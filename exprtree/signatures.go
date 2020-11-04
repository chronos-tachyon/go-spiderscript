package exprtree

import (
	"fmt"
	"sort"
)

// GenericSignatureBuilder
// {{{

type GenericSignatureBuilder struct {
	interp *Interp
	params []GenericParam
}

func (interp *Interp) GenericSignatureBuilder() *GenericSignatureBuilder {
	return &GenericSignatureBuilder{
		interp: interp,
		params: make([]GenericParam, 0, 4),
	}
}

func (builder *GenericSignatureBuilder) Reset() *GenericSignatureBuilder {
	builder.params = builder.params[:0]
	return builder
}

func (builder *GenericSignatureBuilder) WithType() *GenericSignatureBuilder {
	builder.params = append(builder.params, GenericParam{kind: TypeGenericParam})
	return builder
}

func (builder *GenericSignatureBuilder) WithUInt() *GenericSignatureBuilder {
	builder.params = append(builder.params, GenericParam{kind: IntegerGenericParam})
	return builder
}

func (builder *GenericSignatureBuilder) WithEnum(t *Type) *GenericSignatureBuilder {
	checkNotNil("t", t)
	if kind := t.Chase().Kind(); kind != EnumKind {
		panic(fmt.Errorf("BUG: t.Chase().Kind() is %v, expected EnumKind", kind))
	}
	builder.params = append(builder.params, GenericParam{kind: EnumGenericParam, type_: t})
	return builder
}

func (builder *GenericSignatureBuilder) Build() *GenericSignature {
	interp := builder.interp
	params := builder.params
	return interp.registerGenSig(&GenericSignature{
		interp: interp,
		params: cloneGenericParams(params),
	})
}

// }}}

// GenericSignature
// {{{

type GenericSignature struct {
	interp *Interp
	params []GenericParam
	id     GenericSignatureID
}

func (g *GenericSignature) Interp() *Interp {
	return g.interp
}

func (g *GenericSignature) ID() GenericSignatureID {
	return g.id
}

func (g *GenericSignature) NumParams() uint {
	return uint(len(g.params))
}

func (g *GenericSignature) Param(index uint) GenericParam {
	return g.params[index]
}

func (g *GenericSignature) Params() []GenericParam {
	return cloneGenericParams(g.params)
}

func (g *GenericSignature) String() string {
	return g.stringImpl(false)
}

func (g *GenericSignature) GoString() string {
	return g.stringImpl(true)
}

func (g *GenericSignature) stringImpl(goMode bool) string {
	estimatedLen := uint(2)
	if goMode {
		estimatedLen = 18
	}

	length := g.NumParams()
	if length != 0 {
		estimatedLen += 2*length - 2
	}

	strings := make([]string, length)
	for index := uint(0); index < length; index++ {
		str := g.Param(index).String()
		estimatedLen += uint(len(str))
		strings[index] = str
	}

	buf := takeBuffer()
	buf.Grow(int(estimatedLen))
	defer giveBuffer(buf)

	if goMode {
		buf.WriteString("GenericSignature(")
	} else {
		buf.WriteByte('[')
	}

	for index, str := range strings {
		if index > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(str)
	}

	if goMode {
		buf.WriteByte(')')
	} else {
		buf.WriteByte(']')
	}

	return checkEstimatedLength(buf, estimatedLen)
}

var _ fmt.Stringer = (*GenericSignature)(nil)
var _ fmt.GoStringer = (*GenericSignature)(nil)

// }}}

// GenericParam
// {{{

type GenericParam struct {
	type_ *Type
	kind  GenericParamKind
}

func (param GenericParam) Kind() GenericParamKind {
	return param.kind
}

func (param GenericParam) Type() *Type {
	return param.type_
}

func (param GenericParam) String() string {
	switch param.Kind() {
	case TypeGenericParam:
		return "type"

	case IntegerGenericParam:
		return "uint"

	case EnumGenericParam:
		return param.Type().CanonicalName()

	default:
		panic(fmt.Errorf("BUG: GenericParam.Kind() is %v, which is not implemented", param.Kind()))
	}
}

var _ fmt.Stringer = GenericParam{}

// }}}

// FunctionSignatureBuilder
// {{{

type FunctionSignatureBuilder struct {
	interp *Interp
	ret    *Type
	pos    []FunctionArg
	named  map[string]FunctionArg
}

func (interp *Interp) FunctionSignatureBuilder() *FunctionSignatureBuilder {
	return &FunctionSignatureBuilder{
		interp: interp,
		ret:    interp.VoidType(),
		pos:    make([]FunctionArg, 0, 4),
		named:  make(map[string]FunctionArg, 4),
	}
}

func (builder *FunctionSignatureBuilder) Reset() *FunctionSignatureBuilder {
	builder.ret = builder.interp.VoidType()
	builder.pos = builder.pos[:0]
	for key := range builder.named {
		delete(builder.named, key)
	}
	return builder
}

func (builder *FunctionSignatureBuilder) WithReturn(t *Type) *FunctionSignatureBuilder {
	checkNotNil("t", t)
	builder.ret = t
	return builder
}

func (builder *FunctionSignatureBuilder) WithPositionalArg(t *Type) *FunctionSignatureBuilder {
	checkNotNil("t", t)
	builder.pos = append(builder.pos, FunctionArg{type_: t, repeat: false})
	return builder
}

func (builder *FunctionSignatureBuilder) WithRepeatedPositionalArg(t *Type) *FunctionSignatureBuilder {
	checkNotNil("t", t)
	builder.pos = append(builder.pos, FunctionArg{type_: t, repeat: true})
	return builder
}

func (builder *FunctionSignatureBuilder) WithNamedArg(name string, t *Type) *FunctionSignatureBuilder {
	if !reSymbolName.MatchString(name) {
		panic(fmt.Errorf("BUG: name is invalid; got %q", name))
	}
	checkNotNil("t", t)
	builder.named[name] = FunctionArg{type_: t, repeat: false}
	return builder
}

func (builder *FunctionSignatureBuilder) Build() *FunctionSignature {
	interp := builder.interp
	ret := builder.ret
	pos := builder.pos
	named := builder.named
	return interp.registerFuncSig(&FunctionSignature{
		interp: interp,
		ret:    ret,
		pos:    cloneFunctionArgs(pos),
		named:  cloneFunctionArgMap(named),
	})
}

// }}}

// FunctionSignature
// {{{

type FunctionSignature struct {
	interp *Interp
	ret    *Type
	pos    []FunctionArg
	named  map[string]FunctionArg
	id     FunctionSignatureID
}

func (f *FunctionSignature) Interp() *Interp {
	return f.interp
}

func (f *FunctionSignature) ID() FunctionSignatureID {
	return f.id
}

func (f *FunctionSignature) Return() *Type {
	return f.ret
}

func (f *FunctionSignature) NumPositionalArgs() uint {
	return uint(len(f.pos))
}

func (f *FunctionSignature) PositionalArg(index uint) FunctionArg {
	return f.pos[index]
}

func (f *FunctionSignature) PositionalArgs() []FunctionArg {
	return cloneFunctionArgs(f.pos)
}

func (f *FunctionSignature) NumNamedArgs() uint {
	return uint(len(f.named))
}

func (f *FunctionSignature) ArgNames() []string {
	out := make([]string, 0, len(f.named))
	for name := range f.named {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func (f *FunctionSignature) NamedArg(name string) FunctionArg {
	return f.named[name]
}

func (f *FunctionSignature) NamedArgs() map[string]FunctionArg {
	return cloneFunctionArgMap(f.named)
}

func (f *FunctionSignature) String() string {
	return f.stringImpl(false)
}

func (f *FunctionSignature) GoString() string {
	return f.stringImpl(true)
}

func (f *FunctionSignature) stringImpl(goMode bool) string {
	m := uint(len(f.pos))
	n := uint(len(f.named))
	mn := m + n
	o := uint(len(f.ret.CanonicalName()))
	estimatedLen := 4 + o
	if goMode {
		estimatedLen += 15
	}
	if mn != 0 {
		estimatedLen += 2*m + 4*n
		if !goMode {
			estimatedLen -= 2
		}
		for _, arg := range f.pos {
			estimatedLen += uint(len(arg.Type().CanonicalName()))
			if arg.repeat {
				estimatedLen += 3
			}
		}
		for name, arg := range f.named {
			estimatedLen += uint(len(name)) + uint(len(arg.Type().CanonicalName()))
		}
	}

	buf := takeBuffer()
	buf.Grow(int(estimatedLen))
	defer giveBuffer(buf)

	if goMode {
		buf.WriteString("FunctionSignature")
	}
	buf.WriteByte('(')

	first := true
	if goMode {
		buf.WriteString(f.ret.CanonicalName())
		first = false
	}

	if m != 0 {
		for index := uint(0); index < m; index++ {
			arg := f.pos[index]
			if !first {
				buf.WriteString(", ")
			}
			if arg.repeat {
				buf.WriteString("...")
			}
			buf.WriteString(arg.Type().CanonicalName())
			first = false
		}
	}

	if n != 0 {
		for name, arg := range f.named {
			if !first {
				buf.WriteString(", ")
			}
			buf.WriteString(name)
			buf.WriteString(": ")
			buf.WriteString(arg.Type().CanonicalName())
			first = false
		}
	}

	buf.WriteByte(')')

	if !goMode {
		buf.WriteString(": ")
		buf.WriteString(f.ret.CanonicalName())
	}

	return checkEstimatedLength(buf, estimatedLen)
}

var _ fmt.Stringer = (*FunctionSignature)(nil)
var _ fmt.GoStringer = (*FunctionSignature)(nil)

// }}}

// FunctionArg
// {{{

type FunctionArg struct {
	type_  *Type
	repeat bool
}

func (arg FunctionArg) Type() *Type {
	return arg.type_
}

func (arg FunctionArg) IsRepeated() bool {
	return arg.repeat
}

// }}}
