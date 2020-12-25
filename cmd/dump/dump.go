package dump

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/cmd"
)

var (
	flags = []cli.Flag{
		cmd.SourceFlag,
	}
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:        "dump",
		Action:      exec,
		Usage:       cmd.Usage,
		HideHelp:    true,
		Description: "dump data about a library",
		Flags:       flags,
		// Before: before,
	}
}

func exec(context *cli.Context) error {
	src, err := cmd.OpenSource(context)
	if err != nil {
		return err
	}
	defer src.Close()

	for _, playlist := range src.Playlists() {
		tracks := playlist.Tracks()
		logrus.Infof("%v (%d items)", playlist.Path(), len(tracks))

		// for tcount, track := range tracks {
		// 	logrus.Infof("  - %2d %s (%s)", tcount, track.Name(), track.FilePath())
		// }
	}

	// count := 0
	// notfound := 0
	// errorsc := 0
	//
	// start := time.Now()

	// err = src.ForEachTrack(func(index int, total int, track music.Track) error {
	// 	count++
	//
	//
	//
	// 	return nil
	// })

	// logrus.Infof("parsed %d files, %d skipped, %d errors, %d not found, duration: %s",
	// 	count, count-notfound, errorsc, notfound, time.Since(start))
	return err
}
