package value

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
)

// String
// {{{

type String struct {
	Segments []string
	Formats  []*FormatSpecification

	PrecomputedStringLength   uint
	PrecomputedGoStringLength uint
}

func NewEmptyString() *String {
	return NewSimpleString("")
}

func NewSimpleString(str string) *String {
	return &String{
		Segments:                  []string{str},
		Formats:                   nil,
		PrecomputedStringLength:   2 + uint(len(str)),
		PrecomputedGoStringLength: 17 + uint(len(str)),
	}
}

func (sv *String) String() string {
	return util.StringImpl(sv)
}

func (sv *String) GoString() string {
	return util.GoStringImpl(sv)
}

func (sv *String) EstimateStringLength() uint {
	return sv.PrecomputedStringLength
}

func (sv *String) EstimateGoStringLength() uint {
	return sv.PrecomputedGoStringLength
}

func (sv *String) WriteStringTo(out *strings.Builder) {
	out.WriteByte('"')
	i := 0
	for i < len(sv.Formats) {
		writeQuotedSegment(out, sv.Segments[i])
		sv.Formats[i].WriteStringTo(out)
		i++
	}
	writeQuotedSegment(out, sv.Segments[i])
	out.WriteByte('"')
}

func (sv *String) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("&value.String{")
	i := 0
	for i < len(sv.Formats) {
		out.WriteString(strconv.Quote(sv.Segments[i]))
		out.WriteByte(',')
		sv.Formats[i].WriteGoStringTo(out)
		out.WriteByte(',')
		i++
	}
	out.WriteString(strconv.Quote(sv.Segments[i]))
	out.WriteByte('}')
}

func (sv *String) SetEstimatedStringLength(length uint) {
	sv.PrecomputedStringLength = length
}

func (sv *String) SetEstimatedGoStringLength(length uint) {
	sv.PrecomputedGoStringLength = length
}

func (sv *String) Parse(input []rune) error {
	*sv = String{
		Segments: make([]string, 0, 4),
		Formats:  make([]*FormatSpecification, 0, 3),
	}

	var segmentBuffer strings.Builder
	segmentBuffer.Grow(len(input))

	var format *FormatSpecification
	var formatBuffer strings.Builder
	var formatIndex uint
	formatBuffer.Grow(20)

	var escRunes []rune = make([]rune, 0, 10)
	var escRadix util.Radix
	var escLength, escMinLength, escMaxLength, escMaxValue uint

	flushEscape := func() error {
		u64, err := escRadix.Parse(string(escRunes[2:]))
		if err != nil {
			return err
		}
		if escLength < escMinLength || escLength > escMaxLength || u64 > uint64(escMaxValue) {
			return &StringEscapeParseError{
				Escape:    escRunes,
				Value:     u64,
				MinLength: escMinLength,
				MaxLength: escMaxLength,
				MaxValue:  escMaxValue,
			}
		}
		segmentBuffer.WriteRune(rune(u64))

		escRunes = escRunes[:0]
		escRadix = nil
		escLength = 0
		escMinLength = 0
		escMaxLength = 0
		escMaxValue = 0
		return nil
	}

	flushSegment := func() {
		format = new(FormatSpecification)
		sv.Segments = append(sv.Segments, segmentBuffer.String())
		sv.Formats = append(sv.Formats, format)
		segmentBuffer.Reset()
	}

	flushFormat := func() {
		util.EstimateLengths(format)
		format = nil
	}

	state := StringWantLeadingQuote
	for index, ch := range input {
		if util.IsFormatFlag(ch) {
			if state == StringDQGotPercent {
				flushSegment()
			}
			if state == StringDQGotPercent || state == StringDQWantFmtFlags {
				format.SetFlag(ch, true)
				state = StringDQWantFmtFlags
				continue
			}
		}
		if (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch == '_') {
			if state == StringDQWantFmtArgMaybeWidth || state == StringDQWantFmtArgMaybePrec || state == StringDQWantFmtArgMaybeValue {
				formatBuffer.WriteByte(byte(ch))
				continue
			}
		}
		if util.IsDecimalDigit(ch) {
			if state == StringDQGotPercent {
				flushSegment()
			}
			if state == StringDQGotPercent || state == StringDQWantFmtFlags || state == StringDQWantFmtWidthDigits {
				formatBuffer.WriteByte(byte(ch))
				state = StringDQWantFmtWidthDigits
				continue
			}
			if state == StringDQWantFmtPrec || state == StringDQWantFmtPrecDigits {
				formatBuffer.WriteByte(byte(ch))
				state = StringDQWantFmtPrecDigits
				continue
			}
		}
		if ch == '[' {
			if state == StringDQGotPercent {
				flushSegment()
			}
			if state == StringDQWantFmtWidthDigits {
				format.HasWidth = true
				format.FixedWidth = util.MustAtoi(formatBuffer.String())
				formatBuffer.Reset()
			}
			if state == StringDQWantFmtPrecDigits {
				format.HasPrecision = true
				format.FixedPrecision = util.MustAtoi(formatBuffer.String())
				formatBuffer.Reset()
			}
			if state == StringDQGotPercent || state == StringDQWantFmtFlags {
				state = StringDQWantFmtArgMaybeWidth
				continue
			}
			if state == StringDQWantFmtDot || state == StringDQWantFmtPrecDigits || state == StringDQWantFmtWidthDigits || state == StringDQWantFmtConv {
				state = StringDQWantFmtArgMaybeValue
				continue
			}
			if state == StringDQWantFmtPrec {
				state = StringDQWantFmtArgMaybePrec
				continue
			}
		}
		if ch == ']' {
			if state == StringDQWantFmtArgMaybeWidth {
				state = StringDQWantFmtArgMaybeWidthHaveRBracket
				continue
			}
			if state == StringDQWantFmtArgMaybePrec {
				state = StringDQWantFmtArgMaybePrecHaveRBracket
				continue
			}
			if state == StringDQWantFmtArgMaybeValue {
				state = StringDQWantFmtArgMaybeValueHaveRBracket
				continue
			}
		}
		if ch == '*' {
			if state == StringDQGotPercent {
				flushSegment()
			}
			if state == StringDQGotPercent || state == StringDQWantFmtFlags {
				formatIndex++
				format.HasWidth = true
				format.WidthIsExternal = true
				format.WidthArgumentIndex = formatIndex
				format.WidthArgumentName = util.Itoa(formatIndex)
				state = StringDQWantFmtDot
				continue
			}
			if state == StringDQWantFmtPrec {
				formatIndex++
				format.HasPrecision = true
				format.PrecisionIsExternal = true
				format.PrecisionArgumentIndex = formatIndex
				format.PrecisionArgumentName = util.Itoa(formatIndex)
				state = StringDQWantFmtConv
				continue
			}
			if state == StringDQWantFmtArgMaybeWidthHaveRBracket {
				name := formatBuffer.String()
				formatBuffer.Reset()

				format.HasWidth = true
				format.WidthIsExternal = true
				format.WidthArgumentIsExplicit = true
				if num, ok := util.TryAtoi(name); ok {
					format.WidthArgumentIndex = num
					format.WidthArgumentName = util.Itoa(num)
					formatIndex = num
				} else {
					format.WidthArgumentIsNamed = true
					format.WidthArgumentName = name
				}

				state = StringDQWantFmtDot
				continue
			}
			if state == StringDQWantFmtArgMaybePrecHaveRBracket {
				name := formatBuffer.String()
				formatBuffer.Reset()

				format.HasPrecision = true
				format.PrecisionIsExternal = true
				format.PrecisionArgumentIsExplicit = true
				if num, ok := util.TryAtoi(name); ok {
					format.PrecisionArgumentIndex = num
					format.PrecisionArgumentName = util.Itoa(num)
					formatIndex = num
				} else {
					format.PrecisionArgumentIsNamed = true
					format.PrecisionArgumentName = name
				}

				state = StringDQWantFmtConv
				continue
			}
		}
		if ch == '.' {
			if state == StringDQGotPercent {
				flushSegment()
			}
			if state == StringDQWantFmtWidthDigits {
				format.HasWidth = true
				format.FixedWidth = util.MustAtoi(formatBuffer.String())
				formatBuffer.Reset()
			}
			if state == StringDQGotPercent || state == StringDQWantFmtFlags || state == StringDQWantFmtWidthDigits || state == StringDQWantFmtDot {
				state = StringDQWantFmtPrec
				continue
			}
		}
		if _, ok := formatConversions[ch]; ok {
			if state == StringDQGotPercent {
				flushSegment()
			}
			if state == StringDQWantFmtWidthDigits {
				format.HasWidth = true
				format.FixedWidth = util.MustAtoi(formatBuffer.String())
				formatBuffer.Reset()
			}
			if state == StringDQWantFmtPrecDigits {
				format.HasPrecision = true
				format.FixedPrecision = util.MustAtoi(formatBuffer.String())
				formatBuffer.Reset()
			}
			if state == StringDQGotPercent || state == StringDQWantFmtFlags || state == StringDQWantFmtWidthDigits || state == StringDQWantFmtDot || state == StringDQWantFmtPrec || state == StringDQWantFmtPrecDigits || state == StringDQWantFmtConv {
				formatIndex++
				format.ValueArgumentIndex = formatIndex
				format.ValueArgumentName = util.Itoa(formatIndex)
				format.Conversion = byte(ch)
				flushFormat()
				state = StringDQReady
				continue
			}
			if state == StringDQWantFmtArgMaybeWidthHaveRBracket || state == StringDQWantFmtArgMaybePrecHaveRBracket || state == StringDQWantFmtArgMaybeValueHaveRBracket {
				name := formatBuffer.String()
				formatBuffer.Reset()

				format.Conversion = byte(ch)
				format.ValueArgumentIsExplicit = true
				if num, ok := util.TryAtoi(name); ok {
					format.ValueArgumentIndex = num
					format.ValueArgumentName = util.Itoa(formatIndex)
					formatIndex = num
				} else {
					format.ValueArgumentIsNamed = true
					format.ValueArgumentName = name
				}

				flushFormat()
				state = StringDQReady
				continue
			}
		}

		if state == StringSQGotBackslash || state == StringDQGotBackslash {
			nextReady := StringSQReady
			nextEscape := StringSQWantEscape
			if state == StringDQGotBackslash {
				nextReady = StringDQReady
				nextEscape = StringDQWantEscape
			}

			if outCh, found := simpleEscapes[ch]; found {
				segmentBuffer.WriteRune(outCh)
				state = nextReady
				continue
			}
			if ch == 'd' {
				state = nextEscape
				escRunes = append(escRunes, '\\', 'd')
				escRadix = util.RadixTable[10]
				escLength = 0
				escMinLength = 1
				escMaxLength = 3
				escMaxValue = 127
				continue
			}
			if ch == 'o' {
				state = nextEscape
				escRunes = append(escRunes, '\\', 'o')
				escRadix = util.RadixTable[8]
				escLength = 0
				escMinLength = 1
				escMaxLength = 3
				escMaxValue = 127
				continue
			}
			if ch == 'x' {
				state = nextEscape
				escRunes = append(escRunes, '\\', 'x')
				escRadix = util.RadixTable[16]
				escLength = 0
				escMinLength = 1
				escMaxLength = 2
				escMaxValue = 127
				continue
			}
			if ch == 'u' {
				state = nextEscape
				escRunes = append(escRunes, '\\', 'u')
				escRadix = util.RadixTable[16]
				escLength = 0
				escMinLength = 4
				escMaxLength = 4
				escMaxValue = uint(unicode.MaxRune)
				continue
			}
			if ch == 'U' {
				state = nextEscape
				escRunes = append(escRunes, '\\', 'U')
				escRadix = util.RadixTable[16]
				escLength = 0
				escMinLength = 4
				escMaxLength = 8
				escMaxValue = uint(unicode.MaxRune)
				continue
			}
		}

		if state == StringSQWantEscape || state == StringDQWantEscape {
			nextReady := StringSQReady
			if state == StringDQWantEscape {
				nextReady = StringDQReady
			}

			lower := unicode.ToLower(ch)
			if escRadix.MatchRune(lower) {
				escRunes = append(escRunes, lower)
				escLength++
				if escLength >= escMaxLength {
					if err := flushEscape(); err != nil {
						return err
					}
					state = nextReady
				}
				continue
			}

			if err := flushEscape(); err != nil {
				return err
			}
			state = nextReady
			// fallthrough
		}

		if state == StringWantLeadingQuote && ch == '\'' {
			state = StringSQReady
			continue
		}
		if state == StringWantLeadingQuote && ch == '"' {
			state = StringDQReady
			continue
		}
		if state == StringDQGotPercent && ch == '%' {
			segmentBuffer.WriteRune(ch)
			continue
		}
		if state == StringSQReady && ch == '\'' {
			state = StringWantEnd
			continue
		}
		if state == StringDQReady && ch == '"' {
			state = StringWantEnd
			continue
		}
		if state == StringSQReady && ch == '\\' {
			state = StringSQGotBackslash
			continue
		}
		if state == StringDQReady && ch == '\\' {
			state = StringDQGotBackslash
			continue
		}
		if state == StringDQReady && ch == '%' {
			state = StringDQGotPercent
			continue
		}

		if state == StringSQReady || state == StringDQReady {
			segmentBuffer.WriteRune(ch)
			continue
		}

		return &StringParseError{
			Input: input,
			Index: uint(index),
			State: state,
		}
	}

	if state != StringWantEnd {
		return &StringParseError{
			Input: input,
			Index: uint(len(input)),
			State: state,
		}
	}

	sv.Segments = append(sv.Segments, segmentBuffer.String())
	sv.PrecomputedGoStringLength = 13
	for _, seg := range sv.Segments {
		sv.PrecomputedStringLength += uint(len(seg))
		sv.PrecomputedGoStringLength += uint(len(seg)) + 3
	}
	for _, format := range sv.Formats {
		sv.PrecomputedStringLength += format.PrecomputedStringLength
		sv.PrecomputedGoStringLength += format.PrecomputedGoStringLength + 1
	}
	return nil
}

var _ Value = (*String)(nil)
var _ util.Estimable = (*String)(nil)

// }}}

// StringParseError
// {{{

type StringParseError struct {
	Input []rune
	Index uint
	State StringParseState
}

func (err *StringParseError) Error() string {
	if err.Index >= uint(len(err.Input)) {
		return fmt.Sprintf("unexpected end of input at index %d [state=%v input=%q]", err.Index, err.State, string(err.Input))
	}
	ch := err.Input[err.Index]
	return fmt.Sprintf("unexpected character %q at index %d [state=%v input=%q]", ch, err.Index, err.State, string(err.Input))
}

var _ error = (*StringParseError)(nil)

// }}}

// StringEscapeParseError
// {{{

type StringEscapeParseError struct {
	Escape    []rune
	Value     uint64
	MinLength uint
	MaxLength uint
	MaxValue  uint
}

func (err *StringEscapeParseError) Error() string {
	escString := string(err.Escape)
	escLength := uint(len(err.Escape)) - 2
	if escLength < err.MinLength {
		return fmt.Sprintf("escape sequence is too short: %s: %d < %d", escString, escLength, err.MinLength)
	}
	if escLength > err.MaxLength {
		return fmt.Sprintf("escape sequence is too long: %s: %d > %d", escString, escLength, err.MaxLength)
	}
	return fmt.Sprintf("escape sequence is out of range: %s: U+%04X > U+%04X", escString, err.Value, err.MaxValue)
}

var _ error = (*StringEscapeParseError)(nil)

// }}}

func writeQuotedSegment(out *strings.Builder, segment string) {
	for _, ch := range segment {
		switch {
		case ch == 0x07:
			out.WriteString("\\a")
		case ch == 0x08:
			out.WriteString("\\b")
		case ch == 0x09:
			out.WriteString("\\t")
		case ch == 0x0a:
			out.WriteString("\\n")
		case ch == 0x0b:
			out.WriteString("\\v")
		case ch == 0x0c:
			out.WriteString("\\f")
		case ch == 0x0d:
			out.WriteString("\\r")
		case ch == 0x1b:
			out.WriteString("\\e")
		case ch == 0x22:
			out.WriteString("\\\"")
		case ch == 0x25:
			out.WriteString("%%")
		case ch == 0x5c:
			out.WriteString("\\\\")
		case ch < 0x80 && unicode.IsControl(ch):
			fmt.Fprintf(out, "\\x%02x", ch)
		case ch < 0x10000 && unicode.IsControl(ch):
			fmt.Fprintf(out, "\\u%04x", ch)
		case unicode.IsControl(ch):
			fmt.Fprintf(out, "\\U%08x", ch)
		case ch > unicode.MaxRune:
			fmt.Fprintf(out, "\\U%08x", ch)
		default:
			out.WriteRune(ch)
		}
	}
}

var simpleEscapes = map[rune]rune{
	0x0009: 0x0009, // \TAB
	0x000a: 0x0020, // \LF -> SP
	0x000d: 0x0020, // \CR -> SP
	0x0020: 0x0020, // \SP
	0x0021: 0x0021, // \!
	0x0022: 0x0022, // \"
	0x0023: 0x0023, // \#
	0x0024: 0x0024, // \$
	0x0025: 0x0025, // \%
	0x0026: 0x0026, // \&
	0x0027: 0x0027, // \'
	0x0028: 0x0028, // \(
	0x0029: 0x0029, // \)
	0x002a: 0x002a, // \*
	0x002b: 0x002b, // \+
	0x002c: 0x002c, // \,
	0x002d: 0x002d, // \-
	0x002e: 0x002e, // \.
	0x002f: 0x002f, // \/
	0x0030: 0x0000, // \0 -> NUL
	0x003a: 0x003a, // \:
	0x003b: 0x003b, // \;
	0x003c: 0x003c, // \<
	0x003d: 0x003d, // \=
	0x003e: 0x003e, // \>
	0x003f: 0x003f, // \?
	0x0040: 0x0040, // \@
	0x005b: 0x005b, // \[
	0x005c: 0x005c, // \\
	0x005d: 0x005d, // \]
	0x005e: 0x005e, // \^
	0x005f: 0x005f, // \_
	0x0060: 0x0060, // \`
	0x0061: 0x0007, // \a -> BEL
	0x0062: 0x0008, // \b -> BS
	0x0065: 0x001b, // \e -> ESC
	0x0066: 0x000c, // \f -> FF
	0x006e: 0x000a, // \n -> LF
	0x0072: 0x000d, // \r -> CR
	0x0074: 0x0009, // \t -> TAB
	0x0076: 0x000b, // \v -> VT
	0x007b: 0x007b, // \{
	0x007c: 0x007c, // \|
	0x007d: 0x007d, // \}
	0x007e: 0x007e, // \~
}

var formatConversions = map[rune]struct{}{
	'E': {}, // float: -1.234456E+78
	'F': {}, // float: -1.234456
	'G': {}, // float: -1.234456 or -1.234456E+78
	'O': {}, // int: "0o" + base 8
	'T': {}, // type name
	'U': {}, // "U+%04X"
	'X': {}, // int: base 16 uppercase, float: -0X1.23ABCP+20, string: hex bytes uppercase
	'b': {}, // int: base 2, float: -123456p-78
	'c': {}, // rune: literal
	'd': {}, // int: 255
	'e': {}, // float: -1.234456e+78
	'f': {}, // float: -1.234456
	'g': {}, // float: -1.234456 or -1.234456e+78
	'o': {}, // int: base 8
	'p': {}, // pointer: 0xdeadbeef
	'q': {}, // string: quoted literal, rune: quoted literal
	's': {}, // string: literal
	't': {}, // bool: true
	'v': {},
	'x': {}, // int: base 16 lowercase, float: -0x1.23abcp+20, string: hex bytes lowercase
}
