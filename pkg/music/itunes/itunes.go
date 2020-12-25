package itunes

import (
	"fmt"
	"html"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dhowden/itl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type Itunes struct {
	itllib        *itl.Library
	tracks        map[string]itl.Track
	trackPerId    map[int]itl.Track
	playlistPerId map[string]itl.Playlist
	writer        *writer
	info          string
}

func Open(path string) (music.Library, error) {
	i := &Itunes{
		tracks:     map[string]itl.Track{},
		trackPerId: map[int]itl.Track{},
		playlistPerId: map[string]itl.Playlist{},
	}

	if path == "" {
		path = files.ExpandHomePath("~/Music/iTunes/iTunes Music Library.xml")
	}

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
		i.tracks[normalizePath(t.Location)] = t
		i.trackPerId[t.TrackID] = t
	}

	for _, p := range xml.Playlists {
		i.playlistPerId[p.PlaylistPersistentID] = p
	}

	i.info = fmt.Sprintf("iTunes: App Version: %v, Lib Version: %v.%v, Track Count: %d", xml.ApplicationVersion, xml.MajorVersion, xml.MinorVersion, len(xml.Tracks))
	logrus.Infof("sucessfully loaded itunes library in %s", time.Since(start))
	logrus.Info(i)

	i.itllib = &xml

	return i, nil
}

func (i *Itunes) Close() {
	if i.writer != nil {
		i.writer.Close()
	}
	logrus.Info("iTunes library closed")
}

func (i *Itunes) AddFile(path string) error {
	return i.writer.addFile(path)
}

func (i *Itunes) Track(filename string) music.Track {
	if t, ok := i.tracks[files.NormalizePath(filename)]; ok {
		return i.newTrack(t)
	}
	return nil
}

func (i *Itunes) ForEachTrack(fct music.EachTrackFunc) error {
	count := 0
	for _, it := range i.tracks {
		if !strings.HasPrefix(it.Location, "file://") {
			continue
		}
		count++
		if e := fct(count, len(i.tracks), i.newTrack(it)); e != nil {
			return e
		}
	}
	return nil
}

func (i *Itunes) Crates() []music.Tracklist {
	// in itunes, a crate is also a playlist
	return i.Playlists()
}

func (i *Itunes) Playlists() []music.Tracklist {
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

func (i *Itunes) getCreateWriter() *writer {
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

func (i *Itunes) String() string {
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
