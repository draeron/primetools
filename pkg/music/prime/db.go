package prime

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type PrimeDB struct {
	UUID     string
	origin   string
	info     string
	sql      *sqlx.DB
	total    int
	trackIds map[string]trackEntry
	lib      *Library
}

func OpenDB(path string, lib *Library) (*PrimeDB, error) {
	p := &PrimeDB{
		trackIds: map[string]trackEntry{},
		lib:      lib,
	}

	logrus.Infof("opening PRIME database located at '%s'", path)

	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	p.sql = db

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	p.origin, _ = filepath.Abs(filepath.Dir(path))
	p.origin = files.NormalizePath(p.origin)

	info := infoEntry{}
	err = db.Unsafe().Get(&info, `SELECT * from Information LIMIT 1`)
	if err != nil {
		db.Close()
		return nil, errors.Wrapf(err, "failed to fetch prime database information")
	}
	p.UUID = info.UUID

	if err = p.buildIdsMap(); err != nil {
		p.Close()
		return nil, err
	}

	p.info = fmt.Sprintf("PRIME: Database Version: %v.%v.%v, Track Count: %d, Path: %v", info.SchemaVersionMajor, info.SchemaVersionMinor, info.SchemaVersionPatch, len(p.trackIds), p.origin)
	logrus.Info(p.info)

	return p, nil
}

func (l *PrimeDB) IsExported() bool {
	query := `select count(trackId) from CopiedTrack`
	count := 0
	err := l.sql.Get(&count, query)
	if err != nil {
		logrus.Errorf("could not determine if PRIME src '%s' is an export: %v", l.origin, err)
	}
	return count != 0
}

func (l *PrimeDB) fetchList(listType ListType) []music.Tracklist {
	list, err := l.fetchListEntries(listType)
	if err != nil {
		logrus.Errorf("failed to fetch crates from PRIME database: %v", err)
		return nil
	}

	out := []music.Tracklist{}
	for _, it := range list {
		out = append(out, newList(l, it))
	}
	return out
}

func (l *PrimeDB) fetchListEntries(listType ListType) ([]listEntry, error) {
	list := []listEntry{}

	query := `SELECT id, type, title, path, trackCount FROM List WHERE type = ? AND isFolder = false`
	err := l.sql.Select(&list, query, listType)
	if err != nil {
		return nil, errors.Wrapf(err, "")
	}

	return list, err
}

func (l *PrimeDB) Close() {
	if l.sql != nil {
		l.sql.Close()
	}
	logrus.Infof("PRIME library '%s' closed", l.origin)
}

func (l *PrimeDB) buildIdsMap() error {
	entries := []trackEntry{}
	err := l.sql.Unsafe().Select(&entries, `SELECT * from Track`)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch track ids")
	}
	for _, e := range entries {
		fpath := e.Path.String
		if !filepath.IsAbs(fpath) {
			fpath = files.NormalizePath(l.origin + "/" + fpath)
		}
		fpath = files.RemoveAccent(fpath)
		if _, ok := l.trackIds[fpath]; ok {
			logrus.Warnf("duplicate entry in sqlite for path '%s'", fpath)
		}
		l.trackIds[fpath] = e
	}
	return nil
}

func (l *PrimeDB) Track(filename string) music.Track {
	// if filename == "m:/techno/-= ambient =-/arutani/arutani - the mermaid girl ft. ăvem.mp3" {
	// 	println("qawewqeq")
	// }

	entry, ok := l.trackIds[files.RemoveAccent(filename)]
	if !ok {
		return nil
	}

	query := `SELECT * FROM Track WHERE id = ?`
	err := l.sql.Unsafe().Get(&entry, query, entry.Id)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logrus.Errorf("failed query: %v", err)
		}
		return nil
	} else {
		return newTrack(l, entry)
	}
}

func (l *PrimeDB) Tracks() ([]music.Track, error) {
	entries := []*trackEntry{}
	query := `SELECT * FROM Track`
	err := l.sql.Unsafe().Select(&entries, query)
	if err != nil {
		return nil, errors.Wrapf(err, "query '%s' failed", query)
	}

	tracks := []music.Track{}

	for _, e := range entries {
		tracks = append(tracks, newTrack(l, *e))
	}

	return tracks, nil
}
