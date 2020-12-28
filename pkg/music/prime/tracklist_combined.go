package prime

import (
	"primetools/pkg/music"
)

type CombinedTrackList struct {
	name string
	path string
	lists []music.Tracklist
}

func newCombinedTracklist(list music.Tracklist) *CombinedTrackList {
	c := CombinedTrackList{
		name: list.Name(),
		path: list.Path(),
		lists: []music.Tracklist{ list },
	}
	return &c
}

func (c *CombinedTrackList) Name() string {
	return c.name
}

func (c *CombinedTrackList) Path() string {
	return c.path
}

func (c *CombinedTrackList) MergeWith(list music.Tracklist) *CombinedTrackList {
	if list.Path() == c.path {
		c.lists = append(c.lists, list)
	}
	return c
}

func (c *CombinedTrackList) Tracks() []music.Track {
	tracks := []music.Track{}
	for _, pl := range c.lists {
		tracks = append(tracks, pl.Tracks()...)
	}
	return tracks
}

