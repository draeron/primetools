package itunes

import (
	"strings"
	"sync"
	"time"

	"github.com/dhowden/itl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/draeron/itunes-win/itunes"
	"primetools/pkg/files"
)

type writer_windows struct {
	app    *itunes.ITunes
	tracks map[string]*itunes.Track
	mutex  sync.Mutex
}

func createWriter() (itunes_writer, error) {
	logrus.Infof("connecting to iTunes through COM interface")

	com, err := itunes.Init()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to init to iTunes application interface: %v", err)
	}
	return &writer_windows{
		app: com,
	}, nil
}

func (w *writer_windows) close() {
	w.app.Exit()
}

func (w *writer_windows) load() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.tracks != nil {
		return
	}

	logrus.Infof("loading iTunes data from COM interface...")

	start := time.Now()

	lib, err := w.app.GetMainPlaylist()
	if err != nil {
		logrus.Errorf("ITunes.GetMainPlaylist: %v", err)
		return
	}

	tracks, err := lib.GetTracks()
	if err != nil {
		logrus.Errorf("ITunes.GetTracks: %v", err)
		return
	}

	w.tracks = map[string]*itunes.Track{}
	count, err := tracks.Count()
	// note: index start at 1
	for i := 1; i < int(count); i++ {
		track, err := tracks.ByIndex(i)
		if err != nil {
			logrus.Errorf("tracks.GetTrackByIndex(%d): %v", i, err)
			continue
		}

		if track.Kind() != itunes.ITTrackKindFile {
			continue
		}

		loc, err := track.GetLocation()
		if err != nil {
			logrus.Errorf("track.GetLocation: %v", err)
			continue
		}
		w.tracks[normalizePath(loc)] = track
	}

	logrus.Infof("iTunes COM data loaded in %s", time.Since(start))
}

func (w *writer_windows) addFile(path string) error {
	if !files.Exists(path) {
		return errors.Errorf("file '%s' doesn't exists", path)
	}

	lib, err := w.app.GetMainPlaylist()
	if err != nil {
		return errors.Wrapf(err, "itunes.GetMainPlaylist()")
	}

	return lib.AddFile(path)
}

func (w *writer_windows) track(pid itunes.PersistentID) (*itunes.Track, error) {
	lib, err := w.app.GetMainPlaylist()
	if err != nil {
		return nil, errors.Wrapf(err, "itunes.GetMainPlaylist()")
	}

	tracks, err := lib.GetTracks()
	if err != nil {
		return nil, errors.Wrapf(err, "playlist.GetMainPlaylist()")
	}

	track, err := tracks.ByPersistentID(pid)
	if err != nil {
		return nil, errors.Wrapf(err, "playlist.GetTrackByPersistentID(%s)", pid)
	}

	if track.IsNil() {
		return nil, errors.Errorf("not track found for pid '%s'", pid)
	}

	return track, nil
}

func (w *writer_windows) setLocation(pid string, path string) error {
	ppid, err := itunes.ParsePersistentID(pid)
	if err != nil {
		return errors.Wrapf(err, "string '%s' is not a valid persistent id", pid)
	}

	track, err := w.track(ppid)
	if err != nil {
		return err
	}
	return track.SetLocation(path)
}

func (w *writer_windows) setRating(pid string, rating int) error {
	ppid, err := itunes.ParsePersistentID(pid)
	if err != nil {
		return errors.Wrapf(err, "string '%s' is not a valid persistent id", pid)
	}

	track, err := w.track(ppid)
	if err != nil {
		return err
	}
	return track.SetRating(rating)
}

func (w *writer_windows) setPlayCount(pid string, count int) error {
	ppid, err := itunes.ParsePersistentID(pid)
	if err != nil {
		return errors.Wrapf(err, "string '%s' is not a valid persistent id", pid)
	}

	track, err := w.track(ppid)
	if err != nil {
		return err
	}
	return track.SetPlayedCount(count)
}

func (w *writer_windows) createPlaylist(path string) (*itl.Playlist, error) {
	return nil, errors.New("disabled for now, not tested enough")
	split := strings.Split(path, "/")

	collection, err := w.app.Playlists()
	if err != nil {
		return nil, errors.Wrapf(err, "fail to get main playlist")
	}

	isLast := func(idx int) bool {
		return idx != len(split)-1
	}

	var previousPlaylist *itunes.Playlist

	for idx, name := range split {
		playlist, err := collection.ByName(name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get playlist with name '%s'", name)
		}

		pathName := strings.Join(split[:idx], ".")

		// if it doesn't exists, create it
		if playlist == nil {
			// if !last, create a folder
			if !isLast(idx) {
				logrus.Infof("creating folder '%s' in itune...", pathName)
				playlist, err = w.app.CreateFolder(name)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to create folder '%s'", pathName)
				}
			} else {
				logrus.Infof("creating playlist '%s' in itune...", pathName)
				playlist, err = w.app.CreatePlaylist(name)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to create folder '%s'", pathName)
				}
			}

			if playlist == nil {
				return nil, errors.New("cannot set parent to nil")
			}

			if previousPlaylist != nil {
				err = playlist.SetParent(previousPlaylist)
				if err != nil {
					return nil, errors.Wrapf(err, "fail to set parent for '%s'", pathName)
				}
			}
		}
		previousPlaylist = playlist
	}

	pid, err := w.app.ObjectPersistentID(previousPlaylist)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch persistent id for playlist")
	}

	parent, err := previousPlaylist.Parent()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get parent for playlist")
	}

	ppid, err := w.app.ObjectPersistentID(parent)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch persistent id for playlist")
	}

	return &itl.Playlist{
		Name:                 previousPlaylist.Name(),
		PlaylistPersistentID: pid.String(),
		ParentPersistentID:   ppid.String(),
		DistinguishedKind:    int(itunes.ITUserPlaylistSpecialKindMusic),
		Visible:              true,
		Music:                true,
	}, nil
}

// func (w *writer_windows) setPlaylistContent(pid string, tracks []itunes.Track) error {
// 	playlist, err := w.app.Playlist(pid)
// 	if err != nil || playlist == nil {
// 		return err
// 	}
//
// 	if smart, _ := playlist.IsSmart(); smart {
// 		return errors.New("cannot set content of smart playlist")
// 	}
//
// }
