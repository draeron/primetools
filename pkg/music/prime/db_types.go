package prime

import (
	"database/sql"
)

type trackEntry struct {
	Id           int            `db:"id"`
	Length       sql.NullInt32  `db:"length"`
	BPM          sql.NullInt32  `db:"bpm"`
	Year         sql.NullInt32  `db:"year"`
	Path         sql.NullString `db:"path"`
	Filename     sql.NullString `db:"filename"`
	Bitrate      sql.NullInt32  `db:"bitrate"`
	Size         sql.NullInt32  `db:"fileBytes"`
	External     sql.NullBool   `db:"isExternalTrack"`
	ExternalId   sql.NullInt32  `db:"idTrackInExternalDatabase"`
	ExternalDbId sql.NullString `db:"uuidOfExternalDatabase"`
	// TrackType  int            `db:"trackType"`
}

type listEntry struct {
	Id               int            `db:"id"`
	Type             ListType       `db:"type"`
	Title            string         `db:"title"`
	Path             sql.NullString `db:"path"`
	TrackCount       int            `db:"trackCount"`
	ExplicitExported bool           `db:"isExplicitlyExported"`
	Folder           bool           `db:"isFolder"`
	Ordering         int            `db:"ordering"`
}

type infoEntry struct {
	Id                 int    `db:"id"`
	UUID               string `db:"uuid"`
	SchemaVersionMajor int    `db:"schemaVersionMajor"`
	SchemaVersionMinor int    `db:"schemaVersionMinor"`
	SchemaVersionPatch int    `db:"schemaVersionPatch"`
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

func (m metaStringEntries) Title() string {
	return m.Get(MetaTitle)
}

func (m metaStringEntries) Album() string {
	return m.Get(MetaAlbum)
}

func (m metaStringEntries) Artist() string {
	return m.Get(MetaArtist)
}

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
