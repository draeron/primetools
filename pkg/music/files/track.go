package files

import (
	"log"
	"path/filepath"
	"time"

	"github.com/bogem/id3v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/djherbis/times.v1"

	"primetools/pkg/music"
)

type Track struct {
	path string
}

func (t Track) String() string {
	return filepath.Base(t.path)
}

func newTrack(path string) music.Track {
	return &Track{path}
}

func (t Track) Rating() music.Rating {
	tags, err := id3v2.Open(t.path, id3v2.Options{})
	if err != nil {
		logrus.Errorf("fail to open id3 tags for file %s: %v", t.path, err)
		return music.Zero
	}

	frames := tags.GetFrames(tags.CommonID("Popularimeter"))
	for _, f := range frames {
		popm, ok := f.(id3v2.PopularimeterFrame)
		if !ok {
			log.Fatal("Couldn't assert POPM frame")
		}

		if popm.Email == "traktor@native-instruments.de" {
			return music.Rating(popm.Rating / 51)
		}
	}
	return music.Zero
}

func (t Track) SetRating(rating music.Rating) error {
	logrus.Error("SetRating for file not implemented")
	return nil
}

func (t Track) Added() time.Time {
	tim, err := times.Stat(t.path)
	if err != nil {
		logrus.Errorf("failed to get creation time for file %s: %v", t.path, err)
		return time.Time{}
	}
	return tim.BirthTime()
}

func (t Track) PlayCount() int {
	panic("implement me")
}

func (t Track) SetPlayCount(count int) error {
	panic("implement me")
}

func (t Track) FilePath() string {
	return t.path
}
