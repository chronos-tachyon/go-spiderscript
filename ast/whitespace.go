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

// HorizontalWhitespace
// {{{

type HorizontalWhitespace struct {
	Count uint
}

func (ws *HorizontalWhitespace) Init() {
	*ws = HorizontalWhitespace{}
}

func (ws *HorizontalWhitespace) Match(p *Parser) bool {
	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.HWS)) {
		return false
	}

	startColumn := tok.Start.RawColumnNumber
	endColumn := tok.End.RawColumnNumber
	for p.Consume(&tok, tokenpredicate.Type(token.HWS)) {
		endColumn = tok.End.RawColumnNumber
	}

	if startColumn > endColumn {
		startColumn = endColumn
	}
	ws.Count = endColumn - startColumn
	return true
}

func (ws *HorizontalWhitespace) String() string {
	return util.StringImpl(ws)
}

func (ws *HorizontalWhitespace) GoString() string {
	return util.GoStringImpl(ws)
}

func (ws *HorizontalWhitespace) EstimateStringLength() uint {
	return uint(ws.Count)
}

func (ws *HorizontalWhitespace) EstimateGoStringLength() uint {
	return 26
}

func (ws *HorizontalWhitespace) WriteStringTo(out *strings.Builder) {
	for i := uint(0); i < ws.Count; i++ {
		out.WriteByte(' ')
	}
}

func (ws *HorizontalWhitespace) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("HorizontalWhitespace{")
	out.WriteString(strconv.FormatUint(uint64(ws.Count), 10))
	out.WriteByte('}')
}

func (ws *HorizontalWhitespace) ComputeStringLengthEstimates() {
}

var _ Node = (*HorizontalWhitespace)(nil)

// }}}

// VerticalWhitespace
// {{{

type VerticalWhitespace struct {
	Count uint
}

func (ws *VerticalWhitespace) Init() {
	*ws = VerticalWhitespace{}
}

func (ws *VerticalWhitespace) Match(p *Parser) bool {
	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.VWS)) {
		return false
	}

	startLine := tok.Start.RawLineNumber
	endLine := tok.End.RawLineNumber
	for p.Consume(&tok, tokenpredicate.Type(token.VWS)) {
		endLine = tok.End.RawLineNumber
	}

	if startLine > endLine {
		startLine = endLine
	}
	ws.Count = endLine - startLine
	return true
}

func (ws *VerticalWhitespace) String() string {
	return util.StringImpl(ws)
}

func (ws *VerticalWhitespace) GoString() string {
	return util.GoStringImpl(ws)
}

func (ws *VerticalWhitespace) EstimateStringLength() uint {
	return uint(ws.Count)
}

func (ws *VerticalWhitespace) EstimateGoStringLength() uint {
	return 24
}

func (ws *VerticalWhitespace) WriteStringTo(out *strings.Builder) {
	for i := uint(0); i < ws.Count; i++ {
		out.WriteByte('\n')
	}
}

func (ws *VerticalWhitespace) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("VerticalWhitespace{")
	out.WriteString(strconv.FormatUint(uint64(ws.Count), 10))
	out.WriteByte('}')
}

func (ws *VerticalWhitespace) ComputeStringLengthEstimates() {
}

var _ Node = (*VerticalWhitespace)(nil)

// }}}

// SingleLineCommentWhitespace
// {{{

type SingleLineCommentWhitespace struct {
	Text string
}

func (ws *SingleLineCommentWhitespace) Init() {
	*ws = SingleLineCommentWhitespace{}
}

func (ws *SingleLineCommentWhitespace) Match(p *Parser) bool {
	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.SingleLineComment)) {
		return false
	}
	p.Consume(nil, tokenpredicate.Type(token.VWS))

	text, err := value.ExtractLiteral(tok.Parsed, "//<error>")
	if err != nil {
		p.EmitError(err)
	}

	if !strings.HasPrefix(text, "//") {
		panic(fmt.Errorf("expected *value.Literal that starts with \"//\", got %q", text))
	}

	text = text[2:]
	ws.Text = text
	return true
}

func (ws *SingleLineCommentWhitespace) String() string {
	return util.StringImpl(ws)
}

func (ws *SingleLineCommentWhitespace) GoString() string {
	return util.GoStringImpl(ws)
}

func (ws *SingleLineCommentWhitespace) EstimateStringLength() uint {
	return 3 + uint(len(ws.Text))
}

func (ws *SingleLineCommentWhitespace) EstimateGoStringLength() uint {
	return 31 + uint(len(ws.Text))
}

func (ws *SingleLineCommentWhitespace) WriteStringTo(out *strings.Builder) {
	out.WriteString("//")
	out.WriteString(ws.Text)
	out.WriteByte('\n')
}

func (ws *SingleLineCommentWhitespace) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("SingleLineCommentWhitespace{")
	out.WriteString(strconv.Quote(ws.Text))
	out.WriteByte('}')
}

func (ws *SingleLineCommentWhitespace) ComputeStringLengthEstimates() {
}

var _ Node = (*SingleLineCommentWhitespace)(nil)

// }}}

// MultiLineCommentWhitespace
// {{{

type MultiLineCommentWhitespace struct {
	Text string
}

func (ws *MultiLineCommentWhitespace) Init() {
	*ws = MultiLineCommentWhitespace{}
}

func (ws *MultiLineCommentWhitespace) Match(p *Parser) bool {
	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.MultiLineComment)) {
		return false
	}
	p.Consume(nil, tokenpredicate.Type(token.VWS))

	text, err := value.ExtractLiteral(tok.Parsed, "/*<error>*/")
	if err != nil {
		p.EmitError(err)
	}

	if !strings.HasPrefix(text, "/*") {
		panic(fmt.Errorf("expected *value.Literal that starts with \"/*\", got %q", text))
	}

	if !strings.HasSuffix(text, "*/") {
		panic(fmt.Errorf("expected *value.Literal that starts with \"*/\", got %q", text))
	}

	text = text[2:]
	text = text[:len(text)-2]
	ws.Text = text
	return true
}

func (ws *MultiLineCommentWhitespace) String() string {
	return util.StringImpl(ws)
}

func (ws *MultiLineCommentWhitespace) GoString() string {
	return util.GoStringImpl(ws)
}

func (ws *MultiLineCommentWhitespace) EstimateStringLength() uint {
	return 4 + uint(len(ws.Text))
}

func (ws *MultiLineCommentWhitespace) EstimateGoStringLength() uint {
	return 30 + uint(len(ws.Text))
}

func (ws *MultiLineCommentWhitespace) WriteStringTo(out *strings.Builder) {
	out.WriteString("/*")
	out.WriteString(ws.Text)
	out.WriteString("*/")
}

func (ws *MultiLineCommentWhitespace) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("MultiLineCommentWhitespace{")
	out.WriteString(strconv.Quote(ws.Text))
	out.WriteByte('}')
}

func (ws *MultiLineCommentWhitespace) ComputeStringLengthEstimates() {
}

var _ Node = (*MultiLineCommentWhitespace)(nil)

// }}}

// StatementTerminatorWhitespace
// {{{

type StatementTerminatorWhitespace struct {
}

func (ws *StatementTerminatorWhitespace) Init() {
	*ws = StatementTerminatorWhitespace{}
}

func (ws *StatementTerminatorWhitespace) Match(p *Parser) bool {
	return p.Consume(nil, tokenpredicate.Type(token.Semicolon))
}

func (ws *StatementTerminatorWhitespace) String() string {
	return util.StringImpl(ws)
}

func (ws *StatementTerminatorWhitespace) GoString() string {
	return util.GoStringImpl(ws)
}

func (ws *StatementTerminatorWhitespace) EstimateStringLength() uint {
	return 1
}

func (ws *StatementTerminatorWhitespace) EstimateGoStringLength() uint {
	return 31
}

func (ws *StatementTerminatorWhitespace) WriteStringTo(out *strings.Builder) {
	out.WriteByte(';')
}

func (ws *StatementTerminatorWhitespace) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("StatementTerminatorWhitespace{}")
}

func (ws *StatementTerminatorWhitespace) ComputeStringLengthEstimates() {
}

var _ Node = (*StatementTerminatorWhitespace)(nil)

// }}}

// InternalWhitespaceRun
// {{{

type InternalWhitespaceRun struct {
	Children []Node

	PrecomputedStringLength   uint
	PrecomputedGoStringLength uint
}

func (ws *InternalWhitespaceRun) Init() {
	*ws = InternalWhitespaceRun{}
}

func (ws *InternalWhitespaceRun) Match(p *Parser) bool {
	if ws.Children == nil {
		ws.Children = make([]Node, 0, 16)
	}

	atLeastOne := false
	for {
		var ws0 HorizontalWhitespace
		if ws0.Match(p) {
			atLeastOne = true
			ws.Children = append(ws.Children, &ws0)
			continue
		}

		var ws1 VerticalWhitespace
		if ws1.Match(p) {
			atLeastOne = true
			ws.Children = append(ws.Children, &ws1)
			continue
		}

		var ws2 SingleLineCommentWhitespace
		if ws2.Match(p) {
			atLeastOne = true
			ws.Children = append(ws.Children, &ws2)
			continue
		}

		var ws3 MultiLineCommentWhitespace
		if ws3.Match(p) {
			atLeastOne = true
			ws.Children = append(ws.Children, &ws3)
			continue
		}

		break
	}
	return atLeastOne
}

func (ws *InternalWhitespaceRun) String() string {
	return util.StringImpl(ws)
}

func (ws *InternalWhitespaceRun) GoString() string {
	return util.GoStringImpl(ws)
}

func (ws *InternalWhitespaceRun) EstimateStringLength() uint {
	return ws.PrecomputedStringLength
}

func (ws *InternalWhitespaceRun) EstimateGoStringLength() uint {
	return ws.PrecomputedGoStringLength
}

func (ws *InternalWhitespaceRun) WriteStringTo(out *strings.Builder) {
	for _, item := range ws.Children {
		item.WriteStringTo(out)
	}
}

func (ws *InternalWhitespaceRun) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("InternalWhitespaceRun{")
	for index, item := range ws.Children {
		if index != 0 {
			out.WriteByte(',')
		}
		item.WriteGoStringTo(out)
	}
	out.WriteByte('}')
}

func (ws *InternalWhitespaceRun) SetEstimatedStringLength(length uint) {
	ws.PrecomputedStringLength = length
}

func (ws *InternalWhitespaceRun) SetEstimatedGoStringLength(length uint) {
	ws.PrecomputedGoStringLength = length
}

func (ws *InternalWhitespaceRun) ComputeStringLengthEstimates() {
	util.EstimateLengths(ws)
}

var _ Node = (*InternalWhitespaceRun)(nil)

// }}}

// TerminalWhitespaceRun
// {{{

type TerminalWhitespaceRun struct {
	Children []Node

	PrecomputedStringLength   uint
	PrecomputedGoStringLength uint
}

func (ws *TerminalWhitespaceRun) Init() {
	*ws = TerminalWhitespaceRun{}
}

func (ws *TerminalWhitespaceRun) Match(p *Parser) bool {
	if ws.Children == nil {
		ws.Children = make([]Node, 0, 16)
	}

	mark := p.Mark()
	defer p.Forget(mark)

	for {
		if p.Peek(nil, tokenpredicate.Type(token.EOF)) {
			return true
		}

		var ws0 StatementTerminatorWhitespace
		if ws0.Match(p) {
			ws.Children = append(ws.Children, &ws0)
			return true
		}

		var ws1 VerticalWhitespace
		if ws1.Match(p) {
			ws.Children = append(ws.Children, &ws1)
			return true
		}

		var ws2 SingleLineCommentWhitespace
		if ws2.Match(p) {
			ws.Children = append(ws.Children, &ws2)
			return true
		}

		var ws3 HorizontalWhitespace
		if ws3.Match(p) {
			ws.Children = append(ws.Children, &ws3)
			continue
		}

		var ws4 MultiLineCommentWhitespace
		if ws4.Match(p) {
			ws.Children = append(ws.Children, &ws4)
			continue
		}

		break
	}

	p.Rewind(mark)
	return false
}

func (ws *TerminalWhitespaceRun) String() string {
	return util.StringImpl(ws)
}

func (ws *TerminalWhitespaceRun) GoString() string {
	return util.GoStringImpl(ws)
}

func (ws *TerminalWhitespaceRun) EstimateStringLength() uint {
	return ws.PrecomputedStringLength
}

func (ws *TerminalWhitespaceRun) EstimateGoStringLength() uint {
	return ws.PrecomputedGoStringLength
}

func (ws *TerminalWhitespaceRun) WriteStringTo(out *strings.Builder) {
	for _, item := range ws.Children {
		item.WriteStringTo(out)
	}
}

func (ws *TerminalWhitespaceRun) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("TerminalWhitespaceRun{")
	for index, item := range ws.Children {
		if index != 0 {
			out.WriteByte(',')
		}
		item.WriteGoStringTo(out)
	}
	out.WriteByte('}')
}

func (ws *TerminalWhitespaceRun) SetEstimatedStringLength(length uint) {
	ws.PrecomputedStringLength = length
}

func (ws *TerminalWhitespaceRun) SetEstimatedGoStringLength(length uint) {
	ws.PrecomputedGoStringLength = length
}

func (ws *TerminalWhitespaceRun) ComputeStringLengthEstimates() {
	util.EstimateLengths(ws)
}

var _ Node = (*TerminalWhitespaceRun)(nil)

// }}}
