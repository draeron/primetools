package test

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/cmd"
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
	src := cmd.OpenSource(context)
	defer src.Close()

	_, err := src.CreatePlaylist("un/deux/trois")
	if err != nil {
		logrus.Errorf("Err: %v", err)
		panic(err)
	}

	_, err = src.CreatePlaylist("un/deuxio")
	if err != nil {
		logrus.Errorf("Err: %v", err)
		panic(err)
	}

	return nil
}
