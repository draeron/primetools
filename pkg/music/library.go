package music

import (
	"fmt"
)

type Library interface {
	Close()
	AddFile(path string) error
	Track(filename string) Track

	// try to match a track with the same metadata
	Matches(track Track) Tracks

	Playlists() []Tracklist
	Crates() []Tracklist

	MoveTrack(track Track, newpath string) error

	ForEachTrack(fct EachTrackFunc) error

	fmt.Stringer
}

type EachTrackFunc func(index int, total int, track Track) error

