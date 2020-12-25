package prime

type MetaStringType int
type MetaIntType int

type ListType int

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

	ListPlayList = ListType(1)
	ListHistory = ListType(2)
	ListPrepare = ListType(3)
	ListCrate = ListType(4)
)

func (m MetaStringType) String() string {
	switch m {
	case MetaTitle: return "title"
	case MetaArtist: return "artist"
	case MetaAlbum: return "album"
	case MetaGenre: return "genre"
	case MetaComment: return "comment"
	case MetaPublisher: return "publisher"
	case MetaComposer: return "composer"
	case MetaAlbumArtist: return "album-artist"
	case MetaDuration: return "duration"
	case MetaIsInHistoryList: return "is-in-history-list"
	case MetaFileExtension: return "file-extension"
	}
	panic("unknown string meta type")
}

func (m MetaIntType) String() string {
	switch m {
	case MetaLastPlayed: return "played"
	case MetaModified: return "modified"
	case MetaAdded: return "added"
	case MetaKey: return "key"
	case MetaRating: return "rating"
	}
	panic("unknown int meta type")
}

func (l ListType) String() string {
	switch l {
	case ListPlayList: return "playlist"
	case ListHistory: return "history"
	case ListPrepare: return "prepare"
	case ListCrate: return "crate"
	}
	panic("unknown list type")
}
