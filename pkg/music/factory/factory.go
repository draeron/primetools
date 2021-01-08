package factory

import (
	"strings"

	"github.com/pkg/errors"

	"primetools/pkg/enums"
	"primetools/pkg/music"
	"primetools/pkg/music/files"
	"primetools/pkg/music/itunes"
	"primetools/pkg/music/prime"
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

	ltype, err := enums.ParseSourceType(typed)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid target type '%s'", typed)
	}

	switch ltype {
	case enums.ITunes:
		return itunes.Open(path)
	case enums.PRIME:
		return prime.Open(path)
	case enums.File:
		return files.Open(path), nil
	default:
		return nil, errors.Errorf("invalid library type: %v", typed)
	}
}
