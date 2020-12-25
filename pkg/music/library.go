package music

import (
	"fmt"
	"time"
)

type Library interface {
	Close()
	AddFile(path string) error
	Track(filename string) Track

	Playlists() []Tracklist
	Crates() []Tracklist

	ForEachTrack(fct EachTrackFunc) error

	fmt.Stringer
}

type EachTrackFunc func(index int, total int, track Track) error

type Track interface {
	Title() string
	Album() string
	Year() int

	Rating() Rating
	SetRating(rating Rating) error

	Modified() time.Time
	SetModified(modified time.Time) error

	Added() time.Time
	SetAdded(added time.Time) error

	PlayCount() int
	SetPlayCount(count int) error

	FilePath() string

	fmt.Stringer
}

type Tracklist interface {
	Name() string
	Path() string
	Tracks() []Track
}
