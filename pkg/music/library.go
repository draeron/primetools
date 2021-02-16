package music

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Library interface {
	Close()
	AddFile(path string) (Track, error)
	Track(filename string) Track

	// try to match a track with the same metadata
	Matches(track Track) Tracks

	Playlists() []Tracklist
	Crates() []Tracklist

	CreatePlaylist(path string) (Tracklist, error)
	CreateCrate(path string) (Tracklist, error)

	MoveTrack(track Track, newpath string) error

	ForEachTrack(fct EachTrackFunc) error

	// return a list of file extension supported by this library
	SupportedExtensions() FileExtensions

	fmt.Stringer
}

type FileExtensions []string

type EachTrackFunc func(index int, total int, track Track) error


func (f FileExtensions) Contains(file string) bool {
	ext := filepath.Ext(file)
	for _, it := range f {
		if strings.ToLower(it)== strings.ToLower(ext) {
			return true
		}
	}
	return false
}
