package value

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
)

// Number
// {{{

type Number struct {
	Sign             byte
	RadixSymbol      byte
	ExponentSymbol   byte
	ExponentSign     byte
	IntegralDigits   []byte
	FractionalDigits []byte
	ExponentDigits   []byte
}

func NewZero() *Number {
	return &Number{
		Sign:           '+',
		IntegralDigits: []byte{'0'},
	}
}

func (nv *Number) String() string {
	return util.StringImpl(nv)
}

func (nv *Number) GoString() string {
	return util.GoStringImpl(nv)
}

func (nv *Number) EstimateStringLength() uint {
	return 6 + uint(len(nv.IntegralDigits)) + uint(len(nv.FractionalDigits)) + uint(len(nv.ExponentDigits))
}

func (nv *Number) EstimateGoStringLength() uint {
	return 26 + uint(len(nv.IntegralDigits)) + uint(len(nv.FractionalDigits)) + uint(len(nv.ExponentDigits))
}

func (nv *Number) WriteStringTo(out *strings.Builder) {
	if nv.Sign != 0 {
		out.WriteByte(nv.Sign)
	}
	if nv.RadixSymbol != 0 {
		out.WriteByte('0')
		out.WriteByte(nv.RadixSymbol)
	}
	out.Write(nv.IntegralDigits)
	if nv.FractionalDigits != nil {
		out.WriteByte('.')
		out.Write(nv.FractionalDigits)
	}
	if nv.ExponentSymbol != 0 {
		out.WriteByte(nv.ExponentSymbol)
		if nv.ExponentSign != 0 {
			out.WriteByte(nv.ExponentSign)
		}
		out.Write(nv.ExponentDigits)
	}
}

func (nv *Number) WriteGoStringTo(out *strings.Builder) {
	bytefn := func(ch byte) byte {
		if ch == 0 {
			return '_'
		}
		return ch
	}

	strfn := func(str []byte) []byte {
		if str == nil {
			return util.NilBytes
		}
		return str
	}

	out.WriteString("&value.Number{")
	out.WriteByte(bytefn(nv.Sign))
	out.WriteByte(',')
	out.WriteByte(bytefn(nv.RadixSymbol))
	out.WriteByte(',')
	out.Write(strfn(nv.IntegralDigits))
	out.WriteByte(',')
	out.Write(strfn(nv.FractionalDigits))
	out.WriteByte(',')
	out.WriteByte(bytefn(nv.ExponentSymbol))
	out.WriteByte(',')
	out.WriteByte(bytefn(nv.ExponentSign))
	out.WriteByte(',')
	out.Write(strfn(nv.ExponentDigits))
	out.WriteByte('}')
}

func (nv *Number) Parse(input []rune) error {
	*nv = Number{
		Sign:           '+',
		IntegralDigits: make([]byte, 0, len(input)),
		ExponentSign:   '+',
	}

	state := NumberWantSign
	bufferedZero := false
	for index, ch := range input {
		lower := unicode.ToLower(ch)

		if lower == '_' {
			continue
		}

		if lower == '+' || lower == '-' {
			if state == NumberWantSign {
				nv.Sign = byte(lower)
				state = NumberWantZero
				continue
			}
			if state == NumberWantExponentSign {
				nv.ExponentSign = byte(lower)
				state = NumberWantExponentDigits
				continue
			}
		}

		if lower == '.' {
			if state == NumberWantSign || state == NumberWantZero || state == NumberWantRadixSymbol || state == NumberWantIntegralDigits {
				if bufferedZero {
					nv.IntegralDigits = append(nv.IntegralDigits, '0')
					bufferedZero = false
				}
				nv.FractionalDigits = make([]byte, 0, len(input)-index)
				state = NumberWantFractionalDigits
				continue
			}
		}

		if lower == 'b' || lower == 'o' || lower == 'x' {
			if state == NumberWantRadixSymbol {
				nv.RadixSymbol = byte(lower)
				state = NumberWantIntegralDigits
				bufferedZero = false
				continue
			}
		}

		if lower == '0' {
			if state == NumberWantSign || state == NumberWantZero {
				state = NumberWantRadixSymbol
				bufferedZero = true
				continue
			}
		}

		if util.IsLegalForRadix(nv.RadixSymbol, lower) {
			if state == NumberWantSign || state == NumberWantZero || state == NumberWantRadixSymbol || state == NumberWantIntegralDigits {
				if bufferedZero {
					nv.IntegralDigits = append(nv.IntegralDigits, '0')
					bufferedZero = false
				}
				nv.IntegralDigits = append(nv.IntegralDigits, byte(lower))
				state = NumberWantIntegralDigits
				continue
			}
			if state == NumberWantFractionalDigits {
				nv.FractionalDigits = append(nv.FractionalDigits, byte(lower))
				continue
			}
		}

		if util.IsDecimalDigit(lower) {
			if state == NumberWantExponentSign || state == NumberWantExponentDigits {
				nv.ExponentDigits = append(nv.ExponentDigits, byte(lower))
				state = NumberWantExponentDigits
				continue
			}
		}

		if lower == 'e' || lower == 'p' {
			if state == NumberWantIntegralDigits || state == NumberWantFractionalDigits {
				nv.ExponentSymbol = byte(lower)
				nv.ExponentDigits = make([]byte, 0, len(input)-index)
				state = NumberWantExponentSign
				continue
			}
		}

		return &NumberParseError{
			Input: input,
			Index: uint(index),
			State: state,
		}
	}

	if state == NumberWantSign || state == NumberWantZero || state == NumberWantExponentSign {
		return &NumberParseError{
			Input: input,
			Index: uint(len(input)),
			State: state,
		}
	}

	if len(nv.IntegralDigits) == 0 {
		nv.IntegralDigits = append(nv.IntegralDigits, '0')
	}

	if nv.FractionalDigits != nil && len(nv.FractionalDigits) == 0 {
		nv.FractionalDigits = nil
	}

	nv.IntegralDigits = util.TrimLeadingZeroes(nv.IntegralDigits)
	nv.FractionalDigits = util.TrimTrailingZeroes(nv.FractionalDigits)
	nv.ExponentDigits = util.TrimLeadingZeroes(nv.ExponentDigits)

	return nil
}

func (nv *Number) IsZero() bool {
	if util.IsAllByte('0', nv.IntegralDigits) {
		if util.IsAllByte('0', nv.FractionalDigits) {
			return true
		}
	}
	return false
}

func (nv *Number) AsUint32() (uint32, error) {
	if nv.ExponentSymbol != 0 || len(nv.FractionalDigits) != 0 || len(nv.ExponentDigits) != 0 {
		return 0, fmt.Errorf("value.Number: %#v is float", nv)
	}

	if nv.IsZero() {
		return 0, nil
	}

	if nv.Sign == '-' {
		return 0, fmt.Errorf("value.Number: %#v is negative", nv)
	}

	radix := 10
	switch nv.RadixSymbol {
	case 'b':
		radix = 2
	case 'o':
		radix = 8
	case 'x':
		radix = 16
	}

	str := string(nv.IntegralDigits)
	u64, err := strconv.ParseUint(str, radix, 32)
	if err != nil {
		return 0, fmt.Errorf("value.Number: cannot parse %#v as uint32: %w", nv, err)
	}
	return uint32(u64), nil
}

func (nv *Number) AsUint64() (uint64, error) {
	if nv.ExponentSymbol != 0 || len(nv.FractionalDigits) != 0 || len(nv.ExponentDigits) != 0 {
		return 0, fmt.Errorf("value.Number: %#v is float", nv)
	}

	if nv.IsZero() {
		return 0, nil
	}

	if nv.Sign == '-' {
		return 0, fmt.Errorf("value.Number: %#v is negative", nv)
	}

	radix := 10
	switch nv.RadixSymbol {
	case 'b':
		radix = 2
	case 'o':
		radix = 8
	case 'x':
		radix = 16
	}

	str := string(nv.IntegralDigits)
	u64, err := strconv.ParseUint(str, radix, 64)
	if err != nil {
		return 0, fmt.Errorf("value.Number: cannot parse %#v as uint64: %w", nv, err)
	}
	return u64, nil
}

func (nv *Number) AsInt32() (int32, error) {
	if nv.ExponentSymbol != 0 || len(nv.FractionalDigits) != 0 || len(nv.ExponentDigits) != 0 {
		return 0, fmt.Errorf("value.Number: %#v is float", nv)
	}

	if nv.IsZero() {
		return 0, nil
	}

	radix := 10
	switch nv.RadixSymbol {
	case 'b':
		radix = 2
	case 'o':
		radix = 8
	case 'x':
		radix = 16
	}

	str := string(nv.IntegralDigits)
	if nv.Sign == '-' {
		var sb strings.Builder
		sb.Grow(1 + len(nv.IntegralDigits))
		sb.WriteByte('-')
		sb.WriteString(str)
		str = sb.String()
	}

	s64, err := strconv.ParseInt(str, radix, 32)
	if err != nil {
		return 0, fmt.Errorf("value.Number: cannot parse %#v as int32: %w", nv, err)
	}
	return int32(s64), nil
}

func (nv *Number) AsInt64() (int64, error) {
	if nv.ExponentSymbol != 0 || len(nv.FractionalDigits) != 0 || len(nv.ExponentDigits) != 0 {
		return 0, fmt.Errorf("value.Number: %#v is float", nv)
	}

	if nv.IsZero() {
		return 0, nil
	}

	radix := 10
	switch nv.RadixSymbol {
	case 'b':
		radix = 2
	case 'o':
		radix = 8
	case 'x':
		radix = 16
	}

	str := string(nv.IntegralDigits)
	if nv.Sign == '-' {
		var sb strings.Builder
		sb.Grow(1 + len(nv.IntegralDigits))
		sb.WriteByte('-')
		sb.WriteString(str)
		str = sb.String()
	}

	s64, err := strconv.ParseInt(str, radix, 64)
	if err != nil {
		return 0, fmt.Errorf("value.Number: cannot parse %#v as int64: %w", nv, err)
	}
	return s64, nil
}

var _ Value = (*Number)(nil)

// }}}

// NumberParseError
// {{{

type NumberParseError struct {
	Input []rune
	Index uint
	State NumberParseState
}

func (err *NumberParseError) Error() string {
	if err.Index >= uint(len(err.Input)) {
		return fmt.Sprintf("unexpected end of input at index %d [state=%v input=%q]", err.Index, err.State, string(err.Input))
	}
	ch := err.Input[err.Index]
	return fmt.Sprintf("unexpected character %q at index %d [state=%v input=%q]", ch, err.Index, err.State, string(err.Input))
}

var _ error = (*NumberParseError)(nil)

// }}}
