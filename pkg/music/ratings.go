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

// func (r Rating) String() string {
// 	switch r {
// 	case Zero: return ""
// 	case OneStar: return "*"
// 	case TwoStar: return "**"
// 	case ThreeStar: return "***"
// 	case FourStar: return "****"
// 	case FiveStar: return "*****"
// 	}
// 	panic("invalid rating value")
// }
