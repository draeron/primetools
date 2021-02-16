package music

import (
	"sort"
)

type Tracklist interface {
	Name() string
	Path() string
	Tracks() Tracks

	SetTracks(Tracks) error
}

type Tracks []Track

/*
	Return a sorted de-deduplicated (based on path)
*/
func (t Tracks) Dedupe() Tracks {
	mappe := map[string]Track{}

	for _, it := range t {
		mappe[it.FilePath()] = it
	}

	out := Tracks{}
	for _, it := range mappe {
		out = append(out, it)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].FilePath() < out[j].FilePath()
	})

	return out
}

func (t Tracks) Titles() (titles []string) {
	for _, it := range t {
		titles = append(titles, it.Title())
	}
	return titles
}

func (t Tracks) HasFile(path string) bool {
	for _, it := range t {
		if it.FilePath() == path {
			return true
		}
	}
	return false
}

func (t Tracks) Filepaths() (paths []string) {
	for _, it := range t {
		paths = append(paths, it.FilePath())
	}
	return paths
}
