package factory

import (
	"github.com/pkg/errors"

	"primetools/pkg/enums"
	"primetools/pkg/music"
	"primetools/pkg/music/files"
	"primetools/pkg/music/itunes"
	"primetools/pkg/music/prime"
	"primetools/pkg/music/rekordbox"
)

/*
	Arg format is type:<extra>, where extra is a comma seperated
*/
func Open(libtype enums.LibraryType, path string) (music.Library, error) {
	switch libtype {
	case enums.ITunes:
		return itunes.Open(path)
	case enums.PRIME:
		return prime.Open(path)
	case enums.File:
		return files.Open(path), nil
	case enums.Rekordbox:
		return rekordbox.Open(path)
	default:
		return nil, errors.Errorf("invalid library type: %v", libtype)
	}
}
