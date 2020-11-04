package token

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
	"github.com/chronos-tachyon/go-spiderscript/value"
)

type Position struct {
	Path            string
	RuneOffset      uint
	RawLineNumber   uint
	RawColumnNumber uint
	ConsumeNextLF   bool
}

func (pos *Position) Init(path string) {
	*pos = Position{Path: path}
}

func (pos *Position) Advance(ch rune) {
	pos.RuneOffset++

	if pos.ConsumeNextLF {
		pos.ConsumeNextLF = false
		if ch == '\n' {
			return
		}
	}

	switch ch {
	case '\t':
		pos.RawColumnNumber += (8 - (pos.RawColumnNumber & 7))

	case '\n', '\v', '\f':
		pos.RawLineNumber++
		pos.RawColumnNumber = 0

	case '\r':
		pos.RawLineNumber++
		pos.RawColumnNumber = 0
		pos.ConsumeNextLF = true

	default:
		pos.RawColumnNumber++
	}
}

func (pos Position) LineNumber() uint {
	return pos.RawLineNumber + 1
}

func (pos Position) ColumnNumber() uint {
	return pos.RawColumnNumber + 1
}

func (pos Position) String() string {
	return util.StringImpl(pos)
}

func (pos Position) GoString() string {
	return util.GoStringImpl(pos)
}

func (pos Position) EstimateStringLength() uint {
	return 24 + uint(len(pos.Path))
}

func (pos Position) EstimateGoStringLength() uint {
	return 64 + uint(len(pos.Path))
}

func (pos Position) WriteStringTo(out *strings.Builder) {
	fmt.Fprintf(out, "%q(L%d,C%d,+%d)", pos.Path, pos.LineNumber(), pos.ColumnNumber(), pos.RuneOffset)
}

func (pos Position) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("Position{")
	out.WriteString(strconv.Quote(pos.Path))
	out.WriteByte(',')
	out.WriteString(util.Itoa(pos.RuneOffset))
	out.WriteByte(',')
	out.WriteString(util.Itoa(pos.RawLineNumber))
	out.WriteByte(',')
	out.WriteString(util.Itoa(pos.RawColumnNumber))
	out.WriteByte(',')
	out.WriteString(strconv.FormatBool(pos.ConsumeNextLF))
	out.WriteByte('}')
}

var _ value.Value = Position{}
var _ value.Value = (*Position)(nil)
