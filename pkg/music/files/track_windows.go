package files

import (
	"os"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

func (t Track) SetAdded(added time.Time) error {
	fd, err := syscall.Open(t.path, os.O_RDWR, 0755)
	if err != nil {
		return errors.Wrapf(err, "could not open file %s", t.path)
	}
	defer syscall.Close(fd)

	// st, err := os.Stat(t.path)
	// if err != nil {
	// 	return err
	// }

	ctime := syscall.NsecToFiletime(int64(added.Nanosecond()))
	// mtime := syscall.NsecToFiletime(int64(st.ModTime().Nanosecond()))

	return syscall.SetFileTime(fd, &ctime, &ctime, &ctime)
}
