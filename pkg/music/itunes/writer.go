package itunes

import (
	"github.com/dhowden/itl"
)

type itunes_writer interface {
	close()
	load()
	addFile(path string) error
	setRating(pid string, rating int) error
	setPlayCount(pid string, count int) error
	setLocation(pid string, path string) error
	createPlaylist(name string) (*itl.Playlist, error)
	// setPlaylistContent(pid string, tracks []itunes.Track) error
}
