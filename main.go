package main

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/cmd"
	"primetools/cmd/fix"
	"primetools/cmd/sync"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: false,
		ForceColors:   true,
	})

	// app := cli.NewApp()
	app := &cli.App{
		Name:    "primetools",
		Usage:   cmd.Usage,
		Version: "0.1.0",
		Action:  run,
		Commands: []*cli.Command{
			sync.Cmd(),
			fix.Cmd(),
		},
	}
	app.Setup()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	return nil
}
