package ast

import (
	"strconv"
	"strings"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
	"github.com/chronos-tachyon/go-spiderscript/token"
	"github.com/chronos-tachyon/go-spiderscript/tokenpredicate"
	"github.com/chronos-tachyon/go-spiderscript/value"
)

func MatchPragma(ptr *Node, p *Parser) bool {
	pragma0 := &VersionPragma{}
	if pragma0.Match(p) {
		*ptr = pragma0
		return true
	}

	pragma1 := &GenericPragma{}
	pragma1.Init()
	if pragma1.Match(p) {
		*ptr = pragma1
		return true
	}
	return false
}

// VersionPragma
// {{{

type VersionPragma struct {
	Major uint32
	Minor uint32
	Patch uint32
}

func (pragma *VersionPragma) Init() {
	*pragma = VersionPragma{}
}

func (pragma *VersionPragma) Match(p *Parser) bool {
	mark := p.Mark()
	defer p.Forget(mark)

	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.Pragma)) {
		p.Rewind(mark)
		return false
	}

	name, err := value.ExtractLiteral(tok.Parsed, "#<error>")
	if err != nil {
		p.EmitError(err)
	}
	if name != "#version" {
		p.Rewind(mark)
		return false
	}

	p.DropWhile(tokenpredicate.Type(token.HWS))

	if !p.Consume(nil, tokenpredicate.Type(token.LParen)) {
		p.Rewind(mark)
		return false
	}

	p.DropWhile(tokenpredicate.Type(token.HWS))

	var major LiteralNumber
	var minor LiteralNumber
	var patch LiteralNumber

	if !major.Match(p) {
		p.Rewind(mark)
		return false
	}

	p.DropWhile(tokenpredicate.Type(token.HWS))

	if p.Consume(nil, tokenpredicate.Type(token.Comma)) {
		p.DropWhile(tokenpredicate.Type(token.HWS))
		if !minor.Match(p) {
			p.Rewind(mark)
			return false
		}
		p.DropWhile(tokenpredicate.Type(token.HWS))
	}

	if p.Consume(nil, tokenpredicate.Type(token.Comma)) {
		p.DropWhile(tokenpredicate.Type(token.HWS))
		if !patch.Match(p) {
			p.Rewind(mark)
			return false
		}
		p.DropWhile(tokenpredicate.Type(token.HWS))
	}

	if !p.Consume(nil, tokenpredicate.Type(token.RParen)) {
		p.Rewind(mark)
		return false
	}

	pragma.Major, err = major.Value.AsUint32()
	if err != nil {
		p.EmitError(err)
	}

	pragma.Minor, err = minor.Value.AsUint32()
	if err != nil {
		p.EmitError(err)
	}

	pragma.Patch, err = patch.Value.AsUint32()
	if err != nil {
		p.EmitError(err)
	}

	return true
}

func (pragma *VersionPragma) String() string {
	return util.StringImpl(pragma)
}

func (pragma *VersionPragma) GoString() string {
	return util.GoStringImpl(pragma)
}

func (pragma *VersionPragma) EstimateStringLength() uint {
	return 42
}

func (pragma *VersionPragma) EstimateGoStringLength() uint {
	return 47
}

func (pragma *VersionPragma) WriteStringTo(out *strings.Builder) {
	out.WriteString("#version(")
	out.WriteString(strconv.FormatUint(uint64(pragma.Major), 10))
	out.WriteByte(',')
	out.WriteString(strconv.FormatUint(uint64(pragma.Minor), 10))
	out.WriteByte(',')
	out.WriteString(strconv.FormatUint(uint64(pragma.Patch), 10))
	out.WriteByte(')')
}

func (pragma *VersionPragma) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("VersionPragma{")
	out.WriteString(strconv.FormatUint(uint64(pragma.Major), 10))
	out.WriteByte(',')
	out.WriteString(strconv.FormatUint(uint64(pragma.Minor), 10))
	out.WriteByte(',')
	out.WriteString(strconv.FormatUint(uint64(pragma.Patch), 10))
	out.WriteByte('}')
}

func (pragma *VersionPragma) ComputeStringLengthEstimates() {
}

var _ Node = (*VersionPragma)(nil)

// }}}

// GenericPragma
// {{{

type GenericPragma struct {
	Name        string
	Expressions []Node
}

func (pragma *GenericPragma) Init() {
	*pragma = GenericPragma{
		Name:        "#<error>",
		Expressions: make([]Node, 0, 4),
	}
}

func (pragma *GenericPragma) Match(p *Parser) bool {
	mark := p.Mark()
	defer p.Forget(mark)

	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.Pragma)) {
		p.Rewind(mark)
		return false
	}

	name, err := value.ExtractLiteral(tok.Parsed, "#<error>")
	if err != nil {
		p.EmitError(err)
	}
	pragma.Name = name

	p.Forget(mark)
	mark = p.Mark()

	p.DropWhile(tokenpredicate.Type(token.HWS))

	if !p.Consume(nil, tokenpredicate.Type(token.LParen)) {
		pragma.Expressions = nil
		p.Rewind(mark)
		return true
	}

	p.DropWhile(tokenpredicate.Type(token.HWS))

	if p.Consume(nil, tokenpredicate.Type(token.RParen)) {
		pragma.Expressions = nil
		return true
	}

	var expr Node
	if !MatchPragmaExpr(&expr, p) {
		pragma.Expressions = nil
		p.Rewind(mark)
		return true
	}
	pragma.Expressions = append(pragma.Expressions, expr)

	for {
		p.DropWhile(tokenpredicate.Type(token.HWS))

		if p.Consume(nil, tokenpredicate.Type(token.RParen)) {
			return true
		}

		if p.Consume(nil, tokenpredicate.Type(token.Comma)) {
			expr = nil
			p.DropWhile(tokenpredicate.Type(token.HWS))
			if MatchPragmaExpr(&expr, p) {
				pragma.Expressions = append(pragma.Expressions, expr)
				continue
			}
		}

		break
	}

	p.Rewind(mark)
	pragma.Expressions = nil
	return true
}

func (pragma *GenericPragma) String() string {
	return util.StringImpl(pragma)
}

func (pragma *GenericPragma) GoString() string {
	return util.GoStringImpl(pragma)
}

func (pragma *GenericPragma) EstimateStringLength() uint {
	sum := 1 + uint(len(pragma.Name)) + uint(len(pragma.Expressions))
	for _, expr := range pragma.Expressions {
		sum += expr.EstimateStringLength()
	}
	return sum
}

func (pragma *GenericPragma) EstimateGoStringLength() uint {
	sum := 17 + uint(len(pragma.Name)) + uint(len(pragma.Expressions))
	for _, expr := range pragma.Expressions {
		sum += expr.EstimateGoStringLength()
	}
	return sum
}

func (pragma *GenericPragma) WriteStringTo(out *strings.Builder) {
	out.WriteString(pragma.Name)
	out.WriteByte('(')
	for index, expr := range pragma.Expressions {
		if index > 0 {
			out.WriteByte(',')
		}
		expr.WriteStringTo(out)
	}
	out.WriteByte(')')
}

func (pragma *GenericPragma) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("GenericPragma{")
	out.WriteString(strconv.Quote(pragma.Name))
	for _, expr := range pragma.Expressions {
		out.WriteByte(',')
		expr.WriteGoStringTo(out)
	}
	out.WriteByte('}')
}

func (pragma *GenericPragma) ComputeStringLengthEstimates() {
	for _, expr := range pragma.Expressions {
		expr.ComputeStringLengthEstimates()
	}
}

var _ Node = (*GenericPragma)(nil)

// }}}
