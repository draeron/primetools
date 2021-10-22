package music

type MarshallTracklist struct {
	Name   string `json:"name"`
	Path   string `json:"path,omitempty"`
	Count  int    `json:"count"`
	Tracks []MarshalTrack
}

func NewMarshallTracklist(list Tracklist) interface{} {
	tracks := []MarshalTrack{}
	for _, track := range list.Tracks() {
		val := NewMarchalTrack(track)
		tracks = append(tracks, val)
	}
	return MarshallTracklist{
		Name:   list.Name(),
		Path:   list.Path(),
		Count:  list.Count(),
		Tracks: tracks,
	}
}
