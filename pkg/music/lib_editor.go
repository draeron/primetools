package music

type LibraryEditor interface {
	Library
	
	AddFile(path string) (Track, error)
	CreatePlaylist(path string) (Tracklist, error)
	CreateCrate(path string) (Tracklist, error)
	MoveTrack(track Track, newpath string) error

	// return a list of file extension supported by this library
	SupportedExtensions() FileExtensions
}
