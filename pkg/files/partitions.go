package files

import (
	"strings"

	"github.com/deepakjois/gousbdrivedetector"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/disk"
)

func DiskPartitions() ([]string, error) {
	usbs := map[string]bool{}
	if drives, err := usbdrivedetector.Detect(); err == nil {
		// fmt.Printf("%d USB Devices Found\n", len(drives))
		for _, d := range drives {
			usbs[strings.TrimSuffix(d, "\\")] = true
		}
	} else {
		return nil, errors.Wrap(err, "failed to detect usb drives")
	}

	parts, err := disk.Partitions(true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list partitions")
	}

	out := []string{}

	for _, part := range parts {
		if _, ok := usbs[part.Mountpoint]; !ok {
			out = append(out, part.Mountpoint)
		}
	}

	return out, nil
}
