package itunes

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tdsh/itunes-win/itunes"
)

type writer struct {
	app *itunes.ITunes
	tracks map[string]*itunes.Track
	mutex sync.Mutex
}

func createWriter() (*writer, error) {
	logrus.Infof("connecting to iTunes through COM interface")

	com, err := itunes.Init()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to init to iTunes application interface: %v", err)
	}
	return &writer{
		app: com,
	}, nil
}

func (w *writer) Close() {
	w.app.Exit()
}

func (w *writer) load() {
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

	// t, err := tracks.GetTrackByName("m:\\techno\\-= prog.tek =-\\animal trainer\\animal trainer - euphorie .mp3")
	// if t != nil {
	// 	loc, _ := t.GetLocation()
	// 	println(loc)
	// }

	w.tracks = map[string]*itunes.Track{}
	count, err := tracks.GetCount()
	// note: index start at 1
	for i := 1; i < int(count); i++ {
		track, err := tracks.GetTrackByIndex(i)
		if err != nil {
			logrus.Errorf("tracks.GetTrackByIndex(%d): %v", i, err)
			// todo: print error?
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
		w.tracks[convertPath(loc)] = track
	}

	logrus.Infof("iTunes COM data loaded in %s", time.Since(start))
}

func (w *writer) rating(location string, rating int) error {
	w.load()
	t, ok := w.tracks[location]
	if !ok {
		return errors.Errorf("track with location '%s' was not found in iTunes", location)
	}

	return t.SetAlbumRating(rating * 20)
}

func (w *writer) playCount(location string, count int) error {
	w.load()
	location = convertPath(location)
	t, ok := w.tracks[location]
	if !ok {
		return errors.Errorf("track with location '%s' was not found in iTunes", location)
	}
	return t.SetPlayedCount(count)
}
