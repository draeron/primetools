package itunes

import (
	"encoding/json"
	"html"
	"time"

	"github.com/dhowden/itl"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type Track struct {
	itrack itl.Track
	lib    *Library
}

func (i *Library) newTrack(track itl.Track) *Track {
	return &Track{
		itrack: track,
	}
}

func (t *Track) String() string {
	return html.UnescapeString(t.itrack.Name)
}

func (t *Track) PlayCount() int {
	return t.itrack.PlayCount
}

func (t *Track) Rating() music.Rating {
	if !t.itrack.RatingComputed {
		return music.Rating(t.itrack.Rating / 20)
	}
	return 0
}

func (t *Track) FilePath() string {
	path := files.ConvertUrlFilePath(t.itrack.Location)
	return path
}

func (t *Track) Added() time.Time {
	return t.itrack.DateAdded.UTC()
}

func (t *Track) Modified() time.Time {
	return t.itrack.DateModified.UTC()
}

func (t *Track) SetModified(added time.Time) error {
	return errors.New("cannot set modified date in iTunes")
}

func (t *Track) SetAdded(added time.Time) error {
	return errors.New("cannot set added date in iTunes")
}

func (t *Track) SetRating(rating music.Rating) error {
	return t.lib.getCreateWriter().setRating(t.itrack.PersistentID, int(rating)*20)
}

func (t *Track) SetPlayCount(count int) error {
	return t.lib.getCreateWriter().setPlayCount(t.itrack.PersistentID, count)
}

func (t *Track) Title() string {
	return html.UnescapeString(t.itrack.Name)
}

func (t *Track) Album() string {
	return html.UnescapeString(t.itrack.Album)
}

func (t *Track) Artist() string {
	return html.UnescapeString(t.itrack.Artist)
}

func (t *Track) Size() int64 {
	return int64(t.itrack.Size)
}

func (t *Track) Year() int {
	return t.itrack.Year
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
