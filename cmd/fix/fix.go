package fix

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/cmd"
	"primetools/pkg/enums"
	"primetools/pkg/files"
	"primetools/pkg/music"
	"primetools/pkg/music/factory"
)

var (
	flags = []cli.Flag{
		cmd.SourceFlag,
		cmd.SourcePathFlag,
		cmd.DryrunFlag,
		&cli.BoolFlag{
			Name:        "yes",
			Aliases:     []string{"y"},
			Usage:       "Do not prompt for write confirmation",
			Destination: &opts.accept,
		},
		&cli.PathFlag{
			Name:        "search-path",
			Aliases:     []string{"p"},
			Usage:       "path to search for music file",
			Destination: &opts.searchPath,
		},
	}

	opts = struct {
		accept     bool
		searchPath string
	}{}
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:        "fix",
		Usage:       cmd.Usage,
		HideHelp:    true,
		Description: fmt.Sprintf("try to fix problem database [%s]", strings.Join(enums.FixTypeNames(), ", ")),
		Flags:       flags,
		Action: func(context *cli.Context) error {
			return errors.Errorf("unknown fix type: %s", context.Args().First())
		},
		Subcommands: cmd.SubCmds(enums.FixTypeNames(), exec, flags, func(cmd *cli.Command) {
			cmd.Before = func(context *cli.Context) error {
				if !files.Exists(opts.searchPath) {
					return errors.Errorf("search path '%s' doesn't exists", opts.searchPath)
				}
				return nil
			}
		}),
	}
}

func exec(context *cli.Context) error {
	src := cmd.OpenSource(context)
	defer src.Close()

	typ, err := enums.ParseFixType(strings.ToLower(context.Command.Name))
	if err != nil {
		return err
	}

	logrus.Infof("scanning for %s files...", typ)

	switch typ {
	case enums.Duplicate:
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
		file, err := factory.Open(enums.File, opts.searchPath)
		if err != nil {
			return err
		}

		err = src.ForEachTrack(func(index int, total int, track music.Track) error {
			if !files.Exists(track.FilePath()) {
				logrus.Warnf("file '%s' is missing from disk", track.FilePath())

				matches := file.Matches(track)

				if len(matches) > 0 {
					logrus.Infof("found %d matching tracks", len(matches))

					var match music.Track

					if cmd.IsDryRun(context) {
						logrus.Infof("[DRY] would be changed to '%s'", matches[0])
						return nil
					}

					if !opts.accept {
						sprompt := promptui.Select{
							Label: fmt.Sprintf("Please select path to use as replacement for '%s'", track.FilePath()),
							Items: matches.Filepaths(),
						}
						idx, _, err := sprompt.Run()
						if err != nil {
							return err
						}
						match = matches[idx]
						logrus.Infof("moving '%s' to '%s'", track.FilePath(), matches[idx])
					} else {
						match = matches[0]
					}

					if match == nil {
						return errors.New("invalid match selection")
					}

					return src.MoveTrack(track, match.FilePath())
				} else {
					logrus.Errorf("could not find a match for '%v'", track)
				}
			}
			return nil
		})
	}

	return err
}
