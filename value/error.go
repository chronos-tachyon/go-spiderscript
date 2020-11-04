package value

import (
	"fmt"
	"strings"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
)

type Error struct {
	Err error

	PrecomputedStringLength   uint
	PrecomputedGoStringLength uint
}

func (ev *Error) String() string {
	return util.StringImpl(ev)
}

func (ev *Error) GoString() string {
	return util.GoStringImpl(ev)
}

func (ev *Error) EstimateStringLength() uint {
	return ev.PrecomputedStringLength
}

func (ev *Error) EstimateGoStringLength() uint {
	return ev.PrecomputedGoStringLength
}

func (ev *Error) WriteStringTo(out *strings.Builder) {
	fmt.Fprintf(out, "error: %v", ev.Err)
}

func (ev *Error) WriteGoStringTo(out *strings.Builder) {
	fmt.Fprintf(out, "&value.Error{%#v}", ev.Err)
}

func (ev *Error) SetEstimatedStringLength(length uint) {
	ev.PrecomputedStringLength = length
}

func (ev *Error) SetEstimatedGoStringLength(length uint) {
	ev.PrecomputedGoStringLength = length
}

var _ Value = (*Error)(nil)
var _ util.Estimable = (*Error)(nil)
