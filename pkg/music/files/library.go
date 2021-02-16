package files

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type FileLibrary struct {
	basePath  string
	cache     map[string]*Track
	hashCache map[string]*Track
}

func Open(path string) *FileLibrary {
	f := &FileLibrary{
		basePath:  path,
		cache:     map[string]*Track{},
		hashCache: map[string]*Track{},
	}

	logrus.Infof("file library created from folder '%s'", path)
	return f
}

func (f *FileLibrary) Close() {
}

func (f *FileLibrary) SupportedExtensions() music.FileExtensions {
	return music.FileExtensions{} // no import supported
}

func (f *FileLibrary) Crates() []music.Tracklist {
	logrus.Warn("there cannot be any crates in files")
	return nil
}

func (f *FileLibrary) Playlists() []music.Tracklist {
	logrus.Warn("there cannot be any playlist in files")
	return nil
}

func (f *FileLibrary) AddFile(path string) (music.Track, error) {
	logrus.Warn("cannot add files to files, nothing to do")
	return nil, nil
}

func (f *FileLibrary) CreatePlaylist(path string) (music.Tracklist, error) {
	return nil, errors.New("cannot create playlist on file library")
}

func (f *FileLibrary) CreateCrate(path string) (music.Tracklist, error) {
	return nil, errors.New("cannot create creates on file library")
}

func (l *FileLibrary) MoveTrack(track music.Track, newpath string) error {
	logrus.Warn("cannot move files, nothing to do")
	return nil
}

func (f *FileLibrary) Track(filename string) music.Track {
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

func (f *FileLibrary) Matches(track music.Track) (matches music.Tracks) {
	if track == nil {
		return
	}

	if found := f.Track(track.FilePath()); found != nil {
		if files.Exists(found.FilePath()) {
			matches = append(matches, found)
		}
	}

	hash := music.TrackHash(track)
	if match, ok := f.hashCache[hash]; ok {
		matches = append(matches, match)
	}

	if f.basePath != "" {
		matches = append(matches, f.MatchInDirectory(track, f.basePath)...)
	} else {
		matches = append(matches, f.MatchInDirectory(track, path.Dir(track.FilePath()))...)
	}

	return matches.Dedupe()
}

func (f *FileLibrary) MatchInDirectory(track music.Track, dir string) (matches music.Tracks) {
	hash := music.TrackHash(track)

	logrus.Infof("scanning for music files at '%s'", dir)

	files.WalkMusicFiles(dir, func(path string, directoryEntry *godirwalk.Dirent) error {
		cached := f.cache[path]
		if cached == nil {
			cached = newTrack(path)
			f.cache[path] = cached
		}

		// only check for files with the same file extension
		if filepath.Ext(cached.path) != filepath.Ext(track.FilePath()) {
			return nil
		}

		if files.Size(cached.path) != track.Size() {
			return nil
		}

		// println(track.String(), "\n", music.TrackMeta(track))
		// println(cached.String(), "\n", music.TrackMeta(cached))

		ithash := music.TrackHash(cached)
		f.hashCache[ithash] = cached
		if ithash == hash {
			matches = append(matches, cached)
		}
		return nil
	})

	return
}

func (f *FileLibrary) ForEachTrack(fct music.EachTrackFunc) error {
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

func (f *FileLibrary) String() string {
	return fmt.Sprintf("Files: %s", f.basePath)
}
