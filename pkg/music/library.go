package music

import (
	"fmt"
)

type Library interface {
	Close()
	AddFile(path string) error
	Track(filename string) Track

	Playlists() []Tracklist
	Crates() []Tracklist

	ForEachTrack(fct EachTrackFunc) error

	fmt.Stringer
}

type EachTrackFunc func(index int, total int, track Track) error

type Tracklist interface {
	Name() string
	Path() string
	Tracks() []Track
}

func TracklistToMarshal(list Tracklist) interface{} {
	tracks := []trackJson{}
	for _, track := range list.Tracks() {
		val := TrackToMarshalObject(track)
		tracks = append(tracks, val)
	}
	return struct {
		Name string `json:"name"`
		Path string `json:"path,omitempty"`
		Tracks []trackJson
	}{
		Name:   list.Name(),
		Path:   list.Path(),
		Tracks: tracks,
	}
}
