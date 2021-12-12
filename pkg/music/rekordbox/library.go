package rekordbox

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"primetools/pkg/files"
	"primetools/pkg/music"
)

type Library struct {
	xml  XmlLibrary
	info string

	keyToTrack  map[int]music.Track
	pathToTrack map[string]music.Track
	hashCache   map[string]music.Tracks
}

func Open(path string) (*Library, error) {
	lib := &Library{
		keyToTrack:  map[int]music.Track{},
		pathToTrack: map[string]music.Track{},
	}

	start := time.Now()

	logrus.Infof("opening rekordbox xml '%s'", path)

	if !files.Exists(path) {
		return nil, errors.Errorf("rekordbox library file '%s' doesn't exists", path)
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read content of file '%s'", path)
	}

	err = xml.Unmarshal(content, &lib.xml)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parsed rekordbox file '%s' content", path)
	}

	if lib.xml.Product.Name == "" || lib.xml.Product.Version == "" || lib.xml.Product.Company == "" {
		return nil, errors.Errorf("rekordbox library file looks invalid, empty product id")
	}

	for _, it := range lib.xml.Tracks {
		track := newTrack(it)
		lib.keyToTrack[it.TrackID] = track
		lib.pathToTrack[track.FilePath()] = track
	}

	lib.info = fmt.Sprintf("%v: Version: %v, Company: %s, Track Count: %d", lib.xml.Product.Name, lib.xml.Product.Version, lib.xml.Product.Company, len(lib.xml.Tracks))
	logrus.Infof("sucessfully loaded rekordbox library in %s", time.Since(start))

	return lib, nil
}

func (l *Library) Close() {
}

func (l *Library) Track(filename string) music.Track {
	if track, ok := l.pathToTrack[filename]; ok {
		return track
	}
	return nil
}

func (l *Library) Matches(track music.Track) (matches music.Tracks) {
	if track == nil {
		return
	}

	if found := l.Track(track.FilePath()); found != nil {
		if files.Exists(found.FilePath()) {
			matches = append(matches, found)
		}
	}

	if l.hashCache == nil {
		start := time.Now()
		logrus.Info("constructing track hashes from Rekordbox library metadata")

		l.hashCache = map[string]music.Tracks{}

		err := l.ForEachTrack(func(index int, total int, track music.Track) error {
			h := music.TrackHash(track)
			if dupe, ok := l.hashCache[h]; ok {
				list := dupe.Filepaths()
				logrus.Warnf("duplicate metadata for '%s': \n  %s", track.String(), strings.Join(append(list, track.FilePath()), "\n  "))
			}
			if t, ok := track.(*Track); ok {
				l.hashCache[h] = append(l.hashCache[h], t)
			}
			return nil
		})
		if err != nil {
			logrus.Error("%v", err)
		}
		logrus.Infof("processed %d tracks in %v", len(l.hashCache), time.Since(start))
	}

	hash := music.TrackHash(track)
	if match, ok := l.hashCache[hash]; ok {
		matches = append(matches, match...)
	}

	h := music.TrackHash(track)
	return l.hashCache[h]
}

func (l *Library) Playlists() []music.Tracklist {
	return l.tracklists()
}

func (l *Library) Crates() []music.Tracklist {
	logrus.Warnf("Rekordbox doesn't have crates, getting playlists insteads")
	return l.tracklists()
}

func (l *Library) tracklists() []music.Tracklist {
	list := []music.Tracklist{}
	for _, node := range l.xml.Nodes {
		list = append(list, l.flatten(node, "")...)
	}
	return list
}

func (l *Library) flatten(node XmlPlaylistNode, parentName string) (lists []music.Tracklist) {
	pat := node.Name
	if parentName != "ROOT" {
		pat = path.Join(parentName, node.Name)
	}

	if node.Type == 1 {
		if len(node.Tracks) > 0 {
			lists = append(lists, &TrackList{
				lib:  l,
				path: pat,
				xml:  node,
			})
		}
	} else {
		for _, it := range node.Childs {
			lists = append(lists, l.flatten(it, pat)...)
		}
	}
	return
}

func (l *Library) ForEachTrack(fct music.EachTrackFunc) error {
	count := len(l.xml.Tracks)
	for idx, it := range l.xml.Tracks {
		track := &Track{it}
		err := fct(idx, count, track)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Library) String() string {
	return l.info
}

func (l *Library) trackByKey(key int) music.Track {
	if track, ok := l.keyToTrack[key]; ok {
		return track
	} else {
		return nil
	}
}
