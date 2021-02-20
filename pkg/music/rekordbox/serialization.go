package rekordbox

import (
	"github.com/sirupsen/logrus"

	"primetools/pkg/music"
)

type XmlLibrary struct {
	Product struct {
		Name    string `xml:"Name,attr"`
		Version string `xml:"Version,attr"`
		Company string `xml:"Company,attr"`
	} `xml:"PRODUCT"`

	Tracks []XmlTrack        `xml:"COLLECTION>TRACK"`
	Nodes  []XmlPlaylistNode `xml:"PLAYLISTS>NODE"`
}

type XmlTrack struct {
	TrackID   int    `xml:"TrackID,attr"`
	Name      string `xml:"Name,attr"`
	Album     string `xml:"Album,attr"`
	Artist    string `xml:"Artist,attr"`
	Genre     string `xml:"Genre,attr"`
	Year      int    `xml:"Year,attr"`
	Size      int64  `xml:"Size,attr"`
	Rating    int    `xml:"Rating,attr"`
	DateAdded string `xml:"DateAdded,attr"`
	PlayCount int    `xml:"PlayCount,attr"`
	Location  string `xml:"Location,attr"`
}

type XmlPlaylistNode struct {
	Type   int    `xml:"Type,attr"`
	Name   string `xml:"Name,attr"`
	Tracks []struct {
		Key int `xml:"Key,attr"`
	} `xml:"TRACK"`
	Childs []XmlPlaylistNode `xml:"NODE"`
}

func (x XmlPlaylistNode) toTracks(library *Library) (tracks music.Tracks) {
		for _, it := range x.Tracks {
			if track := library.trackByKey(it.Key); track != nil {
				tracks = append(tracks, track)
			} else {
				logrus.Warnf("playlist '%s' refer to a invalid track key %d", x.Name, it.Key)
			}
		}
		return
}
