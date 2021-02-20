package files

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/karrick/godirwalk"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"gopkg.in/djherbis/times.v1"
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

func WalkMusicFiles(root string, walkFunc godirwalk.WalkFunc) error {
	// root = filepath.ToSlash(root) + "/"
	return godirwalk.Walk(root, &godirwalk.Options{
		FollowSymbolicLinks: false,
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

func ModifiedTime(path string) time.Time {
	tim, err := times.Stat(path)
	if err != nil {
		logrus.Errorf("failed to get creation time for file %s: %v", path, err)
		return time.Time{}
	}
	return tim.ModTime()
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

// file://localhost/m:/Techno/-=%20Ambient%20=-/Bluetech/2005%20-%20Sines%20And%20Singularities/01%20-%20Enter%20The%20Lovely.mp3
func ConvertUrlFilePath(path string) string {
	path = strings.Replace(path, "file://localhost/", "", 1)
	path, _ = url.PathUnescape(path)
	path = html.UnescapeString(path)
	path = NormalizePath(path)
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

func ReadFrom(path string, data interface{}) error {

	fd, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "fail to open path '%s'", path)
	}

	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return errors.Wrapf(err, "fail to read '%s'", path)
	}

	switch filepath.Ext(path) {
	case ".yml", ".yaml":
		err = yamlv2.Unmarshal(content, data)
	case ".json":
		err = json.Unmarshal(content, data)
	case ".toml":
		err = toml.Unmarshal(content, data)
	case ".bson":
		err = bson.Unmarshal(content, data)
	default:
		return errors.Errorf("unsupported file format %s", filepath.Ext(path))
	}

	if err != nil {
		return errors.Wrapf(err, "fail to parse '%s' content", path)
	}

	return nil
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
