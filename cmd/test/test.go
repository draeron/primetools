package test

import (
	"encoding/xml"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"primetools/cmd"
	"primetools/pkg/music/traktor"

	"github.com/grafov/m3u8"
	"github.com/urfave/cli/v2"
)

var (
	flags = []cli.Flag{
		cmd.SourceFlag,
		cmd.SourcePathFlag,
	}
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:        "test",
		Description: "command to test code paths",
		Usage:       cmd.Usage,
		HideHelp:    true,
		Hidden:      true,
		Flags:       flags,
		Action:      exec,
	}
}

func exec(context *cli.Context) error {
	root := os.DirFS("C:\\Users\\draeron\\Documents\\Native Instruments\\Traktor 2.11.3\\History")
	err := fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) != ".nml" || d.IsDir() {
			return nil
		}

		file, err := root.Open(path)
		if err != nil {
			return err
		}

		decoder := xml.NewDecoder(file)
		xmllib := traktor.XmlLibrary{}

		err = decoder.Decode(&xmllib)
		if err != nil {
			return err
		}

		playlist := xmllib.GetHistoryPlaylist()
		if playlist == nil {
			return nil
		}

		m3uplaylist, err := m3u8.NewMediaPlaylist(1000, 1000)
		if err != nil {
			return err
		}
		for _, entry := range playlist.Entries {
			if !entry.ExtendedData.PlayedPublic {
				continue
			}
			trackfilepath := strings.ReplaceAll(entry.PrimaryKey.Key, "/:", "/")
			trackfilepath = strings.Replace(trackfilepath, "Stuff/", "O:/", 1)
			seconds, _ := strconv.ParseFloat(entry.ExtendedData.Duration, 64)
			err = m3uplaylist.Append(trackfilepath, seconds, filepath.Base(trackfilepath))
			if err != nil {
				return err
			}
		}

		outfile, err := os.Create(filepath.Join("O:\\Traktor History", filepath.Base(path)+".m3u8"))
		if err != nil {
			return err
		}
		m3uplaylist.Close()

		content := m3uplaylist.Encode()
		if content != nil {
			outfile.Write(content.Bytes())
		}
		outfile.Close()

		return nil
	})

	// lib, _ := itunes.Open("")
	//
	// lib.ForEachTrack(func(index int, total int, track music.Track) error {
	//
	// 	if time.Since(track.Added()) < time.Hour * 24 {
	// 		println(track.FilePath())
	// 	}
	//
	// 	return nil
	// })

	return err
}
