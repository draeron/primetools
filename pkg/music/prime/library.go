package prime

import (
	fpath "path"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type Library struct {
	main      *PrimeDB
	dbs       map[string]*PrimeDB
	hashCache map[string]music.Tracks
}

func Open(path string) (music.Library, error) {
	p := &Library{
		dbs: map[string]*PrimeDB{},
	}

	if path == "" {
		path = files.ExpandHomePath("~/Music/Engine Library/m.db")
	} else if files.IsDir(path) && !strings.Contains(path, "Engine Library") {
		if !strings.HasSuffix(path, "/") || !strings.HasSuffix(path, "\\") {
			path = path + "/"
		}
		path += "Engine Library/m.db"
	}

	var err error
	p.main, err = OpenDB(path, p)
	if err != nil {
		return nil, err
	}

	if !p.main.IsExported() {
		disks, err := files.DiskPartitions()
		if err != nil {
			p.Close()
			return nil, err
		}
		for _, disk := range disks {
			fpath := fpath.Join(disk, "Engine Library", "m.db")
			if files.Exists(fpath) {
				db, err := OpenDB(fpath, p)
				if err != nil {
					p.Close()
					return nil, err
				}
				p.dbs[db.UUID] = db
			}
		}
	}

	return p, nil
}

func (l *Library) Close() {
	l.main.Close()
	for _, it := range l.dbs {
		it.Close()
	}
	logrus.Info("PRIME library closed")
}

func (l *Library) Playlists() []music.Tracklist {
	list := l.main.fetchList(ListPlayList)
	for _, db := range l.dbs {
		list = append(list, db.fetchList(ListPlayList)...)
	}
	return list
}

func (l *Library) Crates() []music.Tracklist {
	mape := map[string]*CombinedTrackList{}

	for _, p := range l.main.fetchList(ListCrate) {
		mape[p.Path()] = newCombinedTracklist(p)
	}

	for _, db := range l.dbs {
		for _, p := range db.fetchList(ListCrate) {
			if _, ok := mape[p.Path()]; ok {
				mape[p.Path()].MergeWith(p)
			} else {
				mape[p.Path()] = newCombinedTracklist(p)
			}
		}
	}

	list := []music.Tracklist{}
	for _, it := range mape {
		list = append(list, it)
	}
	return list
}

func (l *Library) AddFile(path string) error {
	msg := "AddFile is not implemented in Library library"
	logrus.Warnf(msg)
	return errors.New(msg)
}

func (l *Library) MoveTrack(track music.Track, newpath string) error {
	itrack, ok := track.(*Track)
	if !ok {
		panic("invalid track type parameter")
	}

	return itrack.SetPath(newpath)
}

func (l *Library) Track(filename string) music.Track {
	if track := l.main.Track(filename); track != nil {
		return track
	}

	for _, db := range l.dbs {
		if track := db.Track(filename); track != nil {
			return track
		}
	}
	return nil
}

func (l *Library) Matches(track music.Track) music.Tracks {
	if l.hashCache == nil {
		start := time.Now()
		logrus.Info("constructing track hashes from PRIME library metadata")

		l.hashCache = map[string]music.Tracks{}
		err := l.ForEachTrack(func(index int, total int, track music.Track) error {
			h := music.TrackHash(track)
			if _, ok := l.hashCache[h]; ok {
				logrus.Warnf("duplicate metadata entries: ")
			}
			if t, ok := track.(*Track); ok {
				l.hashCache[h] = append(l.hashCache[h], t)
			}
			return nil
		})
		if err != nil {
			logrus.Error("%v", err)
		}
		logrus.Infof("processed %d tracks in %v", len(l.hashCache), time.Since(start))
	}

	h := music.TrackHash(track)

	return l.hashCache[h]
}

func (l *Library) ForEachTrack(fct music.EachTrackFunc) error {
	tracks, err := l.main.Tracks()
	if err != nil {
		return err
	}

	for _, db := range l.dbs {
		if list, err := db.Tracks(); err == nil {
			tracks = append(tracks, list...)
		} else {
			return err
		}
	}

	for idx, track := range tracks {
		if err != fct(idx, len(tracks), track) {
			return err
		}
	}
	return nil
}

func (l *Library) String() string {
	return l.main.info
}
