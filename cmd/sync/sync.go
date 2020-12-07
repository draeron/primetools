package sync

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"primetools/pkg/files"
	"primetools/pkg/music/factory"
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:   "sync",
		Action: exec,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "source",
				Aliases:     []string{"src", "s"},
				Required:    true,
				DefaultText: "itune",
			},
		},
	}
}

func exec(context *cli.Context) error {

	source := context.String("source")
	if source == "" {
		return errors.New("source cannot be empty")
	}

	lib, err := factory.Open(source)
	if err != nil {
		return errors.Wrapf(err, "fail to open source: %v", err)
	}
	defer lib.Close()

	track := lib.Track(files.SanitizePath("M:\\Techno\\-= Prog.Tek =-\\Animal Trainer\\Animal Trainer - Euphorie .mp3"))

	if track != nil {
		c := track.PlayCount()
		err = track.SetPlayCount(c+1)
		if err != nil {
			logrus.Errorf("failed to write playcount: %v", err)
		}
	}

	return errors.New("not implemented")
}
