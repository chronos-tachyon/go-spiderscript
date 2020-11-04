package value

import (
	"strconv"
	"strings"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
)

type Literal struct {
	Text string
}

func NewLiteral(text string) *Literal {
	return &Literal{Text: text}
}

func (lv *Literal) String() string {
	return util.StringImpl(lv)
}

func (lv *Literal) GoString() string {
	return util.GoStringImpl(lv)
}

func (lv *Literal) EstimateStringLength() uint {
	return uint(len(lv.Text))
}

func (lv *Literal) EstimateGoStringLength() uint {
	return 18 + uint(len(lv.Text))
}

func (lv *Literal) WriteStringTo(out *strings.Builder) {
	out.WriteString(lv.Text)
}

func (lv *Literal) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("&value.Literal{")
	out.WriteString(strconv.Quote(lv.Text))
	out.WriteByte('}')
}

var _ Value = (*Literal)(nil)
