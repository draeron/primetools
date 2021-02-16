// +build !windows

package itunes

import (
	"fmt"

	"github.com/dhowden/itl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/draeron/itunes-win/itunes"
)

type writer_notimplemented struct{}

func (w writer_notimplemented) addFile(path string) (*itl.Track, error) {
	return nil, errors.New("not implemented")
}

func (w writer_notimplemented) createPlaylist(name string) (*itl.Playlist, error) {
	return nil, errors.New("not implemented")
}

func createWriter() (itunes_writer, error) {
	logrus.Warn("itunes writing implemented on this platform")
	return &writer_notimplemented{}, nil
}

func (w writer_notimplemented) close() {
}

func (w writer_notimplemented) load() {
}

func (w writer_notimplemented) setLocation(pid string, path string) error {
	return errors.New("not implemented")
}

func (w writer_notimplemented) track(pid string) (*itunes.Track, error) {
	return nil, fmt.Errorf("not found")
}

func (w writer_notimplemented) setRating(pid string, rating int) error {
	return fmt.Errorf("setRating not implemented on this platform")
}

func (w writer_notimplemented) setPlayCount(pid string, count int) error {
	return fmt.Errorf("setPlayCount not implemented on this platform")
}
