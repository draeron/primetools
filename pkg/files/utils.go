package files

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/karrick/godirwalk"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	yamlv2 "gopkg.in/yaml.v2"

	"primetools/pkg/enums"
)

func ExpandHomePath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		logrus.WithError(err).Errorf("could not resolve user home dir: %v", err)
	}

	path = strings.Replace(path, "~", home, -1)
	return path
}

func Exists(path string) bool {
	// if runtime.GOOS == "windows" {
	// 	path = filepath.FromSlash(path)
	// }
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsMusicFile(path string) bool {
	return filepath.Ext(path) == ".mp3" || filepath.Ext(path) == ".flac"
}

func Size(path string) int64 {
	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return stat.Size()
}

func WalkMusicFiles(root string, walkFunc godirwalk.WalkFunc) {
	// root = filepath.ToSlash(root) + "/"

	godirwalk.Walk(root, &godirwalk.Options{
		ErrorCallback: func(s string, err error) godirwalk.ErrorAction {
			logrus.Warnf("cannot walk '%s': %v", s, err)
			return godirwalk.SkipNode
		},
		Callback: func(osPathname string, directoryEntry *godirwalk.Dirent) error {
			if IsMusicFile(osPathname) {
				return walkFunc(osPathname, directoryEntry)
			}
			return nil
		},
		AllowNonDirectory: true,
	})
}

func IsDir(path string) bool {
	st, err := os.Stat(path)
	if st == nil || err != nil {
		return false
	}
	return st.IsDir()
}

func RemoveAccent(path string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	path, _, _ = transform.String(t, path)
	return path
}

/*
	Find the absolute path, with forward slash and on windows, with lowercase
*/
func NormalizePath(path string) string {
	path, _ = filepath.Abs(path)
	path = filepath.ToSlash(path)
	if runtime.GOOS == "windows" {
		path = strings.ToLower(path)
	}
	return path
}

func WriteTo(opath string, format enums.FormatType, data interface{}) error {

	var err error
	content := []byte{}
	ext := path.Ext(opath)

	if ext == "" && format == enums.Auto {
		format = enums.Yaml
	}

	switch {
	case ext == ".yaml", ext == ".yml", format == enums.Yaml:
		yamlv2.FutureLineWrap()
		content, err = yamlv2.Marshal(data)
	case ext == ".bson":
		content, err = bson.Marshal(data)
	case ext == ".toml":
		content, err = toml.Marshal(data)
	case ext == ".json", format == enums.Json:
		content, err = json.MarshalIndent(data, "", "  ")
	default:
		content = []byte(fmt.Sprintf("%v", data))
	}

	if err != nil {
		return errors.Wrap(err, "fail to marshal struct")
	}

	var fh *os.File
	if opath == "-" {
		fh = os.Stdout
	} else if len(opath) > 0 {
		fh, err = os.Create(opath)
		if err != nil {
			return errors.Wrapf(err, "fail to open file '%s", opath)
		}
		logrus.Infof("opened file '%s' for writing", opath)
	} else {
		return errors.Errorf("cannot save to file, opath is empty")
	}

	_, err = fh.Write(content)
	return errors.Wrapf(err, "fail to write into file '%s", opath)
}
