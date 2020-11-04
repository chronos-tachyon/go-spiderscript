package token

import (
	"fmt"
)

type State uint

const (
	StateZero State = iota
	StateIllegalRune

	StateDone
	StateUnreadAndDone

	StateReadyAtStartOfFile
	StateReady

	StateInShebang
	StateInSingleLineComment
	StateInMultiLineComment
	StateInMultiLineCommentWithStar

	StateInExclaim
	StateInHashAtStartOfFile
	StateInHash
	StateInPercent
	StateInAmpersand
	StateInAmpersandAmpersand
	StateInStar
	StateInStarStar
	StateInPlus
	StateInMinus
	StateInDot
	StateInDotDot
	StateInSlash
	StateInColon
	StateInLess
	StateInLessEqual
	StateInLessLess
	StateInLessLessBar
	StateInEqual
	StateInGreater
	StateInGreaterGreater
	StateInGreaterGreaterBar
	StateInQuestion
	StateInCaret
	StateInCaretCaret
	StateInBar
	StateInBarBar
	StateInTilde

	StateInPragma
	StateInIdentifier
	StateInNumber
	StateInSingleQuoteString
	StateInSingleQuoteStringWithBackslash
	StateInDoubleQuoteString
	StateInDoubleQuoteStringWithBackslash
	StateInTickQuoteString
	StateInTickQuoteStringWithBackslash
	StateInRegexSlash
	StateInRegexSlashWithBackslash
	StateInRegexExclaim
	StateInRegexExclaimWithBackslash
	StateInRegexAt
	StateInRegexAtWithBackslash
	StateInRegexBrace
	StateInRegexBraceWithBackslash
	StateInRegexFlags
	StateInPEG
	StateInPEGWithBackslash
)

var stateNames = []string{
	"StateZero",
	"StateIllegalRune",
	"StateDone",
	"StateUnreadAndDone",
	"StateReadyAtStartOfFile",
	"StateReady",
	"StateInShebang",
	"StateInSingleLineComment",
	"StateInMultiLineComment",
	"StateInMultiLineCommentWithStar",
	"StateInExclaim",
	"StateInHashAtStartOfFile",
	"StateInHash",
	"StateInPercent",
	"StateInAmpersand",
	"StateInAmpersandAmpersand",
	"StateInStar",
	"StateInStarStar",
	"StateInPlus",
	"StateInMinus",
	"StateInDot",
	"StateInDotDot",
	"StateInSlash",
	"StateInColon",
	"StateInLess",
	"StateInLessEqual",
	"StateInLessLess",
	"StateInLessLessBar",
	"StateInEqual",
	"StateInGreater",
	"StateInGreaterGreater",
	"StateInGreaterGreaterBar",
	"StateInQuestion",
	"StateInCaret",
	"StateInCaretCaret",
	"StateInBar",
	"StateInBarBar",
	"StateInTilde",
	"StateInPragma",
	"StateInIdentifier",
	"StateInNumber",
	"StateInSingleQuoteString",
	"StateInSingleQuoteStringWithBackslash",
	"StateInDoubleQuoteString",
	"StateInDoubleQuoteStringWithBackslash",
	"StateInTickQuoteString",
	"StateInTickQuoteStringWithBackslash",
	"StateInRegexSlash",
	"StateInRegexSlashWithBackslash",
	"StateInRegexExclaim",
	"StateInRegexExclaimWithBackslash",
	"StateInRegexAt",
	"StateInRegexAtWithBackslash",
	"StateInRegexBrace",
	"StateInRegexBraceWithBackslash",
	"StateInRegexFlags",
	"StateInPEG",
	"StateInPEGWithBackslash",
}

func (enum State) String() string {
	if uint(enum) >= uint(len(stateNames)) {
		return fmt.Sprintf("State(%d)", uint(enum))
	}
	return stateNames[enum]
}

func (enum State) GoString() string {
	return enum.String()
}

var _ fmt.Stringer = State(0)
var _ fmt.GoStringer = State(0)
