package exprtree

import (
	"encoding/binary"
	"fmt"
)

const (
	MaxAlignShift = 12
	MaxStructSize = 0x8000
)

// RuntimeOS
// {{{

type RuntimeOS uint8

const (
	InvalidRuntimeOS RuntimeOS = iota
	LINUX
)

var runtimeOSNames = []string{
	"InvalidRuntimeOS",
	"LINUX",
}

func (os RuntimeOS) String() string {
	if uint(os) >= uint(len(runtimeOSNames)) {
		return fmt.Sprintf("RuntimeOS(%d)", uint(os))
	}
	return runtimeOSNames[os]
}

func (os RuntimeOS) GoString() string {
	return os.String()
}

var _ fmt.Stringer = RuntimeOS(0)
var _ fmt.GoStringer = RuntimeOS(0)

// }}}

// RuntimeCPU
// {{{

type RuntimeCPU uint8

const (
	InvalidRuntimeCPU RuntimeCPU = iota
	X86_64
	X86
	X32
	ARM64
	ARM
)

var runtimeCPUNames = []string{
	"InvalidRuntimeCPU",
	"X86_64",
	"X86",
	"X32",
	"ARM64",
	"ARM",
}

var runtimeCPUDataModels = []DataModel{
	InvalidDataModel,
	LP64,
	ILP32,
	ILP32,
	LP64,
	ILP32,
}

func (cpu RuntimeCPU) String() string {
	if uint(cpu) >= uint(len(runtimeCPUNames)) {
		return fmt.Sprintf("RuntimeCPU(%d)", uint(cpu))
	}
	return runtimeCPUNames[cpu]
}

func (cpu RuntimeCPU) GoString() string {
	return cpu.String()
}

func (cpu RuntimeCPU) DataModel() DataModel {
	if uint(cpu) >= uint(len(runtimeCPUDataModels)) {
		return InvalidDataModel
	}
	return runtimeCPUDataModels[cpu]
}

func (cpu RuntimeCPU) ByteOrder() binary.ByteOrder {
	return binary.LittleEndian
}

var _ fmt.Stringer = RuntimeCPU(0)
var _ fmt.GoStringer = RuntimeCPU(0)

// }}}

// DataModel
// {{{

type DataModel uint8

const (
	InvalidDataModel DataModel = iota
	LP64
	ILP32
)

var dataModelNames = []string{
	"InvalidDataModel",
	"LP64",
	"ILP32",
}

func (model DataModel) String() string {
	if uint(model) >= uint(len(dataModelNames)) {
		return fmt.Sprintf("DataModel(%d)", uint(model))
	}
	return dataModelNames[model]
}

func (model DataModel) GoString() string {
	return model.String()
}

var _ fmt.Stringer = DataModel(0)
var _ fmt.GoStringer = DataModel(0)

// }}}

// SymbolKind
// {{{

type SymbolKind uint8

const (
	InvalidSymbolKind SymbolKind = iota
	SimpleSymbol
	UnboundGenericTypeSymbol
	BoundGenericTypeSymbol
	SimpleFunctionSymbol
	UnboundGenericFunctionSymbol
	BoundGenericFunctionSymbol
)

var symbolKindNames = []string{
	"InvalidSymbolKind",
	"SimpleSymbol",
	"UnboundGenericTypeSymbol",
	"BoundGenericTypeSymbol",
	"SimpleFunctionSymbol",
	"UnboundGenericFunctionSymbol",
	"BoundGenericFunctionSymbol",
}

func (kind SymbolKind) String() string {
	if uint(kind) >= uint(len(symbolKindNames)) {
		return fmt.Sprintf("SymbolKind(%d)", uint(kind))
	}
	return symbolKindNames[kind]
}

func (kind SymbolKind) GoString() string {
	return kind.String()
}

var _ fmt.Stringer = SymbolKind(0)
var _ fmt.GoStringer = SymbolKind(0)

// }}}

// TypeKind
// {{{

type TypeKind uint8

const (
	InvalidTypeKind TypeKind = iota

	ReflectedTypeKind

	U8Kind
	U16Kind
	U32Kind
	U64Kind

	S8Kind
	S16Kind
	S32Kind
	S64Kind

	F16Kind
	F32Kind
	F64Kind

	C32Kind
	C64Kind
	C128Kind

	StringKind
	ErrorKind

	EnumKind
	BitfieldKind
	StructKind
	UnionKind
	InterfaceKind
	FunctionKind

	PointerKind
	ArrayKind
	SliceKind

	MutableKind
	ConstKind
	NamedKind
)

var typeKindNames = []string{
	"InvalidTypeKind",
	"ReflectedTypeKind",
	"U8Kind",
	"U16Kind",
	"U32Kind",
	"U64Kind",
	"S8Kind",
	"S16Kind",
	"S32Kind",
	"S64Kind",
	"F16Kind",
	"F32Kind",
	"F64Kind",
	"C32Kind",
	"C64Kind",
	"C128Kind",
	"StringKind",
	"ErrorKind",
	"EnumKind",
	"BitfieldKind",
	"StructKind",
	"UnionKind",
	"InterfaceKind",
	"FunctionKind",
	"PointerKind",
	"ArrayKind",
	"SliceKind",
	"MutableKind",
	"ConstKind",
	"NamedKind",
}

func (kind TypeKind) String() string {
	if uint(kind) >= uint(len(typeKindNames)) {
		return fmt.Sprintf("TypeKind(%d)", uint(kind))
	}
	return typeKindNames[kind]
}

func (kind TypeKind) GoString() string {
	return kind.String()
}

var _ fmt.Stringer = TypeKind(0)
var _ fmt.GoStringer = TypeKind(0)

// }}}

// GenericParamKind
// {{{

type GenericParamKind uint32

const (
	InvalidGenericParamKind GenericParamKind = iota
	TypeGenericParam
	IntegerGenericParam
	EnumGenericParam
)

var typeParamKindNames = []string{
	"InvalidGenericParamKind",
	"TypeGenericParam",
	"IntegerGenericParam",
	"EnumGenericParam",
}

func (kind GenericParamKind) String() string {
	if uint(kind) >= uint(len(typeParamKindNames)) {
		return fmt.Sprintf("GenericParamKind(%d)", uint(kind))
	}
	return typeParamKindNames[kind]
}

func (kind GenericParamKind) GoString() string {
	return kind.String()
}

var _ fmt.Stringer = GenericParamKind(0)
var _ fmt.GoStringer = GenericParamKind(0)

// }}}

// StatementContext
// {{{

type StatementContext uint32

const (
	InvalidStatementContext StatementContext = iota
	EnumStatementContext
	BitfieldStatementContext
	StructStatementContext
	UnionStatementContext
	InterfaceStatementContext
	FunctionStatementContext
)

var statementContextNames = []string{
	"InvalidStatementContext",
	"EnumStatementContext",
	"BitfieldStatementContext",
	"StructStatementContext",
	"UnionStatementContext",
	"InterfaceStatementContext",
	"FunctionStatementContext",
}

func (kind StatementContext) String() string {
	if uint(kind) >= uint(len(statementContextNames)) {
		return fmt.Sprintf("StatementContext(%d)", uint(kind))
	}
	return statementContextNames[kind]
}

func (kind StatementContext) GoString() string {
	return kind.String()
}

var _ fmt.Stringer = StatementContext(0)
var _ fmt.GoStringer = StatementContext(0)

// }}}

// StatementKind
// {{{

type StatementKind uint32

const (
	InvalidStatementKind StatementKind = iota

	AlignPragmaStatement
	MinimumSizePragmaStatement
	PreserveFieldOrderPragmaStatement

	OmitNewPragmaStatement
	OmitCopyPragmaStatement
	OmitMovePragmaStatement
	OmitHashPragmaStatement
	OmitComparePragmaStatement
	OmitToStringPragmaStatement
	OmitToReprPragmaStatement

	StaticConstantStatement
	InstanceConstantStatement
	StaticFieldStatement

	EnumKindStatement
	EnumValueStatement
	EnumAliasStatement

	BitfieldKindStatement
	BitfieldValueStatement
	BitfieldAliasStatement

	StructFieldStatement

	UnionTagStatement
	UnionFieldStatement

	InterfaceFieldStatement
	InterfacePropertyStatement
	InterfaceMethodStatement
)

var statementKindNames = []string{
	"InvalidStatementKind",
	"AlignPragmaStatement",
	"MinimumSizePragmaStatement",
	"PreserveFieldOrderPragmaStatement",
	"OmitNewPragmaStatement",
	"OmitCopyPragmaStatement",
	"OmitMovePragmaStatement",
	"OmitHashPragmaStatement",
	"OmitComparePragmaStatement",
	"OmitToStringPragmaStatement",
	"OmitToReprPragmaStatement",
	"StaticConstantStatement",
	"InstanceConstantStatement",
	"StaticFieldStatement",
	"EnumKindStatement",
	"EnumValueStatement",
	"EnumAliasStatement",
	"BitfieldKindStatement",
	"BitfieldValueStatement",
	"BitfieldAliasStatement",
	"StructFieldStatement",
	"UnionTagStatement",
	"UnionFieldStatement",
	"InterfaceFieldStatement",
	"InterfacePropertyStatement",
	"InterfaceMethodStatement",
}

func (kind StatementKind) String() string {
	if uint(kind) >= uint(len(statementKindNames)) {
		return fmt.Sprintf("StatementKind(%d)", uint(kind))
	}
	return statementKindNames[kind]
}

func (kind StatementKind) GoString() string {
	return kind.String()
}

var _ fmt.Stringer = StatementKind(0)
var _ fmt.GoStringer = StatementKind(0)

// }}}

// TraversalOrder
// {{{

type TraversalOrder uint8

const (
	InvalidTraversalOrder TraversalOrder = iota
	PreOrder
	InOrder
	PostOrder
)

var traversalOrderNames = []string{
	"InvalidTraversalOrder",
	"PreOrder",
	"InOrder",
	"PostOrder",
}

func (order TraversalOrder) String() string {
	if uint(order) >= uint(len(traversalOrderNames)) {
		return fmt.Sprintf("TraversalOrder(%d)", uint(order))
	}
	return traversalOrderNames[order]
}

func (order TraversalOrder) GoString() string {
	return order.String()
}

var _ fmt.Stringer = TraversalOrder(0)
var _ fmt.GoStringer = TraversalOrder(0)

// }}}
