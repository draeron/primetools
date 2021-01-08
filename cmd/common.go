package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/pkg/enums"
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
	SourceFlag = &cli.GenericFlag{
		Name:    Source,
		Aliases: []string{"s"},
		Value: enums.ITunes.ToCliGeneric(),
	}

	TargetFlag = &cli.GenericFlag{
		Name:    Target,
		Aliases: []string{"t"},
		// Required: true,
		Value: enums.PRIME.ToCliGeneric(),
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

func open(context *cli.Context, flag string) music.Library {
	if context.String(flag) == "" {
		logrus.Errorf("--%s cannot be empty", flag)
		os.Exit(1)
	}

	lib, err := factory.Open(context.String(flag))
	if err != nil {
		logrus.Errorf("fail to open %s: %v", flag, err)
		os.Exit(1)
	}
	return lib
}

func OpenTarget(context *cli.Context) music.Library {
	return open(context, Target)
}

func OpenSource(context *cli.Context) music.Library {
	return open(context, Source)
}
