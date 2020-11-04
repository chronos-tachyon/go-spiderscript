package util

import (
	"strings"
)

type Estimable interface {
	WriteStringTo(out *strings.Builder)
	WriteGoStringTo(out *strings.Builder)
	SetEstimatedStringLength(length uint)
	SetEstimatedGoStringLength(length uint)
}

func EstimateLengths(value Estimable) {
	var sb strings.Builder
	sb.Grow(64)
	value.WriteStringTo(&sb)
	value.SetEstimatedStringLength(uint(len(sb.String())))
	sb.Reset()
	value.WriteGoStringTo(&sb)
	value.SetEstimatedGoStringLength(uint(len(sb.String())))
}
