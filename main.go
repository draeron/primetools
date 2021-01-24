package main

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/cmd"
	"primetools/cmd/dump"
	"primetools/cmd/fix"
	_import "primetools/cmd/import"
	"primetools/cmd/sync"
	"primetools/cmd/test"
	"primetools/pkg/options"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: false,
		ForceColors:   true,
	})
	log.SetOutput(logrus.New().Writer())

	// app := cli.NewApp()
	app := &cli.App{
		Name:    "primetools",
		Usage:   cmd.Usage,
		Version: "0.1.0",
		Commands: []*cli.Command{
			sync.Cmd(),
			fix.Cmd(),
			dump.Cmd(),
			_import.Cmd(),
			test.Cmd(),
		},
		Before: func(context *cli.Context) error {
			if context.Bool(cmd.Dryrun) {
				options.SetDryRun()
			}
			return nil
		},
	}
	app.Setup()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
