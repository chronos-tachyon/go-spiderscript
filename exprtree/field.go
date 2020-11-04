package exprtree

import (
	"fmt"
)

// Layout
// {{{

type Layout struct {
	fields     []*Field
	bySymbol   map[*Symbol]*Field
	byOffset   map[uint]*Field
	alignShift uint8
	minSize    uint16
	padSize    uint16
}

func (layout *Layout) AlignShift() uint {
	return uint(layout.alignShift)
}

func (layout *Layout) AlignBytes() uint {
	return uint(1) << layout.alignShift
}

func (layout *Layout) MinimumSize() uint {
	return uint(layout.minSize)
}

func (layout *Layout) PaddedSize() uint {
	return uint(layout.padSize)
}

func (layout *Layout) Fields() []*Field {
	return cloneFields(layout.fields)
}

func (layout *Layout) FieldBySymbol(sym *Symbol) (*Field, bool) {
	field, found := layout.bySymbol[sym]
	return field, found
}

func (layout *Layout) FieldByOffset(offset uint) (*Field, bool) {
	field, found := layout.byOffset[offset]
	return field, found
}

// }}}

// Field
// {{{

type Field struct {
	sym    *Symbol
	offset uint
	length uint
}

func (field *Field) Check() {
	if field == nil {
		panic(fmt.Errorf("BUG: *Field is nil"))
	}
}

func (field *Field) Interp() *Interp {
	return field.Symbol().Interp()
}

func (field *Field) Symbol() *Symbol {
	return field.sym
}

func (field *Field) CanonicalName() string {
	return field.Symbol().CanonicalName()
}

func (field *Field) MangledName() string {
	return field.Symbol().MangledName()
}

func (field *Field) Type() *Type {
	return field.Symbol().Type()
}

func (field *Field) Offset() uint {
	return field.offset
}

func (field *Field) Length() uint {
	return field.length
}

func (field *Field) BindTo(mem *Memory) *Value {
	value := &Value{
		field:  field,
		memory: mem,
	}
	return value
}

// }}}
