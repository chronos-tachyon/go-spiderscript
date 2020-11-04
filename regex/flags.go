package regex

import (
	"fmt"
	"strings"
)

type Flag uint16

const (
	FlagI Flag = (1 << iota)
	FlagM
	FlagS
	FlagX
)

const flagBits = 16

var flagNames = []byte{'i', 'm', 's', 'x'}
var flagGoNames = []string{
	"FlagI",
	"FlagM",
	"FlagS",
	"FlagX",
}

func (flags Flag) String() string {
	var sb strings.Builder
	sb.Grow(4)
	for bitIndex := uint(0); bitIndex < flagBits; bitIndex++ {
		bit := Flag(1) << bitIndex
		if (flags & bit) == bit {
			if bitIndex < uint(len(flagNames)) {
				sb.WriteByte(flagNames[bitIndex])
			}
		}
	}
	return sb.String()
}

func (flags Flag) GoString() string {
	if flags == 0 {
		return "0"
	}

	var sb strings.Builder
	sb.Grow(24)
	first := true
	for bitIndex := uint(0); bitIndex < flagBits; bitIndex++ {
		bit := Flag(1) << bitIndex
		if (flags & bit) == bit {
			if !first {
				sb.WriteByte('|')
			}
			first = false
			if bitIndex < uint(len(flagGoNames)) {
				sb.WriteString(flagGoNames[bitIndex])
			} else {
				fmt.Fprintf(&sb, "Flag(%#04x)", uint16(bit))
			}
		}
	}
	return sb.String()
}

var _ fmt.Stringer = Flag(0)
var _ fmt.GoStringer = Flag(0)
