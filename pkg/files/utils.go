package files

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
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
	_, err := os.Stat(path)
	return os.IsExist(err)
}

func SanitizePath(path string) string {
	path = filepath.ToSlash(path)
	if runtime.GOOS == "windows" {
		path = strings.ToLower(path)
	}
	return path
}
