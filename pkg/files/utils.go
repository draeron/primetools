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

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
		content, err = yaml.Marshal(data)
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
	} else {
		return errors.Errorf("cannot save to file, opath is empty")
	}

	_, err = fh.Write(content)
	return errors.Wrapf(err, "fail to write into file '%s", opath)
}
