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
		cmd.SourcePathFlag,
		cmd.TargetFlag,
		cmd.TargetPathFlag,
		cmd.DryrunFlag,
		&cli.BoolFlag{
			Name: "force",
			Aliases: []string{"f"},
			Usage: "force update (don't do any comparaison",
			Destination: &opts.force,
		},
	}

	opts = struct{
		force bool
	}{}
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name: "sync",
		Action: func(context *cli.Context) error {
			return errors.Errorf("unknown sync type: %s", context.Args().First())
		},
		Usage:       cmd.Usage,
		HideHelp:    true,
		Description: fmt.Sprintf("sync assets from a source to a destination [%s]", strings.Join(enums.SyncTypeNames(), ", ")),
		Subcommands: cmd.SubCmds(enums.SyncTypeNames(), exec, flags, nil),
		Flags:       flags,
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

	err := tgt.ForEachTrack(func(index int, total int, track music.Track) error {
		count++

		// if strings.Contains(track.FilePath(), "hunger") {
		// 	logrus.Print(track.String())
		// }

		srct := src.Track(track.FilePath())
		if srct == nil {
			logrus.Warnf("not match found for '%s' in %v", track, cmd.SourceFlag.Value)
			notfound++
			return nil
		}

		stype, err := enums.ParseSyncType(context.Command.Name)
		if err != nil {
			return err
		}

		switch stype {
		case enums.Modified:
			if opts.force || srct.Modified().String() != track.Modified().String() {
				changed++
				msg := fmt.Sprintf("updating modified for '%s': %v => %v", track, track.Modified().Format(time.RFC822), srct.Modified().Format(time.RFC822))
				if !cmd.IsDryRun(context) {
					logrus.Info(msg)
					err := track.SetModified(srct.Modified())
					if err != nil {
						logrus.Errorf("failed to sync modified date for '%s': %v", srct.Title(), err)
						errorsc++
					}
				} else {
					logrus.Info("[DRY] ", msg)
				}
			}
		case enums.Added:
			// left := srct.Added().String()
			// right := track.Added().String()
			// println(left, right)
			if opts.force || srct.Added().String() != track.Added().String() {
				changed++
				msg := fmt.Sprintf("updating added for '%s': %v => %v", track, track.Added().Format(time.RFC822), srct.Added().Format(time.RFC822))
				if !cmd.IsDryRun(context) {
					logrus.Info(msg)
					err := track.SetAdded(srct.Added())
					if err != nil {
						errorsc++
						logrus.Errorf("failed to sync added date for '%s': %v", srct.Title(), err)
					}
				} else {
					logrus.Info("[DRY] ", msg)
				}
			}
		case enums.PlayCount:
			if opts.force || srct.PlayCount() != track.PlayCount() {
				changed++
				msg := fmt.Sprintf("updating play count for '%s': %v => %v", track, srct.PlayCount(), srct.PlayCount())
				if !cmd.IsDryRun(context) {
					logrus.Info(msg)
					err := track.SetPlayCount(srct.PlayCount())
					if err != nil {
						errorsc++
						logrus.Errorf("failed to sync playcount for '%s': %v", srct.Title(), err)
					}
				} else {
					logrus.Info("[DRY] ", msg)
				}
			}
		case enums.Ratings:
			if opts.force || srct.Rating() != track.Rating() {
				changed++
				msg := fmt.Sprintf("updating rating for '%s': %v => %v", track, track.Rating(), srct.Rating())
				if !cmd.IsDryRun(context) {
					logrus.Info(msg)
					err := track.SetRating(srct.Rating())
					if err != nil {
						errorsc++
						logrus.Errorf("failed to sync rating for '%s': %v", srct.Title(), err)
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
