package prime

import (
	"database/sql"
)

type trackEntry struct {
	Id        int           `db:"id"`
	Length    int           `db:"length"`
	BPM       sql.NullInt32 `db:"bpm"`
	Year      sql.NullInt32 `db:"year"`
	Path      string        `db:"path"`
	Filename  string        `db:"filename"`
	Bitrate   int           `db:"bitrate"`
	TrackType int           `db:"trackType"`
}

type infoEntry struct {
	Id                 int `db:"id"`
	SchemaVersionMajor int `db:"schemaVersionMajor"`
	SchemaVersionMinor int `db:"schemaVersionMinor"`
	SchemaVersionPatch int `db:"schemaVersionPatch"`
}

type metaStringEntry struct {
	Id   int            `db:"id"`
	Type MetaStringType `db:"type"`
	Text sql.NullString `db:"text"`
}

type metaIntEntry struct {
	Id    int           `db:"id"`
	Type  MetaIntType   `db:"type"`
	Value sql.NullInt64 `db:"value"`
}

type metaStringEntries []metaStringEntry
type metaIntEntries []metaIntEntry

func (m metaStringEntries) Get(typed MetaStringType) string {
	for _, it := range m {
		if it.Type == typed {
			return it.Text.String
		}
	}
	return ""
}

func (m metaIntEntries) Get(typed MetaIntType) int64 {
	for _, it := range m {
		if it.Type == typed {
			return it.Value.Int64
		}
	}
	return 0
}