package rekordbox

import (
	"encoding/json"

	"github.com/pkg/errors"

	"primetools/pkg/music"
)

type TrackList struct {
	lib *Library
	path string
	xml XmlPlaylistNode
}

func (t *TrackList) Name() string {
	return t.xml.Name
}

func (t *TrackList) Path() string {
	return t.path
}

func (t *TrackList) Tracks() (tracks music.Tracks) {
	return t.xml.toTracks(t.lib)
}

func (t *TrackList) SetTracks(tracks music.Tracks) error {
	return errors.New("Tracklist.SetTracks operation is not supported for rekordbox")
}

func (t *TrackList) MarshalYAML() (interface{}, error) {
	return music.NewMarshallTracklist(t), nil
}

func (t *TrackList) MarshalJSON() ([]byte, error) {
	return json.Marshal(music.NewMarshallTracklist(t))
}
