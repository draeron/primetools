package test

import (
	"log"

	"primetools/pkg/music"
	"primetools/pkg/music/traktor"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

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
	lib, err := traktor.Open("C:\\Users\\draeron\\Documents\\Native Instruments\\Traktor 2.11.3\\collection.nml")
	if err != nil {
		log.Fatalf("%+v", err)
	}

	println(lib.String())

	lib.ForEachTrack(func(index int, total int, track music.Track) error {
		if index > 3 {
			return nil
		}

		bytes, _ := yaml.Marshal(track)
		logrus.Info(string(bytes))

		return nil
	})

	// lib, _ := itunes.Open("")
	//
	// lib.ForEachTrack(func(index int, total int, track music.Track) error {
	//
	// 	if time.Since(track.Added()) < time.Hour * 24 {
	// 		println(track.FilePath())
	// 	}
	//
	// 	return nil
	// })

	return nil
}
