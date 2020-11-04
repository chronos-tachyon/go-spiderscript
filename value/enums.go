package value

import (
	"fmt"
)

// NumberParseState
// {{{

type NumberParseState uint

const (
	NumberWantSign NumberParseState = iota
	NumberWantZero
	NumberWantRadixSymbol
	NumberWantIntegralDigits
	NumberWantFractionalDigits
	NumberWantExponentSign
	NumberWantExponentDigits
)

var numberValueParseStateNames = []string{
	"NumberWantSign",
	"NumberWantZero",
	"NumberWantRadixSymbol",
	"NumberWantIntegralDigits",
	"NumberWantFractionalDigits",
	"NumberWantExponentSign",
	"NumberWantExponentDigits",
}

func (enum NumberParseState) String() string {
	if uint(enum) >= uint(len(numberValueParseStateNames)) {
		return fmt.Sprintf("NumberParseState(%d)", uint(enum))
	}
	return numberValueParseStateNames[enum]
}

func (enum NumberParseState) GoString() string {
	return enum.String()
}

var _ fmt.Stringer = NumberParseState(0)
var _ fmt.GoStringer = NumberParseState(0)

// }}}

// StringParseState
// {{{

type StringParseState uint

const (
	StringWantLeadingQuote StringParseState = iota

	StringSQReady
	StringSQGotBackslash
	StringSQWantEscape

	StringDQReady
	StringDQGotBackslash
	StringDQWantEscape

	StringDQGotPercent
	StringDQWantFmtFlags
	StringDQWantFmtWidthDigits
	StringDQWantFmtDot
	StringDQWantFmtPrec
	StringDQWantFmtPrecDigits
	StringDQWantFmtConv
	StringDQWantFmtArgMaybeWidth
	StringDQWantFmtArgMaybeWidthHaveRBracket
	StringDQWantFmtArgMaybePrec
	StringDQWantFmtArgMaybePrecHaveRBracket
	StringDQWantFmtArgMaybeValue
	StringDQWantFmtArgMaybeValueHaveRBracket

	StringWantEnd
)

var stringValueParseStateNames = []string{
	"StringWantLeadingQuote",
	"StringSQReady",
	"StringSQGotBackslash",
	"StringSQWantEscape",
	"StringDQReady",
	"StringDQGotBackslash",
	"StringDQWantEscape",
	"StringDQGotPercent",
	"StringDQWantFmtFlags",
	"StringDQWantFmtWidthDigits",
	"StringDQWantFmtDot",
	"StringDQWantFmtPrec",
	"StringDQWantFmtPrecDigits",
	"StringDQWantFmtConv",
	"StringDQWantFmtArgMaybeWidth",
	"StringDQWantFmtArgMaybeWidthHaveRBracket",
	"StringDQWantFmtArgMaybePrec",
	"StringDQWantFmtArgMaybePrecHaveRBracket",
	"StringDQWantFmtArgMaybeValue",
	"StringDQWantFmtArgMaybeValueHaveRBracket",
	"StringWantEnd",
}

func (enum StringParseState) String() string {
	if uint(enum) >= uint(len(stringValueParseStateNames)) {
		return fmt.Sprintf("StringParseState(%d)", uint(enum))
	}
	return stringValueParseStateNames[enum]
}

func (enum StringParseState) GoString() string {
	return enum.String()
}

var _ fmt.Stringer = StringParseState(0)
var _ fmt.GoStringer = StringParseState(0)

// }}}
