package rekordbox

import (
	"encoding/json"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type Track struct {
	xml XmlTrack
}

func newTrack(src XmlTrack) music.Track {
	return &Track{src }
}

func (t Track) Title() string {
	return t.xml.Name
}

func (t Track) Album() string {
	return t.xml.Album
}

func (t Track) Artist() string {
	return t.xml.Artist
}

func (t Track) Year() int {
	return t.xml.Year
}

func (t Track) Rating() music.Rating {
	return music.Rating(t.xml.Rating / 51)
}

func (t Track) SetRating(rating music.Rating) error {
	return errors.New("SetRating operation is not supported for rekordbox")
}

func (t Track) Modified() time.Time {
	loc := t.FilePath()
	return files.ModifiedTime(loc)
}

func (t Track) SetModified(modified time.Time) error {
	return errors.New("SetModified operation is not supported for rekordbox")
}

func (t Track) Added() time.Time {
	added, err := time.Parse("2006-01-02", t.xml.DateAdded)
	if err != nil {
		logrus.Error(err)
	}
	return added
}

func (t Track) SetAdded(added time.Time) error {
	return errors.New("SetAdded operation is not supported for rekordbox")
}

func (t Track) PlayCount() int {
	return t.xml.PlayCount
}

func (t Track) SetPlayCount(count int) error {
	return errors.New("SetPlayCount operation is not supported for rekordbox")
}

func (t Track) FilePath() string {
	return files.ConvertUrlFilePath(t.xml.Location)
}

func (t Track) Size() int64 {
	return t.xml.Size
}

func (t Track) String() string {
	return t.xml.Name
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
