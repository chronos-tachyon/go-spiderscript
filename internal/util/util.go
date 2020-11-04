package util

import (
	"strconv"
)

var (
	NilBytes   = []byte{'<', 'n', 'i', 'l', '>'}
	RegexRunes = []rune{'#', 'r', 'x'}
	PEGRunes   = []rune{'#', 'p', 'e', 'g'}
)

func Itoa64(num uint64) string {
	return strconv.FormatUint(num, 10)
}

func Itoa(num uint) string {
	return Itoa64(uint64(num))
}

func TryAtoi(str string) (uint, bool) {
	u64, err := strconv.ParseUint(str, 10, 32)
	ok := true
	if err != nil {
		u64 = 0
		ok = false
	}
	return uint(u64), ok
}

func MustAtoi(str string) uint {
	u64, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint(u64)
}
