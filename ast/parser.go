package ast

import (
	"errors"

	"github.com/chronos-tachyon/go-spiderscript/token"
	"github.com/chronos-tachyon/go-spiderscript/tokenpredicate"
)

var ErrBadParse = errors.New("bad parse")

type Parser struct {
	errors     []error
	tokens     []token.Token
	tokenIndex uint
}

type Mark struct {
	savedErrorLength uint
	savedTokenIndex  uint
}

func NewParser(lexer *token.Lexer) *Parser {
	p := new(Parser)
	p.Init(lexer)
	return p
}

func (p *Parser) Init(lexer *token.Lexer) {
	*p = Parser{
		errors: make([]error, 0, 4),
		tokens: make([]token.Token, 0, 1024),
	}
	for lexer.HasNext() {
		p.tokens = append(p.tokens, lexer.Token())
	}
}

func (p *Parser) Errors() []error {
	if len(p.errors) == 0 {
		return nil
	}
	return p.errors
}

func (p *Parser) EmitError(err error) {
	p.errors = append(p.errors, err)
}

func (p *Parser) Mark() Mark {
	return Mark{
		savedErrorLength: uint(len(p.errors)),
		savedTokenIndex:  p.tokenIndex,
	}
}

func (p *Parser) Rewind(mark Mark) {
	p.errors = p.errors[:mark.savedErrorLength]
	p.tokenIndex = mark.savedTokenIndex
}

func (p *Parser) Forget(mark Mark) {
	_ = mark
}

func (p *Parser) Parse(node Node) bool {
	node.Init()
	if node.Match(p) {
		return true
	}
	node.Init()
	return false
}

func (p *Parser) Peek(out *token.Token, pred tokenpredicate.TokenPredicate) bool {
	if p.tokenIndex < uint(len(p.tokens)) {
		tok := p.tokens[p.tokenIndex]
		if pred.MatchToken(tok) {
			if out != nil {
				*out = tok
			}
			return true
		}
	}
	if out != nil {
		*out = token.Token{}
	}
	return false
}

func (p *Parser) Consume(out *token.Token, pred tokenpredicate.TokenPredicate) bool {
	ok := p.Peek(out, pred)
	if ok {
		p.tokenIndex++
	}
	return ok
}

func (p *Parser) ConsumeWhile(out *[]token.Token, pred tokenpredicate.TokenPredicate) bool {
	atLeastOne := false
	for p.tokenIndex < uint(len(p.tokens)) {
		tok := p.tokens[p.tokenIndex]
		if !pred.MatchToken(tok) {
			break
		}
		p.tokenIndex++
		if out != nil {
			*out = append(*out, tok)
		}
		atLeastOne = true
	}
	return atLeastOne
}

func (p *Parser) ConsumeUntil(out *[]token.Token, pred tokenpredicate.TokenPredicate) bool {
	atLeastOne := false
	for p.tokenIndex < uint(len(p.tokens)) {
		tok := p.tokens[p.tokenIndex]
		if pred.MatchToken(tok) {
			break
		}
		p.tokenIndex++
		if out != nil {
			*out = append(*out, tok)
		}
		atLeastOne = true
	}
	return atLeastOne
}

func (p *Parser) DropWhile(pred tokenpredicate.TokenPredicate) bool {
	return p.ConsumeWhile(nil, pred)
}

func (p *Parser) DropUntil(pred tokenpredicate.TokenPredicate) bool {
	return p.ConsumeUntil(nil, pred)
}
