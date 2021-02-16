package prime

import (
	"encoding/json"

	"github.com/pkg/errors"

	"primetools/pkg/music"
)

type Crate struct {
	name  string
	path  string
	lists []music.Tracklist
	src   *Library
}

func newCrate(list music.Tracklist, src *Library) *Crate {
	c := Crate{
		name:  list.Name(),
		path:  list.Path(),
		lists: []music.Tracklist{list},
		src:   src,
	}
	return &c
}

func (c *Crate) Name() string {
	return c.name
}

func (c *Crate) Path() string {
	return c.path
}

func (c *Crate) MergeWith(list music.Tracklist) *Crate {
	if list.Path() == c.path {
		c.lists = append(c.lists, list)
	}
	return c
}

func (t *Crate) MarshalYAML() (interface{}, error) {
	return music.NewMarshallTracklist(t), nil
}

func (t *Crate) MarshalJSON() ([]byte, error) {
	return json.Marshal(music.NewMarshallTracklist(t))
}

/*
	For each track, split them based on their origin DB
*/
func (c *Crate) SetTracks(tracks music.Tracks) error {

	savetrack := func(db *PrimeDB, tracks music.Tracks) error {
		newlist := []music.Track{}

		for _, it := range tracks {
			track, ok := it.(*Track)
			if !ok {
				return errors.New("cannot save track object which are not from the same library")
			}

			if track.src == db {
				newlist = append(newlist, track)
			}
		}

		list, err := createListIn(db, c.path, ListCrate)
		if err != nil {
			return err
		}

		if list == nil {
			return errors.Errorf("failed to create the crate '%s' into PRIME db '%s'", c.path, db.origin)
		}

		return list.SetTracks(newlist)
	}

	err := savetrack(c.src.main, tracks)
	if err != nil {
		return err
	}

	for _, db := range c.src.dbs {
		err = savetrack(db, tracks)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Crate) Tracks() music.Tracks {
	tracks := []music.Track{}
	for _, pl := range c.lists {
		tracks = append(tracks, pl.Tracks()...)
	}
	return tracks
}
