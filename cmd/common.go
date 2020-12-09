package cmd

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"primetools/pkg/music"
	"primetools/pkg/music/factory"
)

const (
	Source = "source"
	Target = "target"
	Dryrun = "dryrun"

	Usage = "the swiss knife of Denon's Engine PRIME"
)

var (
	SourceFlag = &cli.StringFlag{
		Name:     Source,
		Aliases:  []string{"s"},
		// Required: true,
		Value:    "itunes",
	}

	TargetFlag = &cli.StringFlag{
		Name:     Target,
		Aliases:  []string{"t"},
		// Required: true,
		Value:    "prime",
	}

	DryrunFlag = &cli.BoolFlag{
		Name:    Dryrun,
		Aliases: []string{"ro"},
	}
)

func CheckSourceAndTarget(context *cli.Context) error {
	if context.String(Source) == "" {
		return errors.New("--source cannot be empty")
	}

	if context.String(Target) == "" {
		return errors.New("--target cannot be empty")
	}
	return nil
}

func open(context *cli.Context, flag string) (music.Library, error) {
	if context.String(flag) == "" {
		return nil, errors.Errorf("--%s cannot be empty", flag)
	}

	return factory.Open(context.String(flag))
}

func OpenTarget(context *cli.Context) (music.Library, error) {
	return open(context, Target)
}

func OpenSource(context *cli.Context) (music.Library, error) {
	return open(context, Source)
}
