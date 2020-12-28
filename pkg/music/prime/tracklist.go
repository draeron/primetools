package prime

import (
	"strings"

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
	path := strings.Replace(t.entry.Path.String, ";", "/", -1)
	path = strings.TrimSuffix(path, "/")
	return path
}

func (t *TrackList) Tracks() []music.Track {
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
