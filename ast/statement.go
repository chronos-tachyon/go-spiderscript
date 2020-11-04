package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
	"github.com/chronos-tachyon/go-spiderscript/token"
	"github.com/chronos-tachyon/go-spiderscript/tokenpredicate"
	"github.com/chronos-tachyon/go-spiderscript/value"
)

func MatchStatement(ptr *Node, p *Parser) bool {
	p.DropWhile(tokenpredicate.Type(token.HWS))

	stmt0 := &ShebangLineStatement{}
	if stmt0.Match(p) {
		*ptr = stmt0
		return true
	}

	stmt1 := &EmptyStatement{}
	if stmt1.Match(p) {
		*ptr = stmt1
		return true
	}

	stmt2 := &PragmaStatement{}
	if stmt2.Match(p) {
		*ptr = stmt2
		return true
	}

	stmt3 := &ImportStatement{}
	if stmt3.Match(p) {
		*ptr = stmt3
		return true
	}

	stmt99 := &AnyStatement{}
	stmt99.Match(p)
	*ptr = stmt99
	return true
}

// ShebangLineStatement
// {{{

type ShebangLineStatement struct {
	Text string
	VWS  VerticalWhitespace
}

func (stmt *ShebangLineStatement) Init() {
	*stmt = ShebangLineStatement{}
}

func (stmt *ShebangLineStatement) Match(p *Parser) bool {
	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.ShebangLine)) {
		return false
	}

	text, err := value.ExtractLiteral(tok.Parsed, "#!<error>")
	if err != nil {
		p.EmitError(err)
	}

	if !strings.HasPrefix(text, "#!") {
		panic(fmt.Errorf("expected *value.Literal that starts with \"#!\", got %q", text))
	}

	text = text[2:]
	stmt.Text = text
	stmt.VWS.Match(p)
	return true
}

func (stmt *ShebangLineStatement) String() string {
	return util.StringImpl(stmt)
}

func (stmt *ShebangLineStatement) GoString() string {
	return util.GoStringImpl(stmt)
}

func (stmt *ShebangLineStatement) EstimateStringLength() uint {
	return 2 + uint(len(stmt.Text)) + stmt.VWS.EstimateStringLength()
}

func (stmt *ShebangLineStatement) EstimateGoStringLength() uint {
	return 25 + uint(len(stmt.Text)) + stmt.VWS.EstimateGoStringLength()
}

func (stmt *ShebangLineStatement) WriteStringTo(out *strings.Builder) {
	out.WriteString("#!")
	out.WriteString(stmt.Text)
	stmt.VWS.WriteStringTo(out)
}

func (stmt *ShebangLineStatement) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("ShebangLineStatement{")
	out.WriteString(strconv.Quote(stmt.Text))
	out.WriteByte(',')
	stmt.VWS.WriteGoStringTo(out)
	out.WriteByte('}')
}

func (stmt *ShebangLineStatement) ComputeStringLengthEstimates() {
	stmt.VWS.ComputeStringLengthEstimates()
}

var _ Node = (*ShebangLineStatement)(nil)

// }}}

// EmptyStatement
// {{{

type EmptyStatement struct {
	WS TerminalWhitespaceRun
}

func (stmt *EmptyStatement) Init() {
	*stmt = EmptyStatement{}
}

func (stmt *EmptyStatement) Match(p *Parser) bool {
	return stmt.WS.Match(p)
}

func (stmt *EmptyStatement) String() string {
	return util.StringImpl(stmt)
}

func (stmt *EmptyStatement) GoString() string {
	return util.GoStringImpl(stmt)
}

func (stmt *EmptyStatement) EstimateStringLength() uint {
	return stmt.WS.EstimateStringLength()
}

func (stmt *EmptyStatement) EstimateGoStringLength() uint {
	return 16 + stmt.WS.EstimateGoStringLength()
}

func (stmt *EmptyStatement) WriteStringTo(out *strings.Builder) {
	stmt.WS.WriteStringTo(out)
}

func (stmt *EmptyStatement) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("EmptyStatement{")
	stmt.WS.WriteGoStringTo(out)
	out.WriteByte('}')
}

func (stmt *EmptyStatement) ComputeStringLengthEstimates() {
	stmt.WS.ComputeStringLengthEstimates()
}

var _ Node = (*EmptyStatement)(nil)

// }}}

// PragmaStatement
// {{{

type PragmaStatement struct {
	Pragma Node
	WS0    InternalWhitespaceRun
	WS1    TerminalWhitespaceRun
}

func (stmt *PragmaStatement) Init() {
	*stmt = PragmaStatement{}
}

func (stmt *PragmaStatement) Match(p *Parser) bool {
	mark := p.Mark()
	defer p.Forget(mark)

	stmt.WS0.Match(p)

	if !MatchPragma(&stmt.Pragma, p) {
		p.Rewind(mark)
		return false
	}

	if !stmt.WS1.Match(p) {
		p.Rewind(mark)
		return false
	}

	return true
}

func (stmt *PragmaStatement) String() string {
	return util.StringImpl(stmt)
}

func (stmt *PragmaStatement) GoString() string {
	return util.GoStringImpl(stmt)
}

func (stmt *PragmaStatement) EstimateStringLength() uint {
	sum := 0 + stmt.Pragma.EstimateStringLength()
	sum += stmt.WS0.EstimateStringLength()
	sum += stmt.WS1.EstimateStringLength()
	return sum
}

func (stmt *PragmaStatement) EstimateGoStringLength() uint {
	sum := 19 + stmt.Pragma.EstimateGoStringLength()
	sum += stmt.WS0.EstimateGoStringLength()
	sum += stmt.WS1.EstimateGoStringLength()
	return sum
}

func (stmt *PragmaStatement) WriteStringTo(out *strings.Builder) {
	stmt.WS0.WriteStringTo(out)
	stmt.Pragma.WriteStringTo(out)
	stmt.WS1.WriteStringTo(out)
}

func (stmt *PragmaStatement) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("PragmaStatement{")
	stmt.Pragma.WriteGoStringTo(out)
	out.WriteByte(',')
	stmt.WS0.WriteGoStringTo(out)
	out.WriteByte(',')
	stmt.WS1.WriteGoStringTo(out)
	out.WriteByte('}')
}

func (stmt *PragmaStatement) ComputeStringLengthEstimates() {
	stmt.Pragma.ComputeStringLengthEstimates()
	stmt.WS0.ComputeStringLengthEstimates()
	stmt.WS1.ComputeStringLengthEstimates()
}

var _ Node = (*PragmaStatement)(nil)

// }}}

// ImportStatement
// {{{

type ImportStatement struct {
	ModuleName ModuleName
	Pragmas    []Node
	WS2        []Node
	WS0        InternalWhitespaceRun
	WS1        InternalWhitespaceRun
	WS3        TerminalWhitespaceRun
}

func (stmt *ImportStatement) Init() {
	*stmt = ImportStatement{}
}

func (stmt *ImportStatement) Match(p *Parser) bool {
	mark := p.Mark()
	defer p.Forget(mark)

	stmt.WS0.Match(p)

	if !p.Consume(nil, tokenpredicate.Type(token.KeywordImport)) {
		p.Rewind(mark)
		return false
	}

	if !stmt.WS1.Match(p) {
		p.Rewind(mark)
		return false
	}

	if !stmt.ModuleName.Match(p) {
		p.Rewind(mark)
		return false
	}

	for {
		if stmt.WS3.Match(p) {
			return true
		}

		var ws InternalWhitespaceRun
		ws.Match(p)

		var pragma Node
		if MatchPragma(&pragma, p) {
			stmt.Pragmas = append(stmt.Pragmas, pragma)
			stmt.WS2 = append(stmt.WS2, &ws)
			continue
		}

		p.Rewind(mark)
		return false
	}
}

func (stmt *ImportStatement) String() string {
	return util.StringImpl(stmt)
}

func (stmt *ImportStatement) GoString() string {
	return util.GoStringImpl(stmt)
}

func (stmt *ImportStatement) EstimateStringLength() uint {
	sum := 6 + stmt.ModuleName.EstimateStringLength()
	sum += sumStringLengthEstimates(stmt.Pragmas)
	sum += stmt.WS0.EstimateStringLength()
	sum += stmt.WS1.EstimateStringLength()
	sum += sumStringLengthEstimates(stmt.WS2)
	sum += stmt.WS3.EstimateStringLength()
	return sum
}

func (stmt *ImportStatement) EstimateGoStringLength() uint {
	sum := 26 + stmt.ModuleName.EstimateGoStringLength()
	sum += sumGoStringLengthEstimates(stmt.Pragmas)
	sum += stmt.WS0.EstimateGoStringLength()
	sum += stmt.WS1.EstimateGoStringLength()
	sum += sumGoStringLengthEstimates(stmt.WS2)
	sum += stmt.WS3.EstimateGoStringLength()
	return sum
}

func (stmt *ImportStatement) WriteStringTo(out *strings.Builder) {
	stmt.WS0.WriteStringTo(out)
	out.WriteString("import")
	stmt.WS1.WriteStringTo(out)
	stmt.ModuleName.WriteStringTo(out)
	for i := uint(0); i < uint(len(stmt.Pragmas)); i++ {
		stmt.WS2[i].WriteStringTo(out)
		stmt.Pragmas[i].WriteStringTo(out)
	}
	stmt.WS3.WriteStringTo(out)
}

func (stmt *ImportStatement) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("ImportStatement{")
	stmt.ModuleName.WriteGoStringTo(out)
	out.WriteByte(',')
	out.WriteByte('[')
	writeGoStringsTo(out, stmt.Pragmas)
	out.WriteByte(']')
	out.WriteByte(',')
	stmt.WS0.WriteGoStringTo(out)
	out.WriteByte(',')
	stmt.WS1.WriteGoStringTo(out)
	out.WriteByte(',')
	out.WriteByte('[')
	writeGoStringsTo(out, stmt.WS2)
	out.WriteByte(']')
	out.WriteByte(',')
	stmt.WS3.WriteGoStringTo(out)
	out.WriteByte('}')
}

func (stmt *ImportStatement) ComputeStringLengthEstimates() {
}

var _ Node = (*ImportStatement)(nil)

// }}}

// AnyStatement
// {{{

type AnyStatement struct {
	Token token.Token
}

func (stmt *AnyStatement) Init() {
	*stmt = AnyStatement{}
}

func (stmt *AnyStatement) Match(p *Parser) bool {
	p.Consume(&stmt.Token, tokenpredicate.Any())
	return true
}

func (stmt *AnyStatement) String() string {
	return util.StringImpl(stmt)
}

func (stmt *AnyStatement) GoString() string {
	return util.GoStringImpl(stmt)
}

func (stmt *AnyStatement) EstimateStringLength() uint {
	return 11 + stmt.Token.EstimateStringLength()
}

func (stmt *AnyStatement) EstimateGoStringLength() uint {
	return 10 + stmt.Token.EstimateGoStringLength()
}

func (stmt *AnyStatement) WriteStringTo(out *strings.Builder) {
	out.WriteString("AnyStatement[")
	stmt.Token.WriteStringTo(out)
	out.WriteString("]\n")
}

func (stmt *AnyStatement) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("AnyStatement{")
	stmt.Token.WriteGoStringTo(out)
	out.WriteByte('}')
}

func (stmt *AnyStatement) ComputeStringLengthEstimates() {
}

var _ Node = (*AnyStatement)(nil)

// }}}
