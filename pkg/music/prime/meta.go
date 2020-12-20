package prime

type MetaStringType int
type MetaIntType int

const (
	// these are types in MetaData table
	MetaTitle           = MetaStringType(1)
	MetaArtist          = MetaStringType(2)
	MetaAlbum           = MetaStringType(3)
	MetaGenre           = MetaStringType(4)
	MetaComment         = MetaStringType(5)
	MetaPublisher       = MetaStringType(6)
	MetaComposer        = MetaStringType(7)
	MetaAlbumArtist     = MetaStringType(8)
	MetaDuration        = MetaStringType(10)
	MetaIsInHistoryList = MetaStringType(12)
	MetaFileExtension   = MetaStringType(13)

	// these are types in MetaDataInteger table
	MetaLastPlayed = MetaIntType(1)
	MetaModified   = MetaIntType(2)
	MetaAdded      = MetaIntType(3)
	MetaKey        = MetaIntType(4)
	MetaRating     = MetaIntType(5)
)
