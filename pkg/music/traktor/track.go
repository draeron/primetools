package traktor

import (
	"encoding/json"
	"time"

	"primetools/pkg/files"
	"primetools/pkg/music"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
)

type Track struct {
	xml XmlTrack
}

func newTrack(src XmlTrack) music.Track {
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
	return time.Time(t.xml.ModifiedDate).Year()
}

func (t Track) Rating() music.Rating {
	return music.Rating(t.xml.Info.Ranking / 51)
}

func (t Track) SetRating(rating music.Rating) error {
	return errors.New("SetRating operation is not supported for traktor")
}

func (t Track) Modified() time.Time {
	date := time.Time(t.xml.ModifiedDate)
	date.Add(time.Duration(t.xml.ModifiedTime))
	return date
}

func (t Track) SetModified(modified time.Time) error {
	return errors.New("SetModified operation is not supported for traktor")
}

func (t Track) Added() time.Time {
	return time.Time(t.xml.Info.ImportDate)
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
	return files.ConvertUrlFilePath(t.xml.Filepath())
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
