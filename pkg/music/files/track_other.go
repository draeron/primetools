// +build !windows

package files

import (
	"time"

	"github.com/pkg/errors"
)

func (t Track) SetAdded(added time.Time) error {
	return errors.New("not implemented")
}
