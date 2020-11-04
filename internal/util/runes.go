package util

import (
	"unicode"
)

func IsIllegal(ch rune) bool {
	switch {
	case ch <= 0x08:
		return true

	case ch >= 0x0e && ch <= 0x1f:
		return true

	case ch >= 0x7f && ch <= 0xbf:
		return true

	case ch >= 0x100 && unicode.IsControl(ch):
		return true

	default:
		return false
	}
}

func IsHWS(ch rune) bool {
	switch {
	case ch == 0x0009:
		return true

	case ch == 0x0020:
		return true

	case ch == 0x00a0:
		return true

	case ch == 0x1680:
		return true

	case ch >= 0x2000 && ch <= 0x200b:
		return true

	case ch == 0x202f:
		return true

	case ch == 0x205f:
		return true

	case ch == 0x3000:
		return true

	case ch == 0xfeff:
		return true

	default:
		return false
	}
}

func IsVWS(ch rune) bool {
	switch {
	case ch >= 0x000a && ch <= 0x000d:
		return true

	case ch == 0x0085:
		return true

	case ch == 0x2028:
		return true

	case ch == 0x2029:
		return true

	default:
		return false
	}
}

func IsIdentStart(ch rune) bool {
	switch {
	case ch >= 'A' && ch <= 'Z':
		return true

	case ch >= 'a' && ch <= 'z':
		return true

	case ch == '_':
		return true

	case ch == '$':
		return true

	case unicode.IsLetter(ch):
		return true

	case unicode.IsMark(ch):
		return true

	default:
		return false
	}
}

func IsIdentContinue(ch rune) bool {
	switch {
	case ch >= '0' && ch <= '9':
		return true

	case ch >= 'A' && ch <= 'Z':
		return true

	case ch >= 'a' && ch <= 'z':
		return true

	case ch == '_':
		return true

	case ch == '$':
		return true

	case unicode.IsNumber(ch):
		return true

	case unicode.IsLetter(ch):
		return true

	case unicode.IsMark(ch):
		return true

	default:
		return false
	}
}

func IsNumberContinue(ch rune) bool {
	switch {
	case ch >= '0' && ch <= '9':
		return true

	case ch >= 'A' && ch <= 'Z':
		return true

	case ch >= 'a' && ch <= 'z':
		return true

	case ch == '+':
		return true

	case ch == '-':
		return true

	case ch == '.':
		return true

	case ch == '_':
		return true

	default:
		return false
	}
}

func IsBinaryDigit(ch rune) bool {
	return (ch == '0' || ch == '1')
}

func IsOctalDigit(ch rune) bool {
	return (ch >= '0' && ch <= '7')
}

func IsDecimalDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func IsHexDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'F') || (ch >= 'a' && ch <= 'f')
}

func IsLegalForRadix(radix byte, ch rune) bool {
	switch radix {
	case 'b':
		return IsBinaryDigit(ch)

	case 'o':
		return IsOctalDigit(ch)

	case 'x':
		return IsHexDigit(ch)

	default:
		return IsDecimalDigit(ch)
	}
}

func IsFormatFlag(ch rune) bool {
	switch ch {
	case '\'':
		return true
	case '#':
		return true
	case '+':
		return true
	case '-':
		return true
	case ' ':
		return true
	case '0':
		return true
	default:
		return false
	}
}
