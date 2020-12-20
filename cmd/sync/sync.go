package sync

import (
	"fmt"
	"time"

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
	Modified  = SyncType("modified")
	PlayCount = SyncType("playcount")
)

var (
	SyncTypes = []SyncType{
		Ratings,
		Added,
		Modified,
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
		Name: "sync",
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
			newType(Modified),
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
		Flags:     flags,
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

	count := 0
	notfound := 0
	changed := 0
	errorsc := 0

	start := time.Now()

	err = src.ForEachTrack(func(index int, total int, track music.Track) error {
		count++

		// if strings.Contains(track.String(), "Mermaid") {
		// 	logrus.Print(track.String())
		// }

		tt := tgt.Track(track.FilePath())
		if tt == nil {
			notfound++
			return nil
		}

		switch SyncType(context.Command.Name) {
		case Modified:
			if tt.Modified() != track.Modified() {
				changed++
				msg := fmt.Sprintf("updating modified for '%s': %v => %v", track, tt.Modified().Format(time.RFC822), track.Modified().Format(time.RFC822))
				if !context.Bool(cmd.Dryrun) {
					logrus.Info(msg)
					err := tt.SetModified(track.Modified())
					if err != nil {
						errorsc++
						return err
					}
				} else {
					logrus.Info("[DRY] ", msg)
				}
			}
		case Added:
			if tt.Added() != track.Added() {
				changed++
				msg := fmt.Sprintf("updating added for '%s': %v => %v", track, tt.Added().Format(time.RFC822), track.Added().Format(time.RFC822))
				if !context.Bool(cmd.Dryrun) {
					logrus.Info(msg)
					err := tt.SetAdded(track.Added())
					if err != nil {
						errorsc++
						return err
					}
				} else {
					logrus.Info("[DRY] ", msg)
				}
			}
		case PlayCount:
			if tt.PlayCount() != track.PlayCount() {
				changed++
				msg := fmt.Sprintf("updating play count for '%s': %v => %v", track, tt.PlayCount(), track.PlayCount())
				if !context.Bool(cmd.Dryrun) {
					logrus.Info(msg)
					err := tt.SetPlayCount(track.PlayCount())
					if err != nil {
						errorsc++
						return err
					}
				} else {
					logrus.Info("[DRY] ", msg)
				}
			}
		case Ratings:
			if tt.Rating() != track.Rating() {
				changed++
				msg := fmt.Sprintf("updating rating for '%s': %v => %v", track, tt.Rating(), track.Rating())
				if !context.Bool(cmd.Dryrun) {
					logrus.Info(msg)
					err := tt.SetRating(track.Rating())
					if err != nil {
						errorsc++
						return err
					}
				} else {
					logrus.Info("[DRY] ", msg)
				}
			}

		}

		return nil
	})

	logrus.Infof("processed %d files, %d updated, %d skipped, %d errors, %d not found, duration: %s",
		count, changed, count-changed-notfound, errorsc, notfound, time.Since(start))
	return err
}
