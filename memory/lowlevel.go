package memory

import (
	"errors"
	"fmt"
)

var ErrNotImplemented = errors.New("not implemented")

type HugePagesMode byte

const (
	HugePagesOff HugePagesMode = iota
	HugePages2M
	HugePages1G
)

func (mode HugePagesMode) String() string {
	switch mode {
	case HugePages2M:
		return "HugePages2M"
	case HugePages1G:
		return "HugePages1G"
	default:
		return "HugePagesOff"
	}
}

func (mode HugePagesMode) GoString() string {
	return mode.String()
}

func (mode HugePagesMode) PageSize() uintptr {
	switch mode {
	case HugePages2M:
		return (1 << 21)
	case HugePages1G:
		return (1 << 30)
	default:
		return (1 << 12)
	}
}

var _ fmt.Stringer = HugePagesOff
var _ fmt.GoStringer = HugePagesOff

func mallocAlign(x uintptr, pageSize uintptr) uintptr {
	mask := pageSize - 1
	return (x + mask) & ^mask
}
