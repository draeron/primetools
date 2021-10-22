package enginedj

import (
	"context"
	"encoding/json"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"primetools/pkg/music"
)

type TrackList struct {
	entry playlistEntry
	src   *EngineDJDB
	path  string
}

func newList(src *EngineDJDB, entry playlistEntry) TrackList {
	return TrackList{
		entry: entry,
		src:   src,
	}
}

func (t TrackList) Name() string {
	return t.entry.Title.String
}

func (t TrackList) Path() string {
	if t.path == "" {
		t.path = path.Join(t.fetchParentPath(), t.entry.Title.String)
	}

	return t.path
}

func (t TrackList) fetchParentPath() string {
	if t.entry.ParentListId.Valid && t.entry.ParentListId.Int32 != 0 {
		list, err := t.src.fetchListWith(int(t.entry.ParentListId.Int32))
		if err != nil {
			logrus.Errorf("failed to fetch playlist: %v", err)
		}
		return list.Path()
	}
	return ""
}

func (t TrackList) String() string {
	names := []string{}
	for _, track := range t.Tracks() {
		names = append(names, track.String())
	}
	return "[" + strings.Join(names, ",") + "]"
}

func (t *TrackList) setParent(parent *TrackList) {
	logrus.Infof("updating parent for playlist '%s'", t.Path())
	var err error

	// delete previous entries
	_, err = t.src.sql.Exec(`DELETE FROM PlaylistAllParent WHERE childListId = ?`, t.entry.Id)
	if err != nil {
		logrus.Errorf("failed to update playlist '%s': %v", t.Path(), err)
		return
	}
	_, err = t.src.sql.Exec(`DELETE FROM ListHierarchy WHERE childIdList = ?`, t.entry.Id)
	if err != nil {
		logrus.Errorf("failed to update playlist '%s': %v", t.Path(), err)
		return
	}

	if parent == nil {
		parent = t
	}

	_, err = t.src.sql.Exec(`INSERT INTO ListParentList (listOriginId, listParentId) VALUES (?,?)`,
		t.entry.Id, parent.entry.Id)
	if err != nil {
		logrus.Errorf("failed to update playlist '%s': %v", t.Path(), err)
		return
	}

	if t.entry.Id != parent.entry.Id {
		_, err = t.src.sql.Exec(`INSERT INTO ListHierarchy (listId, listIdChild, listTypeChild) VALUES (?,?)`,
			parent.entry.Id, t.entry.Id)
		if err != nil {
			logrus.Errorf("failed to update playlist '%s': %v", t.Path(), err)
		}
	}
}

func (t *TrackList) SetTracks(tracks music.Tracks) error {
	logrus.Infof("updating tracklist for playlist '%s' in db '%s' with %d entries", t.Path(), t.src.origin, len(tracks))

	tx, err := t.src.sql.BeginTxx(context.TODO(), nil)
	if err != nil {
		return errors.Wrapf(err, "failed start db transaction for playlist '%s'", t.Path())
	}

	query := `DELETE FROM ListTrackList WHERE listId = ?`
	_, err = tx.Exec(query, t.entry.Id)
	if err != nil {
		tx.Rollback()
		return errors.Wrapf(err, "failed deleting previous track from playlist '%s'", t.Path())
	}

	for idx, track := range tracks {
		tr, ok := track.(*Track)
		if !ok {
			panic("feed tracks from the same library")
		}

		query = `INSERT INTO ListTrackList (listId, trackId, trackIdInOriginDatabase, databaseUuid, trackNumber) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(query, t.entry.Id, tr.entry.Id, tr.entry.OriginTrackId, tr.entry.OriginDatabaseUuid, idx+1)

		if err != nil {
			tx.Rollback()
			return errors.Wrapf(err, "failed to add track to playlist '%s'", t.Path())
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrapf(err, "fail to commit transaction for playlist '%s'", t.Path())
	}
	return nil
}

func (t *TrackList) Count() int {
	count := 0
	query := `SELECT COUNT(DISTINCT(trackId)) FROM PlaylistEntity WHERE listId = ?`
	_ = t.src.sql.Get(&count, query, t.entry.Id)
	return count
}

func (t *TrackList) Tracks() music.Tracks {
	tracks := []trackEntry{}
	query := `select * from PlaylistEntity JOIN Track ON Track.id = PlaylistEntity.trackId WHERE PlaylistEntity.listId = ? ORDER BY playOrder`

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
