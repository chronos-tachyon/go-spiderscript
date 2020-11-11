package exprtree

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// SymbolData
// {{{

type SymbolData struct {
	Kind     SymbolKind
	Name     string
	Type     *Type
	Generic  GenericSymbolData
	Function FunctionSymbolData

	HasCanonicalNameOverride bool
	CanonicalNameOverride    string

	HasMangledNameOverride bool
	MangledNameOverride    string

	ignored struct{}
}

type GenericSymbolData struct {
	Signature   *GenericSignature
	ParamNames  []string
	ParamValues []interface{}

	ignored struct{}
}

type FunctionSymbolData struct {
	Signature       *FunctionSignature
	PositionalNames []string

	ignored struct{}
}

// }}}

// SymbolName
// {{{

type SymbolName struct {
	generic  *GenericSymbolName
	function *FunctionSymbolName
	hname    string
	lname    string
	cname    string
	mname    string
	kind     SymbolKind
}

func (sn *SymbolName) Kind() SymbolKind {
	return sn.kind
}

func (sn *SymbolName) HumanName() string {
	return sn.hname
}

func (sn *SymbolName) LocalName() string {
	return sn.lname
}

func (sn *SymbolName) CanonicalName() string {
	return sn.cname
}

func (sn *SymbolName) MangledName() string {
	return sn.mname
}

func (sn *SymbolName) Generic() *GenericSymbolName {
	return sn.generic
}

func (sn *SymbolName) Function() *FunctionSymbolName {
	return sn.function
}

// GenericSymbolName
// {{{

type GenericSymbolName struct {
	sig    *GenericSignature
	names  []string
	values []interface{}
}

func (gsn *GenericSymbolName) Signature() *GenericSignature {
	return gsn.sig
}

func (gsn *GenericSymbolName) ParamNames() []string {
	return cloneStrings(gsn.names)
}

func (gsn *GenericSymbolName) ParamValues() []interface{} {
	return cloneAnys(gsn.values)
}

// }}}

// FunctionSymbolName
// {{{

type FunctionSymbolName struct {
	sig        *FunctionSignature
	posNames   []string
	namedNames []string
}

func (fsn *FunctionSymbolName) Signature() *FunctionSignature {
	return fsn.sig
}

func (fsn *FunctionSymbolName) PositionalArgumentNames() []string {
	return cloneStrings(fsn.posNames)
}

func (fsn *FunctionSymbolName) NamedArgumentNames() []string {
	return cloneStrings(fsn.namedNames)
}

// }}}

// }}}

// Symbol
// {{{

type Symbol struct {
	symtab   *SymbolTable
	generic  *GenericSymbolName
	function *FunctionSymbolName
	type_    *Type
	hname    string
	lname    string
	cname    string
	mname    string
	ctv      interface{}
	id       SymbolID
	kind     SymbolKind
}

func (sym *Symbol) Interp() *Interp {
	return sym.SymbolTable().Interp()
}

func (sym *Symbol) SymbolTable() *SymbolTable {
	return sym.symtab
}

func (sym *Symbol) ID() SymbolID {
	return sym.id
}

func (sym *Symbol) Kind() SymbolKind {
	return sym.kind
}

func (sym *Symbol) HumanName() string {
	return sym.hname
}

func (sym *Symbol) LocalName() string {
	return sym.lname
}

func (sym *Symbol) CanonicalName() string {
	return sym.cname
}

func (sym *Symbol) MangledName() string {
	return sym.mname
}

func (sym *Symbol) Generic() *GenericSymbolName {
	return sym.generic
}

func (sym *Symbol) Function() *FunctionSymbolName {
	return sym.function
}

func (sym *Symbol) Type() *Type {
	return sym.type_
}

func (sym *Symbol) CompileTimeValue() interface{} {
	var ctv interface{}
	locked(sym.symtab.rmu, func() {
		ctv = sym.ctv
	})
	return ctv
}

func (sym *Symbol) SetCompileTimeValue(ctv interface{}) {
	locked(sym.symtab.wmu, func() {
		sym.ctv = ctv
	})
}

func (sym *Symbol) String() string {
	return sym.CanonicalName()
}

func (sym *Symbol) GoString() string {
	kname := sym.Kind().String()
	mname := sym.MangledName()
	cname := sym.CanonicalName()
	tname := sym.type_.CanonicalName()
	estimatedLen := 18 + uint(len(kname)) + uint(len(mname)) + uint(len(cname)) + uint(len(tname))

	buf := takeBuffer()
	buf.Grow(int(estimatedLen))
	defer giveBuffer(buf)

	buf.WriteString("Symbol(")
	buf.WriteString(kname)
	buf.WriteString(", ")
	buf.WriteString(strconv.Quote(mname))
	buf.WriteString(", ")
	buf.WriteString(strconv.Quote(cname))
	buf.WriteString(", ")
	buf.WriteString(tname)
	buf.WriteByte(')')

	return checkEstimatedLength(buf, estimatedLen)
}

var _ fmt.Stringer = (*Symbol)(nil)
var _ fmt.GoStringer = (*Symbol)(nil)

// }}}

// Interface: Resolver
// {{{

type Resolver interface {
	All(out map[string]*Symbol)
	Get(name string) (*Symbol, bool)
	Put(mangledName string, sym *Symbol) error
}

// ResolverList
// {{{

type ResolverList []Resolver

func (list ResolverList) All(out map[string]*Symbol) {
	if out == nil {
		panic(fmt.Errorf("out is nil"))
	}

	index := uint(len(list))
	for index > 0 {
		index--
		list[index].All(out)
	}
}

func (list ResolverList) Get(name string) (*Symbol, bool) {
	for _, r := range list {
		if sym, found := r.Get(name); found {
			return sym, found
		}
	}
	return nil, false
}

func (list ResolverList) Put(name string, sym *Symbol) error {
	if len(list) == 0 {
		return fmt.Errorf("ResolverList has no scopes")
	}
	first := list[0]
	return first.Put(name, sym)
}

var _ Resolver = ResolverList(nil)

// }}}

// SymbolTable
// {{{

type SymbolTable struct {
	interp       *Interp
	wmu          sync.Locker
	rmu          sync.Locker
	canonPrefix  string
	manglePrefix string
	tbl          map[string]*Symbol
}

func (symtab *SymbolTable) Init(interp *Interp, wmu sync.Locker, rmu sync.Locker, canonPrefix string, manglePrefix string) {
	if interp == nil {
		panic(fmt.Errorf("BUG: *Interp is nil"))
	}
	if wmu == nil {
		panic(fmt.Errorf("BUG: wmu is nil"))
	}
	if rmu == nil {
		panic(fmt.Errorf("BUG: rmu is nil"))
	}

	*symtab = SymbolTable{
		interp:       interp,
		wmu:          wmu,
		rmu:          rmu,
		canonPrefix:  canonPrefix,
		manglePrefix: manglePrefix,
		tbl:          make(map[string]*Symbol, 16),
	}
}

func (symtab *SymbolTable) Interp() *Interp {
	return symtab.interp
}

func (symtab *SymbolTable) All(out map[string]*Symbol) {
	if out == nil {
		panic(fmt.Errorf("BUG: out is nil"))
	}

	locked(symtab.rmu, func() {
		for name, sym := range symtab.tbl {
			out[name] = sym
		}
	})
}

func (symtab *SymbolTable) Get(name string) (*Symbol, bool) {
	var sym *Symbol
	var found bool
	locked(symtab.rmu, func() {
		sym, found = symtab.tbl[name]
	})

	if found {
		return sym, found
	}

	sym, found = symtab.resolveSynthetic(name)
	if found {
		locked(symtab.wmu, func() {
			sym2, found2 := symtab.tbl[name]
			if found2 {
				sym = sym2
			} else {
				symtab.tbl[name] = sym
			}
		})
	}
	return sym, found
}

func (symtab *SymbolTable) Put(name string, sym *Symbol) error {
	if sym == nil {
		panic(fmt.Errorf("BUG: *Symbol is nil"))
	}

	var old *Symbol
	locked(symtab.wmu, func() {
		old = symtab.tbl[name]
		if old == nil {
			symtab.tbl[name] = sym
		}
	})

	if old != nil {
		return &DuplicateSymbolError{Name: name, Old: old, New: sym}
	}
	return nil
}

func (symtab *SymbolTable) NewSymbol(data SymbolData) (*Symbol, error) {
	sn, err := NewSymbolName(data, symtab.canonPrefix, symtab.manglePrefix)
	if err != nil {
		return nil, err
	}

	if data.Type == nil {
		return nil, fmt.Errorf("SymbolData.Type is nil")
	}

	sym := &Symbol{
		symtab:   symtab,
		id:       symtab.interp.allocateSymbol(),
		kind:     data.Kind,
		hname:    sn.HumanName(),
		lname:    sn.LocalName(),
		cname:    sn.CanonicalName(),
		mname:    sn.MangledName(),
		generic:  sn.Generic(),
		function: sn.Function(),
		type_:    data.Type,
		ctv:      nil,
	}
	symtab.interp.registerSymbol(sym)

	err = symtab.Put(sym.LocalName(), sym)
	if err != nil {
		return nil, err
	}

	return sym, nil
}

func (symtab *SymbolTable) NewGenSym(t *Type) *Symbol {
	checkNotNil("t", t)

	symID := symtab.interp.allocateSymbol()
	symName := fmt.Sprintf("__G%08x", uint32(symID))
	sym := &Symbol{
		symtab:   symtab,
		id:       symID,
		kind:     SimpleSymbol,
		hname:    symName,
		lname:    symName,
		cname:    symName,
		mname:    symName,
		generic:  nil,
		function: nil,
		type_:    t,
		ctv:      nil,
	}
	symtab.interp.registerSymbol(sym)

	if err := symtab.Put(sym.LocalName(), sym); err != nil {
		panic(fmt.Errorf("BUG: %w", err))
	}

	return sym
}

var _ Resolver = (*SymbolTable)(nil)

// }}}

// }}}

func NewSymbolName(data SymbolData, canonPrefix string, manglePrefix string) (*SymbolName, error) {
	expectGeneric := false
	expectGenericValues := false
	expectFunction := false

	switch data.Kind {
	case SimpleSymbol:
		// pass

	case UnboundGenericTypeSymbol:
		expectGeneric = true

	case BoundGenericTypeSymbol:
		expectGeneric = true
		expectGenericValues = true

	case SimpleFunctionSymbol:
		expectFunction = true

	case UnboundGenericFunctionSymbol:
		expectGeneric = true
		expectFunction = true

	case BoundGenericFunctionSymbol:
		expectGeneric = true
		expectGenericValues = true
		expectFunction = true

	default:
		return nil, fmt.Errorf("SymbolData.Kind is %v, which is not implemented", data.Kind)
	}

	if !reSymbolName.MatchString(data.Name) {
		return nil, fmt.Errorf("SymbolData.Name is %q, which is not a valid symbol name", data.Name)
	}

	sn := &SymbolName{
		kind:  data.Kind,
		hname: data.Name,
		lname: "???",
		cname: "???",
		mname: "???",
	}

	if expectGeneric {
		gsn := &GenericSymbolName{
			sig:    data.Generic.Signature,
			names:  cloneStrings(data.Generic.ParamNames),
			values: cloneAnys(data.Generic.ParamValues),
		}
		length := gsn.sig.NumParams()
		sn.generic = gsn

		if x := uint(len(gsn.names)); x != length {
			return nil, fmt.Errorf("SymbolData.Generic.ParamNames has length %d, but expected length %d", x, length)
		}

		if expectGenericValues {
			if x := uint(len(gsn.values)); x != length {
				return nil, fmt.Errorf("SymbolData.Generic.ParamValues has length %d, but expected length %d", x, length)
			}
		} else {
			if x := uint(len(gsn.values)); x != 0 {
				return nil, fmt.Errorf("SymbolData.Kind is %v, but SymbolData.Generic.ParamValues is non-nil", data.Kind)
			}
		}

		for index := uint(0); index < length; index++ {
			param := gsn.sig.Param(index)
			paramName := gsn.names[index]
			var paramValue interface{}
			if expectGenericValues {
				paramValue = gsn.values[index]
			}

			if !reSymbolName.MatchString(paramName) {
				return nil, fmt.Errorf("SymbolData.Generic.ParamNames[%d]: invalid symbol name %q", index, paramName)
			}

			switch param.Kind() {
			case TypeGenericParam:
				if expectGenericValues {
					ptr, ok := paramValue.(*Type)
					if !ok {
						return nil, fmt.Errorf("SymbolData.Generic.ParamValues[%d]: expected *Type, got %T", index, paramValue)
					}
					if ptr == nil {
						return nil, fmt.Errorf("SymbolData.Generic.ParamValues[%d].(*Type): pointer is nil", index)
					}
				}

			case IntegerGenericParam:
				if expectGenericValues {
					_, ok := paramValue.(uint64)
					if !ok {
						return nil, fmt.Errorf("SymbolData.Generic.ParamValues[%d]: expected uint64, got %T", index, paramValue)
					}
				}

			case EnumGenericParam:
				if expectGenericValues {
					outerType := param.Type()
					innerType := outerType.Chase()
					e := innerType.Details().(*Enum)

					ptr, ok := paramValue.(*EnumItem)
					if !ok {
						return nil, fmt.Errorf("SymbolData.Generic.ParamValues[%d]: expected *EnumItem, got %T", index, paramValue)
					}
					if ptr == nil {
						return nil, fmt.Errorf("SymbolData.Generic.ParamValues[%d].(*EnumItem): pointer is nil", index)
					}
					if ptr.parent != e {
						// TODO: better error message
						return nil, fmt.Errorf("SymbolData.Generic.ParamValues[%d].(*EnumItem): wrong enum kind", index)
					}
				}
			}
		}
	} else {
		if data.Generic.Signature != nil {
			return nil, fmt.Errorf("SymbolData.Kind is %v, but SymbolData.Generic.Signature is non-nil", data.Kind)
		}
		if data.Generic.ParamNames != nil {
			return nil, fmt.Errorf("SymbolData.Kind is %v, but SymbolData.Generic.ParamNames is non-nil", data.Kind)
		}
		if data.Generic.ParamValues != nil {
			return nil, fmt.Errorf("SymbolData.Kind is %v, but SymbolData.Generic.ParamValues is non-nil", data.Kind)
		}
	}

	if expectFunction {
		fsn := &FunctionSymbolName{
			sig:        data.Function.Signature,
			posNames:   cloneStrings(data.Function.PositionalNames),
			namedNames: data.Function.Signature.ArgNames(),
		}

		posLength := fsn.sig.NumPositionalArgs()
		sn.function = fsn

		if x := uint(len(fsn.posNames)); x != posLength {
			return nil, fmt.Errorf("SymbolData.Function: Signature.Positional has length %d, but PositionalNames has length %d", posLength, x)
		}

		for index := uint(0); index < posLength; index++ {
			argName := fsn.posNames[index]
			if !reSymbolName.MatchString(argName) {
				return nil, fmt.Errorf("SymbolData.Function.PositionalNames[%d]: invalid symbol name %q", index, argName)
			}
		}
	} else {
		if data.Function.Signature != nil {
			return nil, fmt.Errorf("SymbolData.Kind is %v, but SymbolData.Function.Signature is non-nil", data.Kind)
		}
		if data.Function.PositionalNames != nil {
			return nil, fmt.Errorf("SymbolData.Kind is %v, but SymbolData.Function.PositionalNames is non-nil", data.Kind)
		}
	}

	lBuf := takeBuffer()
	cBuf := takeBuffer()
	mBuf := takeBuffer()
	defer func() {
		giveBuffer(mBuf)
		giveBuffer(cBuf)
		giveBuffer(lBuf)
	}()

	cBuf.WriteString(canonPrefix)
	mBuf.WriteString(manglePrefix)

	switch sn.kind {
	case SimpleSymbol:
		cBuf.WriteString(sn.hname)

		writeNameTo(mBuf, sn.hname)

		lBuf.WriteString(sn.hname)

	case UnboundGenericTypeSymbol:
		cBuf.WriteString(sn.hname)
		writeUnboundGenericCanonicalNameTo(cBuf, sn.generic)

		writeNameTo(mBuf, sn.hname)
		writeUnboundGenericMangledNameTo(mBuf, sn.generic)

		lBuf.WriteString("__")
		writeNameTo(lBuf, sn.hname)
		writeUnboundGenericMangledNameTo(lBuf, sn.generic)

	case BoundGenericTypeSymbol:
		cBuf.WriteString(sn.hname)
		writeBoundGenericCanonicalNameTo(cBuf, sn.generic)

		writeNameTo(mBuf, sn.hname)
		writeBoundGenericMangledNameTo(mBuf, sn.generic)

		lBuf.WriteString("__")
		writeNameTo(lBuf, sn.hname)
		writeBoundGenericMangledNameTo(lBuf, sn.generic)

	case SimpleFunctionSymbol:
		cBuf.WriteString(sn.hname)
		writeFunctionCanonicalNameTo(cBuf, sn.function)

		writeNameTo(mBuf, sn.hname)
		writeFunctionMangledNameTo(mBuf, sn.function)

		lBuf.WriteString("__")
		writeNameTo(lBuf, sn.hname)
		writeFunctionMangledNameTo(lBuf, sn.function)

	case UnboundGenericFunctionSymbol:
		cBuf.WriteString(sn.hname)
		writeUnboundGenericCanonicalNameTo(cBuf, sn.generic)
		writeFunctionCanonicalNameTo(cBuf, sn.function)

		writeNameTo(mBuf, sn.hname)
		writeUnboundGenericMangledNameTo(mBuf, sn.generic)
		writeFunctionMangledNameTo(mBuf, sn.function)

		lBuf.WriteString("__")
		writeNameTo(lBuf, sn.hname)
		writeUnboundGenericMangledNameTo(lBuf, sn.generic)
		writeFunctionMangledNameTo(lBuf, sn.function)

	case BoundGenericFunctionSymbol:
		cBuf.WriteString(sn.hname)
		writeBoundGenericCanonicalNameTo(cBuf, sn.generic)
		writeFunctionCanonicalNameTo(cBuf, sn.function)

		writeNameTo(mBuf, sn.hname)
		writeBoundGenericMangledNameTo(mBuf, sn.generic)
		writeFunctionMangledNameTo(mBuf, sn.function)

		lBuf.WriteString("__")
		writeNameTo(lBuf, sn.hname)
		writeBoundGenericMangledNameTo(lBuf, sn.generic)
		writeFunctionMangledNameTo(lBuf, sn.function)
	}

	mBuf.WriteByte('Z')

	sn.lname = lBuf.String()
	sn.cname = cBuf.String()
	sn.mname = mBuf.String()

	if data.HasCanonicalNameOverride {
		sn.cname = data.CanonicalNameOverride
	}

	if data.HasMangledNameOverride {
		sn.mname = data.MangledNameOverride
	}

	return sn, nil
}

func (symtab *SymbolTable) resolveSynthetic(name string) (*Symbol, bool) {
	// TODO:
	// 1. recognize name =~ bound generic
	// 2. check if the unbound generic version of the name exists in the SymbolTable
	// 3. instantiate the generic by binding all provided compile-time parameters
	return nil, false
}

func writeUnboundGenericCanonicalNameTo(buf *strings.Builder, gsn *GenericSymbolName) {
	buf.WriteString("[")
	length := gsn.sig.NumParams()
	for index := uint(0); index < length; index++ {
		param := gsn.sig.Param(index)
		paramName := gsn.names[index]
		if index != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(paramName)
		buf.WriteString(": ")
		buf.WriteString(param.String())
	}
	buf.WriteString("]")
}

func writeBoundGenericCanonicalNameTo(buf *strings.Builder, gsn *GenericSymbolName) {
	buf.WriteString("#[")
	length := gsn.sig.NumParams()
	for index := uint(0); index < length; index++ {
		param := gsn.sig.Param(index)
		paramValue := gsn.values[index]
		if index != 0 {
			buf.WriteString(", ")
		}
		switch param.Kind() {
		case TypeGenericParam:
			t := paramValue.(*Type)
			buf.WriteString(t.CanonicalName())

		case IntegerGenericParam:
			u64 := paramValue.(uint64)
			writeUint64(buf, u64)

		case EnumGenericParam:
			item := paramValue.(*EnumItem)
			buf.WriteString(item.Name())
		}
	}
	buf.WriteString("]")
}

func writeFunctionCanonicalNameTo(buf *strings.Builder, fsn *FunctionSymbolName) {
	buf.WriteString("#(")

	posLength := fsn.sig.NumPositionalArgs()
	namedLength := fsn.sig.NumNamedArgs()
	first := true

	for index := uint(0); index < posLength; index++ {
		argName := fsn.posNames[index]
		argData := fsn.sig.PositionalArg(index)
		if first {
			buf.WriteString(", ")
			first = false
		}
		buf.WriteString(argName)
		buf.WriteString(": ")
		if argData.IsRepeated() {
			buf.WriteString("...")
		}
		buf.WriteString(argData.Type().CanonicalName())
	}

	for index := uint(0); index < namedLength; index++ {
		argName := fsn.namedNames[index]
		argData := fsn.sig.NamedArg(argName)

		if first {
			buf.WriteString(", ")
			first = false
		}
		buf.WriteString(argName)
		buf.WriteString(": ")
		if argData.IsRepeated() {
			buf.WriteString("...")
		}
		buf.WriteString(argData.Type().CanonicalName())
	}

	buf.WriteString("): ")
	buf.WriteString(fsn.sig.Return().CanonicalName())
}

func writeUnboundGenericMangledNameTo(buf *strings.Builder, gsn *GenericSymbolName) {
	writeGenericMangledNameTo(buf, gsn, false)
}

func writeBoundGenericMangledNameTo(buf *strings.Builder, gsn *GenericSymbolName) {
	writeGenericMangledNameTo(buf, gsn, true)
}

func writeGenericMangledNameTo(buf *strings.Builder, gsn *GenericSymbolName, isBound bool) {
	length := gsn.sig.NumParams()

	genericByte := byte('U')
	if isBound {
		genericByte = 'B'
	}

	buf.WriteByte(genericByte)
	writeUint(buf, length)

	for index := uint(0); index < length; index++ {
		param := gsn.sig.Param(index)
		switch param.Kind() {
		case TypeGenericParam:
			buf.WriteByte('t')
		case IntegerGenericParam:
			buf.WriteByte('u')
		case EnumGenericParam:
			buf.WriteByte('e')
		}
	}

	for index := uint(0); index < length; index++ {
		param := gsn.sig.Param(index)
		if param.Kind() == EnumGenericParam {
			tname := param.Type().MangledName()
			buf.WriteString(tname[2:])
		}
	}

	if isBound {
		for index := uint(0); index < length; index++ {
			param := gsn.sig.Param(index)
			paramValue := gsn.values[index]

			switch param.Kind() {
			case TypeGenericParam:
				t := paramValue.(*Type)
				tname := t.MangledName()
				buf.WriteString(tname[2:])

			case IntegerGenericParam:
				u64 := paramValue.(uint64)
				writeUint64(buf, u64)
				buf.WriteByte('z')

			case EnumGenericParam:
				item := paramValue.(*EnumItem)
				name := item.Name()
				writeUint(buf, uint(len(name)))
				buf.WriteString(name)
			}
		}
	}
}

func writeFunctionMangledNameTo(buf *strings.Builder, fsn *FunctionSymbolName) {
	posLength := fsn.sig.NumPositionalArgs()
	namedLength := fsn.sig.NumNamedArgs()

	buf.WriteByte('F')

	tname := fsn.sig.Return().MangledName()
	buf.WriteString(tname[2:])

	writeUint(buf, posLength)
	for index := uint(0); index < posLength; index++ {
		argData := fsn.sig.PositionalArg(index)
		t := argData.Type()
		// FIXME: builtin::List[T] if argData.IsRepeated()
		tname := t.MangledName()

		buf.WriteString(tname[2:])
	}

	writeUint(buf, namedLength)
	for index := uint(0); index < namedLength; index++ {
		argName := fsn.namedNames[index]
		argData := fsn.sig.NamedArg(argName)
		t := argData.Type()
		// FIXME: builtin::List[T] if argData.IsRepeated()
		tname := t.MangledName()

		buf.WriteByte('A')
		writeUint(buf, uint(len(argName)))
		buf.WriteString(argName)
		buf.WriteString(tname[2:])
	}
}
