package itunes

import (
	"time"

	"github.com/dhowden/itl"
	"github.com/pkg/errors"

	"primetools/pkg/music"
)

type Track struct {
	itrack itl.Track
	writer *writer
}

func (t Track) String() string {
	return t.itrack.Name
}

func (t Track) PlayCount() int {
	return t.itrack.PlayCount
}

func (t Track) Rating() music.Rating {
	if !t.itrack.RatingComputed {
		return music.Rating(t.itrack.Rating / 20)
	}
	return 0
}

func (t Track) FilePath() string {
	return normalizePath(t.itrack.Location)
}

func (t Track) Added() time.Time {
	return t.itrack.DateAdded
}

func (t Track) SetAdded(added time.Time) error {
	return errors.New("cannot set added date in iTunes")
}

func (t Track) SetRating(rating music.Rating) error {
	return t.writer.setRating(t.itrack.PersistentID, int(rating)*20)
}

func (t Track) SetPlayCount(count int) error {
	return t.writer.setPlayCount(t.itrack.PersistentID, count)
}
