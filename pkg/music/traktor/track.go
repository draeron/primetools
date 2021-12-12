package traktor

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"primetools/pkg/music"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Track struct {
	xml XmlTrack
}

func newTrack(src XmlTrack) *Track {
	return &Track{src}
}

func (t Track) Title() string {
	return t.xml.Title
}

func (t Track) Album() string {
	return t.xml.Album.Title
}

func (t Track) Artist() string {
	return t.xml.Artist
}

func (t Track) Year() int {
	return t.Modified().Year()
}

func (t Track) Rating() music.Rating {
	return music.Rating(t.xml.Info.Ranking / 51)
}

func (t Track) SetRating(rating music.Rating) error {
	return errors.New("SetRating operation is not supported for traktor")
}

func (t Track) Modified() time.Time {
	modified, err := time.Parse("2006/1/2", t.xml.ModifiedDate)
	if err == nil {
		seconds, err := strconv.ParseInt(t.xml.ModifiedTime, 10, 32)
		if err != nil {
			modified = modified.Add(time.Second * time.Duration(seconds))
		}
		return modified
	} else {
		stat, err := os.Stat(t.FilePath())
		if err == nil {
			return stat.ModTime()
		}
	}
	return time.Now()
}

func (t Track) SetModified(modified time.Time) error {
	return errors.New("SetModified operation is not supported for traktor")
}

func (t Track) Added() time.Time {
	date, err := time.Parse(DateFormat, t.xml.Info.ImportDate)
	if err != nil {
		logrus.Warnf("failed to parse import date '%s'", t.xml.Info.ImportDate)
	}
	return date
}

func (t Track) SetAdded(added time.Time) error {
	return errors.New("SetAdded operation is not supported for traktor")
}

func (t Track) PlayCount() int {
	return t.xml.Info.PlayCount
}

func (t Track) SetPlayCount(count int) error {
	return errors.New("SetPlayCount operation is not supported for traktor")
}

func (t Track) FilePath() string {
	return t.xml.Filepath()
}

func (t Track) Size() int64 {
	return t.xml.Info.FileSize
}

func (t Track) String() string {
	return t.xml.Title
}

func (t Track) MarshalYAML() (interface{}, error) {
	return music.NewMarchalTrack(t), nil
}

func (t Track) MarshalTOML() ([]byte, error) {
	return json.Marshal(music.NewMarchalTrack(t))
}

func (t Track) MarshalJSON() ([]byte, error) {
	return toml.Marshal(music.NewMarchalTrack(t))
}
