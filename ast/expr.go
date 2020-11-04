package ast

import (
	"strings"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
	"github.com/chronos-tachyon/go-spiderscript/token"
	"github.com/chronos-tachyon/go-spiderscript/tokenpredicate"
)

func MatchExpr(ptr *Node, p *Parser) bool {
	var dummy Node
	if ptr == nil {
		ptr = &dummy
	}
	*ptr = nil
	return matchExpr0(ptr, p)
}

func matchExpr0(ptr *Node, p *Parser) bool {
	return matchExpr1(ptr, p)
}

func matchExpr1(ptr *Node, p *Parser) bool {
	return matchExpr2(ptr, p)
}

func matchExpr2(ptr *Node, p *Parser) bool {
	return matchExpr3(ptr, p)
}

func matchExpr3(ptr *Node, p *Parser) bool {
	return matchExpr4(ptr, p)
}

func matchExpr4(ptr *Node, p *Parser) bool {
	return matchExpr5(ptr, p)
}

func matchExpr5(ptr *Node, p *Parser) bool {
	return matchExpr6(ptr, p)
}

func matchExpr6(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr7(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	for {
		var ws0 InternalWhitespaceRun
		ws0.Match(p)

		var tok token.Token
		if !p.Consume(&tok, tokenpredicate.Type(token.QuestionColon)) {
			break
		}

		op := binaryOperatorMap[tok.Type]

		var ws1 InternalWhitespaceRun
		ws1.Match(p)

		var right Node
		if !matchExpr7(&right, p) {
			break
		}

		p.Forget(mark)
		mark = p.Mark()
		left = &BinaryOperatorExpr{
			Operator: op,
			Operand0: left,
			Operand1: right,
			WS0:      ws0,
			WS1:      ws1,
		}
	}

	p.Rewind(mark)
	*ptr = left
	return true
}

func matchExpr7(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr8(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	for {
		var ws0 InternalWhitespaceRun
		ws0.Match(p)

		var tok token.Token
		if !p.Consume(&tok, tokenpredicate.Type(token.DotDot)) {
			break
		}

		op := binaryOperatorMap[tok.Type]

		var ws1 InternalWhitespaceRun
		ws1.Match(p)

		var right Node
		if !matchExpr8(&right, p) {
			break
		}

		p.Forget(mark)
		mark = p.Mark()
		left = &BinaryOperatorExpr{
			Operator: op,
			Operand0: left,
			Operand1: right,
			WS0:      ws0,
			WS1:      ws1,
		}
	}

	p.Rewind(mark)
	*ptr = left
	return true
}

func matchExpr8(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr9(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	for {
		var ws0 InternalWhitespaceRun
		ws0.Match(p)

		var tok token.Token
		if !p.Consume(&tok, tokenpredicate.Type(token.BarBar)) {
			break
		}

		op := binaryOperatorMap[tok.Type]

		var ws1 InternalWhitespaceRun
		ws1.Match(p)

		var right Node
		if !matchExpr9(&right, p) {
			break
		}

		p.Forget(mark)
		mark = p.Mark()
		left = &BinaryOperatorExpr{
			Operator: op,
			Operand0: left,
			Operand1: right,
			WS0:      ws0,
			WS1:      ws1,
		}
	}

	p.Rewind(mark)
	*ptr = left
	return true
}

func matchExpr9(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr10(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	for {
		var ws0 InternalWhitespaceRun
		ws0.Match(p)

		var tok token.Token
		if !p.Consume(&tok, tokenpredicate.Type(token.CaretCaret)) {
			break
		}

		op := binaryOperatorMap[tok.Type]

		var ws1 InternalWhitespaceRun
		ws1.Match(p)

		var right Node
		if !matchExpr10(&right, p) {
			break
		}

		p.Forget(mark)
		mark = p.Mark()
		left = &BinaryOperatorExpr{
			Operator: op,
			Operand0: left,
			Operand1: right,
			WS0:      ws0,
			WS1:      ws1,
		}
	}

	p.Rewind(mark)
	*ptr = left
	return true
}

func matchExpr10(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr11(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	for {
		var ws0 InternalWhitespaceRun
		ws0.Match(p)

		var tok token.Token
		if !p.Consume(&tok, tokenpredicate.Type(token.AmpersandAmpersand)) {
			break
		}

		op := binaryOperatorMap[tok.Type]

		var ws1 InternalWhitespaceRun
		ws1.Match(p)

		var right Node
		if !matchExpr11(&right, p) {
			break
		}

		p.Forget(mark)
		mark = p.Mark()
		left = &BinaryOperatorExpr{
			Operator: op,
			Operand0: left,
			Operand1: right,
			WS0:      ws0,
			WS1:      ws1,
		}
	}

	p.Rewind(mark)
	*ptr = left
	return true
}

func matchExpr11(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr12(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	for {
		var ws0 InternalWhitespaceRun
		ws0.Match(p)

		var tok token.Token
		if !p.Consume(&tok,
			tokenpredicate.Or(
				tokenpredicate.Type(token.BangEqual),
				tokenpredicate.Type(token.EqualEqual),
				tokenpredicate.Type(token.Greater),
				tokenpredicate.Type(token.GreaterEqual),
				tokenpredicate.Type(token.Less),
				tokenpredicate.Type(token.LessEqual),
				tokenpredicate.Type(token.LessEqualGreater))) {
			break
		}

		op := binaryOperatorMap[tok.Type]

		var ws1 InternalWhitespaceRun
		ws1.Match(p)

		var right Node
		if !matchExpr12(&right, p) {
			break
		}

		p.Forget(mark)
		mark = p.Mark()
		left = &BinaryOperatorExpr{
			Operator: op,
			Operand0: left,
			Operand1: right,
			WS0:      ws0,
			WS1:      ws1,
		}
	}

	p.Rewind(mark)
	*ptr = left
	return true
}

func matchExpr12(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr13(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	for {
		var ws0 InternalWhitespaceRun
		ws0.Match(p)

		var tok token.Token
		if !p.Consume(&tok, tokenpredicate.Type(token.Bar)) {
			break
		}

		op := binaryOperatorMap[tok.Type]

		var ws1 InternalWhitespaceRun
		ws1.Match(p)

		var right Node
		if !matchExpr13(&right, p) {
			break
		}

		p.Forget(mark)
		mark = p.Mark()
		left = &BinaryOperatorExpr{
			Operator: op,
			Operand0: left,
			Operand1: right,
			WS0:      ws0,
			WS1:      ws1,
		}
	}

	p.Rewind(mark)
	*ptr = left
	return true
}

func matchExpr13(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr14(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	for {
		var ws0 InternalWhitespaceRun
		ws0.Match(p)

		var tok token.Token
		if !p.Consume(&tok, tokenpredicate.Type(token.Caret)) {
			break
		}

		op := binaryOperatorMap[tok.Type]

		var ws1 InternalWhitespaceRun
		ws1.Match(p)

		var right Node
		if !matchExpr14(&right, p) {
			break
		}

		p.Forget(mark)
		mark = p.Mark()
		left = &BinaryOperatorExpr{
			Operator: op,
			Operand0: left,
			Operand1: right,
			WS0:      ws0,
			WS1:      ws1,
		}
	}

	p.Rewind(mark)
	*ptr = left
	return true
}

func matchExpr14(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr15(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	for {
		var ws0 InternalWhitespaceRun
		ws0.Match(p)

		var tok token.Token
		if !p.Consume(&tok, tokenpredicate.Type(token.Ampersand)) {
			break
		}

		op := binaryOperatorMap[tok.Type]

		var ws1 InternalWhitespaceRun
		ws1.Match(p)

		var right Node
		if !matchExpr15(&right, p) {
			break
		}

		p.Forget(mark)
		mark = p.Mark()
		left = &BinaryOperatorExpr{
			Operator: op,
			Operand0: left,
			Operand1: right,
			WS0:      ws0,
			WS1:      ws1,
		}
	}

	p.Rewind(mark)
	*ptr = left
	return true
}

func matchExpr15(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr16(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	var ws0 InternalWhitespaceRun
	ws0.Match(p)

	var tok token.Token
	if !p.Consume(&tok,
		tokenpredicate.Or(
			tokenpredicate.Type(token.LessLess),
			tokenpredicate.Type(token.LessLessBar),
			tokenpredicate.Type(token.GreaterGreater),
			tokenpredicate.Type(token.GreaterGreaterBar))) {
		p.Rewind(mark)
		*ptr = left
		return true
	}

	op := binaryOperatorMap[tok.Type]

	var ws1 InternalWhitespaceRun
	ws1.Match(p)

	var right Node
	if !matchExpr15(&right, p) {
		p.Rewind(mark)
		*ptr = left
		return true
	}

	left = &BinaryOperatorExpr{
		Operator: op,
		Operand0: left,
		Operand1: right,
		WS0:      ws0,
		WS1:      ws1,
	}
	*ptr = left
	return true
}

func matchExpr16(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr17(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	for {
		var ws0 InternalWhitespaceRun
		ws0.Match(p)

		var tok token.Token
		if !p.Consume(&tok,
			tokenpredicate.Or(
				tokenpredicate.Type(token.Plus),
				tokenpredicate.Type(token.Minus))) {
			break
		}

		op := binaryOperatorMap[tok.Type]

		var ws1 InternalWhitespaceRun
		ws1.Match(p)

		var right Node
		if !matchExpr17(&right, p) {
			break
		}

		p.Forget(mark)
		mark = p.Mark()
		left = &BinaryOperatorExpr{
			Operator: op,
			Operand0: left,
			Operand1: right,
			WS0:      ws0,
			WS1:      ws1,
		}
	}

	p.Rewind(mark)
	*ptr = left
	return true
}

func matchExpr17(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr18(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	for {
		var ws0 InternalWhitespaceRun
		ws0.Match(p)

		var tok token.Token
		if !p.Consume(&tok,
			tokenpredicate.Or(
				tokenpredicate.Type(token.Star),
				tokenpredicate.Type(token.Slash),
				tokenpredicate.Type(token.Percent),
				tokenpredicate.Type(token.SlashPercent))) {
			break
		}

		op := binaryOperatorMap[tok.Type]

		var ws1 InternalWhitespaceRun
		ws1.Match(p)

		var right Node
		if !matchExpr18(&right, p) {
			break
		}

		p.Forget(mark)
		mark = p.Mark()
		left = &BinaryOperatorExpr{
			Operator: op,
			Operand0: left,
			Operand1: right,
			WS0:      ws0,
			WS1:      ws1,
		}
	}

	p.Rewind(mark)
	*ptr = left
	return true
}

func matchExpr18(ptr *Node, p *Parser) bool {
	var left Node
	if !matchExpr19(&left, p) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	var ws0 InternalWhitespaceRun
	ws0.Match(p)

	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.StarStar)) {
		p.Rewind(mark)
		*ptr = left
		return true
	}

	op := binaryOperatorMap[tok.Type]

	var ws1 InternalWhitespaceRun
	ws1.Match(p)

	var right Node
	if !matchExpr18(&right, p) {
		p.Rewind(mark)
		*ptr = left
		return true
	}

	left = &BinaryOperatorExpr{
		Operator: op,
		Operand0: left,
		Operand1: right,
		WS0:      ws0,
		WS1:      ws1,
	}
	*ptr = left
	return true
}

func matchExpr19(ptr *Node, p *Parser) bool {
	return matchExpr20(ptr, p)
}

func matchExpr20(ptr *Node, p *Parser) bool {
	return matchExpr21(ptr, p)
}

func matchExpr21(ptr *Node, p *Parser) bool {
	return matchExpr22(ptr, p)
}

func matchExpr22(ptr *Node, p *Parser) bool {
	return matchExpr23(ptr, p)
}

func matchExpr23(ptr *Node, p *Parser) bool {
	return matchExpr24(ptr, p)
}

func matchExpr24(ptr *Node, p *Parser) bool {
	return matchExpr25(ptr, p)
}

func matchExpr25(ptr *Node, p *Parser) bool {
	return false
}

// UnaryOperatorExpr
// {{{

type UnaryOperatorExpr struct {
	WS       InternalWhitespaceRun
	Operand  Node
	Operator UnaryOperator
}

func (expr *UnaryOperatorExpr) Init() {
	*expr = UnaryOperatorExpr{}
}

func (expr *UnaryOperatorExpr) Match(p *Parser) bool {
	return false
}

func (expr *UnaryOperatorExpr) String() string {
	return util.StringImpl(expr)
}

func (expr *UnaryOperatorExpr) GoString() string {
	return util.GoStringImpl(expr)
}

func (expr *UnaryOperatorExpr) EstimateStringLength() uint {
	sum := 0 + uint(len(expr.Operator.String()))
	sum += expr.Operand.EstimateStringLength()
	sum += expr.WS.EstimateStringLength()
	return sum
}

func (expr *UnaryOperatorExpr) EstimateGoStringLength() uint {
	sum := 21 + uint(len(expr.Operator.GoString()))
	sum += expr.Operand.EstimateGoStringLength()
	sum += expr.WS.EstimateGoStringLength()
	return sum
}

func (expr *UnaryOperatorExpr) WriteStringTo(out *strings.Builder) {
	if expr.Operator.IsPrefix() {
		out.WriteString(expr.Operator.String())
		expr.WS.WriteStringTo(out)
		expr.Operand.WriteStringTo(out)
	} else {
		expr.Operand.WriteStringTo(out)
		expr.WS.WriteStringTo(out)
		out.WriteString(expr.Operator.String())
	}
}

func (expr *UnaryOperatorExpr) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("UnaryOperatorExpr{")
	out.WriteString(expr.Operator.GoString())
	out.WriteByte(',')
	expr.Operand.WriteGoStringTo(out)
	out.WriteByte(',')
	expr.WS.WriteGoStringTo(out)
	out.WriteByte('}')
}

func (expr *UnaryOperatorExpr) ComputeStringLengthEstimates() {
	expr.WS.ComputeStringLengthEstimates()
	expr.Operand.ComputeStringLengthEstimates()
}

var _ Node = (*UnaryOperatorExpr)(nil)

// }}}

// BinaryOperatorExpr
// {{{

type BinaryOperatorExpr struct {
	WS0      InternalWhitespaceRun
	WS1      InternalWhitespaceRun
	Operand0 Node
	Operand1 Node
	Operator BinaryOperator
}

func (expr *BinaryOperatorExpr) Init() {
	*expr = BinaryOperatorExpr{}
}

func (expr *BinaryOperatorExpr) Match(p *Parser) bool {
	return false
}

func (expr *BinaryOperatorExpr) String() string {
	return util.StringImpl(expr)
}

func (expr *BinaryOperatorExpr) GoString() string {
	return util.GoStringImpl(expr)
}

func (expr *BinaryOperatorExpr) EstimateStringLength() uint {
	sum := 0 + uint(len(expr.Operator.String()))
	sum += expr.Operand0.EstimateStringLength()
	sum += expr.Operand1.EstimateStringLength()
	sum += expr.WS0.EstimateStringLength()
	sum += expr.WS1.EstimateStringLength()
	return sum
}

func (expr *BinaryOperatorExpr) EstimateGoStringLength() uint {
	sum := 24 + uint(len(expr.Operator.GoString()))
	sum += expr.Operand0.EstimateGoStringLength()
	sum += expr.Operand1.EstimateGoStringLength()
	sum += expr.WS0.EstimateGoStringLength()
	sum += expr.WS1.EstimateGoStringLength()
	return sum
}

func (expr *BinaryOperatorExpr) WriteStringTo(out *strings.Builder) {
	expr.Operand0.WriteStringTo(out)
	expr.WS0.WriteStringTo(out)
	out.WriteString(expr.Operator.String())
	expr.WS1.WriteStringTo(out)
	expr.Operand1.WriteStringTo(out)
}

func (expr *BinaryOperatorExpr) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("BinaryOperatorExpr{")
	out.WriteString(expr.Operator.GoString())
	out.WriteByte(',')
	expr.Operand0.WriteGoStringTo(out)
	out.WriteByte(',')
	expr.Operand1.WriteGoStringTo(out)
	out.WriteByte(',')
	expr.WS0.WriteGoStringTo(out)
	out.WriteByte(',')
	expr.WS1.WriteGoStringTo(out)
	out.WriteByte('}')
}

func (expr *BinaryOperatorExpr) ComputeStringLengthEstimates() {
	expr.Operand0.ComputeStringLengthEstimates()
	expr.Operand1.ComputeStringLengthEstimates()
	expr.WS0.ComputeStringLengthEstimates()
	expr.WS1.ComputeStringLengthEstimates()
}

var _ Node = (*BinaryOperatorExpr)(nil)

// }}}

// TernaryOperatorExpr
// {{{

type TernaryOperatorExpr struct {
	WS0      InternalWhitespaceRun
	WS1      InternalWhitespaceRun
	WS2      InternalWhitespaceRun
	WS3      InternalWhitespaceRun
	Operand0 Node
	Operand1 Node
	Operand2 Node
	Operator TernaryOperator
}

func (expr *TernaryOperatorExpr) Init() {
	*expr = TernaryOperatorExpr{}
}

func (expr *TernaryOperatorExpr) Match(p *Parser) bool {
	return false
}

func (expr *TernaryOperatorExpr) String() string {
	return util.StringImpl(expr)
}

func (expr *TernaryOperatorExpr) GoString() string {
	return util.GoStringImpl(expr)
}

func (expr *TernaryOperatorExpr) EstimateStringLength() uint {
	sum := 0 + uint(len(expr.Operator.FirstSymbol())) + uint(len(expr.Operator.SecondSymbol()))
	sum += expr.Operand0.EstimateStringLength()
	sum += expr.Operand1.EstimateStringLength()
	sum += expr.Operand2.EstimateStringLength()
	sum += expr.WS0.EstimateStringLength()
	sum += expr.WS1.EstimateStringLength()
	sum += expr.WS2.EstimateStringLength()
	sum += expr.WS3.EstimateStringLength()
	return sum
}

func (expr *TernaryOperatorExpr) EstimateGoStringLength() uint {
	sum := 28 + uint(len(expr.Operator.GoString()))
	sum += expr.Operand0.EstimateGoStringLength()
	sum += expr.Operand1.EstimateGoStringLength()
	sum += expr.Operand2.EstimateGoStringLength()
	sum += expr.WS0.EstimateGoStringLength()
	sum += expr.WS1.EstimateGoStringLength()
	sum += expr.WS2.EstimateGoStringLength()
	sum += expr.WS3.EstimateGoStringLength()
	return sum
}

func (expr *TernaryOperatorExpr) WriteStringTo(out *strings.Builder) {
	isInverted := expr.Operator.IsInverted()
	if isInverted {
		expr.Operand1.WriteStringTo(out)
		expr.WS1.WriteStringTo(out)
		out.WriteString(expr.Operator.FirstSymbol())
		expr.WS0.WriteStringTo(out)
		expr.Operand0.WriteStringTo(out)
	} else {
		expr.Operand0.WriteStringTo(out)
		expr.WS0.WriteStringTo(out)
		out.WriteString(expr.Operator.FirstSymbol())
		expr.WS1.WriteStringTo(out)
		expr.Operand1.WriteStringTo(out)
	}
	expr.WS2.WriteStringTo(out)
	out.WriteString(expr.Operator.SecondSymbol())
	expr.WS3.WriteStringTo(out)
	expr.Operand2.WriteStringTo(out)
}

func (expr *TernaryOperatorExpr) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("TernaryOperatorExpr{")
	out.WriteString(expr.Operator.GoString())
	out.WriteByte(',')
	expr.Operand0.WriteGoStringTo(out)
	out.WriteByte(',')
	expr.Operand1.WriteGoStringTo(out)
	out.WriteByte(',')
	expr.Operand2.WriteGoStringTo(out)
	out.WriteByte(',')
	expr.WS0.WriteGoStringTo(out)
	out.WriteByte(',')
	expr.WS1.WriteGoStringTo(out)
	out.WriteByte(',')
	expr.WS2.WriteGoStringTo(out)
	out.WriteByte(',')
	expr.WS3.WriteGoStringTo(out)
	out.WriteByte('}')
}

func (expr *TernaryOperatorExpr) ComputeStringLengthEstimates() {
	expr.Operand0.ComputeStringLengthEstimates()
	expr.Operand1.ComputeStringLengthEstimates()
	expr.Operand2.ComputeStringLengthEstimates()
	expr.WS0.ComputeStringLengthEstimates()
	expr.WS1.ComputeStringLengthEstimates()
	expr.WS2.ComputeStringLengthEstimates()
	expr.WS3.ComputeStringLengthEstimates()
}

var _ Node = (*TernaryOperatorExpr)(nil)

// }}}
