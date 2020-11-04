package exprtree

import (
	"sync"
)

// Type
// {{{

type Type struct {
	mu         sync.RWMutex
	interp     *Interp
	sym        *Symbol
	id         TypeID
	kind       TypeKind
	alignShift uint8
	minSize    uint16
	padSize    uint16
	static     SymbolTable
	instance   SymbolTable
	details    interface{}
}

func (t *Type) ID() TypeID {
	return t.id
}

func (t *Type) Interp() *Interp {
	return t.interp
}

func (t *Type) Symbol() *Symbol {
	return t.sym
}

func (t *Type) CanonicalName() string {
	return t.Symbol().CanonicalName()
}

func (t *Type) MangledName() string {
	return t.Symbol().MangledName()
}

func (t *Type) Kind() TypeKind {
	return t.kind
}

func (t *Type) AlignShift() uint {
	return uint(t.alignShift)
}

func (t *Type) AlignBytes() uint {
	return uint(1) << t.AlignShift()
}

func (t *Type) MinimumBytes() uint {
	return uint(t.minSize)
}

func (t *Type) PaddedBytes() uint {
	return uint(t.padSize)
}

func (t *Type) StaticSymbols() *SymbolTable {
	return &t.static
}

func (t *Type) InstanceSymbols() *SymbolTable {
	return &t.instance
}

func (t *Type) Details() interface{} {
	return t.details
}

func (t *Type) Chase() *Type {
	for {
		switch t.kind {
		case MutableKind:
			fallthrough
		case ConstKind:
			fallthrough
		case NamedKind:
			t = t.Details().(*Type)

		default:
			return t
		}
	}
}

func (t *Type) Elem() *Type {
	switch t.kind {
	case PointerKind:
		return t.Details().(*Type)

	default:
		return nil
	}
}

func (t *Type) Is(other *Type) bool {
	for {
		if t == other {
			return true
		}

		switch t.kind {
		case MutableKind:
			fallthrough
		case ConstKind:
			fallthrough
		case NamedKind:
			t = t.Details().(*Type)

		default:
			return false
		}
	}
}

// }}}
