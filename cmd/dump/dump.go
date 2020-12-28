package dump

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/cmd"
	"primetools/pkg/enums"
	"primetools/pkg/music"
)

var (
	flags = []cli.Flag{
		cmd.SourceFlag,
	}
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:        "dump",
		Usage:       cmd.Usage,
		HideHelp:    true,
		Description: "dump data about a library",
		Flags:       flags,
		Action: func(context *cli.Context) error {
			return errors.Errorf("unknown sync type: %s", context.Args().First())
		},
		Subcommands: []*cli.Command{
			newSub(enums.Tracks),
			newSub(enums.Playlists),
			newSub(enums.Crates),
		},
		// Before: before,
	}
}

func newSub(syncType enums.ObjectType) *cli.Command {
	return &cli.Command{
		Name:      strings.ToLower(syncType.String()),
		UsageText: "",
		Action:    exec,
		Flags:     flags,
		// Hidden: true,
		HideHelpCommand: true,
	}
}

func exec(context *cli.Context) error {
	src, err := cmd.OpenSource(context)
	if err != nil {
		return err
	}
	defer src.Close()

	typ, err := enums.ParseObjectType(strings.ToLower(context.Command.Name))
	if err != nil {
		return err
	}

	switch typ {
	case enums.Playlists,	enums.Crates:
		playlists := []music.Tracklist{}

		if typ == enums.Playlists {
			playlists = src.Playlists()
		} else {
			playlists = src.Crates()
		}

		for _, playlist := range playlists {
			tracks := playlist.Tracks()
			logrus.Infof("%v (%d items)", playlist.Path(), len(tracks))
			// for tcount, track := range tracks {
			// 	logrus.Infof("  - %2d %s (%s)", tcount, track.Name(), track.FilePath())
			// }
		}
	case enums.Tracks:
		logrus.Info("Tracks in library:")
		err = src.ForEachTrack(func(index int, total int, track music.Track) error {
			logrus.Infof("- '%s' | Rating: %v, Added: %s", track.Title(), track.Rating(), track.Added().Local().Format(time.ANSIC))
			return nil
		})
		return err
	}
	return err
}
