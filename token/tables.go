package token

import (
	"github.com/chronos-tachyon/go-spiderscript/internal/runepredicate"
	"github.com/chronos-tachyon/go-spiderscript/internal/util"
)

type stateTableRow struct {
	CurrentState statePredicate
	CurrentRune  runepredicate.RunePredicate
	CurrentRaw   []rune
	Delta        int8
	OnlyIfZero   bool
	NextState    State
	TokenType    Type
}

var stateTable = []stateTableRow{
	{
		CurrentState: anyState(),
		CurrentRune:  runepredicate.Func(util.IsIllegal),
		NextState:    StateIllegalRune,
	},

	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Func(util.IsHWS),
		NextState:    StateDone,
		TokenType:    HWS,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Func(util.IsVWS),
		NextState:    StateDone,
		TokenType:    VWS,
	},

	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('!'),
		NextState:    StateInExclaim,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('"'),
		NextState:    StateInDoubleQuoteString,
	},
	{
		CurrentState: exactState(StateReadyAtStartOfFile),
		CurrentRune:  runepredicate.Exactly('#'),
		NextState:    StateInHashAtStartOfFile,
	},
	{
		CurrentState: exactState(StateReady),
		CurrentRune:  runepredicate.Exactly('#'),
		NextState:    StateInHash,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('%'),
		NextState:    StateInPercent,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('&'),
		NextState:    StateInAmpersand,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('\''),
		NextState:    StateInSingleQuoteString,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('('),
		NextState:    StateDone,
		TokenType:    LParen,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly(')'),
		NextState:    StateDone,
		TokenType:    RParen,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('*'),
		NextState:    StateInStar,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('+'),
		NextState:    StateInPlus,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly(','),
		NextState:    StateDone,
		TokenType:    Comma,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('-'),
		NextState:    StateInMinus,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('.'),
		NextState:    StateInDot,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('/'),
		NextState:    StateInSlash,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly(':'),
		NextState:    StateInColon,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly(';'),
		NextState:    StateDone,
		TokenType:    Semicolon,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('<'),
		NextState:    StateInLess,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateInEqual,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('>'),
		NextState:    StateInGreater,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('?'),
		NextState:    StateInQuestion,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('@'),
		NextState:    StateDone,
		TokenType:    At,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('['),
		NextState:    StateDone,
		TokenType:    LBracket,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly(']'),
		NextState:    StateDone,
		TokenType:    RBracket,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('^'),
		NextState:    StateInCaret,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('{'),
		NextState:    StateDone,
		TokenType:    LBrace,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('|'),
		NextState:    StateInBar,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('}'),
		NextState:    StateDone,
		TokenType:    RBrace,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Exactly('~'),
		NextState:    StateInTilde,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Func(util.IsIdentStart),
		NextState:    StateInIdentifier,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Func(util.IsDecimalDigit),
		NextState:    StateInNumber,
	},
	{
		CurrentState: someState(StateReadyAtStartOfFile, StateReady),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateDone,
		TokenType:    Invalid,
	},

	{
		CurrentState: exactState(StateInExclaim),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    BangEqual,
	},
	{
		CurrentState: exactState(StateInExclaim),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Bang,
	},

	{
		CurrentState: exactState(StateInPercent),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    PercentEqual,
	},
	{
		CurrentState: exactState(StateInPercent),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Percent,
	},

	{
		CurrentState: exactState(StateInAmpersand),
		CurrentRune:  runepredicate.Exactly('&'),
		NextState:    StateInAmpersandAmpersand,
	},
	{
		CurrentState: exactState(StateInAmpersand),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    AmpersandEqual,
	},
	{
		CurrentState: exactState(StateInAmpersand),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Ampersand,
	},
	{
		CurrentState: exactState(StateInAmpersandAmpersand),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    AmpersandAmpersandEqual,
	},
	{
		CurrentState: exactState(StateInAmpersandAmpersand),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    AmpersandAmpersand,
	},

	{
		CurrentState: exactState(StateInStar),
		CurrentRune:  runepredicate.Exactly('*'),
		NextState:    StateInStarStar,
	},
	{
		CurrentState: exactState(StateInStar),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    StarEqual,
	},
	{
		CurrentState: exactState(StateInStar),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Star,
	},
	{
		CurrentState: exactState(StateInStarStar),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    StarStarEqual,
	},
	{
		CurrentState: exactState(StateInStarStar),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    StarStar,
	},

	{
		CurrentState: exactState(StateInPlus),
		CurrentRune:  runepredicate.Exactly('+'),
		NextState:    StateDone,
		TokenType:    PlusPlus,
	},
	{
		CurrentState: exactState(StateInPlus),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    PlusEqual,
	},
	{
		CurrentState: exactState(StateInPlus),
		CurrentRune:  runepredicate.Func(util.IsDecimalDigit),
		NextState:    StateInNumber,
	},
	{
		CurrentState: exactState(StateInPlus),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Plus,
	},

	{
		CurrentState: exactState(StateInMinus),
		CurrentRune:  runepredicate.Exactly('-'),
		NextState:    StateDone,
		TokenType:    MinusMinus,
	},
	{
		CurrentState: exactState(StateInMinus),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    MinusEqual,
	},
	{
		CurrentState: exactState(StateInMinus),
		CurrentRune:  runepredicate.Func(util.IsDecimalDigit),
		NextState:    StateInNumber,
	},
	{
		CurrentState: exactState(StateInMinus),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Minus,
	},

	{
		CurrentState: exactState(StateInDot),
		CurrentRune:  runepredicate.Exactly('.'),
		NextState:    StateInDotDot,
	},
	{
		CurrentState: exactState(StateInDot),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Dot,
	},

	{
		CurrentState: exactState(StateInDotDot),
		CurrentRune:  runepredicate.Exactly('.'),
		NextState:    StateDone,
		TokenType:    DotDotDot,
	},
	{
		CurrentState: exactState(StateInDotDot),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    DotDot,
	},

	{
		CurrentState: exactState(StateInSlash),
		CurrentRune:  runepredicate.Exactly('*'),
		NextState:    StateInMultiLineComment,
	},
	{
		CurrentState: exactState(StateInSlash),
		CurrentRune:  runepredicate.Exactly('/'),
		NextState:    StateInSingleLineComment,
	},
	{
		CurrentState: exactState(StateInSlash),
		CurrentRune:  runepredicate.Exactly('%'),
		NextState:    StateDone,
		TokenType:    SlashPercent,
	},
	{
		CurrentState: exactState(StateInSlash),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    SlashEqual,
	},
	{
		CurrentState: exactState(StateInSlash),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Slash,
	},

	{
		CurrentState: exactState(StateInColon),
		CurrentRune:  runepredicate.Exactly(':'),
		NextState:    StateDone,
		TokenType:    ColonColon,
	},
	{
		CurrentState: exactState(StateInColon),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    ColonEqual,
	},
	{
		CurrentState: exactState(StateInColon),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Colon,
	},

	{
		CurrentState: exactState(StateInLess),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateInLessEqual,
	},
	{
		CurrentState: exactState(StateInLess),
		CurrentRune:  runepredicate.Exactly('<'),
		NextState:    StateInLessLess,
	},
	{
		CurrentState: exactState(StateInLess),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Less,
	},
	{
		CurrentState: exactState(StateInLessEqual),
		CurrentRune:  runepredicate.Exactly('>'),
		NextState:    StateDone,
		TokenType:    LessEqualGreater,
	},
	{
		CurrentState: exactState(StateInLessEqual),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    LessEqual,
	},
	{
		CurrentState: exactState(StateInLessLess),
		CurrentRune:  runepredicate.Exactly('|'),
		NextState:    StateInLessLessBar,
	},
	{
		CurrentState: exactState(StateInLessLess),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    LessLessEqual,
	},
	{
		CurrentState: exactState(StateInLessLess),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    LessLess,
	},
	{
		CurrentState: exactState(StateInLessLessBar),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    LessLessBarEqual,
	},
	{
		CurrentState: exactState(StateInLessLessBar),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    LessLessBar,
	},

	{
		CurrentState: exactState(StateInEqual),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    EqualEqual,
	},
	{
		CurrentState: exactState(StateInEqual),
		CurrentRune:  runepredicate.Exactly('>'),
		NextState:    StateDone,
		TokenType:    EqualGreater,
	},
	{
		CurrentState: exactState(StateInEqual),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Equal,
	},

	{
		CurrentState: exactState(StateInGreater),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    GreaterEqual,
	},
	{
		CurrentState: exactState(StateInGreater),
		CurrentRune:  runepredicate.Exactly('>'),
		NextState:    StateInGreaterGreater,
	},
	{
		CurrentState: exactState(StateInGreater),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Greater,
	},
	{
		CurrentState: exactState(StateInGreaterGreater),
		CurrentRune:  runepredicate.Exactly('|'),
		NextState:    StateInGreaterGreaterBar,
	},
	{
		CurrentState: exactState(StateInGreaterGreater),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    GreaterGreaterEqual,
	},
	{
		CurrentState: exactState(StateInGreaterGreater),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    GreaterGreater,
	},
	{
		CurrentState: exactState(StateInGreaterGreaterBar),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    GreaterGreaterBarEqual,
	},
	{
		CurrentState: exactState(StateInGreaterGreaterBar),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    GreaterGreaterBar,
	},

	{
		CurrentState: exactState(StateInQuestion),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    QuestionEqual,
	},
	{
		CurrentState: exactState(StateInQuestion),
		CurrentRune:  runepredicate.Exactly(':'),
		NextState:    StateDone,
		TokenType:    QuestionColon,
	},
	{
		CurrentState: exactState(StateInQuestion),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Question,
	},

	{
		CurrentState: exactState(StateInCaret),
		CurrentRune:  runepredicate.Exactly('^'),
		NextState:    StateInCaretCaret,
	},
	{
		CurrentState: exactState(StateInCaret),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    CaretEqual,
	},
	{
		CurrentState: exactState(StateInCaret),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Caret,
	},
	{
		CurrentState: exactState(StateInCaretCaret),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    CaretCaretEqual,
	},
	{
		CurrentState: exactState(StateInCaretCaret),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    CaretCaret,
	},

	{
		CurrentState: exactState(StateInBar),
		CurrentRune:  runepredicate.Exactly('|'),
		NextState:    StateInBarBar,
	},
	{
		CurrentState: exactState(StateInBar),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    BarEqual,
	},
	{
		CurrentState: exactState(StateInBar),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Bar,
	},
	{
		CurrentState: exactState(StateInBarBar),
		CurrentRune:  runepredicate.Exactly('='),
		NextState:    StateDone,
		TokenType:    BarBarEqual,
	},
	{
		CurrentState: exactState(StateInBarBar),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    BarBar,
	},

	{
		CurrentState: exactState(StateInTilde),
		CurrentRune:  runepredicate.Exactly('~'),
		NextState:    StateDone,
		TokenType:    TildeTilde,
	},
	{
		CurrentState: exactState(StateInTilde),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Tilde,
	},

	{
		CurrentState: exactState(StateInHashAtStartOfFile),
		CurrentRune:  runepredicate.Exactly('!'),
		NextState:    StateInShebang,
	},
	{
		CurrentState: someState(StateInHashAtStartOfFile, StateInHash),
		CurrentRune:  runepredicate.Func(util.IsIdentStart),
		NextState:    StateInPragma,
	},
	{
		CurrentState: someState(StateInHashAtStartOfFile, StateInHash),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Hash,
	},

	{
		CurrentState: exactState(StateInShebang),
		CurrentRune:  runepredicate.Func(util.IsVWS),
		NextState:    StateUnreadAndDone,
		TokenType:    ShebangLine,
	},
	{
		CurrentState: exactState(StateInShebang),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateInShebang,
	},

	{
		CurrentState: exactState(StateInSingleLineComment),
		CurrentRune:  runepredicate.Func(util.IsVWS),
		NextState:    StateUnreadAndDone,
		TokenType:    SingleLineComment,
	},
	{
		CurrentState: exactState(StateInSingleLineComment),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateInSingleLineComment,
	},

	{
		CurrentState: exactState(StateInMultiLineCommentWithStar),
		CurrentRune:  runepredicate.Exactly('/'),
		NextState:    StateDone,
		TokenType:    MultiLineComment,
	},
	{
		CurrentState: someState(StateInMultiLineComment, StateInMultiLineCommentWithStar),
		CurrentRune:  runepredicate.Exactly('*'),
		NextState:    StateInMultiLineCommentWithStar,
	},
	{
		CurrentState: someState(StateInMultiLineComment, StateInMultiLineCommentWithStar),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateInMultiLineComment,
	},

	{
		CurrentState: exactState(StateInIdentifier),
		CurrentRune:  runepredicate.Func(util.IsIdentContinue),
		NextState:    StateInIdentifier,
	},
	{
		CurrentState: exactState(StateInIdentifier),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Identifier,
	},

	{
		CurrentState: exactState(StateInPragma),
		CurrentRune:  runepredicate.Func(util.IsIdentContinue),
		NextState:    StateInPragma,
	},
	{
		CurrentState: exactState(StateInPragma),
		CurrentRune:  runepredicate.Exactly('/'),
		CurrentRaw:   util.RegexRunes,
		NextState:    StateInRegexSlash,
	},
	{
		CurrentState: exactState(StateInPragma),
		CurrentRune:  runepredicate.Exactly('!'),
		CurrentRaw:   util.RegexRunes,
		NextState:    StateInRegexExclaim,
	},
	{
		CurrentState: exactState(StateInPragma),
		CurrentRune:  runepredicate.Exactly('@'),
		CurrentRaw:   util.RegexRunes,
		NextState:    StateInRegexAt,
	},
	{
		CurrentState: exactState(StateInPragma),
		CurrentRune:  runepredicate.Exactly('{'),
		CurrentRaw:   util.RegexRunes,
		Delta:        1,
		NextState:    StateInRegexBrace,
	},
	{
		CurrentState: exactState(StateInPragma),
		CurrentRune:  runepredicate.Exactly('{'),
		CurrentRaw:   util.PEGRunes,
		Delta:        1,
		NextState:    StateInPEG,
	},
	{
		CurrentState: exactState(StateInPragma),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Pragma,
	},

	{
		CurrentState: exactState(StateInNumber),
		CurrentRune:  runepredicate.Func(util.IsNumberContinue),
		NextState:    StateInNumber,
	},
	{
		CurrentState: exactState(StateInNumber),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Number,
	},

	{
		CurrentState: exactState(StateInDoubleQuoteString),
		CurrentRune:  runepredicate.Exactly('"'),
		NextState:    StateDone,
		TokenType:    String,
	},
	{
		CurrentState: exactState(StateInDoubleQuoteString),
		CurrentRune:  runepredicate.Exactly('\\'),
		NextState:    StateInDoubleQuoteStringWithBackslash,
	},
	{
		CurrentState: someState(StateInDoubleQuoteString, StateInDoubleQuoteStringWithBackslash),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateInDoubleQuoteString,
	},

	{
		CurrentState: exactState(StateInSingleQuoteString),
		CurrentRune:  runepredicate.Exactly('\''),
		NextState:    StateDone,
		TokenType:    String,
	},
	{
		CurrentState: exactState(StateInSingleQuoteString),
		CurrentRune:  runepredicate.Exactly('\\'),
		NextState:    StateInSingleQuoteStringWithBackslash,
	},
	{
		CurrentState: someState(StateInSingleQuoteString, StateInSingleQuoteStringWithBackslash),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateInSingleQuoteString,
	},

	{
		CurrentState: exactState(StateInRegexSlash),
		CurrentRune:  runepredicate.Exactly('/'),
		NextState:    StateInRegexFlags,
	},
	{
		CurrentState: exactState(StateInRegexSlash),
		CurrentRune:  runepredicate.Exactly('\\'),
		NextState:    StateInRegexSlashWithBackslash,
	},
	{
		CurrentState: someState(StateInRegexSlash, StateInRegexSlashWithBackslash),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateInRegexSlash,
	},

	{
		CurrentState: exactState(StateInRegexExclaim),
		CurrentRune:  runepredicate.Exactly('!'),
		NextState:    StateInRegexFlags,
	},
	{
		CurrentState: exactState(StateInRegexExclaim),
		CurrentRune:  runepredicate.Exactly('\\'),
		NextState:    StateInRegexExclaimWithBackslash,
	},
	{
		CurrentState: someState(StateInRegexExclaim, StateInRegexExclaimWithBackslash),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateInRegexExclaim,
	},

	{
		CurrentState: exactState(StateInRegexAt),
		CurrentRune:  runepredicate.Exactly('@'),
		NextState:    StateInRegexFlags,
	},
	{
		CurrentState: exactState(StateInRegexAt),
		CurrentRune:  runepredicate.Exactly('\\'),
		NextState:    StateInRegexAtWithBackslash,
	},
	{
		CurrentState: someState(StateInRegexAt, StateInRegexAtWithBackslash),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateInRegexAt,
	},

	{
		CurrentState: exactState(StateInRegexBrace),
		CurrentRune:  runepredicate.Exactly('{'),
		Delta:        1,
		NextState:    StateInRegexBrace,
	},
	{
		CurrentState: exactState(StateInRegexBrace),
		CurrentRune:  runepredicate.Exactly('}'),
		Delta:        -1,
		OnlyIfZero:   true,
		NextState:    StateInRegexFlags,
	},
	{
		CurrentState: exactState(StateInRegexBrace),
		CurrentRune:  runepredicate.Exactly('}'),
		Delta:        -1,
		NextState:    StateInRegexBrace,
	},
	{
		CurrentState: exactState(StateInRegexBrace),
		CurrentRune:  runepredicate.Exactly('\\'),
		NextState:    StateInRegexBraceWithBackslash,
	},
	{
		CurrentState: someState(StateInRegexBrace, StateInRegexBraceWithBackslash),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateInRegexBrace,
	},

	{
		CurrentState: exactState(StateInRegexFlags),
		CurrentRune:  runepredicate.OneOf('i', 'm', 's', 'x'),
		NextState:    StateInRegexFlags,
	},
	{
		CurrentState: exactState(StateInRegexFlags),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateUnreadAndDone,
		TokenType:    Regex,
	},

	{
		CurrentState: exactState(StateInPEG),
		CurrentRune:  runepredicate.Exactly('{'),
		Delta:        1,
		NextState:    StateInPEG,
	},
	{
		CurrentState: exactState(StateInPEG),
		CurrentRune:  runepredicate.Exactly('}'),
		Delta:        -1,
		OnlyIfZero:   true,
		NextState:    StateDone,
		TokenType:    PEG,
	},
	{
		CurrentState: exactState(StateInPEG),
		CurrentRune:  runepredicate.Exactly('}'),
		Delta:        -1,
		NextState:    StateInPEG,
	},
	{
		CurrentState: exactState(StateInPEG),
		CurrentRune:  runepredicate.Exactly('\\'),
		NextState:    StateInPEGWithBackslash,
	},
	{
		CurrentState: someState(StateInPEG, StateInPEGWithBackslash),
		CurrentRune:  runepredicate.Any(),
		NextState:    StateInPEG,
	},
}

var earlyEOFTable = map[State]Type{
	StateInShebang:           ShebangLine,
	StateInSingleLineComment: SingleLineComment,
	StateInPragma:            Pragma,
	StateInIdentifier:        Identifier,
	StateInNumber:            Number,
}

var keywordTable = map[string]Type{
	"_":         KeywordPlaceholder,
	"alias":     KeywordAlias,
	"async":     KeywordAsync,
	"await":     KeywordAwait,
	"bitfield":  KeywordBitfield,
	"case":      KeywordCase,
	"const":     KeywordConst,
	"coroutine": KeywordCoroutine,
	"else":      KeywordElse,
	"for":       KeywordFor,
	"foreach":   KeywordForEach,
	"func":      KeywordFunc,
	"generator": KeywordGenerator,
	"goto":      KeywordGoto,
	"if":        KeywordIf,
	"import":    KeywordImport,
	"interface": KeywordInterface,
	"let":       KeywordLet,
	"lock":      KeywordLock,
	"method":    KeywordMethod,
	"null":      KeywordNull,
	"operator":  KeywordOperator,
	"property":  KeywordProperty,
	"return":    KeywordReturn,
	"static":    KeywordStatic,
	"struct":    KeywordStruct,
	"switch":    KeywordSwitch,
	"throw":     KeywordThrow,
	"type":      KeywordType,
	"union":     KeywordUnion,
	"var":       KeywordVar,
	"while":     KeywordWhile,
	"with":      KeywordWith,
	"yield":     KeywordYield,
}
