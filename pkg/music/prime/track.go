package prime

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type Track struct {
	entry       trackEntry
	metaStrings metaStringEntries
	metaInts    metaIntEntries
	src         *PrimeDB
	mutex       sync.Mutex
}

func newTrack(src *PrimeDB, entry trackEntry) *Track {
	return &Track{
		src:         src,
		entry:       entry,
		metaStrings: metaStringEntries{},
		metaInts:    metaIntEntries{},
	}
}

func (t *Track) Rating() music.Rating {
	t.readMetaInts()
	rate := t.metaInts.Get(MetaRating)
	return music.Rating(rate / 20)
}

func (t *Track) SetRating(rating music.Rating) error {
	return t.writeMetaIntCascade(MetaRating, int64(rating*20))
}

func (t *Track) Added() time.Time {
	t.readMetaInts()
	atime := t.metaInts.Get(MetaAdded)
	if atime == 0 {
		return time.Now().UTC()
	} else {
		return time.Unix(atime, 0).UTC()
	}
}

func (t *Track) Modified() time.Time {
	t.readMetaInts()
	atime := t.metaInts.Get(MetaCreated)
	if atime == 0 {
		return time.Now().UTC()
	} else {
		return time.Unix(atime, 0).UTC()
	}
}

func (t *Track) SetAdded(added time.Time) error {
	return t.writeMetaIntCascade(MetaAdded, added.Unix())
}

func (t *Track) SetModified(added time.Time) error {
	return t.writeMetaIntCascade(MetaCreated, added.Unix())
}

func (t *Track) PlayCount() int {
	// logrus.Warnf("PlayCount is not implemented in Library library")
	return 0
}

func (t *Track) SetPlayCount(count int) error {
	msg := "PlayCount is not implemented in Library library"
	logrus.Warnf(msg)
	return errors.New(msg)
}

func (t *Track) FilePath() string {
	return files.NormalizePath(t.src.origin + "/" + t.entry.Path.String)
}

func (t *Track) Title() string {
	t.readMetaString()
	return t.metaStrings.Title()
}

func (t *Track) Album() string {
	t.readMetaString()
	return t.metaStrings.Album()
}

func (t *Track) Year() int {
	return int(t.entry.Year.Int32)
}

func (t *Track) String() string {
	t.readMetaString()
	if title := t.metaStrings.Get(MetaTitle); title != "" {
		return title
	} else {
		return t.entry.Filename.String
	}
}

func (t *Track) writeMetaIntCascade(meta MetaIntType, value int64) error {
	if err := writeMetaInt(t.src.sql, t.entry.Id, meta, value); err != nil {
		return err
	}
	if t.entry.External.Bool && t.entry.ExternalDbId.Valid && t.entry.ExternalId.Valid {
		if db, ok := t.src.lib.dbs[t.entry.ExternalDbId.String]; ok {
			return writeMetaInt(db.sql, int(t.entry.ExternalId.Int32), meta, value)
		} else {
			// todo: log cannot find DB
		}
	}
	return nil
}

func writeMetaInt(sql *sqlx.DB, trackId int, meta MetaIntType, value int64) error {
	// query := `UPDATE MetaDataInteger SET value = ? WHERE id = ? AND type = ?`
	query := `INSERT OR REPLACE INTO MetaDataInteger (value, id, type) VALUES (?, ?, ?)`

	res, err := sql.Exec(query, value, trackId, meta)
	if err != nil {
		return errors.Wrapf(err, "running query '%s'", query)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "getting affected row count")
	}
	logrus.Debugf("Query '%s' row affected: %d", query, count)

	return nil
}

func (t *Track) Size() int64 {
	return int64(t.entry.Size.Int32)
}

func (t *Track) MarshalYAML() (interface{}, error) {
	return music.TrackToMarshalObject(t), nil
}

func (t *Track) MarshalJSON() ([]byte, error) {
	return json.Marshal(music.TrackToMarshalObject(t))
}

func (t *Track) readMetaString() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if len(t.metaStrings) > 0 {
		return
	}

	query := `select * from MetaData WHERE id = ?`
	err := t.src.sql.Select(&t.metaStrings, query, t.entry.Id)
	if err != nil {
		logrus.Errorf("failed to read meta strings from sqlite: %v", err)
	}
}

func (t *Track) readMetaInts() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if len(t.metaInts) > 0 {
		return
	}

	query := `select * from MetaDataInteger WHERE id = ?`
	err := t.src.sql.Select(&t.metaInts, query, t.entry.Id)
	if err != nil {
		logrus.Errorf("failed to read meta ints from sqlite: %v", err)
	}
}
