package prime

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type Prime struct {
	origin   string
	info     string
	sql      *sqlx.DB
	total    int
	trackIds map[string]int
}

func Open(path string) (music.Library, error) {
	p := &Prime{
		trackIds: map[string]int{},
	}

	if path == "" {
		path = files.ExpandHomePath("~/Music/Engine Library/m.db")
	} else if files.IsDir(path) && !strings.Contains(path, "Engine Library") {
		if !strings.HasSuffix(path, "/") || !strings.HasSuffix(path, "\\") {
			path = path + "/"
		}
		path += "Engine Library/m.db"
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

	// load file index... because case and shit makes compares sucks balls
	entries := []trackEntry{}
	err = db.Select(&entries, `SELECT id, path from Track`)
	if err != nil {
		db.Close()
		return nil, errors.Wrapf(err, "failed to fetch track ids")
	}
	for _, e := range entries {
		fpath := files.NormalizePath(p.origin + "/" + e.Path.String)
		fpath = files.RemoveAccent(fpath)
		p.trackIds[fpath] = e.Id
	}

	p.info = fmt.Sprintf("PRIME: Database Version: %v.%v.%v, Track Count: %d", info.SchemaVersionMajor, info.SchemaVersionMinor, info.SchemaVersionPatch, len(p.trackIds))
	logrus.Info(p.info)

	return p, nil
}

func (l *Prime) Close() {
	if l.sql != nil {
		l.sql.Close()
	}
	logrus.Info("PRIME library closed")
}

func (l *Prime) Playlists() []music.Tracklist {
	return l.fetchList(ListPlayList)
}

func (l *Prime) Crates() []music.Tracklist {
	return l.fetchList(ListCrate)
}

func (l *Prime) fetchList(listType ListType) []music.Tracklist {
	list, err := l.fetchListEntries(listType)
	if err != nil {
		logrus.Errorf("failed to fetch crates from PRIME database: %v", err)
		return nil
	}

	out := []music.Tracklist{}
	for _, it := range list {
		out = append(out, newList(l.sql, it))
	}
	return out
}

func (l *Prime) fetchListEntries(listType ListType) ([]listEntry, error) {
	list := []listEntry{}

	query := `SELECT id, type, title, path, trackCount FROM List WHERE type = ? AND isFolder = false`
	err := l.sql.Select(&list, query, listType)
	if err != nil {
		return nil, errors.Wrapf(err, "")
	}

	return list, err
}

func (l *Prime) AddFile(path string) error {
	msg := "AddFile is not implemented in Prime library"
	logrus.Warnf(msg)
	return errors.New(msg)
}

func (l *Prime) Track(filename string) music.Track {
	// if filename == "m:/techno/-= ambient =-/arutani/arutani - the mermaid girl ft. ăvem.mp3" {
	// 	println("qawewqeq")
	// }

	id, ok := l.trackIds[files.RemoveAccent(filename)]
	if !ok {
		return nil
	}

	entry := trackEntry{}
	query := `SELECT * FROM Track WHERE id = ?`
	err := l.sql.Unsafe().Get(&entry, query, id)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logrus.Errorf("failed query: %v", err)
		}
		return nil
	} else {
		return newTrack(l.sql, entry)
	}
}

func (l *Prime) ForEachTrack(fct music.EachTrackFunc) error {
	query := `SELECT * FROM Track`
	rows, err := l.sql.Unsafe().Query(query)
	if err != nil {
		return errors.Wrapf(err, "query '%s' failed", query)
	}

	cur := 0

	for rows.Next() {
		entry := trackEntry{}
		rows.NextResultSet()
		err = rows.Scan(&entry)
		if err != nil {
			return errors.Wrap(err, "scan failed")
		}

		t := newTrack(l.sql, entry)

		err = fct(l.total, cur, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Prime) String() string {
	return l.info
}
