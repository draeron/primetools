package music

type Library interface {
	Close()

	Track(filename string) Track
}

type Rating int

const (
	OneStar = Rating(iota)
	TwoStar
	ThreeStar
	FourStar
	FiveStar
)

type Track interface {
	Rating() int
	SetRating(rating Rating) error

	PlayCount() int
	SetPlayCount(count int) error

	FilePath() string
}

type Playlist interface {
	Tracks() []Track
}
