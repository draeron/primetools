package factory

import (
	"strings"

	"github.com/pkg/errors"

	"primetools/pkg/music"
	"primetools/pkg/music/files"
	"primetools/pkg/music/itunes"
)

const (
	ITunes = "itunes"
	Prime  = "prime"
	Files  = "files"
)

/*
	Arg format is type:<extra>, where extra is a comma seperated
*/
func Open(arg string) (music.Library, error) {
	if arg == "" {
		return nil, errors.New("arg cannot be empty")
	}

	args := strings.Split(arg, ";")
	typed := args[0]

	path := ""
	if len(args) > 1 {
		path = args[1]
	}

	switch typed {
	case ITunes:
		return itunes.Open(path)
	case Prime:
		return nil, nil
	case Files:
		return files.Open(path), nil
	default:
		return nil, errors.Errorf("invalid library type: %v", typed)
	}
}
