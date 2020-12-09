package music

import (
	"fmt"
	"time"
)

type Library interface {
	Close()
	AddFile(path string) error
	Track(filename string) Track

	// Channel to iterate through all tracks
	ForEachTrack(fct EachTrackFunc) error

	fmt.Stringer
}

type EachTrackFunc func(index int, total int, track Track) error

type Track interface {
	Rating() Rating
	SetRating(rating Rating) error

	Added() time.Time
	SetAdded(added time.Time) error

	PlayCount() int
	SetPlayCount(count int) error

	FilePath() string

	fmt.Stringer
}

type Playlist interface {
	Tracks() []Track
}
