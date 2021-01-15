package files

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/bogem/id3v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/djherbis/times.v1"

	"primetools/pkg/music"
)

type Track struct {
	path   string
	title  string
	album  string
	artist string
	year   int
	mutex  sync.Mutex
	loaded bool
}

func newTrack(path string) *Track {
	return &Track{
		path: path,
	}
}

func (t *Track) String() string {
	t.readMetadata()
	return t.title
}

func (t *Track) Rating() music.Rating {
	tags, err := id3v2.Open(t.path, id3v2.Options{})
	if err != nil {
		logrus.Errorf("fail to open id3 tags for file %s: %v", t.path, err)
		return music.Zero
	}
	defer tags.Close()

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

func (t *Track) SetRating(rating music.Rating) error {
	logrus.Error("SetRating for file not implemented")
	return nil
}

func (t *Track) Modified() time.Time {
	tim, err := times.Stat(t.path)
	if err != nil {
		logrus.Errorf("failed to get creation time for file %s: %v", t.path, err)
		return time.Time{}
	}
	return tim.ModTime()
}

func (t *Track) SetModified(modified time.Time) error {
	msg := "SetRating for file not implemented"
	logrus.Error(msg)
	return errors.New(msg)
}

func (t *Track) Added() time.Time {
	tim, err := times.Stat(t.path)
	if err != nil {
		logrus.Errorf("failed to get creation time for file %s: %v", t.path, err)
		return time.Time{}
	}
	return tim.BirthTime()
}

func (t *Track) PlayCount() int {
	panic("Playcount not implemented for files")
}

func (t *Track) SetPlayCount(count int) error {
	panic("Playcount not implemented for files")
}

func (t *Track) FilePath() string {
	return t.path
}

func (t *Track) Title() string {
	t.readMetadata()
	return t.title
}

func (t *Track) Artist() string {
	t.readMetadata()
	return t.artist
}

func (t *Track) Album() string {
	t.readMetadata()
	return t.album
}

func (t *Track) Year() int {
	t.readMetadata()
	return t.year
}

func (t *Track) Size() int64 {
	st, err := os.Stat(t.path)
	if err != nil {
		logrus.Errorf("could not read stat for file '%s': %v", t.path, err)
	}
	return st.Size()
}

func (t *Track) MarshalYAML() (interface{}, error) {
	return music.TrackToMarshalObject(t), nil
}

func (t *Track) MarshalJSON() (interface{}, error) {
	return json.Marshal(music.TrackToMarshalObject(t))
}

func (t *Track) readMetadata() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.loaded {
		return
	}
	t.loaded = true

	tags, err := id3v2.Open(t.path, id3v2.Options{
		Parse: true,
		ParseFrames: []string{
			"Title", "Artist", "Year", "Genre", "POPM", "Album", "TALB",
		},
	})
	if err != nil {
		logrus.Warnf("could not open id3 tags for file '%s': %v", t.path, err)
		return
	}
	defer tags.Close()

	if !tags.HasFrames() {
		logrus.Warnf("file '%s' doesn't have any id3 meta data", t.path)
		return
	}

	t.title = tags.Title()
	t.album = tags.Album()
	t.artist = tags.Artist()

	yearstr := tags.Year()
	if yearstr != "" && len(yearstr) >= 4 {
		t.year, err = strconv.Atoi(yearstr[:4])
		if err != nil {
			logrus.Errorf("could not parse year tags in file '%s': %v", t.path, err)
		}
	}
}
