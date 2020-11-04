package ast

import (
	"strings"
	"sync"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
	"github.com/chronos-tachyon/go-spiderscript/token"
	"github.com/chronos-tachyon/go-spiderscript/tokenpredicate"
	"github.com/chronos-tachyon/go-spiderscript/value"
)

// ModuleName
// {{{

type ModuleName struct {
	Names []string
}

func (mod *ModuleName) Init() {
	*mod = ModuleName{}
}

func (mod *ModuleName) Match(p *Parser) bool {
	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.Identifier)) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	if mod.Names == nil {
		mod.Names = make([]string, 0, 4)
	}

	word, err := value.ExtractLiteral(tok.Parsed, "<error>")
	if err != nil {
		p.EmitError(err)
	}
	mod.Names = append(mod.Names, word)

	for {
		if !p.Consume(nil, tokenpredicate.Type(token.ColonColon)) {
			p.Rewind(mark)
			break
		}

		if !p.Consume(&tok, tokenpredicate.Type(token.Identifier)) {
			p.Rewind(mark)
			break
		}

		p.Forget(mark)
		mark = p.Mark()

		word, err = value.ExtractLiteral(tok.Parsed, "<error>")
		if err != nil {
			p.EmitError(err)
		}
		mod.Names = append(mod.Names, word)
	}

	return true
}

func (mod *ModuleName) String() string {
	return util.StringImpl(mod)
}

func (mod *ModuleName) GoString() string {
	return util.GoStringImpl(mod)
}

func (mod *ModuleName) EstimateStringLength() uint {
	if mod.IsEmpty() {
		return 0
	}

	var sum uint
	for _, word := range mod.Names {
		sum += 2 + uint(len(word))
	}
	if len(mod.Names) != 0 {
		sum -= 2
	}
	return sum
}

func (mod *ModuleName) EstimateGoStringLength() uint {
	if mod.IsEmpty() {
		return 3
	}

	return 14 + mod.EstimateStringLength()
}

func (mod *ModuleName) WriteStringTo(out *strings.Builder) {
	if mod.IsEmpty() {
		return
	}

	for index, word := range mod.Names {
		if index > 0 {
			out.WriteString("::")
		}
		out.WriteString(word)
	}
}

func (mod *ModuleName) WriteGoStringTo(out *strings.Builder) {
	if mod.IsEmpty() {
		out.WriteString("nil")
		return
	}

	out.WriteString("ModuleName{\"")
	mod.WriteStringTo(out)
	out.WriteString("\"}")
}

func (mod *ModuleName) ComputeStringLengthEstimates() {
}

func (mod *ModuleName) IsEmpty() bool {
	return mod == nil || len(mod.Names) == 0
}

var moduleNameInternMutex sync.Mutex
var moduleNameInternCache map[string]*ModuleName

func (mod *ModuleName) Intern() *ModuleName {
	if mod.IsEmpty() {
		return nil
	}

	key := mod.String()

	moduleNameInternMutex.Lock()
	defer moduleNameInternMutex.Unlock()

	if moduleNameInternCache == nil {
		moduleNameInternCache = make(map[string]*ModuleName, 32)
	}

	if existing, found := moduleNameInternCache[key]; found {
		return existing
	}
	moduleNameInternCache[key] = mod
	return mod
}

var _ Node = (*ModuleName)(nil)

// }}}

// Identifier
// {{{

type Identifier struct {
	Name string
}

func (ident *Identifier) Init() {
	*ident = Identifier{}
}

func (ident *Identifier) Match(p *Parser) bool {
	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.Identifier)) {
		return false
	}

	word, err := value.ExtractLiteral(tok.Parsed, "<error>")
	if err != nil {
		p.EmitError(err)
	}
	ident.Name = word
	return true
}

func (ident *Identifier) String() string {
	return util.StringImpl(ident)
}

func (ident *Identifier) GoString() string {
	return util.GoStringImpl(ident)
}

func (ident *Identifier) EstimateStringLength() uint {
	if ident.IsEmpty() {
		return 1
	}

	return uint(len(ident.Name))
}

func (ident *Identifier) EstimateGoStringLength() uint {
	if ident.IsEmpty() {
		return 3
	}

	return 14 + uint(len(ident.Name))
}

func (ident *Identifier) WriteStringTo(out *strings.Builder) {
	if ident.IsEmpty() {
		out.WriteByte('_')
		return
	}

	out.WriteString(ident.Name)
}

func (ident *Identifier) WriteGoStringTo(out *strings.Builder) {
	if ident.IsEmpty() {
		out.WriteString("nil")
		return
	}

	out.WriteString("Identifier{\"")
	out.WriteString(ident.Name)
	out.WriteString("\"}")
}

func (ident *Identifier) ComputeStringLengthEstimates() {
}

func (ident *Identifier) IsEmpty() bool {
	return ident == nil || len(ident.Name) == 0
}

var identifierInternMutex sync.Mutex
var identifierInternCache map[string]*Identifier

func (ident *Identifier) Intern() *Identifier {
	if ident.IsEmpty() {
		return nil
	}

	key := ident.Name

	identifierInternMutex.Lock()
	defer identifierInternMutex.Unlock()

	if identifierInternCache == nil {
		identifierInternCache = make(map[string]*Identifier, 32)
	}

	if existing, found := identifierInternCache[key]; found {
		return existing
	}
	identifierInternCache[key] = ident
	return ident
}

var _ Node = (*Identifier)(nil)

// }}}

// Symbol
// {{{

type Symbol struct {
	Module     *ModuleName
	Identifier *Identifier
	TypeSig    *TypeSignature
	FuncSig    *FuncSignature
}

func (sym *Symbol) Init() {
	*sym = Symbol{}
}

func (sym *Symbol) Match(p *Parser) bool {
	var tok token.Token
	if !p.Consume(&tok, tokenpredicate.Type(token.Identifier)) {
		return false
	}

	mark := p.Mark()
	defer p.Forget(mark)

	words := make([]string, 0, 4)
	word, err := value.ExtractLiteral(tok.Parsed, "<error>")
	if err != nil {
		p.EmitError(err)
	}
	words = append(words, word)

	for {
		if !p.Consume(nil, tokenpredicate.Type(token.ColonColon)) {
			p.Rewind(mark)
			break
		}

		if !p.Consume(&tok, tokenpredicate.Type(token.Identifier)) {
			p.Rewind(mark)
			break
		}

		p.Forget(mark)
		mark = p.Mark()

		word, err = value.ExtractLiteral(tok.Parsed, "<error>")
		if err != nil {
			p.EmitError(err)
		}
		words = append(words, word)
	}

	n := uint(len(words)) - 1
	sym.Module = &ModuleName{Names: words[:n]}
	sym.Identifier = &Identifier{Name: words[n]}

	var tsig TypeSignature
	if tsig.Match(p) {
		sym.TypeSig = &tsig
		return true
	}

	var fsig FuncSignature
	if fsig.Match(p) {
		sym.FuncSig = &fsig
		return true
	}

	return true
}

func (sym *Symbol) String() string {
	return util.StringImpl(sym)
}

func (sym *Symbol) GoString() string {
	return util.GoStringImpl(sym)
}

func (sym *Symbol) EstimateStringLength() uint {
	if sym.IsEmpty() {
		return 1
	}

	sum := uint(len(sym.Identifier.Name))
	if sym.Module != nil {
		for _, word := range sym.Module.Names {
			sum += 2 + uint(len(word))
		}
	}
	return sum
}

func (sym *Symbol) EstimateGoStringLength() uint {
	if sym.IsEmpty() {
		return 3
	}

	return 9 + sym.Identifier.EstimateGoStringLength() + sym.Module.EstimateGoStringLength()
}

func (sym *Symbol) WriteStringTo(out *strings.Builder) {
	if sym.IsEmpty() {
		out.WriteByte('_')
		return
	}

	if sym.Module != nil {
		for _, word := range sym.Module.Names {
			out.WriteString(word)
			out.WriteString("::")
		}
	}
	out.WriteString(sym.Identifier.Name)
}

func (sym *Symbol) WriteGoStringTo(out *strings.Builder) {
	if sym.IsEmpty() {
		out.WriteString("nil")
		return
	}

	out.WriteString("Symbol{")
	sym.Module.WriteGoStringTo(out)
	out.WriteByte(',')
	sym.Identifier.WriteGoStringTo(out)
	out.WriteRune('}')
}

func (sym *Symbol) ComputeStringLengthEstimates() {
}

func (sym *Symbol) IsEmpty() bool {
	return sym == nil || sym.Identifier.IsEmpty()
}

var symbolInternMutex sync.Mutex
var symbolInternCache map[string]*Symbol

func (sym *Symbol) Intern() *Symbol {
	if sym.IsEmpty() {
		return nil
	}

	sym.Module = sym.Module.Intern()
	sym.Identifier = sym.Identifier.Intern()
	key := sym.String()

	symbolInternMutex.Lock()
	defer symbolInternMutex.Unlock()

	if symbolInternCache == nil {
		symbolInternCache = make(map[string]*Symbol, 32)
	}

	if existing, found := symbolInternCache[key]; found {
		return existing
	}
	symbolInternCache[key] = sym
	return sym
}

var _ Node = (*Symbol)(nil)

// }}}

// TypeSignature
// {{{

type TypeSignature struct {
	Args             []Node
	WS0              []Node
	WS1              []Node
	WS2              Node
	HasTrailingComma bool
}

func (sig *TypeSignature) Init() {
	*sig = TypeSignature{}
}

func (sig *TypeSignature) Match(p *Parser) bool {
	mark := p.Mark()
	defer p.Forget(mark)

	if !p.Consume(nil, tokenpredicate.Type(token.Hash)) {
		p.Rewind(mark)
		return false
	}
	if !p.Consume(nil, tokenpredicate.Type(token.LBracket)) {
		p.Rewind(mark)
		return false
	}

	var ws InternalWhitespaceRun
	ws.Match(p)

	if p.Consume(nil, tokenpredicate.Type(token.RBracket)) {
		sig.WS2 = &ws
		return true
	}

	if sig.Args == nil {
		sig.Args = make([]Node, 0, 4)
	}
	if sig.WS0 == nil {
		sig.WS0 = make([]Node, 0, 4)
	}
	if sig.WS1 == nil {
		sig.WS1 = make([]Node, 0, 4)
	}

	for {
		ws0 := ws
		ws.Init()

		var arg Node
		if !MatchExpr(&arg, p) {
			break
		}

		var ws1 InternalWhitespaceRun
		ws1.Match(p)

		sig.Args = append(sig.Args, arg)
		sig.WS0 = append(sig.WS0, &ws0)
		sig.WS1 = append(sig.WS1, &ws1)

		if p.Consume(nil, tokenpredicate.Type(token.RBracket)) {
			return true
		}

		if !p.Consume(nil, tokenpredicate.Type(token.Comma)) {
			break
		}

		ws.Match(p)

		if p.Consume(nil, tokenpredicate.Type(token.RBracket)) {
			sig.WS2 = &ws
			sig.HasTrailingComma = true
			return true
		}
	}

	p.Rewind(mark)
	return false
}

func (sig *TypeSignature) String() string {
	return util.StringImpl(sig)
}

func (sig *TypeSignature) GoString() string {
	return util.GoStringImpl(sig)
}

func (sig *TypeSignature) EstimateStringLength() uint {
	sum := 3 + sig.WS2.EstimateStringLength() + uint(len(sig.Args))
	sum += sumStringLengthEstimates(sig.Args)
	sum += sumStringLengthEstimates(sig.WS0)
	sum += sumStringLengthEstimates(sig.WS1)
	if len(sig.Args) != 0 {
		sum--
	}
	if sig.HasTrailingComma {
		sum++
	}
	return sum
}

func (sig *TypeSignature) EstimateGoStringLength() uint {
	sum := 24 + sig.WS2.EstimateGoStringLength()
	sum += sumGoStringLengthEstimates(sig.Args)
	sum += sumGoStringLengthEstimates(sig.WS0)
	sum += sumGoStringLengthEstimates(sig.WS1)
	return sum
}

func (sig *TypeSignature) WriteStringTo(out *strings.Builder) {
	out.WriteByte('#')
	out.WriteByte('[')
	for index := uint(0); index < uint(len(sig.Args)); index++ {
		if index != 0 {
			out.WriteByte(',')
		}
		sig.WS0[index].WriteStringTo(out)
		sig.Args[index].WriteStringTo(out)
		sig.WS1[index].WriteStringTo(out)
	}
	if sig.HasTrailingComma {
		out.WriteByte(',')
	}
	sig.WS2.WriteStringTo(out)
}

func (sig *TypeSignature) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("TypeSignature{")
	out.WriteByte('[')
	writeGoStringsTo(out, sig.Args)
	out.WriteByte(']')
	out.WriteByte(',')
	out.WriteByte('[')
	writeGoStringsTo(out, sig.WS0)
	out.WriteByte(']')
	out.WriteByte(',')
	out.WriteByte('[')
	writeGoStringsTo(out, sig.WS1)
	out.WriteByte(']')
	out.WriteByte(',')
	sig.WS2.WriteGoStringTo(out)
	out.WriteByte('}')
}

func (sig *TypeSignature) ComputeStringLengthEstimates() {
	for index := uint(0); index < uint(len(sig.Args)); index++ {
		sig.Args[index].ComputeStringLengthEstimates()
		sig.WS0[index].ComputeStringLengthEstimates()
		sig.WS1[index].ComputeStringLengthEstimates()
	}
	sig.WS2.ComputeStringLengthEstimates()
}

var _ Node = (*TypeSignature)(nil)

// }}}

// FuncSignature
// {{{

type FuncSignature struct {
}

func (sig *FuncSignature) Init() {
	*sig = FuncSignature{}
}

func (sig *FuncSignature) Match(p *Parser) bool {
	mark := p.Mark()
	defer p.Forget(mark)

	if !p.Consume(nil, tokenpredicate.Type(token.Hash)) {
		p.Rewind(mark)
		return false
	}
	if !p.Consume(nil, tokenpredicate.Type(token.LParen)) {
		p.Rewind(mark)
		return false
	}

	p.Rewind(mark)
	return false
}

func (sig *FuncSignature) String() string {
	return util.StringImpl(sig)
}

func (sig *FuncSignature) GoString() string {
	return util.GoStringImpl(sig)
}

func (sig *FuncSignature) EstimateStringLength() uint {
	return 0
}

func (sig *FuncSignature) EstimateGoStringLength() uint {
	return 0
}

func (sig *FuncSignature) WriteStringTo(out *strings.Builder) {
}

func (sig *FuncSignature) WriteGoStringTo(out *strings.Builder) {
}

func (sig *FuncSignature) ComputeStringLengthEstimates() {
}

var _ Node = (*FuncSignature)(nil)

// }}}
