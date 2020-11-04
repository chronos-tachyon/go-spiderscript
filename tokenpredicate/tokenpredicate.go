package tokenpredicate

import (
	"github.com/chronos-tachyon/go-spiderscript/token"
	"github.com/chronos-tachyon/go-spiderscript/value"
)

type TokenPredicate interface {
	MatchToken(token.Token) bool
}

// None
// {{{

type nonePredicate struct{}

func (pred nonePredicate) MatchToken(tok token.Token) bool {
	return false
}

func None() TokenPredicate {
	return nonePredicate{}
}

// }}}

// Any
// {{{

type anyPredicate struct{}

func (pred anyPredicate) MatchToken(tok token.Token) bool {
	return true
}

func Any() TokenPredicate {
	return anyPredicate{}
}

// }}}

// Not
// {{{

type notPredicate struct {
	Child TokenPredicate
}

func (pred notPredicate) MatchToken(tok token.Token) bool {
	return !pred.Child.MatchToken(tok)
}

func Not(child TokenPredicate) TokenPredicate {
	return notPredicate{child}
}

// }}}

// And
// {{{

type andPredicate struct {
	Children []TokenPredicate
}

func (pred andPredicate) MatchToken(tok token.Token) bool {
	for _, child := range pred.Children {
		if !child.MatchToken(tok) {
			return false
		}
	}
	return true
}

func And(children ...TokenPredicate) TokenPredicate {
	return andPredicate{children}
}

// }}}

// Or
// {{{

type orPredicate struct {
	Children []TokenPredicate
}

func (pred orPredicate) MatchToken(tok token.Token) bool {
	for _, child := range pred.Children {
		if child.MatchToken(tok) {
			return true
		}
	}
	return false
}

func Or(children ...TokenPredicate) TokenPredicate {
	return orPredicate{children}
}

// }}}

// Type
// {{{

type typePredicate struct {
	Type token.Type
}

func (pred typePredicate) MatchToken(tok token.Token) bool {
	return tok.Type == pred.Type
}

func Type(tt token.Type) TokenPredicate {
	return typePredicate{tt}
}

// }}}

// Literal
// {{{

type literalPredicate struct {
	Text string
}

func (pred literalPredicate) MatchToken(tok token.Token) bool {
	if lv, ok := tok.Parsed.(*value.Literal); ok {
		if lv.Text == pred.Text {
			return true
		}
	}
	return false
}

func Literal(text string) TokenPredicate {
	return literalPredicate{text}
}

// }}}

// Func
// {{{

type funcPredicate struct {
	Func func(token.Token) bool
}

func (pred funcPredicate) MatchToken(tok token.Token) bool {
	return pred.Func(tok)
}

func Func(fn func(token.Token) bool) TokenPredicate {
	return funcPredicate{fn}
}

// }}}
