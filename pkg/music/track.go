package music

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
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

/*
	Resolve file path and tell if both track point to the same file
*/
func IsSameFile(left Track, right Track) bool {
	var err error

	toabs := func(path string) string {
		if !filepath.IsAbs(path) {
			path, err = filepath.Abs(path)
			if err != nil {
				logrus.Errorf("cannot determine absolute path for '%s': %v", path, err)
				return ""
			}
		}
		return path
	}

	lpath := toabs(left.FilePath())
	rpath := toabs(right.FilePath())
	return lpath == rpath
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
