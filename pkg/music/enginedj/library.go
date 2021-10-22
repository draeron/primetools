package enginedj

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
	main      *EngineDJDB
	dbs       map[string]*EngineDJDB
	hashCache map[string]music.Tracks
}

func Open(path string) (music.Library, error) {
	p := &Library{
		dbs: map[string]*EngineDJDB{},
	}

	if path == "" {
		path = files.ExpandHomePath("~/Music/Engine Library/Database2/m.db")
	} else if files.IsDir(path) && !strings.Contains(path, "Engine Library") {
		if !strings.HasSuffix(path, "/") || !strings.HasSuffix(path, "\\") {
			path = path + "/"
		}
		path += "Engine Library/Database2/m.db"
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
			fpath := fpath.Join(disk, "Engine Library", "Database2", "m.db")
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
	logrus.Info("EngineDJ library closed")
}

func (l *Library) Playlists() []music.Tracklist {
	return l.Crates()
}

func (l *Library) Crates() []music.Tracklist {
	mape := map[string]*Playlist{}

	for _, p := range l.main.fetchLists() {
		mape[p.Path()] = newCrate(p, l)
	}

	for _, db := range l.dbs {
		for _, p := range db.fetchLists() {
			if _, ok := mape[p.Path()]; ok {
				mape[p.Path()].MergeWith(p)
			} else {
				mape[p.Path()] = newCrate(p, l)
			}
		}
	}

	list := []music.Tracklist{}
	for _, it := range mape {
		list = append(list, it)
	}
	return list
}

/*
	Playlist are only created in the main DB
*/
func (l *Library) CreatePlaylist(path string) (music.Tracklist, error) {
	return createListIn(l.main, path)
}

func createListIn(db *EngineDJDB, path string) (music.Tracklist, error) {
	org, err := db.fetchList(path)
	if org != nil {
		return org, err
	}

	split := strings.Split(path, "/")

	var list *TrackList
	var previous *TrackList

	isLast := func(idx int) bool {
		return idx == len(split)-1
	}

	for idx, _ := range split {
		pathname := strings.Join(split[:idx+1], "/")

		list, err = db.fetchList(pathname)
		if err != nil {
			return nil, err
		}

		if list == nil {
			list = db.createList(pathname, !isLast(idx))
			if list == nil {
				return nil, errors.Errorf("failed to create playlist '%s'", pathname)
			}
		} else if !isLast(idx) {
			return nil, errors.Errorf("cannot create folder playlist '%s' since there exists another non folder playlist", pathname)
		}
		list.setParent(previous)
		previous = list
	}

	return list, nil
}

/*
	Crates are create independently in all libraries
*/
func (l *Library) CreateCrate(path string) (music.Tracklist, error) {
	return createListIn(l.main, path)
}

func (i *Library) SupportedExtensions() music.FileExtensions {
	return music.FileExtensions{} // no import supported for now
	// return []string{
	// 	".aac",
	// 	".aiff",
	// 	".aif",
	// 	".flac",
	// 	".mp3",
	// 	".mp4",
	// 	".ogg",
	// 	".wav",
	// 	".alac",
	// }
}

func (l *Library) AddFile(path string) (music.Track, error) {
	msg := "AddFile is not implemented in Library library"
	logrus.Warnf(msg)
	return nil, errors.New(msg)
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

func (l *Library) Matches(track music.Track) (matches music.Tracks) {
	if l.hashCache == nil {
		start := time.Now()
		logrus.Info("constructing track hashes from EngineDJ library metadata")

		l.hashCache = map[string]music.Tracks{}
		err := l.ForEachTrack(func(index int, total int, track music.Track) error {
			h := music.TrackHash(track)
			if dupe, ok := l.hashCache[h]; ok {
				list := dupe.Filepaths()
				logrus.Warnf("duplicate metadata for '%s': \n  %s", track.String(), strings.Join(append(list, track.FilePath()), "\n  "))
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

	hash := music.TrackHash(track)
	if match, ok := l.hashCache[hash]; ok {
		matches = append(matches, match...)
	}

	h := music.TrackHash(track)
	return l.hashCache[h]
}

func (l *Library) uniqueTracks() ([]*Track, error) {
	mainTracks, err := l.main.Tracks()
	if err != nil {
		return nil, err
	}

	mape := map[string]*Track{}
	for _, it := range mainTracks {
		mape[music.TrackHash(it)] = it
	}

	for _, db := range l.dbs {
		if list, err := db.Tracks(); err == nil {
			for _, it := range list {
				hash := music.TrackHash(it)
				// ignore duplicate in sub db which already exists in the main DB since they have a reference to
				// the origin DB
				if dupe, found := mape[hash]; !found || !music.IsSameFile(dupe, it) {
					mape[hash] = it
				}
			}
		} else {
			return nil, err
		}
	}

	list := []*Track{}
	for _, v := range mape {
		list = append(list, v)
	}
	return list, nil
}

func (l *Library) ForEachTrack(fct music.EachTrackFunc) error {
	list, err := l.uniqueTracks()
	if err != nil {
		return err
	}

	for idx, track := range list {
		if err != fct(idx, len(list), track) {
			return err
		}
	}
	return nil
}

func (l *Library) String() string {
	return l.main.info
}
