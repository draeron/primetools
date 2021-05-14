package files

import (
	"encoding/json"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bogem/id3v2"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/djherbis/times.v1"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type Track struct {
	path   string
	title  string
	album  string
	artist string
	rating music.Rating
	year   int
	mutex  sync.Mutex
	loaded bool
}

const TracktorEmail = "traktor@native-instruments.de"

func newTrack(path string) *Track {
	return &Track{
		path: path,
	}
}

func (t *Track) String() string {
	t.readMetadata()
	if t.title == "" {
		return t.FilePath()
	} else {
		return strings.Join([]string{t.artist, t.title}, " - ")
	}
}

func (t *Track) Rating() music.Rating {
	t.readMetadata()
	return t.rating
}

func (t *Track) SetRating(rating music.Rating) error {
	tags, err := id3v2.Open(t.path, id3v2.Options{
		Parse: true,
	})
	if err != nil {
		logrus.Errorf("fail to open id3 tags for file %s: %v", t.path, err)
	}
	defer tags.Close()

	var popframe *id3v2.PopularimeterFrame

	for _, frame := range tags.GetFrames("POPM") {
		if popm, ok := frame.(id3v2.PopularimeterFrame); ok {
			if popm.Email == TracktorEmail {
				popframe = &popm
			}
		}
	}
	if popframe == nil {
		popframe = &id3v2.PopularimeterFrame{
			Email: TracktorEmail,
			Counter: &big.Int{},
		}
	}
	popframe.Rating = uint8(rating) * 51
	tags.AddFrame("POPM", popframe)

	tags.SetVersion(4)
	enc := tags.DefaultEncoding()

	for id, framer := range tags.AllFrames() {
		for _, it := range framer {
			switch frame := it.(type) {
			case id3v2.UserDefinedTextFrame:
				// fmt.Printf("USER:%s: %s\n", id, frame.Value)
				frame.Encoding = enc
				tags.AddFrame(id, frame)
			case id3v2.CommentFrame:
				// fmt.Printf("COMMENT:%s: %s\n", id, frame.Text)
				frame.Encoding = enc
				tags.AddCommentFrame(frame)
			case id3v2.TextFrame:
				// fmt.Printf("TEXT:%s: %s\n", id, frame.Text)
				frame.Encoding = enc
				tags.AddTextFrame(id, enc, frame.Text)
			}
		}
	}

	return tags.Save()
}

func (t *Track) Modified() time.Time {
	return files.ModifiedTime(t.path)
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
	return music.NewMarchalTrack(t), nil
}

func (t *Track) MarshalJSON() ([]byte, error) {
	return json.Marshal(music.NewMarchalTrack(t))
}

func (t *Track) MarshalTOML() ([]byte, error) {
	return toml.Marshal(music.NewMarchalTrack(t))
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
	if t.title == "" {
		logrus.Warnf("file '%s' doesn't have any id3 title data", t.path)
	}

	for _, frame := range tags.GetFrames("POPM") {
		if popm, ok := frame.(id3v2.PopularimeterFrame); ok {
			if popm.Email == TracktorEmail {
				t.rating = music.Rating(popm.Rating / 51)
			}
		}
	}

	yearstr := tags.Year()
	if yearstr != "" && len(yearstr) >= 4 {
		t.year, err = strconv.Atoi(yearstr[:4])
		if err != nil {
			logrus.Errorf("could not parse year tags in file '%s': %v", t.path, err)
		}
	}
}
