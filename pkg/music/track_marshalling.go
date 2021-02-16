package music

import (
	"encoding/json"
	"time"

	"github.com/pelletier/go-toml"
)

/*
	Simple struct used for marshalling
*/
type MarshalTrack struct {
	Title     string
	FilePath  string
	Artist    string    `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	Album     string    `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	Year      int       `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	Modified  time.Time `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	Added     time.Time `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	Rating    Rating    `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	PlayCount int       `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	Size      int64
}

/*
	Construct a MarshalTrack from a Track interface
*/
func NewMarchalTrack(track Track) MarshalTrack {
	return MarshalTrack{
		Title:     track.Title(),
		Album:     track.Album(),
		Artist:    track.Artist(),
		Year:      track.Year(),
		FilePath:  track.FilePath(),
		Added:     track.Added(),
		Modified:  track.Modified(),
		Rating:    track.Rating(),
		PlayCount: track.PlayCount(),
		Size:      track.Size(),
	}
}

/*
	Convert a marshalled track object into a Track interface
*/
func (t MarshalTrack) Interface() Track {
	return &marshalTrackAdapter{
		track: t,
	}
}

func (t MarshalTrack) String() string {
	return TrackMeta(t.Interface())
}

type marshalTrackAdapter struct {
	track MarshalTrack
}

func (m marshalTrackAdapter) Title() string {
	return m.track.Title
}

func (m marshalTrackAdapter) Album() string {
	return m.track.Album
}

func (m marshalTrackAdapter) Artist() string {
	return m.track.Artist
}

func (m marshalTrackAdapter) Year() int {
	return m.track.Year
}

func (m marshalTrackAdapter) Rating() Rating {
	return m.track.Rating
}

func (m *marshalTrackAdapter) SetRating(rating Rating) error {
	m.track.Rating = rating
	return nil
}

func (m marshalTrackAdapter) Modified() time.Time {
	return m.track.Modified
}

func (m *marshalTrackAdapter) SetModified(modified time.Time) error {
	m.track.Modified = modified
	return nil
}

func (m marshalTrackAdapter) Added() time.Time {
	return m.track.Added
}

func (m *marshalTrackAdapter) SetAdded(added time.Time) error {
	m.track.Added = added
	return nil
}

func (m marshalTrackAdapter) PlayCount() int {
	return m.track.PlayCount
}

func (m *marshalTrackAdapter) SetPlayCount(count int) error {
	m.track.PlayCount = count
	return nil
}

func (m marshalTrackAdapter) FilePath() string {
	return m.track.FilePath
}

func (m marshalTrackAdapter) Size() int64 {
	return m.track.Size
}

func (m marshalTrackAdapter) String() string {
	return TrackMeta(&m)
}

func (m marshalTrackAdapter) MarshalYAML() (interface{}, error) {
	return m.track, nil
}

func (m marshalTrackAdapter) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.track)
}

func (m marshalTrackAdapter) MarshalTOML() ([]byte, error) {
	return toml.Marshal(m.track)
}
