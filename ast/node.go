package ast

import (
	"strings"
)

type Node interface {
	Init()

	Match(p *Parser) bool

	String() string
	GoString() string
	EstimateStringLength() uint
	EstimateGoStringLength() uint
	WriteStringTo(out *strings.Builder)
	WriteGoStringTo(out *strings.Builder)
	ComputeStringLengthEstimates()
}

func sumStringLengthEstimates(list []Node) uint {
	var sum uint
	for _, item := range list {
		sum += item.EstimateStringLength()
	}
	return sum
}

func sumGoStringLengthEstimates(list []Node) uint {
	sum := uint(len(list))
	if len(list) != 0 {
		sum--
	}
	for _, item := range list {
		sum += item.EstimateGoStringLength()
	}
	return sum
}

func writeStringsTo(out *strings.Builder, list []Node) {
	for _, item := range list {
		item.WriteStringTo(out)
	}
}

func writeGoStringsTo(out *strings.Builder, list []Node) {
	for i, item := range list {
		if i > 0 {
			out.WriteByte(',')
		}
		item.WriteGoStringTo(out)
	}
}

func sumStringLengthEstimates2D(listOfLists [][]Node) uint {
	var sum uint
	for _, list := range listOfLists {
		for _, item := range list {
			sum += item.EstimateStringLength()
		}
	}
	return sum
}

func sumGoStringLengthEstimates2D(listOfLists [][]Node) uint {
	sum := 3 * uint(len(listOfLists))
	if len(listOfLists) != 0 {
		sum--
	}
	for _, list := range listOfLists {
		sum += uint(len(list))
		if len(list) != 0 {
			sum--
		}
		for _, item := range list {
			sum += item.EstimateGoStringLength()
		}
	}
	return sum
}

func writeStringsTo2D(out *strings.Builder, listOfLists [][]Node) {
	for _, list := range listOfLists {
		for _, item := range list {
			item.WriteStringTo(out)
		}
	}
}

func writeGoStringsTo2D(out *strings.Builder, listOfLists [][]Node) {
	for i, list := range listOfLists {
		if i > 0 {
			out.WriteByte(',')
		}
		out.WriteByte('[')
		for j, item := range list {
			if j > 0 {
				out.WriteByte(',')
			}
			item.WriteGoStringTo(out)
		}
		out.WriteByte(']')
	}
}
