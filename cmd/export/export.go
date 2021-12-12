package export

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
	}
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:        "export",
		Usage:       cmd.Usage,
		Description: fmt.Sprintf("sync assets from a source to a destination [%s]", strings.Join(enums.SyncTypeNames(), ", ")),
		Flags:       flags,
		Action:      exec,
	}
}

func exec(context *cli.Context) error {
	src := cmd.OpenSource(context)
	defer src.Close()

	tgtlib := cmd.CreateTarget(context)
	defer tgtlib.Close()

	target, ok := tgtlib.(music.LibraryExporter)
	if !ok {
		return errors.New("target library type doesn't support export")
	}

	start := time.Now()

	count := 0
	// errorsc := 0

	var err error

	err = src.ForEachTrack(func(index int, total int, track music.Track) error {
		count++
		return target.AddTrack(track)
	})
	if err != nil {
		return err
	}

	for _, it := range src.Playlists() {
		err = target.AddPlaylist(it)
		if err != nil {
			return err
		}
	}

	err = target.Export()
	if err != nil {
		return errors.WithMessagef(err, "failed to export to target library")
	}

	logrus.Infof("Export done in: %v", time.Since(start))

	// logrus.Infof("processed %d files, %d updated, %d skipped, %d errors, %d not found, duration: %s",
	// 	count, changed, count-changed-notfound, errorsc, notfound, time.Since(start))
	return err
}
