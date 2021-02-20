package _import

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/cmd"
	"primetools/pkg/enums"
	"primetools/pkg/files"
	"primetools/pkg/music"
)

var (
	flags = []cli.Flag{
		cmd.TargetFlag,
		cmd.TargetPathFlag,
		cmd.DryrunFlag,
		&cli.PathFlag{
			Name:        "source",
			Aliases:     []string{"s"},
			Destination: &opts.source,
			Usage:       "file to use as source",
		},
		&cli.StringSliceFlag{
			Name:        "name",
			Aliases:     []string{"n"},
			Usage:       "Names of crate/playlist to import, can be glob (*something*), if none is given, will import all object in dump file.",
			Destination: &opts.rules.StringSlice,
		},
		// &cli.BoolFlag{
		// 	Name:        "yes",
		// 	Aliases:     []string{"y"},
		// 	DefaultText: "Do not prompt for write confirmation",
		// 	Destination: &opts.accept,
		// },
		&cli.BoolFlag{
			Name:        "ignore-not-found",
			Aliases:     []string{},
			Usage:       "Ignore track which aren't found in target, otherwise the operation will fail.",
			Destination: &opts.ignoreNotFound,
		},
	}

	opts = struct {
		source         string
		accept         bool
		ignoreNotFound bool
		rules          cmd.RuleSlice
		objType        enums.ObjectType
	}{}

	importTypes = []string{enums.Playlists.String(), enums.Crates.String()}
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:        "import",
		Usage:       cmd.Usage,
		HideHelp:    true,
		Description: "import playlist/crates",
		Subcommands: cmd.SubCmds(importTypes, exec, flags, nil),
		Flags:       flags,
		Action:      exec,
	}
}

func exec(context *cli.Context) error {
	var err error

	opts.objType, err = enums.ParseObjectType(strings.ToLower(context.Command.Name))
	if err != nil {
		return errors.Errorf("unsupported import type %s, valid types are [%v]", context.Command.Name, strings.Join(importTypes, ", "))
	}

	target := cmd.OpenTarget(context)
	defer target.Close()

	if !files.Exists(opts.source) {
		return errors.Errorf("file '%s' doesn't exists or is invalid", opts.source)
	}

	lists := []music.MarshallTracklist{}
	err = files.ReadFrom(opts.source, &lists)
	if err != nil {
		return errors.New("failed reading from " + opts.source)
	}

	if len(lists) == 0 {
		return errors.Errorf("file '%s' doesn't contains any %s", opts.source, opts.objType)
	}

	if err = opts.rules.Compile(); err != nil {
		return err
	}

	for _, list := range lists {
		err = importList(context, list, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func importList(context *cli.Context, list music.MarshallTracklist, target music.Library) error {
	var err error

	if !opts.rules.Match(list.Path) {
		logrus.Infof("%s '%s' doesn't match any rule, skipping", opts.objType, list.Path)
	}

	logrus.Infof("importing %s '%s' from file into target library", opts.objType, list.Path)

	var targetList music.Tracklist

	switch opts.objType {
	case enums.Playlists:
		targetList, err = target.CreatePlaylist(list.Path)
	case enums.Crates:
		targetList, err = target.CreateCrate(list.Path)
	default:
		return errors.Errorf("unsupported type: %s", opts.objType)
	}

	if err != nil {
		return errors.Errorf("failed to create %s '%s': %v", opts.objType, list.Path, err)
	}

	oldCount := len(targetList.Tracks())

	var newList music.Tracks

	for _, track := range list.Tracks {
		matches := target.Matches(track.Interface())

		// todo: add prompt to choose the match ?
		if len(matches) > 0 {
			newList = append(newList, matches[0])
		} else if opts.ignoreNotFound {
			logrus.Warnf("failed to find a match for file '%v' in target library for in %s '%s'", track, opts.objType, list.Path)
		} else {
			return errors.Errorf("failed to find a match for track '%v' in target library for in %s '%s', skipping write", track, opts.objType, list.Path)
		}
	}

	msg := fmt.Sprintf("%s '%s' was updated from %d to %d items", opts.objType, targetList.Path(), oldCount, len(newList))
	if !cmd.IsDryRun(context) {
		err = targetList.SetTracks(newList)
		if err != nil {
			return err
		}
	} else {
		msg = "[DRY] " + msg
	}

	logrus.Infof(msg)
	return nil
}
