package prime

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"primetools/pkg/music"
)

type TrackList struct {
	entry listEntry
	src   *PrimeDB
}

func newList(src *PrimeDB, entry listEntry) *TrackList {
	return &TrackList{
		entry: entry,
		src:   src,
	}
}

func (t TrackList) Name() string {
	return t.entry.Title
}

func (t TrackList) Path() string {
	path := strings.Replace(t.entry.Path.String, "/", "|", -1)
	path = strings.Replace(t.entry.Path.String, ";", "/", -1)
	path = strings.TrimSuffix(path, "/")
	return path
}

func (t TrackList) String() string {
	names := []string{}
	for _, track := range t.Tracks() {
		names = append(names, track.String())
	}
	return "[" + strings.Join(names, ",") + "]"
}

func (t *TrackList) setParent(parent *TrackList) {
	logrus.Infof("updating parent for %s '%s'", t.entry.Type, t.Path())
	var err error

	// delete previous entries
	_, err = t.src.sql.Exec(`DELETE FROM ListParentList WHERE listOriginId = ?`, t.entry.Id)
	if err != nil {
		logrus.Errorf("failed to update %v '%s': %v", t.entry.Type, t.Path(), err)
		return
	}
	_, err = t.src.sql.Exec(`DELETE FROM ListHierarchy WHERE listIdChild = ?`, t.entry.Id)
	if err != nil {
		logrus.Errorf("failed to update %v '%s': %v", t.entry.Type, t.Path(), err)
		return
	}

	if parent == nil {
		parent = t
	}

	_, err = t.src.sql.Exec(`INSERT INTO ListParentList (listOriginId, listOriginType, listParentId, listParentType) VALUES (?,?,?,?)`,
		t.entry.Id, t.entry.Type, parent.entry.Id, parent.entry.Type)
	if err != nil {
		logrus.Errorf("failed to update %v '%s': %v", t.entry.Type, t.Path(), err)
		return
	}

	if t.entry.Id != parent.entry.Id {
		_, err = t.src.sql.Exec(`INSERT INTO ListHierarchy (listId, listType, listIdChild, listTypeChild) VALUES (?,?,?,?)`,
			parent.entry.Id, parent.entry.Type, t.entry.Id, t.entry.Type)
		if err != nil {
			logrus.Errorf("failed to update %v '%s': %v", t.entry.Type, t.Path(), err)
		}
	}
}

func (t *TrackList) SetTracks(tracks music.Tracks) error {
	logrus.Infof("updating tracklist for %v '%s' in db '%s' with %d entries", t.entry.Type, t.Path(), t.src.origin, len(tracks))

	tx, err := t.src.sql.BeginTxx(context.TODO(), nil)
	if err != nil {
		return errors.Wrapf(err, "failed start db transaction for %v '%s'", t.entry.Type, t.Path())
	}

	query := `DELETE FROM ListTrackList WHERE listId = ? and listType = ?`
	_, err = tx.Exec(query, t.entry.Id, t.entry.Type)
	if err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "failed deleting previous track from %v '%s'", t.entry.Type, t.Path())
	}

	for idx, track := range tracks {
		tr, ok := track.(*Track)
		if !ok {
			panic("feed tracks from the same library")
		}

		query = `INSERT INTO ListTrackList (listId, listType, trackId, trackIdInOriginDatabase, databaseUuid, trackNumber) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(query, t.entry.Id, t.entry.Type, tr.entry.Id, tr.entry.ExternalId, tr.entry.ExternalDbId, idx+1)
		if err != nil {
			tx.Rollback()
			return errors.Wrapf(err, "failed to add track to %v '%s'", t.entry.Type, t.Path())
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrapf(err, "fail to commit transaction for %v '%s'", t.entry.Type, t.Path())
	}
	return nil
}

func (t *TrackList) Tracks() music.Tracks {
	tracks := []trackEntry{}
	query := `select * from ListTrackList join Track ON Track.id = ListTrackList.trackId WHERE listId = ? ORDER BY trackNumber`

	err := t.src.sql.Unsafe().Select(&tracks, query, t.entry.Id)
	if err != nil {
		logrus.Errorf("fail to fetch track list for playlist '%s': %v", t.Name(), err)
		return nil
	}

	out := []music.Track{}
	for _, it := range tracks {
		out = append(out, newTrack(t.src, it))
	}
	return out
}

func (t *TrackList) MarshalYAML() (interface{}, error) {
	return music.NewMarshallTracklist(t), nil
}

func (t *TrackList) MarshalJSON() ([]byte, error) {
	return json.Marshal(music.NewMarshallTracklist(t))
}
