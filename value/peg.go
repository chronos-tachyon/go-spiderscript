package value

import (
	"strconv"
	"strings"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
)

type PEG struct {
	Input string
}

func (pv *PEG) String() string {
	return util.StringImpl(pv)
}

func (pv *PEG) GoString() string {
	return util.GoStringImpl(pv)
}

func (pv *PEG) EstimateStringLength() uint {
	return uint(len(pv.Input))
}

func (pv *PEG) EstimateGoStringLength() uint {
	return 14 + uint(len(pv.Input))
}

func (pv *PEG) WriteStringTo(out *strings.Builder) {
	out.WriteString(pv.Input)
}

func (pv *PEG) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("&value.PEG{")
	out.WriteString(strconv.Quote(pv.Input))
	out.WriteByte('}')
}

func (pv *PEG) Parse(input []rune) error {
	pv.Input = string(input)
	return nil
}

var _ Value = (*PEG)(nil)
