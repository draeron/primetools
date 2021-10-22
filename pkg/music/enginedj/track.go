package enginedj

import (
	"encoding/json"
	"path/filepath"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type Track struct {
	entry       trackEntry
	src         *EngineDJDB
	mutex       sync.Mutex
}

func newTrack(src *EngineDJDB, entry trackEntry) *Track {
	return &Track{
		src:         src,
		entry:       entry,
	}
}

func (t *Track) Rating() music.Rating {
	return music.Rating(t.entry.Rating.Int32 / 20)
}

func (t *Track) SetRating(rating music.Rating) error {
	t.entry.Rating.Int32 = int32(rating * 20)
	query := ` UPDATE Track SET rating = ? WHERE id = ?`
	_, err := t.src.sql.Exec(query, t.entry.Rating.Int32, t.entry.Id)
	return errors.Wrapf(err, "failed to set rating %v to track '%s'", rating, t.String())
}

func (t *Track) Added() time.Time {
	return t.entry.Added.Time
}

func (t *Track) Modified() time.Time {
	return t.entry.Added.Time
}

func (t *Track) SetAdded(added time.Time) error {
	t.entry.Added.Time = added
	query := ` UPDATE Track SET dateCreated = ? WHERE id = ?`
	_, err := t.src.sql.Exec(query, t.entry.Added.Time, t.entry.Id)
	return errors.Wrapf(err, "failed to set added date %v to track '%s'", added, t.String())
}

func (t *Track) SetModified(added time.Time) error {
	return errors.New("not implemented")
}

func (t *Track) PlayCount() int {
	return 0
}

func (t *Track) SetPlayCount(count int) error {
	msg := "PlayCount is not implemented in EngineDJ library"
	logrus.Warnf(msg)
	return errors.New(msg)
}

func (t *Track) FilePath() string {
	if filepath.IsAbs(t.entry.Path.String) {
		return t.entry.Path.String
	} else {
		return files.NormalizePath(t.src.origin + "/" + t.entry.Path.String)
	}
}

func (t *Track) SetPath(newpath string) error {
	rpath, err := filepath.Rel(t.src.origin, newpath)
	if err != nil {
		rpath = newpath
	}

	return t.runQuery(func(sql *sqlx.DB, trackId int) error {
		return writeFilepath(sql, trackId, rpath)
	})
}

func (t *Track) Title() string {
	return t.entry.Title.String
}

func (t *Track) Album() string {
	return t.entry.Album.String
}

func (t *Track) Year() int {
	return int(t.entry.Year.Int32)
}

func (t *Track) Artist() string {
	return t.entry.Artist.String
}

func (t *Track) String() string {
	if t.entry.Title.String != "" {
		return t.entry.Title.String
	} else {
		return t.entry.Filename.String
	}
}

func (t *Track) Size() int64 {
	return int64(t.entry.Size.Int32)
}

func (t *Track) MarshalYAML() (interface{}, error) {
	return music.NewMarchalTrack(t), nil
}

func (t *Track) MarshalJSON() ([]byte, error) {
	return json.Marshal(music.NewMarchalTrack(t))
}

func (t *Track) MarshalTOML() ([]byte, error) {
	return toml.Marshal(music.NewMarchalTrack(t))
}

func (t *Track) runQuery(fct func(sql *sqlx.DB, trackId int) error) error {
	err := fct(t.src.sql, t.entry.Id)
	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}
	if t.isExternal() {
		if db, ok := t.src.lib.dbs[t.entry.OriginDatabaseUuid.String]; ok {
			return fct(db.sql, int(t.entry.OriginTrackId.Int32))
		} else {
			// todo: log cannot find DB
		}
	}
	return nil
}

func (t *Track) isExternal() bool {
	return t.entry.OriginDatabaseUuid.String == t.src.UUID
}

func writeFilepath(sql *sqlx.DB, trackId int, newpath string) error {
	query := `UPDATE Track SET path = ?, filename = ? WHERE id = ?`
	fname := filepath.Base(newpath)

	_, err := sql.Exec(query, newpath, fname, trackId)
	if err != nil {
		return errors.Wrapf(err, "failed to location of track '%v' in EngineDJ db: %v", trackId, err)
	}
	logrus.Infof("path for '%v' updated in EngineDJ db", trackId)
	return nil
}
