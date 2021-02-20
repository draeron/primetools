package test

import (
	"time"

	"github.com/urfave/cli/v2"

	"primetools/cmd"
	"primetools/pkg/music"
	"primetools/pkg/music/itunes"
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
	// lib, err := rekordbox.Open("rekorbox.xml")
	// if err != nil {
	// 	panic(err)
	// }
	//
	// println(lib.String())

	lib, _ := itunes.Open("")

	lib.ForEachTrack(func(index int, total int, track music.Track) error {

		if time.Since(track.Added()) < time.Hour * 24 {
			println(track.FilePath())
		}

		return nil
	})

	return nil
}
