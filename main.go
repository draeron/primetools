package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"primetools/cmd/fix"
	"primetools/cmd/sync"
)

func main() {
	app := &cli.App{
		Name:   "primetools",
		Usage:  "the swiss knife of engine prime",
		Action: run,
		Commands: []*cli.Command{
			sync.Cmd(),
			fix.Cmd(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	return nil
}
