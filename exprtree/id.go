package exprtree

import (
	"fmt"
)

// SymbolID
// {{{

type SymbolID uint32

func (id SymbolID) String() string {
	return fmt.Sprintf("symbol #%d", uint32(id))
}

func (id SymbolID) GoString() string {
	return fmt.Sprintf("SymbolID(%d)", uint32(id))
}

var _ fmt.Stringer = SymbolID(0)
var _ fmt.GoStringer = SymbolID(0)

// }}}

// TypeID
// {{{

type TypeID uint32

func (id TypeID) String() string {
	return fmt.Sprintf("type #%d", uint32(id))
}

func (id TypeID) GoString() string {
	return fmt.Sprintf("TypeID(%d)", uint32(id))
}

var _ fmt.Stringer = TypeID(0)
var _ fmt.GoStringer = TypeID(0)

// }}}

// BufferID
// {{{

type BufferID uint32

func (id BufferID) String() string {
	return fmt.Sprintf("buffer #%d", uint32(id))
}

func (id BufferID) GoString() string {
	return fmt.Sprintf("BufferID(%d)", uint32(id))
}

var _ fmt.Stringer = BufferID(0)
var _ fmt.GoStringer = BufferID(0)

// }}}

// ErrorID
// {{{

type ErrorID uint32

func (id ErrorID) String() string {
	return fmt.Sprintf("error #%d", uint32(id))
}

func (id ErrorID) GoString() string {
	return fmt.Sprintf("ErrorID(%d)", uint32(id))
}

var _ fmt.Stringer = ErrorID(0)
var _ fmt.GoStringer = ErrorID(0)

// }}}

// GenericSignatureID
// {{{

type GenericSignatureID uint32

func (id GenericSignatureID) String() string {
	return fmt.Sprintf("generic signature #%d", uint32(id))
}

func (id GenericSignatureID) GoString() string {
	return fmt.Sprintf("GenericSignatureID(%d)", uint32(id))
}

var _ fmt.Stringer = GenericSignatureID(0)
var _ fmt.GoStringer = GenericSignatureID(0)

// }}}

// FunctionSignatureID
// {{{

type FunctionSignatureID uint32

func (id FunctionSignatureID) String() string {
	return fmt.Sprintf("function signature #%d", uint32(id))
}

func (id FunctionSignatureID) GoString() string {
	return fmt.Sprintf("FunctionSignatureID(%d)", uint32(id))
}

var _ fmt.Stringer = FunctionSignatureID(0)
var _ fmt.GoStringer = FunctionSignatureID(0)

// }}}
