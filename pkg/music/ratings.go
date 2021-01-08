package music

const (
	Zero    = Rating(0)
	OneStar = Rating(iota)
	TwoStar
	ThreeStar
	FourStar
	FiveStar
)

type Rating int
