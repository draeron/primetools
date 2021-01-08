package sync

import (
	"fmt"
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
			newSub(enums.Ratings),
			newSub(enums.Added),
			newSub(enums.Modified),
			newSub(enums.PlayCount),
		},
		Flags: flags,
		// Before: before,
	}
}

func newSub(syncType enums.SyncType) *cli.Command {
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
	src := cmd.OpenSource(context)
	defer src.Close()

	tgt := cmd.OpenTarget(context)
	defer tgt.Close()

	count := 0
	notfound := 0
	changed := 0
	errorsc := 0

	start := time.Now()

	err := src.ForEachTrack(func(index int, total int, track music.Track) error {
		count++

		// if strings.Contains(track.String(), "Troja Átma") {
		// 	logrus.Print(track.String())
		// }

		tt := tgt.Track(track.FilePath())
		if tt == nil {
			notfound++
			return nil
		}

		stype, err := enums.ParseSyncType(context.Command.Name)
		if err != nil {
			return err
		}

		switch stype {
		case enums.Modified:
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
		case enums.Added:
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
		case enums.PlayCount:
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
		case enums.Ratings:
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
