package music

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

type Track interface {
	Title() string
	Album() string
	Artist() string
	Year() int

	Rating() Rating
	SetRating(rating Rating) error

	Modified() time.Time
	SetModified(modified time.Time) error

	Added() time.Time
	SetAdded(added time.Time) error

	PlayCount() int
	SetPlayCount(count int) error

	FilePath() string

	Size() int64

	fmt.Stringer
	yaml.Marshaler
	toml.Marshaler
	json.Marshaler
}

type trackJson struct {
	Title     string
	FilePath  string
	Artist    string    `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	Album     string    `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	Year      int       `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	Modified  time.Time `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	Added     time.Time `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	Rating    Rating    `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
	PlayCount int       `json:",omitempty" yaml:",omitempty" toml:",omitempty"`
}

func TrackToMarshalObject(track Track) trackJson {
	return trackJson{
		Title:     track.Title(),
		Album:     track.Album(),
		Artist:    track.Artist(),
		Year:      track.Year(),
		FilePath:  track.FilePath(),
		Added:     track.Added(),
		Modified:  track.Modified(),
		Rating:    track.Rating(),
		PlayCount: track.PlayCount(),
	}
}

func TrackMeta(track Track) string {
	msg := ""
	msg += fmt.Sprintf("Impl: %v\n", reflect.TypeOf(track).Elem().Name())
	msg += fmt.Sprintf("Title: %v\n", track.Title())
	msg += fmt.Sprintf("Album: %v\n", track.Album())
	msg += fmt.Sprintf("Artist: %v\n", track.Artist())
	msg += fmt.Sprintf("Year: %v\n", strconv.Itoa(track.Year()))
	msg += fmt.Sprintf("File: %v\n", track.FilePath())
	msg += fmt.Sprintf("Hash: %v\n", TrackHash(track))
	return msg
}

func TrackHash(track Track) string {
	hash := sha1.New()

	hash.Write([]byte(track.Title()))
	hash.Write([]byte(track.Album()))
	hash.Write([]byte(track.Artist()))
	hash.Write([]byte(strconv.Itoa(track.Year())))

	res := hash.Sum(nil)
	return fmt.Sprintf("%x", res)
}
