package fix

import (
	"strings"

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
		Name:        "fix",
		Usage:       cmd.Usage,
		HideHelp:    true,
		Description: "try to fix problem database",
		Flags:       flags,
		Action: func(context *cli.Context) error {
			return errors.Errorf("unknown fix type: %s", context.Args().First())
		},
		Subcommands: []*cli.Command{
			newSub(enums.Duplicate),
			newSub(enums.Missing),
		},
		// Before: before,
	}
}

func newSub(typ enums.FixType) *cli.Command {
	return &cli.Command{
		Name:      strings.ToLower(typ.String()),
		UsageText: "",
		Action:    exec,
		Flags:     flags,
		// Hidden: true,
		HideHelpCommand: true,
	}
}

func exec(context *cli.Context) error {
	src := cmd.OpenSource(context)
	defer src.Close()

	typ, err := enums.ParseFixType(strings.ToLower(context.Command.Name))
	if err != nil {
		return err
	}

	switch typ {
	case enums.Duplicate:
		logrus.Info("scanning for duplicate file in database")

		tracks := map[string]music.Track{}
		duplicates := map[string][]music.Track{}

		src.ForEachTrack(func(index int, total int, track music.Track) error {
			path := strings.ToLower(track.FilePath())

			if _, ok := tracks[path]; ok {
				duplicates[path] = append(duplicates[path], track, tracks[path])
			}
			tracks[path] = track
			return nil
		})

		if len(duplicates) > 0 {
			logrus.Infof("found %d duplicates in database", len(duplicates))
			for _, dup := range duplicates {
				for _, t := range dup {
					logrus.Infof("%s", t.FilePath())
				}
			}
		} else {
			logrus.Info("no duplicate file were found")
		}
	case enums.Missing:
	}

	return nil
}
