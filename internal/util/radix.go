package util

import (
	"strconv"
)

// Interface: Radix
// {{{

type Radix interface {
	Number() uint
	MatchRune(ch rune) bool
	Parse(str string) (uint64, error)
}

var RadixTable = map[uint]Radix{
	2:  BinaryRadix{},
	8:  OctalRadix{},
	10: DecimalRadix{},
	16: HexadecimalRadix{},
}

// BinaryRadix
// {{{

type BinaryRadix struct{}

func (BinaryRadix) Number() uint {
	return 2
}

func (BinaryRadix) MatchRune(ch rune) bool {
	return IsBinaryDigit(ch)
}

func (BinaryRadix) Parse(str string) (uint64, error) {
	return strconv.ParseUint(str, 2, 64)
}

var _ Radix = BinaryRadix{}

// }}}

// OctalRadix
// {{{

type OctalRadix struct{}

func (OctalRadix) Number() uint {
	return 8
}

func (OctalRadix) MatchRune(ch rune) bool {
	return IsOctalDigit(ch)
}

func (OctalRadix) Parse(str string) (uint64, error) {
	return strconv.ParseUint(str, 8, 64)
}

var _ Radix = OctalRadix{}

// }}}

// DecimalRadix
// {{{

type DecimalRadix struct{}

func (DecimalRadix) Number() uint {
	return 10
}

func (DecimalRadix) MatchRune(ch rune) bool {
	return IsDecimalDigit(ch)
}

func (DecimalRadix) Parse(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}

var _ Radix = DecimalRadix{}

// }}}

// HexadecimalRadix
// {{{

type HexadecimalRadix struct{}

func (HexadecimalRadix) Number() uint {
	return 16
}

func (HexadecimalRadix) MatchRune(ch rune) bool {
	return IsHexDigit(ch)
}

func (HexadecimalRadix) Parse(str string) (uint64, error) {
	return strconv.ParseUint(str, 16, 64)
}

var _ Radix = HexadecimalRadix{}

// }}}
// }}}
