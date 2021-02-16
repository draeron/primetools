package add

import (
	"time"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/cmd"
	"primetools/pkg/files"
	flib "primetools/pkg/music/files"
)

var (
	flags = []cli.Flag{
		cmd.TargetFlag,
		cmd.TargetPathFlag,
		&cli.PathFlag{
			Name:        "search-path",
			Aliases:     []string{"p"},
			Destination: &opts.searchPath,
		},
		&cli.BoolFlag{
			Name:        "rating",
			Usage:       "also import rating which is stored in the file",
			Destination: &opts.rating,
		},
		cmd.DryrunFlag,
	}

	opts = struct {
		searchPath string
		rating     bool
	}{}
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:        "add",
		Usage:       cmd.Usage,
		Description: "Add file to target library",
		Flags:       flags,
		Before: func(context *cli.Context) error {
			if !files.Exists(opts.searchPath) {
				return errors.Errorf("search path '%s' doesn't exists.", opts.searchPath)
			}
			return nil
		},
		Action: exec,
	}
}

func exec(context *cli.Context) error {
	tgt := cmd.OpenTarget(context)
	defer tgt.Close()

	filelib := flib.Open(opts.searchPath)
	defer filelib.Close()

	start := time.Now()

	count := 0
	scanned := 0

	exts := tgt.SupportedExtensions()

	err := files.WalkMusicFiles(opts.searchPath, func(osPathname string, directoryEntry *godirwalk.Dirent) error {
		scanned++

		// skip files which are not supported by the target lib
		if !exts.Contains(osPathname) {
			return nil
		}

		if tgt.Track(osPathname) != nil {
			return nil
		}

		logrus.Infof("file '%s' not found in target library", osPathname)

		count++
		if !cmd.IsDryRun(context) {
			track, err := tgt.AddFile(osPathname)

			if err != nil {
				return err
			}

			// read rating from file and set into target lib
			if opts.rating && track != nil {
				tmp := filelib.Track(osPathname)
				return track.SetRating(tmp.Rating())
			}
		}

		return nil
	})

	logrus.Infof("Scanned %d files and added %d files to library in %s", scanned, count, time.Since(start))
	return err
}
