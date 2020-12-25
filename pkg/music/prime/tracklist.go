package prime

import (
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"primetools/pkg/music"
)

type TrackList struct {
	entry listEntry
	sql   *sqlx.DB
}

func newList(sql *sqlx.DB, entry listEntry) *TrackList {
	return &TrackList{
		entry: entry,
		sql: sql,
	}
}

func (t TrackList) Name() string {
	return t.entry.Title
}

func (t TrackList) Path() string {
	path := strings.Replace(t.entry.Path.String, ";", "/", -1)
	path = strings.TrimSuffix(path, "/")
	return path
}

func (t *TrackList) Tracks() []music.Track {

	tracks := []trackEntry{}
	query := `select * from ListTrackList join Track ON Track.id = ListTrackList.trackId WHERE listId = ? ORDER BY trackNumber`

	err := t.sql.Unsafe().Select(&tracks, query, t.entry.Id)
	if err != nil {
		logrus.Errorf("fail to fetch track list for playlist '%s': %v", t.Name(), err)
		return nil
	}

	out := []music.Track{}
	for _, it := range tracks {
		out = append(out, newTrack(t.sql, it))
	}
	return out
}

