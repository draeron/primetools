package enginedj

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type EngineDJDB struct {
	UUID     string
	origin   string
	info     string
	sql      *sqlx.DB
	total    int
	trackIds map[string]trackEntry
	lib      *Library
}

func OpenDB(path string, lib *Library) (*EngineDJDB, error) {
	p := &EngineDJDB{
		trackIds: map[string]trackEntry{},
		lib:      lib,
	}

	logrus.Infof("opening EngineDJ database located at '%s'", path)

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
		return nil, errors.Wrapf(err, "failed to fetch EngineDJ database information")
	}
	p.UUID = info.UUID

	if err = p.buildIdsMap(); err != nil {
		p.Close()
		return nil, err
	}

	p.info = fmt.Sprintf("EngineDJ: Database Version: %v.%v.%v, Track Count: %d, Path: %v", info.SchemaVersionMajor, info.SchemaVersionMinor, info.SchemaVersionPatch, len(p.trackIds), p.origin)
	logrus.Info(p.info)

	return p, nil
}

func (l *EngineDJDB) IsExported() bool {
	query := `select count(*) from Pack`
	count := 0
	err := l.sql.Get(&count, query)
	if err != nil {
		logrus.Errorf("could not determine if EngineDJ src '%s' is an export: %v", l.origin, err)
	}
	return count != 0
}

func (l *EngineDJDB) createList(path string, folder bool) *TrackList {
	split := strings.Split(path, "/")
	name := split[len(split)-1]
	ppath := strings.Replace(path, "/", ";", -1) + ";"

	logrus.Infof("creating playlist '%s in EngineDJ database", path)

	query := `INSERT INTO Playlist (title, path, isFolder, isExplicitlyExported, id) VALUES (?, ?, ?, 1, (SELECT max(id) + 1 FROM List))`
	_, err := l.sql.Exec(query, name, ppath, folder)
	if err != nil {
		logrus.Errorf("failed to create playlist '%s': %v", path, err)
		return nil
	}

	list, err := l.fetchList(path)
	if err != nil {
		logrus.Errorf("failed to create playlist '%s': %v", path, err)
	}
	return list
}

func (l *EngineDJDB) fetchList(path string) (*TrackList, error) {
	ppath := strings.Replace(path, "/", ";", -1) + ";"
	entry := playlistEntry{}
	err := l.sql.Get(&entry, `SELECT * FROM Playlist WHERE path = ?`, ppath)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, errors.Wrapf(err, "fail to fetch playlist '%s'", path)
		} else {
			return nil, nil
		}
	}
	return &TrackList{
		src:   l,
		entry: entry,
	}, nil
}

func (l *EngineDJDB) fetchLists() []TrackList {
	lists, err := l.fetchListEntries()
	if err != nil {
		logrus.Errorf("failed to fetch playlist from EngineDJ database: %v", err)
		return nil
	}

	out := []TrackList{}
	for _, it := range lists {
		out = append(out, newList(l, it))
	}
	return out
}

func (l *EngineDJDB) fetchListWith(id int) (*TrackList, error) {
	list := playlistEntry{}
	query := `SELECT * FROM Playlist WHERE id = ?`
	err := l.sql.Get(&list, query, id)
	if err != nil {
		return nil, errors.Wrapf(err, "")
	}
	return &TrackList{
		entry: list,
		src:   l,
	}, nil
}

func (l *EngineDJDB) fetchListEntries() ([]playlistEntry, error) {
	lists := []playlistEntry{}
	query := `SELECT * FROM Playlist`
	err := l.sql.Select(&lists, query)
	if err != nil {
		return nil, errors.Wrapf(err, "")
	}

	out := []playlistEntry{}
	// filter out folders
	for _, it := range lists {
		children := 0
		query = `SELECT COUNT(*) from PlaylistAllChildren WHERE PlaylistAllChildren.id = ?`
		err = l.sql.Get(&children, query, it.Id)
		if err != nil {
			logrus.Errorf("failed to get child count for playlist '%s'", it.Title)
		}
		if children == 0 {
			out = append(out, it)
		}
	}

	return out, nil
}

func (l *EngineDJDB) Close() {
	if l.sql != nil {
		l.sql.Close()
	}
	logrus.Infof("EngineDJ library '%s' closed", l.origin)
}

func (l *EngineDJDB) buildIdsMap() error {
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
		// fpath = files.RemoveAccent(fpath)
		if _, ok := l.trackIds[fpath]; ok {
			logrus.Warnf("duplicate entry in sqlite for path '%s'", fpath)
		}
		l.trackIds[fpath] = e
	}
	return nil
}

func (l *EngineDJDB) Track(filename string) music.Track {
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

func (l *EngineDJDB) Tracks() ([]*Track, error) {
	entries := []*trackEntry{}
	query := `SELECT * FROM Track`
	err := l.sql.Unsafe().Select(&entries, query)
	if err != nil {
		return nil, errors.Wrapf(err, "query '%s' failed", query)
	}

	tracks := []*Track{}

	for _, e := range entries {
		tracks = append(tracks, newTrack(l, *e))
	}

	return tracks, nil
}
