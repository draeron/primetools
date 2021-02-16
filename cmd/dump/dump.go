package dump

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/cmd"
	"primetools/pkg/enums"
	"primetools/pkg/files"
	"primetools/pkg/music"
)

const (
	FormatFlag = "format"
	OutputFlag = "output"
)

var (
	flags = []cli.Flag{
		cmd.SourceFlag,
		cmd.SourcePathFlag,
		&cli.PathFlag{
			Name:    OutputFlag,
			Aliases: []string{"o"},
			Value:   "-",
		},
		&cli.GenericFlag{
			Name:    FormatFlag,
			Aliases: []string{"f"},
			Value:   enums.Auto.ToCliGeneric(),
		},
		&cli.StringSliceFlag{
			Name:        "name",
			Aliases:     []string{"n"},
			DefaultText: "Names of crate/playlist to dump, can be glob (ie: *something*), if empty, will dump all objects.",
			Destination: &opts.rules.StringSlice,
		},
	}

	opts = struct {
		rules cmd.RuleSlice
	}{}
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
		Subcommands: cmd.SubCmds(enums.ObjectTypeNames(), exec, flags, nil),
	}
}

func exec(context *cli.Context) error {
	src := cmd.OpenSource(context)
	defer src.Close()

	typ, err := enums.ParseObjectType(strings.ToLower(context.Command.Name))
	if err != nil {
		return err
	}

	if err = opts.rules.Compile(); err != nil {
		return err
	}

	output := context.String(OutputFlag)
	format, _ := context.Generic(FormatFlag).(*enums.FormatType)

	switch typ {
	case enums.Playlists, enums.Crates:
		playlists := []music.Tracklist{}

		if typ == enums.Playlists {
			playlists = src.Playlists()
		} else {
			playlists = src.Crates()
		}

		playlists = filterLists(playlists)

		sort.Slice(playlists, func(i, j int) bool {
			return playlists[i].Path() < playlists[j].Path()
		})

		for _, it := range playlists {
			logrus.Infof("dumping %s '%s'", typ, it.Path())
		}

		err := files.WriteTo(output, *format, playlists)
		return errors.Cause(err)

	case enums.Tracks:
		logrus.Info("Tracks in library:")
		tracks := []music.Track{}
		err = src.ForEachTrack(func(index int, total int, track music.Track) error {
			tracks = append(tracks, track)
			return nil
		})
		if err != nil {
			return errors.Cause(err)
		}
		err := files.WriteTo(output, *format, tracks)
		return errors.Cause(err)
	}
	return err
}

func filterLists(lists []music.Tracklist) []music.Tracklist {
	out := []music.Tracklist{}
	for _, it := range lists {
		if opts.rules.Match(it.Path()) {
			out = append(out, it)
		}
	}
	return out
}
