package traktor

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"

	"primetools/pkg/files"
	"primetools/pkg/music"

	"github.com/pkg/errors"
)

const DateFormat = "2006/1/2"

type XmlLibrary struct {
	XMLName xml.Name `xml:"NML"`

	Version string `xml:"VERSION,attr"`
	Header  struct {
		Company string `xml:"COMPANY,attr"`
		Program string `xml:"PROGRAM,attr"`
	} `xml:"HEAD"`

	MusicFolder struct {
	} `xml:"MUSICFOLDER"`

	Collection struct {
		Count   int        `xml:"ENTRIES,attr"`
		Entries []XmlTrack `xml:"ENTRY"`
	} `xml:"COLLECTION"`

	Stem []struct {
	} `xml:"SETS"`

	Playlists []XmlPlaylistNode `xml:"PLAYLISTS>NODE"`

	SortingOrders []struct {
		Path string `xml:"Path,attr"`
	} `xml:"SORTING ORDERS"`
}

/*
<ENTRY MODIFIED_DATE="2013/12/26" MODIFIED_TIME="72910"
       AUDIO_ID="Ae0AERIiIzMzMzMzQzMzMzRDRVVlZVVVVmZWVVZmVlVmZmZmd3eHZmh3d2Znd3dmaHd4d2d3d3Zmdnd2ZoeHh3Z3d3Z3eIiIiIeIh4iHiIeIiIeHiYmImZiJmJmYiJiZmJZ2ZnZ4eHh4d3d4h4d3h4eHh4eHiYmIiZmJmJmZiJiZmJiYh4ZkRUVFREVVVVVVZWVkVWVlZDNFRDMyREMzMzMyMzREQzWIeIh4d3d4d3eHiIeHh4iGZ4iIiIiIiIiIeHiIiIh4iIiIiIiIiIiIiIiHiIiHh4iIhlRUVEREVEMzMzMzMyMzIyIzMzIiMjMiIiMzRDMyIzMyIREAAAAAAA=="
       TITLE="Enter The Lovely" ARTIST="Bluetech">
    <LOCATION DIR="/:Techno/:-= Ambient =-/:Bluetech/:2005 - Sines And Singularities/:" FILE="01 - Enter The Lovely.mp3"
              VOLUME="M:" VOLUMEID="Stuff"></LOCATION>
    <ALBUM TRACK="1" TITLE="Sines And Singularities"></ALBUM>
    <MODIFICATION_INFO AUTHOR_TYPE="user"></MODIFICATION_INFO>
    <INFO BITRATE="240000" GENRE="Chill" COMMENT="Aleph Zero" COVERARTID="067/DW1AVODG4ZIV4DGQJQJVBWUK14GB" KEY="2m"
          PLAYCOUNT="2" PLAYTIME="494" PLAYTIME_FLOAT="493.087006" RANKING="0" IMPORT_DATE="2009/1/7"
          LAST_PLAYED="2009/1/8" FLAGS="14" FILESIZE="14639"></INFO>
    <TEMPO BPM="120.000000" BPM_QUALITY="100.000000"></TEMPO>
    <LOUDNESS PEAK_DB="5.351240" PERCEIVED_DB="5.709400" ANALYZED_DB="5.709400"></LOUDNESS>
    <MUSICAL_KEY VALUE="16"></MUSICAL_KEY>
    <CUE_V2 NAME="AutoGrid" DISPL_ORDER="0" TYPE="4" START="155.974000" LEN="0.000000" REPEATS="-1" HOTCUE="0"></CUE_V2>
</ENTRY>
*/
type XmlTrack struct {
	ModifiedDate string `xml:"MODIFIED_DATE,attr"`
	ModifiedTime string `xml:"MODIFIED_TIME,attr"`
	AudioId      []byte `xml:"AUDIO_ID,attr"`

	Title  string `xml:"TITLE,attr"`
	Artist string `xml:"ARTIST,attr"`

	Location struct {
		Directory string `xml:"DIR,attr"`
		File      string `xml:"FILE,attr"`
		Volume    string `xml:"VOLUME,attr"`
		VolumeID  string `xml:"VOLUMEID,attr"`
	} `xml:"LOCATION"`

	Album struct {
		Track int    `xml:"TRACK,attr"`
		Title string `xml:"TITLE,attr"`
	} `xml:"ALBUM"`

	Modification struct {
		AuthorType string `xml:"AUTHOR_TYPE,attr"`
	} `xml:"MODIFICATION_INFO"`

	Info struct {
		Bitrate     int     `xml:"BITRATE,attr"`
		Genre       string  `xml:"GENRE,attr"`
		Comment     string  `xml:"COMMENT,attr"`
		CovertArtID string  `xml:"COVERARTID,attr"`
		Key         string  `xml:"KEY,attr"`
		PlayCount   int     `xml:"PLAYCOUNT,attr"`
		Playtime    int     `xml:"PLAYTIME,attr"`
		PlaytimeF   float32 `xml:"PLAYTIME_FLOAT,attr"`
		Ranking     int     `xml:"RANKING,attr"`
		ImportDate  string  `xml:"IMPORT_DATE,attr"`
		LastPlayed  string  `xml:"LAST_PLAYED,attr"`
		Flags       string  `xml:"FLAGS,attr"`
		FileSize    int64   `xml:"FILESIZE,attr"`
	} `xml:"INFO"`

	Tempo struct {
		BPM        float32 `xml:"BPM,attr"`
		BPMQuality float32 `xml:"BPM_QUALITY,attr"`
	} `xml:"TEMPO"`

	Loudness struct {
		PeakDB      float32 `xml:"PEAK_DB,attr"`
		PerceivedDB float32 `xml:"PERCEIVED_DB,attr"`
		AnalyzedDB  float32 `xml:"ANALYZED_DB,attr"`
	} `xml:"LOUDNESS"`

	MusicalKey string `xml:"MUSICAL_KEY>VALUE"`

	Cue []struct {
		Name        string  `xml:"NAME,attr"`
		DiplayOrder int     `xml:"DISPL_ORDER,attr"`
		Type        int     `xml:"TYPE,attr"`
		Start       float32 `xml:"START,attr"`
		Length      float32 `xml:"LEN,attr"`
		Repeats     int     `xml:"REPEATS,attr"`
		Hotcue      int     `xml:"HOTCUE,attr"`
	} `xml:"CUE_V2"`
}

func (x XmlTrack) Filepath() string {
	dir := x.Location.Directory
	dir = strings.Replace(dir, "/:", "/", -1)
	path := filepath.Join(x.Location.Volume, dir, x.Location.File)
	return files.NormalizePath(path)
}

/*
<NODE TYPE="FOLDER" NAME="$ROOT">
	<SUBNODES COUNT="1">
		<NODE TYPE="PLAYLIST" NAME="all.best">
			<PLAYLIST ENTRIES="410" TYPE="LIST" UUID="60f003a8047b46b0be356823c4f09e36">
				<ENTRY>
					<PRIMARYKEY TYPE="TRACK" KEY="M:/:Techno/:-= Prog.Trance =-/:Sun Control Species/:Sun Control Species - Bringing the Rain.mp3"></PRIMARYKEY>
				</ENTRY>
			</PLAYLIST>
		</NODE>
	</SUBNODES>
</NODE>
*/
type XmlPlaylistNode struct {
	Type string `xml:"TYPE,attr"`
	Name string `xml:"NAME,attr"`

	SubNodes *struct {
		Count int               `xml:"COUNT,attr"`
		Nodes []XmlPlaylistNode `xml:"NODE"`
	} `xml:"SUBNODES"`

	Playlist *struct {
		Type    string               `xml:"TYPE,attr"`
		Count   int                  `xml:"ENTRIES,attr"`
		Id      string               `xml:"UUID,attr"`
		Entries []XmlPlaylistEntries `xml:"ENTRY>PRIMARYKEY"`
	} `xml:"PLAYLIST"`
}

/*
<ENTRY>
	<PRIMARYKEY TYPE="TRACK" KEY="M:/:Techno/:-= Prog.Trance =-/:Lish/:2011 - Miles Away/:09 - Lish - Feel Good.mp3"></PRIMARYKEY>
</ENTRY>
*/
type XmlPlaylistEntries struct {
	Type string `xml:"TYPE,attr"`
	Key  string `xml:"KEY,attr"`
}

// func (x XmlPlaylistNode) toTracks(library *Library) (tracks music.Tracks) {
// 	for _, it := range x.Tracks {
// 		if track := library.trackByKey(it.Key); track != nil {
// 			tracks = append(tracks, track)
// 		} else {
// 			logrus.Warnf("playlist '%s' refer to a invalid track key %d", x.Name, it.Key)
// 		}
// 	}
// 	return
// }

func (x *XmlTrack) CopyFromTrack(track music.Track) error {
	filestats, err := os.Stat(track.FilePath())
	if err != nil {
		return errors.Errorf("failed to fetch file info for '%s'", track.FilePath())
	}

	x.Title = track.Title()
	x.Album.Title = track.Album()
	x.Artist = track.Artist()
	x.ModifiedDate = track.Modified().Format(DateFormat)
	// x.Cue = []bytes()
	// x.AudioId = ""

	// x.Loudness.AnalyzedDB =
	// x.Loudness.PeakDB =
	// x.Loudness.PerceivedDB =
	// x.Modification.AuthorType =
	// x.Tempo.BPM = track.be
	// x.Tempo.BPMQuality = 100

	x.Location.Directory = filepath.Dir(track.FilePath())
	x.Location.Directory = strings.Replace(x.Location.Directory, "/", "/:", -1)
	x.Location.File = filepath.Base(track.FilePath())
	// x.Location.Volume = ""
	// x.Location.VolumeID = ""

	x.Info.PlayCount = track.PlayCount()
	x.Info.FileSize = filestats.Size()
	x.Info.Ranking = int(track.Rating()) * 51
	x.Info.ImportDate = track.Added().Format(DateFormat)
	// x.Info.Key = track.Key()
	// x.Info.Genre = track.Genre()
	// x.Info.Comment = track.Comment()
	// x.Info.Bitrate = track.Bitrate()
	// x.Info.Playtime
	// x.Info.PlaytimeF

	return nil
}
