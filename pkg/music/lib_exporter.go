package music

type LibraryExporter interface {
	Library

	AddTrack(track Track) error
	AddPlaylist(tracklist Tracklist) error
	Export() error
}
