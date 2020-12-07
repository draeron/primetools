package factory

import (
	"strings"

	"github.com/pkg/errors"

	"primetools/pkg/music"
	"primetools/pkg/music/itunes"
)

/*
	Arg format is type:<extra>, where extra is a comma seperated
*/
func Open(arg string) (music.Library, error) {
	if arg == "" {
		return nil, errors.New("arg cannot be empty")
	}

	args := strings.Split(arg, ":")
	typed := args[0]

	switch typed {
	case "itunes":
		return itunes.Open("")
	default:
		return nil, errors.Errorf("invalid library type: %v", typed)
	}
}
