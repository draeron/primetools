package music

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Track interface {
	Title() string
	Album() string
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
}

type trackJson struct {
	Title     string    `json:"title"`
	FilePath  string    `json:"file_path"`
	Album     string    `json:"album,omitempty"`
	Year      int       `json:"year,omitempty"`
	Modified  time.Time `json:"modified,omitempty"`
	Added     time.Time `json:"added,omitempty"`
	Rating    Rating    `json:"rating,omitempty"`
	PlayCount int       `json:"play_count,omitempty"`
}

func TrackToMarshalObject(track Track) trackJson {
	return trackJson{
		Title:     track.Title(),
		Album:     track.Album(),
		Year:      track.Year(),
		FilePath:  track.FilePath(),
		Added:     track.Added(),
		Modified:  track.Modified(),
		Rating:    track.Rating(),
		PlayCount: track.PlayCount(),
	}
}

func TrackHash(track Track) string {
	hash := sha1.New()

	hash.Write([]byte(track.Title()))
	hash.Write([]byte(track.Album()))
	hash.Write([]byte(strconv.Itoa(track.Year())))

	res := hash.Sum(nil)
	return fmt.Sprintf("%x", res)
}
