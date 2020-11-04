package exprtree

import (
	"encoding/binary"
	"fmt"
	"sync"
)

type Interp struct {
	mu sync.RWMutex
	cv *sync.Cond

	modulesByName map[string]*Module
	symbolsByID   map[SymbolID]*Symbol
	symbolsByName map[string]*Symbol
	typesByID     map[TypeID]*Type
	typesByName   map[string]*Type
	buffersByID   map[BufferID]*Buffer
	errorsByID    map[ErrorID]*Error
	genSigByID    map[GenericSignatureID]*GenericSignature
	genSigByName  map[string]*GenericSignature
	funcSigByID   map[FunctionSignatureID]*FunctionSignature
	funcSigByName map[string]*FunctionSignature

	builtinModuleSingleton         *Module
	builtinEnumModuleSingleton     *Module
	builtinBitfieldModuleSingleton *Module
	builtinStructModuleSingleton   *Module
	builtinUnionModuleSingleton    *Module

	typeTypeSingleton       *Type
	uInt8TypeSingleton      *Type
	uInt16TypeSingleton     *Type
	uInt32TypeSingleton     *Type
	uInt64TypeSingleton     *Type
	sInt8TypeSingleton      *Type
	sInt16TypeSingleton     *Type
	sInt32TypeSingleton     *Type
	sInt64TypeSingleton     *Type
	float16TypeSingleton    *Type
	float32TypeSingleton    *Type
	float64TypeSingleton    *Type
	complex32TypeSingleton  *Type
	complex64TypeSingleton  *Type
	complex128TypeSingleton *Type
	stringTypeSingleton     *Type
	errorTypeSingleton      *Type
	boolTypeSingleton       *Type
	orderTypeSingleton      *Type
	voidTypeSingleton       *Type
	nullTypeSingleton       *Type
	anyTypeSingleton        *Type

	pointerTypeCache  map[*Type]*Type
	mutableTypeCache  map[*Type]*Type
	constTypeCache    map[*Type]*Type
	enumTypeCache     map[string]*Type
	bitfieldTypeCache map[string]*Type
	structTypeCache   map[string]*Type
	unionTypeCache    map[string]*Type

	lastSymbolID  SymbolID
	lastTypeID    TypeID
	lastBufferID  BufferID
	lastErrorID   ErrorID
	lastGenSigID  GenericSignatureID
	lastFuncSigID FunctionSignatureID

	cpu RuntimeCPU
	os  RuntimeOS
}

func NewSystemInterp() *Interp {
	interp := new(Interp)
	interp.Init(SystemCPU(), SystemOS())
	return interp
}

func NewInterp(cpu RuntimeCPU, os RuntimeOS) *Interp {
	interp := new(Interp)
	interp.Init(cpu, os)
	return interp
}

func (interp *Interp) Init(cpu RuntimeCPU, os RuntimeOS) {
	if cpu == InvalidRuntimeCPU {
		panic(fmt.Errorf("cpu is %v", cpu))
	}

	if os == InvalidRuntimeOS {
		panic(fmt.Errorf("os is %v", os))
	}

	*interp = Interp{
		cv:                sync.NewCond(&interp.mu),
		modulesByName:     make(map[string]*Module, 256),
		symbolsByID:       make(map[SymbolID]*Symbol, 256),
		symbolsByName:     make(map[string]*Symbol, 256),
		typesByID:         make(map[TypeID]*Type, 256),
		typesByName:       make(map[string]*Type, 256),
		buffersByID:       make(map[BufferID]*Buffer, 256),
		errorsByID:        make(map[ErrorID]*Error, 256),
		genSigByID:        make(map[GenericSignatureID]*GenericSignature, 256),
		genSigByName:      make(map[string]*GenericSignature, 256),
		funcSigByID:       make(map[FunctionSignatureID]*FunctionSignature, 256),
		funcSigByName:     make(map[string]*FunctionSignature, 256),
		pointerTypeCache:  make(map[*Type]*Type, 256),
		mutableTypeCache:  make(map[*Type]*Type, 256),
		constTypeCache:    make(map[*Type]*Type, 256),
		enumTypeCache:     make(map[string]*Type, 256),
		bitfieldTypeCache: make(map[string]*Type, 256),
		structTypeCache:   make(map[string]*Type, 256),
		unionTypeCache:    make(map[string]*Type, 256),
		lastSymbolID:      0,
		lastTypeID:        0,
		lastBufferID:      0,
		lastErrorID:       0,
		lastGenSigID:      0,
		lastFuncSigID:     0,
		cpu:               cpu,
		os:                os,
	}

	interp.populateBuiltinTypes()
}

func (interp *Interp) CPU() RuntimeCPU {
	return interp.cpu
}

func (interp *Interp) OS() RuntimeOS {
	return interp.os
}

func (interp *Interp) DataModel() DataModel {
	return interp.cpu.DataModel()
}

func (interp *Interp) ByteOrder() binary.ByteOrder {
	return interp.cpu.ByteOrder()
}

func (interp *Interp) AllModules(out map[string]*Module) {
	checkNotNil("out", out)
	locked(interp.mu.RLocker(), func() {
		for cname, mod := range interp.modulesByName {
			out[cname] = mod
		}
	})
}

func (interp *Interp) ModuleByName(name string) (*Module, bool) {
	var mod *Module
	var found bool
	locked(interp.mu.RLocker(), func() {
		mod, found = interp.modulesByName[name]
	})
	return mod, found
}

func (interp *Interp) NewModule(cname string) (*Module, error) {
	if reReservedModuleNames.MatchString(cname) {
		return nil, fmt.Errorf("module name %q is reserved", cname)
	}

	return interp.newModuleCommon(cname, true)
}

func (interp *Interp) AllSymbols(out map[SymbolID]*Symbol) {
	checkNotNil("out", out)
	locked(interp.mu.RLocker(), func() {
		for id, sym := range interp.symbolsByID {
			out[id] = sym
		}
	})
}

func (interp *Interp) SymbolByID(id SymbolID) (*Symbol, bool) {
	var sym *Symbol
	var found bool
	locked(interp.mu.RLocker(), func() {
		sym, found = interp.symbolsByID[id]
	})
	return sym, found
}

func (interp *Interp) SymbolByMangledName(name string) (*Symbol, bool) {
	var sym *Symbol
	var found bool
	locked(interp.mu.RLocker(), func() {
		sym, found = interp.symbolsByName[name]
	})
	return sym, found
}

func (interp *Interp) AllTypes(out map[TypeID]*Type) {
	checkNotNil("out", out)
	locked(interp.mu.RLocker(), func() {
		for id, t := range interp.typesByID {
			out[id] = t
		}
	})
}

func (interp *Interp) TypeByID(id TypeID) (*Type, bool) {
	var t *Type
	var found bool
	locked(interp.mu.RLocker(), func() {
		t, found = interp.typesByID[id]
	})
	return t, found
}

func (interp *Interp) TypeByMangledName(mname string) (*Type, bool) {
	var t *Type
	var found bool
	locked(interp.mu.RLocker(), func() {
		t, found = interp.typesByName[mname]
	})
	return t, found
}

func (interp *Interp) AllBuffers(out map[BufferID]*Buffer) {
	checkNotNil("out", out)
	locked(interp.mu.RLocker(), func() {
		for id, buf := range interp.buffersByID {
			out[id] = buf
		}
	})
}

func (interp *Interp) BufferByID(id BufferID) (*Buffer, bool) {
	var buf *Buffer
	var found bool
	locked(interp.mu.RLocker(), func() {
		buf, found = interp.buffersByID[id]
	})
	return buf, found
}

func (interp *Interp) AllErrors(out map[ErrorID]*Error) {
	checkNotNil("out", out)
	locked(interp.mu.RLocker(), func() {
		for id, err := range interp.errorsByID {
			out[id] = err
		}
	})
}

func (interp *Interp) ErrorByID(id ErrorID) (*Error, bool) {
	var err *Error
	var found bool
	locked(interp.mu.RLocker(), func() {
		err, found = interp.errorsByID[id]
	})
	return err, found
}

func (interp *Interp) BuiltinModule() *Module         { return interp.builtinModuleSingleton }
func (interp *Interp) BuiltinEnumModule() *Module     { return interp.builtinEnumModuleSingleton }
func (interp *Interp) BuiltinBitfieldModule() *Module { return interp.builtinBitfieldModuleSingleton }
func (interp *Interp) BuiltinStructModule() *Module   { return interp.builtinStructModuleSingleton }
func (interp *Interp) BuiltinUnionModule() *Module    { return interp.builtinUnionModuleSingleton }

func (interp *Interp) TypeType() *Type       { return interp.typeTypeSingleton }
func (interp *Interp) UInt8Type() *Type      { return interp.uInt8TypeSingleton }
func (interp *Interp) UInt16Type() *Type     { return interp.uInt16TypeSingleton }
func (interp *Interp) UInt32Type() *Type     { return interp.uInt32TypeSingleton }
func (interp *Interp) UInt64Type() *Type     { return interp.uInt64TypeSingleton }
func (interp *Interp) SInt8Type() *Type      { return interp.sInt8TypeSingleton }
func (interp *Interp) SInt16Type() *Type     { return interp.sInt16TypeSingleton }
func (interp *Interp) SInt32Type() *Type     { return interp.sInt32TypeSingleton }
func (interp *Interp) SInt64Type() *Type     { return interp.sInt64TypeSingleton }
func (interp *Interp) Float16Type() *Type    { return interp.float16TypeSingleton }
func (interp *Interp) Float32Type() *Type    { return interp.float32TypeSingleton }
func (interp *Interp) Float64Type() *Type    { return interp.float64TypeSingleton }
func (interp *Interp) Complex32Type() *Type  { return interp.complex32TypeSingleton }
func (interp *Interp) Complex64Type() *Type  { return interp.complex64TypeSingleton }
func (interp *Interp) Complex128Type() *Type { return interp.complex128TypeSingleton }
func (interp *Interp) StringType() *Type     { return interp.stringTypeSingleton }
func (interp *Interp) ErrorType() *Type      { return interp.errorTypeSingleton }
func (interp *Interp) BoolType() *Type       { return interp.boolTypeSingleton }
func (interp *Interp) OrderType() *Type      { return interp.orderTypeSingleton }
func (interp *Interp) VoidType() *Type       { return interp.voidTypeSingleton }
func (interp *Interp) NullType() *Type       { return interp.nullTypeSingleton }

func (interp *Interp) SignedType(in *Type) (*Type, error) {
	checkNotNil("in", in)

	wrapMutable := false
	wrapConst := false

	switch in.kind {
	case MutableKind:
		wrapMutable = true
		in = in.details.(*Type)
	case ConstKind:
		wrapConst = true
		in = in.details.(*Type)
	}

	var out *Type
	switch in.kind {
	case U8Kind, S8Kind:
		out = interp.SInt8Type()
	case U16Kind, S16Kind:
		out = interp.SInt16Type()
	case U32Kind, S32Kind:
		out = interp.SInt32Type()
	case U64Kind, S64Kind:
		out = interp.SInt64Type()
	default:
		return nil, fmt.Errorf("illegal application of builtin::Signed[type] with type %s; only primitive integer types are permitted", in.CanonicalName())
	}

	if wrapMutable {
		out, _ = interp.MutableType(out)
	}
	if wrapConst {
		out, _ = interp.ConstType(out)
	}
	return out, nil
}

func (interp *Interp) UnsignedType(in *Type) (*Type, error) {
	checkNotNil("in", in)

	wrapMutable := false
	wrapConst := false

	switch in.kind {
	case MutableKind:
		wrapMutable = true
		in = in.details.(*Type)
	case ConstKind:
		wrapConst = true
		in = in.details.(*Type)
	}

	var out *Type
	switch in.kind {
	case U8Kind, S8Kind:
		out = interp.UInt8Type()
	case U16Kind, S16Kind:
		out = interp.UInt16Type()
	case U32Kind, S32Kind:
		out = interp.UInt32Type()
	case U64Kind, S64Kind:
		out = interp.UInt64Type()
	default:
		return nil, fmt.Errorf("illegal application of builtin::Unsigned#[type] with type %s; only primitive integer types are permitted", in.CanonicalName())
	}

	if wrapMutable {
		out, _ = interp.MutableType(out)
	}
	if wrapConst {
		out, _ = interp.ConstType(out)
	}
	return out, nil
}

func (interp *Interp) PointerType(in *Type) (*Type, error) {
	checkNotNil("in", in)

	var out *Type
	var found bool
	locked(&interp.mu, func() {
		out, found = interp.pointerTypeCache[in]
		if !found {
			interp.pointerTypeCache[in] = nil
			return
		}
		for out == nil {
			interp.cv.Wait()
			out = interp.pointerTypeCache[in]
		}
	})
	if found {
		return out, nil
	}

	g1 := interp.GenericSignatureBuilder().WithType().Build()

	var err error
	out, err = interp.createType(
		interp.BuiltinModule().Symbols(),
		SymbolData{
			Kind: BoundGenericTypeSymbol,
			Name: "Pointer",
			Generic: GenericSymbolData{
				Signature:   g1,
				ParamNames:  []string{"T"},
				ParamValues: []interface{}{in},
			},

			HasCanonicalNameOverride: true,
			CanonicalNameOverride:    "*" + in.CanonicalName(),

			HasMangledNameOverride: true,
			MangledNameOverride:    "_Ap" + in.MangledName()[2:],
		},
		func(t *Type) {
			t.kind = PointerKind
			t.alignShift = 3
			t.minSize = 8
			t.padSize = 8
			t.details = in

			// FIXME: provide auto-pointer-chase wrappers
		})

	if err != nil {
		return nil, err
	}

	locked(&interp.mu, func() {
		interp.pointerTypeCache[in] = out
		interp.cv.Broadcast()
	})

	return out, nil
}

func (interp *Interp) MutableType(in *Type) (*Type, error) {
	checkNotNil("in", in)

	if in.kind == MutableKind {
		return in, nil
	}

	if in.kind == ConstKind {
		in = in.details.(*Type)
	}

	var out *Type
	var found bool
	locked(&interp.mu, func() {
		out, found = interp.mutableTypeCache[in]
		if !found {
			interp.mutableTypeCache[in] = nil
			return
		}
		for out == nil {
			interp.cv.Wait()
			out = interp.mutableTypeCache[in]
		}
	})
	if found {
		return out, nil
	}

	g1 := interp.GenericSignatureBuilder().WithType().Build()

	var err error
	out, err = interp.createType(
		interp.BuiltinModule().Symbols(),
		SymbolData{
			Kind: BoundGenericTypeSymbol,
			Name: "Mutable",
			Generic: GenericSymbolData{
				Signature:   g1,
				ParamNames:  []string{"T"},
				ParamValues: []interface{}{in},
			},

			HasCanonicalNameOverride: true,
			CanonicalNameOverride:    "mutable " + in.CanonicalName(),

			HasMangledNameOverride: true,
			MangledNameOverride:    "_Am" + in.MangledName()[2:],
		},
		func(t *Type) {
			t.kind = MutableKind
			t.alignShift = in.alignShift
			t.minSize = in.minSize
			t.padSize = in.padSize
			t.details = in

			// FIXME: propagate symbols
		})

	if err != nil {
		return nil, err
	}

	locked(&interp.mu, func() {
		interp.mutableTypeCache[in] = out
		interp.cv.Broadcast()
	})

	return out, nil
}

func (interp *Interp) ConstType(in *Type) (*Type, error) {
	checkNotNil("in", in)

	if in.kind == MutableKind {
		return in, nil
	}

	if in.kind == ConstKind {
		return in, nil
	}

	var out *Type
	var found bool
	locked(&interp.mu, func() {
		out, found = interp.constTypeCache[in]
		if !found {
			interp.constTypeCache[in] = nil
			return
		}
		for out == nil {
			interp.cv.Wait()
			out = interp.constTypeCache[in]
		}
	})
	if found {
		return out, nil
	}

	g1 := interp.GenericSignatureBuilder().WithType().Build()

	var err error
	out, err = interp.createType(
		interp.BuiltinModule().Symbols(),
		SymbolData{
			Kind: BoundGenericTypeSymbol,
			Name: "Const",
			Generic: GenericSymbolData{
				Signature:   g1,
				ParamNames:  []string{"T"},
				ParamValues: []interface{}{in},
			},

			HasCanonicalNameOverride: true,
			CanonicalNameOverride:    "const " + in.CanonicalName(),

			HasMangledNameOverride: true,
			MangledNameOverride:    "_Ak" + in.MangledName()[2:],
		},
		func(t *Type) {
			t.kind = ConstKind
			t.alignShift = in.alignShift
			t.minSize = in.minSize
			t.padSize = in.padSize
			t.details = in

			// FIXME: filter out non-const methods
		})

	if err != nil {
		return nil, err
	}

	locked(&interp.mu, func() {
		interp.constTypeCache[in] = out
		interp.cv.Broadcast()
	})

	return out, nil
}

func (interp *Interp) NamedType(symtab *SymbolTable, data SymbolData, in *Type) (*Type, error) {
	checkNotNil("symtab", symtab)
	checkNotNil("in", in)
	return interp.createType(symtab, data, func(t *Type) {
		t.kind = NamedKind
		t.alignShift = in.alignShift
		t.minSize = in.minSize
		t.padSize = in.padSize
		t.details = in
	})
}

func (interp *Interp) populateBuiltinTypes() {
	interp.builtinModuleSingleton = interp.newModuleInternal("builtin", false)
	interp.builtinEnumModuleSingleton = interp.newModuleInternal("builtin::enum", false)
	interp.builtinBitfieldModuleSingleton = interp.newModuleInternal("builtin::bitfield", false)
	interp.builtinStructModuleSingleton = interp.newModuleInternal("builtin::struct", false)
	interp.builtinUnionModuleSingleton = interp.newModuleInternal("builtin::union", false)

	interp.typeTypeSingleton = func() *Type {
		t := new(Type)

		sym, err := interp.BuiltinModule().Symbols().NewSymbol(SymbolData{
			Kind: SimpleSymbol,
			Name: "Type",
			Type: t,

			HasMangledNameOverride: true,
			MangledNameOverride:    "_At",
		})
		checkBug(err)

		id := interp.allocateType()

		*t = Type{
			interp:     interp,
			sym:        sym,
			id:         id,
			kind:       ReflectedTypeKind,
			alignShift: 2,
			minSize:    4,
			padSize:    4,
		}

		cname := sym.CanonicalName()
		mname := sym.MangledName()
		xname := mname[:len(mname)-1]

		canonPrefix := cname + "."
		staticManglePrefix := xname + "S"
		localManglePrefix := xname + "L"

		t.static.Init(interp, &t.mu, t.mu.RLocker(), canonPrefix, staticManglePrefix)
		t.instance.Init(interp, &t.mu, t.mu.RLocker(), canonPrefix, localManglePrefix)

		// TODO: add methods

		interp.registerType(t)
		sym.SetCompileTimeValue(t)
		return t
	}()

	type primitiveTypeRow struct {
		Kind       TypeKind
		AlignShift uint
		Bias       int
		Format     bool
		Name       string
		Mangle     string
		Pointer    **Type
	}

	var primitiveTypeTable = []primitiveTypeRow{
		{U8Kind, 0, 0, true, "UInt%d", "_Au%d", &interp.uInt8TypeSingleton},
		{U16Kind, 1, 0, true, "UInt%d", "_Au%d", &interp.uInt16TypeSingleton},
		{U32Kind, 2, 0, true, "UInt%d", "_Au%d", &interp.uInt32TypeSingleton},
		{U64Kind, 3, 0, true, "UInt%d", "_Au%d", &interp.uInt64TypeSingleton},
		{S8Kind, 0, 0, true, "SInt%d", "_Ai%d", &interp.sInt8TypeSingleton},
		{S16Kind, 1, 0, true, "SInt%d", "_Ai%d", &interp.sInt16TypeSingleton},
		{S32Kind, 2, 0, true, "SInt%d", "_Ai%d", &interp.sInt32TypeSingleton},
		{S64Kind, 3, 0, true, "SInt%d", "_Ai%d", &interp.sInt64TypeSingleton},
		{F16Kind, 1, 0, true, "Float%d", "_Af%d", &interp.float16TypeSingleton},
		{F32Kind, 2, 0, true, "Float%d", "_Af%d", &interp.float32TypeSingleton},
		{F64Kind, 3, 0, true, "Float%d", "_Af%d", &interp.float64TypeSingleton},
		{C32Kind, 2, -1, true, "Complex%d", "_Ac%d", &interp.complex32TypeSingleton},
		{C64Kind, 3, -1, true, "Complex%d", "_Ac%d", &interp.complex64TypeSingleton},
		{C128Kind, 4, -1, true, "Complex%d", "_Ac%d", &interp.complex128TypeSingleton},
		{StringKind, 3, 0, false, "String", "_As", &interp.stringTypeSingleton},
		{ErrorKind, 3, 0, false, "Error", "_Ae", &interp.errorTypeSingleton},
	}

	for _, row := range primitiveTypeTable {
		sizeBytes := uint(1) << row.AlignShift
		sizeBits := 8 * sizeBytes

		oname := row.Name
		mname := row.Mangle
		if row.Format {
			oname = fmt.Sprintf(oname, sizeBits)
			mname = fmt.Sprintf(mname, int(row.AlignShift)+row.Bias)
		}

		out, err := interp.createType(
			interp.BuiltinModule().Symbols(),
			SymbolData{
				Kind: SimpleSymbol,
				Name: oname,

				HasMangledNameOverride: true,
				MangledNameOverride:    mname,
			},
			func(t *Type) {
				t.kind = row.Kind
				t.alignShift = uint8(row.AlignShift)
				t.minSize = uint16(sizeBytes)
				t.padSize = uint16(sizeBytes)
			})

		checkBug(err)
		*row.Pointer = out
	}

	type enumDataRow struct {
		Name       string
		HasMangled bool
		Mangled    string
		List       Statements
		Pointer    **Type
	}

	enumDataTable := []enumDataRow{
		{
			Name:       "Bool",
			HasMangled: true,
			Mangled:    "_Ab",
			List: Statements{
				{Kind: EnumKindStatement, EnumKind: S8Kind},
				{Kind: EnumValueStatement, EnumName: "true", EnumNumber: -1},
				{Kind: EnumAliasStatement, EnumName: "True", EnumAliasOf: "true"},
				{Kind: EnumAliasStatement, EnumName: "TRUE", EnumAliasOf: "true"},
				{Kind: EnumValueStatement, EnumName: "false", EnumNumber: 0},
				{Kind: EnumAliasStatement, EnumName: "False", EnumAliasOf: "false"},
				{Kind: EnumAliasStatement, EnumName: "FALSE", EnumAliasOf: "false"},
			},
			Pointer: &interp.boolTypeSingleton,
		},
		{
			Name:       "Order",
			HasMangled: true,
			Mangled:    "_Ao",
			List: Statements{
				{Kind: EnumKindStatement, EnumKind: S8Kind},
				{Kind: EnumValueStatement, EnumName: "LT", EnumNumber: -1},
				{Kind: EnumValueStatement, EnumName: "EQ", EnumNumber: 0},
				{Kind: EnumValueStatement, EnumName: "GT", EnumNumber: 1},
			},
			Pointer: &interp.orderTypeSingleton,
		},
	}

	for _, row := range enumDataTable {
		in, err := interp.EnumType(row.List)
		checkBug(err)

		out, err := interp.NamedType(
			interp.BuiltinModule().Symbols(),
			SymbolData{
				Kind: SimpleSymbol,
				Name: row.Name,

				HasMangledNameOverride: row.HasMangled,
				MangledNameOverride:    row.Mangled,
			},
			in)

		checkBug(err)
		*row.Pointer = out
	}

	type bitfieldDataRow struct {
		Name       string
		HasMangled bool
		Mangled    string
		List       Statements
		Pointer    **Type
	}

	bitfieldDataTable := []bitfieldDataRow{}

	for _, row := range bitfieldDataTable {
		in, err := interp.BitfieldType(row.List)
		checkBug(err)

		out, err := interp.NamedType(
			interp.BuiltinModule().Symbols(),
			SymbolData{
				Kind: SimpleSymbol,
				Name: row.Name,

				HasMangledNameOverride: row.HasMangled,
				MangledNameOverride:    row.Mangled,
			},
			in)

		checkBug(err)
		*row.Pointer = out
	}

	interp.voidTypeSingleton = func() *Type {
		out, _ := interp.createType(
			interp.BuiltinModule().Symbols(),
			SymbolData{
				Kind:                   SimpleSymbol,
				Name:                   "Void",
				HasMangledNameOverride: true,
				MangledNameOverride:    "_Av",
			},
			func(t *Type) {
				t.kind = StructKind
				t.alignShift = 0
				t.minSize = 0
				t.padSize = 1
				t.details = &Struct{
					list:       Statements{},
					fields:     []*StructField{},
					alignShift: 0,
					minSize:    0,
				}
			})

		locked(&interp.mu, func() {
			interp.structTypeCache[""] = out
		})
		return out
	}()

	interp.nullTypeSingleton = func() *Type {
		in := interp.voidTypeSingleton

		out, err := interp.NamedType(
			interp.BuiltinModule().Symbols(),
			SymbolData{
				Kind: SimpleSymbol,
				Name: "Null",
			},
			in)

		checkBug(err)
		return out
	}()

	type structDataRow struct {
		Name       string
		HasMangled bool
		Mangled    string
		List       Statements
		Pointer    **Type
	}

	structDataTable := []structDataRow{}

	for _, row := range structDataTable {
		in, err := interp.StructType(row.List)
		checkBug(err)

		out, err := interp.NamedType(
			interp.BuiltinModule().Symbols(),
			SymbolData{
				Kind: SimpleSymbol,
				Name: row.Name,

				HasMangledNameOverride: row.HasMangled,
				MangledNameOverride:    row.Mangled,
			},
			in)

		checkBug(err)
		*row.Pointer = out
	}

	type unionDataRow struct {
		Name       string
		HasMangled bool
		Mangled    string
		List       Statements
		Pointer    **Type
	}

	unionDataTable := []unionDataRow{}

	for _, row := range unionDataTable {
		in, err := interp.UnionType(row.List)
		checkBug(err)

		out, err := interp.NamedType(
			interp.BuiltinModule().Symbols(),
			SymbolData{
				Kind: SimpleSymbol,
				Name: row.Name,

				HasMangledNameOverride: row.HasMangled,
				MangledNameOverride:    row.Mangled,
			},
			in)

		checkBug(err)
		*row.Pointer = out
	}

	gsb := interp.GenericSignatureBuilder()
	gsb.Build()
	gsb.WithType().Build()

	fsb := interp.FunctionSignatureBuilder()
	fsb.Build()

	// warm up the cache
	type F func(*Type) (*Type, error)
	dataF := []F{
		interp.PointerType,
		interp.ConstType,
		interp.MutableType,
	}

	type G func() *Type
	dataG := []G{
		interp.TypeType,
		interp.UInt8Type,
		interp.UInt16Type,
		interp.UInt32Type,
		interp.UInt64Type,
		interp.SInt8Type,
		interp.SInt16Type,
		interp.SInt32Type,
		interp.SInt64Type,
		interp.Float16Type,
		interp.Float32Type,
		interp.Float64Type,
		interp.Complex32Type,
		interp.Complex64Type,
		interp.Complex128Type,
		interp.StringType,
		interp.ErrorType,
		interp.BoolType,
		interp.OrderType,
	}

	for _, g := range dataG {
		t := g()
		for _, f := range dataF {
			_, err := f(t)
			checkBug(err)
		}
	}
}

func (interp *Interp) newModuleInternal(cname string, withImports bool) *Module {
	mod, err := interp.newModuleCommon(cname, withImports)
	checkBug(err)
	return mod
}

func (interp *Interp) newModuleCommon(cname string, withImports bool) (*Module, error) {
	mname := MangleModuleName(cname)
	xname := mname[:len(mname)-1]

	mod := &Module{
		interp: interp,
		cname:  cname,
		mname:  mname,
	}

	if withImports {
		mod.imports = make(map[string]*Module, 16)
		mod.imports["this"] = mod
		mod.imports["builtin"] = interp.BuiltinModule()
		mod.imports["builtin::enum"] = interp.BuiltinEnumModule()
		mod.imports["builtin::bitfield"] = interp.BuiltinBitfieldModule()
		mod.imports["builtin::struct"] = interp.BuiltinStructModule()
		mod.imports["builtin::union"] = interp.BuiltinUnionModule()
	}

	canonPrefix := cname + "::"
	manglePrefix := xname + "G"
	mod.symbols.Init(interp, &mod.mu, mod.mu.RLocker(), canonPrefix, manglePrefix)

	var old *Module
	locked(&interp.mu, func() {
		old = interp.modulesByName[cname]
		if old == nil {
			interp.modulesByName[cname] = mod
		}
	})

	if old != nil {
		return nil, fmt.Errorf("module name %q already exists", cname)
	}
	return mod, nil
}

func (interp *Interp) createType(symtab *SymbolTable, data SymbolData, fn func(*Type)) (*Type, error) {
	id := interp.allocateType()

	data.Type = interp.TypeType()
	sym, err := symtab.NewSymbol(data)
	if err != nil {
		return nil, err
	}

	t := &Type{
		interp: interp,
		sym:    sym,
		id:     id,
	}

	cname := sym.CanonicalName()
	mname := sym.MangledName()
	xname := mname[:len(mname)-1]

	canonPrefix := cname + "."
	staticManglePrefix := xname + "S"
	localManglePrefix := xname + "L"

	t.static.Init(interp, &t.mu, t.mu.RLocker(), canonPrefix, staticManglePrefix)
	t.instance.Init(interp, &t.mu, t.mu.RLocker(), canonPrefix, localManglePrefix)

	fn(t)

	interp.registerType(t)
	sym.SetCompileTimeValue(t)
	return t, nil
}

func (interp *Interp) allocateSymbol() SymbolID {
	var id SymbolID
	locked(&interp.mu, func() {
		interp.lastSymbolID++
		id = interp.lastSymbolID
	})
	return id
}

func (interp *Interp) allocateType() TypeID {
	var id TypeID
	locked(&interp.mu, func() {
		interp.lastTypeID++
		id = interp.lastTypeID
	})
	return id
}

func (interp *Interp) allocateBuffer() BufferID {
	var id BufferID
	locked(&interp.mu, func() {
		interp.lastBufferID++
		id = interp.lastBufferID
	})
	return id
}

func (interp *Interp) allocateError() ErrorID {
	var id ErrorID
	locked(&interp.mu, func() {
		interp.lastErrorID++
		id = interp.lastErrorID
	})
	return id
}

func (interp *Interp) registerSymbol(ptr *Symbol) {
	locked(&interp.mu, func() {
		interp.symbolsByID[ptr.ID()] = ptr
		interp.symbolsByName[ptr.MangledName()] = ptr
	})
}

func (interp *Interp) registerType(ptr *Type) {
	locked(&interp.mu, func() {
		interp.typesByID[ptr.ID()] = ptr
		interp.typesByName[ptr.MangledName()] = ptr
	})
}

func (interp *Interp) registerBuffer(ptr *Buffer) {
	locked(&interp.mu, func() {
		interp.buffersByID[ptr.ID()] = ptr
	})
}

func (interp *Interp) registerError(ptr *Error) {
	locked(&interp.mu, func() {
		interp.errorsByID[ptr.ID()] = ptr
	})
}

func (interp *Interp) registerGenSig(ptr *GenericSignature) *GenericSignature {
	name := ptr.String()
	locked(&interp.mu, func() {
		if existing, found := interp.genSigByName[name]; found {
			ptr = existing
			return
		}

		interp.lastGenSigID++
		id := interp.lastGenSigID

		ptr.id = id
		interp.genSigByID[id] = ptr
		interp.genSigByName[name] = ptr
	})
	return ptr
}

func (interp *Interp) registerFuncSig(ptr *FunctionSignature) *FunctionSignature {
	name := ptr.String()
	locked(&interp.mu, func() {
		if existing, found := interp.funcSigByName[name]; found {
			ptr = existing
			return
		}

		interp.lastFuncSigID++
		id := interp.lastFuncSigID

		ptr.id = id
		interp.funcSigByID[id] = ptr
		interp.funcSigByName[name] = ptr
	})
	return ptr
}
