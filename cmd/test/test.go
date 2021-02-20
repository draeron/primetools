package test

import (
	"github.com/urfave/cli/v2"

	"primetools/cmd"
	"primetools/pkg/music/rekordbox"
)

var (
	flags = []cli.Flag{
		cmd.SourceFlag,
		cmd.SourcePathFlag,
	}
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:        "test",
		Description: "command to test code paths",
		Usage:       cmd.Usage,
		HideHelp:    true,
		Hidden:      true,
		Flags:       flags,
		Action:      exec,
	}
}

func exec(context *cli.Context) error {
	lib, err := rekordbox.Open("rekorbox.xml")
	if err != nil {
		panic(err)
	}

	println(lib.String())

	return nil
}
