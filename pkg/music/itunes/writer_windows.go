package itunes

import (
	"sync"
	"time"

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

func (w *writer_windows) Close() {
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
	count, err := tracks.GetCount()
	// note: index start at 1
	for i := 1; i < int(count); i++ {
		track, err := tracks.GetTrackByIndex(i)
		if err != nil {
			logrus.Errorf("tracks.GetTrackByIndex(%d): %v", i, err)
			continue
		}

		kind, err := track.GetKind()
		if err != nil {
			logrus.Errorf("track.GetKind: %v", err)
		}
		if kind != "File track" {
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

func (w *writer_windows) track(pid string) (*itunes.Track, error) {
	lib, err := w.app.GetMainPlaylist()
	if err != nil {
		return nil, errors.Wrapf(err, "itunes.GetMainPlaylist()")
	}

	tracks, err := lib.GetTracks()
	if err != nil {
		return nil, errors.Wrapf(err, "playlist.GetMainPlaylist()")
	}

	track, err := tracks.GetTrackByPersistentID(pid)
	if err != nil {
		return nil, errors.Wrapf(err, "playlist.GetTrackByPersistentID(%s)", pid)
	}

	if track.IsNil() {
		return nil, errors.Errorf("not track found for pid '%s'", pid)
	}

	return track, nil
}

func (w *writer_windows) setLocation(pid string, path string) error {
	track, err := w.track(pid)
	if err != nil {
		return err
	}
	return track.SetLocation(path)
}

func (w *writer_windows) setRating(pid string, rating int) error {
	track, err := w.track(pid)
	if err != nil {
		return err
	}
	return track.SetRating(rating)
}

func (w *writer_windows) setPlayCount(pid string, count int) error {
	track, err := w.track(pid)
	if err != nil {
		return err
	}
	return track.SetPlayedCount(count)
}
