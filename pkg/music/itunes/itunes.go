package itunes

import (
	"fmt"
	"os"
	"path/filepath"
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

		location := files.ConvertUrlFilePath(t.Location)

		if strings.HasPrefix(t.Location, "http") {
			// ignore URLs
			continue
		} else if !strings.HasPrefix(t.Location, "file://") {
			tmp := i.newTrack(t)
			i.MoveTrack(tmp, tmp.FilePath())
		}

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
		i.writer.close()
	}
	logrus.Info("iTunes library closed")
}

func (i *Library) SupportedExtensions() music.FileExtensions {
	return music.FileExtensions{
		".aac",
		".aiff",
		".mp3",
		".wav",
		".m4a",
	}
}

func (i *Library) AddFile(path string) (music.Track, error) {
	if filepath.Ext(path) == ".flac" {
		logrus.Warnf("flac file not supported in itunes, ignoring file '%s'", path)
		return nil, nil
	}

	writer := i.getCreateWriter()

	track, err := writer.addFile(path)
	if err != nil {
		return nil, err
	}
	if track == nil {
		return nil, nil
	}

	return &Track{
		itrack: *track,
		lib: i,
	}, nil
}

func (i *Library) CreatePlaylist(path string) (music.Tracklist, error) {
	writer := i.getCreateWriter()
	playlist, err := writer.createPlaylist(path)
	if err != nil {
		return nil, err
	}

	return &Playlist{
		lib:   i,
		plist: *playlist,
	}, nil
}

func (i *Library) CreateCrate(path string) (music.Tracklist, error) {
	return nil, errors.New("crates don't exists in ITunes, use playlists")
}

func (i *Library) MoveTrack(track music.Track, newpath string) error {
	itrack := track.(*Track)
	writer := i.getCreateWriter()
	return writer.setLocation(itrack.itrack.PersistentID, newpath)
}

func (i *Library) Track(filename string) music.Track {
	filename = files.NormalizePath(filename)
	filename = files.RemoveAccent(filename)
	if t, ok := i.trackByLocation[filename]; ok {
		return t
	}
	return nil
}

func (i *Library) Matches(track music.Track) (matches music.Tracks) {
	if track == nil {
		return
	}

	// match for the same filename
	if found := i.Track(track.FilePath()); found != nil {
		if found.Size() == track.Size() && files.Exists(found.FilePath()) {
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
		// if !strings.HasPrefix(it.itrack.Location, "file://") {
		// 	continue
		// }
		count++
		if e := fct(count, len(i.trackByLocation), it); e != nil {
			return e
		}
	}
	return nil
}

func (i *Library) Crates() []music.Tracklist {
	logrus.Warnf("crates don't exists in ITunes, only playlists")
	return nil
}

func (i *Library) Playlists() []music.Tracklist {
	out := []music.Tracklist{}
	for _, playlist := range i.itllib.Playlists {
		// only include music playlists
		if !playlist.Master || playlist.DistinguishedKind != 0 || playlist.Name == "Library" || playlist.Folder {
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
