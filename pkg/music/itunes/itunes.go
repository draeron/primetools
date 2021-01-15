package itunes

import (
	"fmt"
	"html"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dhowden/itl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type Library struct {
	itllib          *itl.Library
	trackByLocation map[string]*Track
	trackById       map[int]*itl.Track
	playlistPerId   map[string]*itl.Playlist
	metaHashes      map[string]*Track
	writer          itunes_writer
	info            string
	mutex           sync.Mutex
}

func Open(path string) (music.Library, error) {
	i := &Library{
		trackByLocation: map[string]*Track{},
		trackById:       map[int]*itl.Track{},
		metaHashes:      map[string]*Track{},
		playlistPerId:   map[string]*itl.Playlist{},
	}

	if path == "" {
		path = files.ExpandHomePath("~/Music/iTunes/iTunes Music Library.xml")
	}
	path = files.NormalizePath(path)

	logrus.Info("opening iTunes xml...")
	logrus.Infof("library resolved at '%s'", path)

	start := time.Now()
	ifile, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open itunes xml file: %v", err)
	}

	xml, err := itl.ReadFromXML(ifile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read itunes xml file: %v", err)
	}

	for _, t := range xml.Tracks {
		if t.Location == "" {
			continue
		}

		location := normalizePath(t.Location)
		if _, ok := i.trackByLocation[location]; ok {
			logrus.Warnf("file '%s' seems to be duplicated in itunes xml", location)
		}
		i.trackByLocation[location] = i.newTrack(t)
		i.trackById[t.TrackID] = &t
	}

	for _, p := range xml.Playlists {
		i.playlistPerId[p.PlaylistPersistentID] = &p
	}

	i.info = fmt.Sprintf("iTunes: App Version: %v, Lib Version: %v.%v, Track Count: %d", xml.ApplicationVersion, xml.MajorVersion, xml.MinorVersion, len(xml.Tracks))
	logrus.Infof("sucessfully loaded itunes library in %s", time.Since(start))
	logrus.Info(i)

	i.itllib = &xml

	return i, nil
}

func (i *Library) Close() {
	if i.writer != nil {
		i.writer.Close()
	}
	logrus.Info("iTunes library closed")
}

func (i *Library) AddFile(path string) error {
	return i.writer.addFile(path)
}

func (i *Library) MoveTrack(track music.Track, newpath string) error {
	itrack := track.(*Track)
	writer := i.getCreateWriter()
	return writer.setLocation(itrack.itrack.PersistentID, newpath)
}

func (i *Library) Track(filename string) music.Track {
	if t, ok := i.trackByLocation[files.NormalizePath(filename)]; ok {
		return t
	}
	return nil
}

func (i *Library) Matches(track music.Track) (matches music.Tracks) {
	if track == nil {
		return
	}

	// match for the same filename
	if found := i.Track(track.FilePath()); track != nil {
		if files.Exists(found.FilePath()) {
			matches = append(matches, found)
		}
	}

	i.mutex.Lock()
	defer i.mutex.Unlock()

	if len(i.metaHashes) == 0 {
		start := time.Now()
		logrus.Infof("constructing track meta cache")

		for _, it := range i.trackByLocation {
			hash := music.TrackHash(it)
			i.metaHashes[hash] = it
		}
		logrus.Infof("processed metadata for %d tracks in %s", len(i.trackByLocation), time.Since(start))
	}

	hash := music.TrackHash(track)

	if found, ok := i.metaHashes[hash]; ok {
		matches = append(matches, found)
	}

	return matches.Dedupe()
}

func (i *Library) ForEachTrack(fct music.EachTrackFunc) error {
	count := 0
	for _, it := range i.trackByLocation {
		if !strings.HasPrefix(it.itrack.Location, "file://") {
			continue
		}
		count++
		if e := fct(count, len(i.trackByLocation), it); e != nil {
			return e
		}
	}
	return nil
}

func (i *Library) Crates() []music.Tracklist {
	// in itunes, a crate is also a playlist
	return i.Playlists()
}

func (i *Library) Playlists() []music.Tracklist {
	out := []music.Tracklist{}
	for _, playlist := range i.itllib.Playlists {
		// only include music playlists
		if playlist.DistinguishedKind != 0 || playlist.Name == "Library" || playlist.Folder {
			continue
		}

		out = append(out, Playlist{
			plist: playlist,
			lib:   i,
		})
	}
	return out
}

func (i *Library) getCreateWriter() itunes_writer {
	if i.writer != nil {
		return i.writer
	}

	var err error
	i.writer, err = createWriter()
	if err != nil {
		logrus.Errorf("failed to init iTunes writer interface, writes operations will fail: %v", err)
		panic(err)
	}
	return i.writer
}

func (i *Library) String() string {
	return i.info
}

// file://localhost/m:/Techno/-=%20Ambient%20=-/Bluetech/2005%20-%20Sines%20And%20Singularities/01%20-%20Enter%20The%20Lovely.mp3
func normalizePath(path string) string {
	path = strings.Replace(path, "file://localhost/", "", 1)
	path, _ = url.PathUnescape(path)
	path = html.UnescapeString(path)
	path = files.NormalizePath(path)
	return path
}
