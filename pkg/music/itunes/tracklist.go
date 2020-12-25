package itunes

import (
	"html"

	"github.com/dhowden/itl"

	"primetools/pkg/music"
)

type Playlist struct {
	plist itl.Playlist
	lib   *Itunes
}

func (p Playlist) Path() string {
	path := ""
	parent := p.plist.ParentPersistentID
	for parent != "" {
		if pplist, ok := p.lib.playlistPerId[parent]; ok {
			path = html.UnescapeString(pplist.Name) + "/" + path
			parent = pplist.ParentPersistentID
		} else {
			break
		}
	}
	return path + p.Name()
}

func (p Playlist) Name() string {
	return html.UnescapeString(p.plist.Name)
}

func (p Playlist) Tracks() []music.Track {
	out := []music.Track{}
	for _, t := range p.plist.PlaylistItems {
		out = append(out, p.lib.newTrack(p.lib.trackPerId[t.TrackID]))
	}
	return out
}
