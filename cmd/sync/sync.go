package sync

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/cmd"
	"primetools/pkg/music"
)

type SyncType string

const (
	Ratings   = SyncType("ratings")
	Added     = SyncType("added")
	PlayCount = SyncType("playcount")
)

var (
	SyncTypes = []SyncType{
		Ratings,
		Added,
		PlayCount,
	}

	flags = []cli.Flag{
		cmd.SourceFlag,
		cmd.TargetFlag,
		cmd.DryrunFlag,
	}
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:        "sync",
		Action: func(context *cli.Context) error {
			return errors.Errorf("unknown sync type: %s", context.Args().First())
		},
		Usage:       cmd.Usage,
		HideHelp:    true,
		Description: "sync assets from a source to a destination",
		Subcommands: []*cli.Command{
			newType(Ratings),
			newType(Added),
			newType(PlayCount),
		},
		Flags: flags,
		// Before: before,
	}
}

func newType(syncType SyncType) *cli.Command {
	return &cli.Command{
		Name:      string(syncType),
		UsageText: "",
		Action:    exec,
		Flags: flags,
		// Hidden: true,
		HideHelpCommand: true,
	}
}

// func before(context *cli.Context) error {
// 	switch SyncType(context.Command.Name) {
// 	case Ratings:
// 	case PlayCount:
// 	case Added:
// 	default:
// 		return errors.New("suported type: " + context.Args().First())
// 	}
// 	return nil
// }

func exec(context *cli.Context) error {
	src, err := cmd.OpenSource(context)
	if err != nil {
		return err
	}
	defer src.Close()

	tgt, err := cmd.OpenTarget(context)
	if err != nil {
		return err
	}
	defer tgt.Close()

	switch SyncType(context.Command.Name) {
	case Added:
		err = src.ForEachTrack(func(index int, total int, track music.Track) error {
			tt := tgt.Track(track.FilePath())
			if tt != nil && tt.Added() != track.Added() {
				logrus.Infof("updating added for '%s': => %v", track, track.Added())
				err := tt.SetAdded(track.Added())
				if err != nil {
					return err
				}
			}
			return nil
		})
	case PlayCount:
		err = src.ForEachTrack(func(index int, total int, track music.Track) error {
			tt := tgt.Track(track.FilePath())
			if tt != nil && tt.PlayCount() != track.PlayCount() {
				logrus.Infof("updating play count for '%s': %v => %v", track, tt.PlayCount(), track.PlayCount())
				err := tt.SetPlayCount(track.PlayCount())
				if err != nil {
					return err
				}
			}
			return nil
		})
	case Ratings:
		err = src.ForEachTrack(func(index int, total int, track music.Track) error {
			tt := tgt.Track(track.FilePath())
			if tt != nil && tt.Rating() != track.Rating() {
				logrus.Infof("updating rating for '%s': %v => %v", track, tt.Rating(), track.Rating())
				err := tt.SetRating(track.Rating())
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	return err
}
