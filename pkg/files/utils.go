package files

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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
