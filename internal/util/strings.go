package util

import (
	"strings"
)

type WriteStringToer interface {
	EstimateStringLength() uint
	WriteStringTo(out *strings.Builder)
}

func StringImpl(value WriteStringToer) string {
	var sb strings.Builder
	sb.Grow(int(value.EstimateStringLength()))
	value.WriteStringTo(&sb)
	return sb.String()
}

type WriteGoStringToer interface {
	EstimateGoStringLength() uint
	WriteGoStringTo(out *strings.Builder)
}

func GoStringImpl(value WriteGoStringToer) string {
	var sb strings.Builder
	sb.Grow(int(value.EstimateGoStringLength()))
	value.WriteGoStringTo(&sb)
	return sb.String()
}
