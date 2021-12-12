package cmd

import (
	"strings"

	"github.com/gobwas/glob"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/pkg/enums"
	"primetools/pkg/music"
	"primetools/pkg/music/factory"
)

const (
	Source     = "source"
	Target     = "target"
	SourcePath = "source-path"
	TargetPath = "target-path"
	Dryrun     = "dryrun"

	Usage = "the swiss knife of Denon's Engine PRIME"
)

var (
	SourceFlag = &cli.GenericFlag{
		Name:    Source,
		Aliases: []string{"s"},
		Value:   enums.ITunes.ToCliGeneric(),
	}

	SourcePathFlag = &cli.PathFlag{
		Name:    SourcePath,
		Aliases: []string{"sp"},
	}

	TargetFlag = &cli.GenericFlag{
		Name:    Target,
		Aliases: []string{"t"},
		// Required: true,
		Value: enums.PRIME.ToCliGeneric(),
	}

	TargetPathFlag = &cli.PathFlag{
		Name:    TargetPath,
		Aliases: []string{"tp"},
	}

	DryrunFlag = &cli.BoolFlag{
		Name:    Dryrun,
		Aliases: []string{"ro"},
		Usage:   "read only mode",
	}
)

type RuleSlice struct {
	rules       []glob.Glob
	StringSlice cli.StringSlice
}

func CheckSourceAndTarget(context *cli.Context) error {
	if context.String(Source) == "" {
		return errors.New("--source cannot be empty")
	}

	if context.String(Target) == "" {
		return errors.New("--target cannot be empty")
	}
	return nil
}

func IsDryRun(context *cli.Context) bool {
	return context.Bool(DryrunFlag.Name)
}

func SubCmds(namenames []string, action cli.ActionFunc, flags []cli.Flag, extra func(cmd *cli.Command)) (subs []*cli.Command) {
	for _, name := range namenames {
		cmd := &cli.Command{
			Name:            strings.ToLower(name),
			UsageText:       "",
			Action:          action,
			Flags:           flags,
			HideHelpCommand: true,
		}
		if extra != nil {
			extra(cmd)
		}
		subs = append(subs, cmd)
	}
	return
}

type opFunc func(libtype enums.LibraryType, path string) (music.Library, error)

func open(context *cli.Context, flag string, pathflag string, op opFunc) music.Library {
	if context.String(flag) == "" {
		logrus.Errorf("--%s cannot be empty", flag)
		logrus.Exit(1)
	}

	ltype, err := enums.ParseLibraryType(context.String(flag))
	if err != nil {
		logrus.Errorf("invalid library type '%s', valid values: [%s]'", context.String(flag), strings.Join(enums.LibraryTypeNames(), ","))
		logrus.Exit(1)
	}

	lib, err := op(ltype, context.String(pathflag))
	if err != nil {
		logrus.Errorf("fail to open %s: %v", flag, err)
		logrus.Exit(1)
	}
	return lib
}

func OpenTarget(context *cli.Context) music.Library {
	return open(context, Target, TargetPath, factory.Open)
}

func CreateTarget(context *cli.Context) music.Library {
	return open(context, Target, TargetPath, factory.Create)
}

func OpenSource(context *cli.Context) music.Library {
	return open(context, Source, SourcePath, factory.Open)
}

func (r *RuleSlice) Compile() error {
	// verify rules
	for _, it := range r.StringSlice.Value() {
		rx, err := glob.Compile(it)
		if err != nil {
			return errors.Errorf("Invalid matching rule '%': %v", it, err)
		}
		r.rules = append(r.rules, rx)
	}
	return nil
}

func (r RuleSlice) Match(name string) bool {
	if len(r.rules) == 0 {
		return true
	}
	for _, rule := range r.rules {
		if rule.Match(name) {
			return true
		}
	}
	return false
}
