package ast

import (
	"fmt"
	"strings"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
	"github.com/chronos-tachyon/go-spiderscript/token"
	"github.com/chronos-tachyon/go-spiderscript/tokenpredicate"
	"github.com/chronos-tachyon/go-spiderscript/value"
)

func MatchValue(ptr *Node, p *Parser) bool {
	num := &LiteralNumber{}
	if num.Match(p) {
		*ptr = num
		return true
	}

	str := &LiteralString{}
	if str.Match(p) {
		*ptr = str
		return true
	}

	ident := &Identifier{}
	if ident.Match(p) {
		*ptr = ident
		return true
	}

	return false
}

func MatchPragmaExpr(ptr *Node, p *Parser) bool {
	num := &LiteralNumber{}
	if num.Match(p) {
		*ptr = num
		return true
	}

	str := &LiteralString{}
	if str.Match(p) {
		*ptr = str
		return true
	}

	ident := &Identifier{}
	if ident.Match(p) {
		*ptr = ident
		return true
	}

	return false
}

// LiteralNumber
// {{{

type LiteralNumber struct {
	Value value.Number
}

func (num *LiteralNumber) Init() {
	*num = LiteralNumber{}
}

func (num *LiteralNumber) Match(p *Parser) bool {
	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.Number)) {
		return false
	}

	switch x := tok.Parsed.(type) {
	case *value.Number:
		num.Value = *x

	case *value.Error:
		p.EmitError(x.Err)
		num.Value = *value.NewZero()

	default:
		panic(fmt.Errorf("expected *value.Number or *value.Error, got %T", tok.Parsed))
	}

	return true
}

func (num *LiteralNumber) String() string {
	return util.StringImpl(num)
}

func (num *LiteralNumber) GoString() string {
	return util.GoStringImpl(num)
}

func (num *LiteralNumber) EstimateStringLength() uint {
	return num.Value.EstimateStringLength()
}

func (num *LiteralNumber) EstimateGoStringLength() uint {
	return num.Value.EstimateGoStringLength()
}

func (num *LiteralNumber) WriteStringTo(out *strings.Builder) {
	num.Value.WriteStringTo(out)
}

func (num *LiteralNumber) WriteGoStringTo(out *strings.Builder) {
	num.Value.WriteGoStringTo(out)
}

func (num *LiteralNumber) ComputeStringLengthEstimates() {
}

var _ Node = (*LiteralNumber)(nil)

// }}}

// LiteralString
// {{{

type LiteralString struct {
	Value value.String
}

func (str *LiteralString) Init() {
	*str = LiteralString{}
}

func (str *LiteralString) Match(p *Parser) bool {
	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.String)) {
		return false
	}

	switch x := tok.Parsed.(type) {
	case *value.String:
		str.Value = *x

	case *value.Error:
		str.Value = *value.NewEmptyString()
		p.EmitError(x.Err)

	default:
		panic(fmt.Errorf("expected *value.String or *value.Error, got %T", tok.Parsed))
	}

	return true
}

func (str *LiteralString) String() string {
	return util.StringImpl(str)
}

func (str *LiteralString) GoString() string {
	return util.GoStringImpl(str)
}

func (str *LiteralString) EstimateStringLength() uint {
	return str.Value.EstimateStringLength()
}

func (str *LiteralString) EstimateGoStringLength() uint {
	return str.Value.EstimateGoStringLength()
}

func (str *LiteralString) WriteStringTo(out *strings.Builder) {
	str.Value.WriteStringTo(out)
}

func (str *LiteralString) WriteGoStringTo(out *strings.Builder) {
	str.Value.WriteGoStringTo(out)
}

func (str *LiteralString) ComputeStringLengthEstimates() {
	util.EstimateLengths(&str.Value)
}

var _ Node = (*LiteralString)(nil)

// }}}
