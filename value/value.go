package value

import (
	"strings"
)

type Value interface {
	String() string
	GoString() string
	EstimateStringLength() uint
	EstimateGoStringLength() uint
	WriteStringTo(out *strings.Builder)
	WriteGoStringTo(out *strings.Builder)
}
