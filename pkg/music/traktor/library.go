package traktor

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"primetools/pkg/files"
	"primetools/pkg/music"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Library struct {
	xml        XmlLibrary
	info       string
	pathHashes map[string]*XmlTrack
}

func Open(path string) (music.Library, error) {
	start := time.Now()

	logrus.Infof("opening traktor xml '%s'", path)

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to open file %s", path)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)

	xmllib := XmlLibrary{}

	err = decoder.Decode(&xmllib)
	if err != nil {
		return nil, errors.WithMessagef(err, "parsing of file '%s' failed", path)
	}

	lib := &Library{
		xml:        xmllib,
		info:       fmt.Sprintf("%v: Version: %v, Company: %s, Track Count: %d", xmllib.Header.Program, xmllib.Version, xmllib.Header.Company, xmllib.Collection.Count),
		pathHashes: map[string]*XmlTrack{},
	}

	for _, it := range xmllib.Collection.Entries {
		lib.pathHashes[it.Filepath()] = &it
	}

	logrus.Infof("sucessfully loaded traktor library in %s", time.Since(start))

	return lib, nil
}

func (l Library) Close() {}

func (l Library) Track(filename string) music.Track {
	filename = files.NormalizePath(filename)

	if track, ok := l.pathHashes[filename]; ok && track != nil {
		return newTrack(*track)
	}

	return nil
}

func (l Library) Matches(track music.Track) music.Tracks {
	// TODO implement me
	panic("implement me")
}

func (l Library) Playlists() []music.Tracklist {
	// TODO implement me
	panic("implement me")
}

func (l Library) Crates() []music.Tracklist {
	// TODO implement me
	panic("implement me")
}

func (l Library) CreatePlaylist(path string) (music.Tracklist, error) {
	// TODO implement me
	panic("implement me")
}

func (l Library) CreateCrate(path string) (music.Tracklist, error) {
	// TODO implement me
	panic("implement me")
}

func (l Library) MoveTrack(track music.Track, newpath string) error {
	return errors.New("not supported")
}

func (l Library) ForEachTrack(fct music.EachTrackFunc) error {
	for idx, track := range l.xml.Collection.Entries {
		if err := fct(idx, l.xml.Collection.Count, newTrack(track)); err != nil {
			return err
		}
	}
	return nil
}

func (l Library) AddFile(path string) (music.Track, error) {
	return nil, errors.New("not supported")
}

func (l Library) SupportedExtensions() music.FileExtensions {
	// import not supported
	return music.FileExtensions{}
}

func (l Library) String() string {
	return l.info
}
