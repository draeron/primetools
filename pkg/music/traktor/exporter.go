package traktor

import (
	"encoding/xml"
	"os"
	"time"

	"primetools/pkg/music"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (l *Library) AddTrack(track music.Track) error {
	xmltrack := &XmlTrack{}

	existing := l.track(track.FilePath())
	if existing != nil {
		xmltrack = &existing.xml
	}

	err := xmltrack.CopyFromTrack(track)
	if err != nil {
		return err
	}

	if existing == nil {
		l.xml.Collection.Count++
		l.xml.Collection.Entries = append(l.xml.Collection.Entries, *xmltrack)
	}
	
	return nil
}

func (l *Library) AddPlaylist(tracklist music.Tracklist) error {
	return nil
}

func (l *Library) Export() error {
	file, err := os.Create(l.path)
	if err != nil {
		return errors.WithMessagef(err, "could not create file '%s'", l.path)
	}

	encoder := xml.NewEncoder(file)

	start := time.Now()
	err = encoder.Encode(&l.xml)
	if err != nil {
		return errors.WithMessage(err, "failed to encode into xml")
	}

	logrus.Infof("successfully exported to file '%s' in %v", l.path, time.Since(start))
	return nil
}
