package value

import (
	"strconv"
	"strings"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
)

type Regex struct {
	Input string
}

func (rx *Regex) String() string {
	return util.StringImpl(rx)
}

func (rx *Regex) GoString() string {
	return util.GoStringImpl(rx)
}

func (rx *Regex) EstimateStringLength() uint {
	return uint(len(rx.Input))
}

func (rx *Regex) EstimateGoStringLength() uint {
	return 16 + uint(len(rx.Input))
}

func (rx *Regex) WriteStringTo(out *strings.Builder) {
	out.WriteString(rx.Input)
}

func (rx *Regex) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("&value.Regex{")
	out.WriteString(strconv.Quote(rx.Input))
	out.WriteByte('}')
}

func (rx *Regex) Parse(input []rune) error {
	rx.Input = string(input)
	return nil
}

var _ Value = (*Regex)(nil)
