package itunes

import (
	"github.com/dhowden/itl"

	"primetools/pkg/music"
)

type Track struct {
	itrack itl.Track
	writer *writer
}

func (t Track) Rating() int {
	if !t.itrack.RatingComputed {
		return t.itrack.Rating / 20
	}
	return 0
}

func (t Track) PlayCount() int {
	return t.itrack.PlayCount
}

func (t Track) FilePath() string {
	return convertPath(t.itrack.Location)
}

func (t Track) SetRating(rating music.Rating) error {
	return t.writer.rating(t.itrack.Location, int(rating) * 20)
}

func (t Track) SetPlayCount(count int) error {
	return t.writer.playCount(t.itrack.Location, count)
}
