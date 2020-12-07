package itunes

import (
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dhowden/itl"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type Itunes struct {
	itllib *itl.Library
	tracks map[string]itl.Track
	writer *writer
}

func Open(path string) (music.Library, error) {
	i := &Itunes{
		tracks: map[string]itl.Track{},
	}

	if path == "" {
		path = files.ExpandHomePath("~/Music/iTunes/iTunes Music Library.xml")
	}

	logrus.Infof("opening iTunes xml at '%s'...", path)

	start := time.Now()
	ifile, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open itunes xml file: %v", err)
	}

	xml, err := itl.ReadFromXML(ifile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read itunes xml file: %v", err)
	}

	for _, t := range xml.Tracks {
		i.tracks[convertPath(t.Location)] = t
	}

	logrus.Infof("sucessfully loaded itunes library in %s", time.Since(start))
	logrus.Infof("iTunes: App Version: %v, Lib Version: %v.%v, Track Count: %d", xml.ApplicationVersion, xml.MajorVersion, xml.MinorVersion, len(xml.Tracks))
	i.itllib = &xml

	i.writer, err = createWriter()
	if err != nil {
		logrus.Errorf("failed to init iTunes writer interface, writes operations will fail: %v", err)
	}

	return i, nil
}

func (i *Itunes) Close() {
	if i.writer != nil {
		i.writer.Close()
	}
}

func (i *Itunes) Track(filename string) music.Track {
	if t, ok := i.tracks[filename]; ok {
		return Track{
			itrack:   t,
			writer:   i.writer,
		}
	}
	return nil
}

// file://localhost/m:/Techno/-=%20Ambient%20=-/Bluetech/2005%20-%20Sines%20And%20Singularities/01%20-%20Enter%20The%20Lovely.mp3
func convertPath(path string) string {
	path = strings.Replace(path, "file://localhost/", "", 1)
	path, _ = url.PathUnescape(path)
	path = files.SanitizePath(path)
	return path
}
