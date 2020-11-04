package operator

import (
	"fmt"
)

type Kind byte

const (
	InvalidKind Kind = iota
	UnaryPrefix
	UnaryPostfix
	UnaryOther
	BinaryInfix
	BinaryOther
	TernaryOther
	AssignStatement
	MutateStatement
)

var kindNames = []string{
	"InvalidKind",
	"UnaryPrefix",
	"UnaryPostfix",
	"UnaryOther",
	"BinaryInfix",
	"BinaryOther",
	"TernaryOther",
	"AssignStatement",
	"MutateStatement",
}

func (kind Kind) String() string {
	if uint(kind) >= uint(len(kindNames)) {
		return fmt.Sprintf("Kind(%d)", uint(kind))
	}
	return kindNames[kind]
}

func (kind Kind) GoString() string {
	return kind.String()
}

func (kind Kind) IsUnary() bool {
	switch kind {
	case UnaryPrefix, UnaryPostfix, UnaryOther:
		return true
	default:
		return false
	}
}

func (kind Kind) IsBinary() bool {
	switch kind {
	case BinaryInfix, BinaryOther:
		return true
	default:
		return false
	}
}

func (kind Kind) IsTernary() bool {
	switch kind {
	case TernaryOther:
		return true
	default:
		return false
	}
}

func (kind Kind) IsAssignment() bool {
	switch kind {
	case AssignStatement:
		return true
	default:
		return false
	}
}

func (kind Kind) IsMutation() bool {
	switch kind {
	case MutateStatement:
		return true
	default:
		return false
	}
}

var _ fmt.Stringer = Kind(0)
var _ fmt.GoStringer = Kind(0)
