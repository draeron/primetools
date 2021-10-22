package enginedj

import (
	"encoding/json"

	"github.com/pkg/errors"

	"primetools/pkg/music"
)

type Playlist struct {
	name  string
	path  string
	lists []TrackList
	src   *Library
}

func newCrate(list TrackList, src *Library) *Playlist {
	c := Playlist{
		name:  list.Name(),
		path:  list.Path(),
		lists: []TrackList{list},
		src:   src,
	}
	return &c
}

func (p *Playlist) Name() string {
	return p.name
}

func (p *Playlist) Path() string {
	return p.path
}

func (p *Playlist) MergeWith(list TrackList) *Playlist {
	if list.Path() == p.path {
		p.lists = append(p.lists, list)
	}
	return p
}

func (p *Playlist) MarshalYAML() (interface{}, error) {
	return music.NewMarshallTracklist(p), nil
}

func (p *Playlist) MarshalJSON() ([]byte, error) {
	return json.Marshal(music.NewMarshallTracklist(p))
}

/*
	For each track, split them based on their origin DB
*/
func (p *Playlist) SetTracks(tracks music.Tracks) error {

	savetrack := func(db *EngineDJDB, tracks music.Tracks) error {
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

		list, err := createListIn(db, p.path)
		if err != nil {
			return err
		}

		if list == nil {
			return errors.Errorf("failed to create the crate '%s' into PRIME db '%s'", p.path, db.origin)
		}

		return list.SetTracks(newlist)
	}

	err := savetrack(p.src.main, tracks)
	if err != nil {
		return err
	}

	for _, db := range p.src.dbs {
		err = savetrack(db, tracks)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Playlist) Count() int {
	count := 0
	for _, list := range p.lists {
		count += list.Count()
	}
	return count
}

func (p *Playlist) Tracks() music.Tracks {
	tracks := []music.Track{}
	for _, pl := range p.lists {
		tracks = append(tracks, pl.Tracks()...)
	}
	return tracks
}
