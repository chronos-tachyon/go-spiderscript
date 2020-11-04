package ast

import (
	"fmt"

	"github.com/chronos-tachyon/go-spiderscript/operators"
	"github.com/chronos-tachyon/go-spiderscript/token"
)

var _ = operators.InvalidOperator

// UnaryOperator
// {{{

type UnaryOperator uint8

const (
	InvalidUnaryOperator UnaryOperator = iota
	OpAddressOf
	OpDerefPointer
	OpPos
	OpNeg
	OpBitwiseNOT
	OpLogicalNOT
	OpSplat
)

var unaryOperatorSymbols = []string{
	"???",
	"&",
	"*",
	"+",
	"-",
	"~",
	"!",
	"...",
}

var unaryOperatorNames = []string{
	"InvalidUnaryOperator",
	"OpAddressOf",
	"OpDerefPointer",
	"OpPos",
	"OpNeg",
	"OpBitwiseNOT",
	"OpLogicalNOT",
	"OpSplat",
}

var unaryOperatorIsPrefixMap = map[UnaryOperator]bool{
	OpSplat: true,
}

func (op UnaryOperator) String() string {
	if uint(op) >= uint(len(unaryOperatorSymbols)) {
		op = 0
	}
	return unaryOperatorSymbols[op]
}

func (op UnaryOperator) GoString() string {
	if uint(op) >= uint(len(unaryOperatorNames)) {
		return fmt.Sprintf("UnaryOperator(%d)", uint(op))
	}
	return unaryOperatorNames[op]
}

func (op UnaryOperator) IsPrefix() bool {
	return unaryOperatorIsPrefixMap[op]
}

var _ fmt.Stringer = UnaryOperator(0)
var _ fmt.GoStringer = UnaryOperator(0)

// }}}

// BinaryOperator
// {{{

type BinaryOperator uint8

const (
	InvalidBinaryOperator BinaryOperator = iota

	OpScope
	OpDeref

	OpPow

	OpMul
	OpDiv
	OpMod
	OpDivMod

	OpAdd
	OpSub

	OpBitwiseLShift
	OpBitwiseRShift
	OpBitwiseLRotate
	OpBitwiseRRotate

	OpBitwiseAND
	OpBitwiseXOR
	OpBitwiseOR

	OpCmpCMP
	OpCmpEQ
	OpCmpNE
	OpCmpLT
	OpCmpLE
	OpCmpGT
	OpCmpGE

	OpLogicalAND
	OpLogicalXOR
	OpLogicalOR

	OpRange

	OpElvis
)

var binaryOperatorSymbols = []string{
	"???",
	"::",
	".",
	"**",
	"*",
	"/",
	"%",
	"/%",
	"+",
	"-",
	"<<",
	">>",
	"<<|",
	">>|",
	"&",
	"^",
	"|",
	"<=>",
	"==",
	"!=",
	"<",
	"<=",
	">",
	">=",
	"&&",
	"^^",
	"||",
	"..",
	"?:",
}

var binaryOperatorNames = []string{
	"InvalidBinaryOperator",
	"OpScope",
	"OpDeref",
	"OpPow",
	"OpMul",
	"OpDiv",
	"OpMod",
	"OpDivMod",
	"OpAdd",
	"OpSub",
	"OpBitwiseLShift",
	"OpBitwiseRShift",
	"OpBitwiseLRotate",
	"OpBitwiseRRotate",
	"OpBitwiseAND",
	"OpBitwiseXOR",
	"OpBitwiseOR",
	"OpCmpCMP",
	"OpCmpEQ",
	"OpCmpNE",
	"OpCmpLT",
	"OpCmpLE",
	"OpCmpGT",
	"OpCmpGE",
	"OpLogicalAND",
	"OpLogicalXOR",
	"OpLogicalOR",
	"OpRange",
	"OpElvis",
}

func (op BinaryOperator) String() string {
	if uint(op) >= uint(len(binaryOperatorSymbols)) {
		op = 0
	}
	return binaryOperatorSymbols[op]
}

func (op BinaryOperator) GoString() string {
	if uint(op) >= uint(len(binaryOperatorNames)) {
		return fmt.Sprintf("BinaryOperator(%d)", uint(op))
	}
	return binaryOperatorNames[op]
}

var _ fmt.Stringer = BinaryOperator(0)
var _ fmt.GoStringer = BinaryOperator(0)

// }}}

// TernaryOperator
// {{{

type TernaryOperator uint8

const (
	InvalidTernaryOperator TernaryOperator = iota
	OpQuestionColon
	OpIfElse
)

var ternaryOperatorFirstSymbols = []string{
	"???",
	"?",
	"if",
}

var ternaryOperatorSecondSymbols = []string{
	"???",
	":",
	"else",
}

var ternaryOperatorNames = []string{
	"InvalidTernaryOperator",
	"OpQuestionColon",
	"OpIfElse",
}

var ternaryOperatorIsInvertedMap = map[TernaryOperator]bool{
	OpIfElse: true,
}

func (op TernaryOperator) FirstSymbol() string {
	if uint(op) >= uint(len(ternaryOperatorFirstSymbols)) {
		op = 0
	}
	return ternaryOperatorFirstSymbols[op]
}

func (op TernaryOperator) SecondSymbol() string {
	if uint(op) >= uint(len(ternaryOperatorSecondSymbols)) {
		op = 0
	}
	return ternaryOperatorSecondSymbols[op]
}

func (op TernaryOperator) String() string {
	return op.GoString()
}

func (op TernaryOperator) GoString() string {
	if uint(op) >= uint(len(ternaryOperatorNames)) {
		return fmt.Sprintf("TernaryOperator(%d)", uint(op))
	}
	return ternaryOperatorNames[op]
}

func (op TernaryOperator) IsInverted() bool {
	return ternaryOperatorIsInvertedMap[op]
}

var _ fmt.Stringer = TernaryOperator(0)
var _ fmt.GoStringer = TernaryOperator(0)

// }}}

var unaryOperatorMap = map[token.Type]UnaryOperator{
	token.Bang:      OpLogicalNOT,
	token.DotDotDot: OpSplat,
	token.Minus:     OpNeg,
	token.Plus:      OpPos,
	token.Tilde:     OpBitwiseNOT,
	token.Ampersand: OpAddressOf,
	token.Star:      OpDerefPointer,
}

var binaryOperatorMap = map[token.Type]BinaryOperator{
	token.Ampersand:          OpBitwiseAND,
	token.AmpersandAmpersand: OpLogicalAND,
	token.BangEqual:          OpCmpNE,
	token.Bar:                OpBitwiseOR,
	token.BarBar:             OpLogicalOR,
	token.Caret:              OpBitwiseXOR,
	token.CaretCaret:         OpLogicalXOR,
	token.ColonColon:         OpScope,
	token.Dot:                OpDeref,
	token.DotDot:             OpRange,
	token.EqualEqual:         OpCmpEQ,
	token.Greater:            OpCmpGT,
	token.GreaterEqual:       OpCmpGE,
	token.GreaterGreater:     OpBitwiseRShift,
	token.GreaterGreaterBar:  OpBitwiseRRotate,
	token.Less:               OpCmpLT,
	token.LessEqual:          OpCmpLE,
	token.LessEqualGreater:   OpCmpCMP,
	token.LessLess:           OpBitwiseLShift,
	token.LessLessBar:        OpBitwiseLRotate,
	token.Minus:              OpSub,
	token.Percent:            OpMod,
	token.Plus:               OpAdd,
	token.QuestionColon:      OpElvis,
	token.Slash:              OpDiv,
	token.SlashPercent:       OpDivMod,
	token.Star:               OpMul,
	token.StarStar:           OpPow,
}

var ternaryOperatorMap = map[token.Type]TernaryOperator{
	token.Question:  OpQuestionColon,
	token.KeywordIf: OpIfElse,
}

var ternaryOperatorSecondTokenMap = map[TernaryOperator]token.Type{
	OpQuestionColon: token.Colon,
	OpIfElse:        token.KeywordElse,
}
