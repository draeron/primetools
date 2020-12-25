package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type Files struct {
	basePath string
	cache    map[string]*Track
}

func Open(path string) music.Library {
	f := Files{
		basePath: path,
		cache:    map[string]*Track{},
	}

	logrus.Infof("scanning files from folder '%s'", path)
	return f
}

func (f Files) Close() {
}

func (f Files) Crates() []music.Tracklist {
	logrus.Warn("there cannot be any crates in files")
	return nil
}

func (f Files) Playlists() []music.Tracklist {
	logrus.Warn("there cannot be any playlist in files")
	return nil
}

func (f Files) AddFile(path string) error {
	logrus.Warn("cannot add files to files, nothing to do")
	return nil
}

func (f Files) Track(filename string) music.Track {
	if !files.Exists(filename) {
		return nil
	}

	ppath := files.NormalizePath(filename)

	if cached, ok := f.cache[ppath]; ok {
		return cached
	}

	t := &Track{path: ppath}
	f.cache[ppath] = t
	return t
}

func (f Files) ForEachTrack(fct music.EachTrackFunc) error {
	paths := []string{}
	err := filepath.Walk(f.basePath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".mp3" && filepath.Ext(path) != ".flac" {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	for i, it := range paths {
		if e := fct(i, len(paths), newTrack(it)); e != nil {
			return e
		}
	}
	return err
}

func (f Files) String() string {
	return fmt.Sprintf("Files: %s", f.basePath)
}
