package test

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/cmd"
	"primetools/pkg/music"
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

	var err error

	// fcontent, _ := ioutil.ReadFile("crates.yaml")
	//
	// tracks := []music.TracklistJson{}
	// err = yaml.Unmarshal(fcontent, &tracks)
	// if err != nil {
	// 	panic(err)
	// }

	tracks := []music.Track{}
	src.ForEachTrack(func(index int, total int, track music.Track) error {
		tracks = append(tracks, track)

		if len(tracks) > 9 {
			return errors.New("end")
		} else {
			return nil
		}
	})

	crate, err := src.CreatePlaylist("un/deux/trois")
	if err != nil {
		logrus.Errorf("Err: %v", err)
		panic(err)
	}

	if len(crate.Tracks()) > 2 {
		crate.SetTracks(crate.Tracks()[:1])
	} else {
		crate.SetTracks(tracks[:5])
	}

	playlist, err := src.CreatePlaylist("un/deuxio")
	if err != nil {
		logrus.Errorf("Err: %v", err)
		panic(err)
	}
	if len(playlist.Tracks()) > 2 {
		playlist.SetTracks(crate.Tracks()[:1])
	} else {
		playlist.SetTracks(tracks[5:])
	}

	return nil
}
