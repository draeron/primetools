package _import

import (
	"github.com/urfave/cli/v2"

	"primetools/cmd"
	"primetools/pkg/enums"
)

var (
	flags = []cli.Flag{
		cmd.SourceFlag,
		cmd.SourcePathFlag,
		cmd.TargetFlag,
		cmd.TargetPathFlag,
		&cli.StringSliceFlag{
			Name: "playlist",
		},
		&cli.StringSliceFlag{
			Name: "crate",
		},
	}
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:        "import",
		Usage:       cmd.Usage,
		HideHelp:    true,
		Description: "import a playlist",
		Subcommands: cmd.SubCmds(enums.SyncTypeNames(), exec, flags),
		Flags:       flags,
		Action:      exec,
	}
}

func exec(context *cli.Context) error {
	src := cmd.OpenSource(context)
	defer src.Close()

	tgt := cmd.OpenTarget(context)
	defer tgt.Close()

	return nil
}
