package fix

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:   "fix",
		Action: exec,
	}
}

func exec(context *cli.Context) error {
	return errors.New("not implemented")
}
