package token

import (
	"strings"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
	"github.com/chronos-tachyon/go-spiderscript/value"
)

type Token struct {
	Type   Type
	Raw    []rune
	Parsed value.Value
	Start  Position
	End    Position

	PrecomputedStringLength   uint
	PrecomputedGoStringLength uint
}

func (tok Token) String() string {
	return util.StringImpl(tok)
}

func (tok Token) GoString() string {
	return util.GoStringImpl(tok)
}

func (tok Token) EstimateStringLength() uint {
	return tok.PrecomputedStringLength
}

func (tok Token) EstimateGoStringLength() uint {
	return tok.PrecomputedGoStringLength
}

func (tok Token) WriteStringTo(out *strings.Builder) {
	out.WriteString(tok.Type.String())
	if tok.Parsed != nil {
		out.WriteByte('[')
		tok.Parsed.WriteStringTo(out)
		out.WriteByte(']')
	}
}

func (tok Token) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("Token{")
	out.WriteString(tok.Type.GoString())
	out.WriteByte(',')
	if tok.Parsed == nil {
		out.Write(util.NilBytes)
	} else {
		tok.Parsed.WriteGoStringTo(out)
	}
	out.WriteByte(',')
	tok.Start.WriteGoStringTo(out)
	out.WriteByte(',')
	tok.End.WriteGoStringTo(out)
	out.WriteByte('}')
}

func (tok *Token) SetEstimatedStringLength(length uint) {
	tok.PrecomputedStringLength = length
}

func (tok *Token) SetEstimatedGoStringLength(length uint) {
	tok.PrecomputedGoStringLength = length
}

var _ value.Value = Token{}
var _ value.Value = (*Token)(nil)
var _ util.Estimable = (*Token)(nil)
