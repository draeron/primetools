package itunes

import (
	"github.com/draeron/itunes-win/itunes"
)

type itunes_writer interface {
	Close()
	load()
	addFile(path string) error
	track(pid string) (*itunes.Track, error)
	setRating(pid string, rating int) error
	setPlayCount(pid string, count int) error
	setLocation(pid string, path string) error
}
