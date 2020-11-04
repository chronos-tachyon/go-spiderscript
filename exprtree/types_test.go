package exprtree

import (
	"testing"
)

func testTypeSingleton(t *testing.T, name string, fn func() *Type) *Type {
	t.Helper()
	type1 := fn()
	type2 := fn()

	if type1 == nil {
		t.Fatalf("%s returned nil", name)
	}

	if type1 != type2 {
		t.Errorf("%s did not return a singleton: %p vs %p", name, type1, type2)
	}

	return type1
}

func testTypeCanonicalName(t *testing.T, name string, type_ *Type, expectedCanonicalName string) {
	t.Helper()
	actualCanonicalName := type_.CanonicalName()
	if expectedCanonicalName != actualCanonicalName {
		t.Errorf("%s.CanonicalName(): expected %q, actual %q", name, expectedCanonicalName, actualCanonicalName)
	}
}

func testTypeMangledName_Exact(t *testing.T, name string, type_ *Type, expectedMangledName string) {
	t.Helper()
	actualMangledName := type_.MangledName()
	if expectedMangledName != actualMangledName {
		t.Errorf("%s.MangledName(): expected %q, actual %q", name, expectedMangledName, actualMangledName)
	}
}

func testTypeMangledName_Global(t *testing.T, name string, type_ *Type, mod *Module, symbolName string) {
	t.Helper()
	expectedMangledName := MangleGlobalSymbolName(mod, symbolName)
	actualMangledName := type_.MangledName()
	if expectedMangledName != actualMangledName {
		t.Errorf("%s.MangledName(): expected %q, actual %q", name, expectedMangledName, actualMangledName)
	}
}

func testTypePadding(t *testing.T, name string, type_ *Type) {
	t.Helper()

	maxAlignShift := uint(8)
	if actualAlignShift := type_.AlignShift(); actualAlignShift > maxAlignShift {
		t.Errorf("%s.AlignShift(): expected 0 ≤ x ≤ %d, got %d", name, maxAlignShift, actualAlignShift)
	}

	k := type_.AlignBytes()
	x := type_.MinimumBytes()
	y := type_.PaddedBytes()

	i := uint(1)
	z := k
	for z < x {
		i++
		z += k
	}

	if x > y {
		t.Errorf("%s.MinimumBytes(): expected x ≥ %d, actual %d", name, y, x)
	}

	if y != z {
		t.Errorf("%s.PaddedBytes(): expected %d [%d*%d], actual %d", name, z, i, k, y)
	}
}

func TestType_VoidType(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "VoidType()", interp.VoidType)
	testTypeCanonicalName(t, "VoidType()", type_, "builtin::Void")
	testTypeMangledName_Exact(t, "VoidType()", type_, "_Av")
	testTypePadding(t, "VoidType()", type_)
}

func TestType_NullType(t *testing.T) {
	interp := GlobalTestInterp()
	builtinModule := interp.BuiltinModule()
	type_ := testTypeSingleton(t, "NullType()", interp.NullType)
	testTypeCanonicalName(t, "NullType()", type_, "builtin::Null")
	testTypeMangledName_Global(t, "NullType()", type_, builtinModule, "Null")
	testTypePadding(t, "NullType()", type_)
}

func TestType_TypeType(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "TypeType()", interp.TypeType)
	testTypeCanonicalName(t, "TypeType()", type_, "builtin::Type")
	testTypeMangledName_Exact(t, "TypeType()", type_, "_At")
	testTypePadding(t, "TypeType()", type_)
}

func TestType_UInt8(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "UInt8Type()", interp.UInt8Type)
	testTypeCanonicalName(t, "UInt8Type()", type_, "builtin::UInt8")
	testTypeMangledName_Exact(t, "UInt8Type()", type_, "_Au0")
	testTypePadding(t, "UInt8Type()", type_)
}

func TestType_UInt16(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "UInt16Type()", interp.UInt16Type)
	testTypeCanonicalName(t, "UInt16Type()", type_, "builtin::UInt16")
	testTypeMangledName_Exact(t, "UInt16Type()", type_, "_Au1")
	testTypePadding(t, "UInt16Type()", type_)
}

func TestType_UInt32(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "UInt32Type()", interp.UInt32Type)
	testTypeCanonicalName(t, "UInt32Type()", type_, "builtin::UInt32")
	testTypeMangledName_Exact(t, "UInt32Type()", type_, "_Au2")
	testTypePadding(t, "UInt32Type()", type_)
}

func TestType_UInt64(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "UInt64Type()", interp.UInt64Type)
	testTypeCanonicalName(t, "UInt64Type()", type_, "builtin::UInt64")
	testTypeMangledName_Exact(t, "UInt64Type()", type_, "_Au3")
	testTypePadding(t, "UInt64Type()", type_)
}

func TestType_SInt8(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "SInt8Type()", interp.SInt8Type)
	testTypeCanonicalName(t, "SInt8Type()", type_, "builtin::SInt8")
	testTypeMangledName_Exact(t, "SInt8Type()", type_, "_Ai0")
	testTypePadding(t, "SInt8Type()", type_)
}

func TestType_SInt16(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "SInt16Type()", interp.SInt16Type)
	testTypeCanonicalName(t, "SInt16Type()", type_, "builtin::SInt16")
	testTypeMangledName_Exact(t, "SInt16Type()", type_, "_Ai1")
	testTypePadding(t, "SInt16Type()", type_)
}

func TestType_SInt32(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "SInt32Type()", interp.SInt32Type)
	testTypeCanonicalName(t, "SInt32Type()", type_, "builtin::SInt32")
	testTypeMangledName_Exact(t, "SInt32Type()", type_, "_Ai2")
	testTypePadding(t, "SInt32Type()", type_)
}

func TestType_SInt64(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "SInt64Type()", interp.SInt64Type)
	testTypeCanonicalName(t, "SInt64Type()", type_, "builtin::SInt64")
	testTypeMangledName_Exact(t, "SInt64Type()", type_, "_Ai3")
	testTypePadding(t, "SInt64Type()", type_)
}

func TestType_Float16(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "Float16Type()", interp.Float16Type)
	testTypeCanonicalName(t, "Float16Type()", type_, "builtin::Float16")
	testTypeMangledName_Exact(t, "Float16Type()", type_, "_Af1")
	testTypePadding(t, "Float16Type()", type_)
}

func TestType_Float32(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "Float32Type()", interp.Float32Type)
	testTypeCanonicalName(t, "Float32Type()", type_, "builtin::Float32")
	testTypeMangledName_Exact(t, "Float32Type()", type_, "_Af2")
	testTypePadding(t, "Float32Type()", type_)
}

func TestType_Float64(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "Float64Type()", interp.Float64Type)
	testTypeCanonicalName(t, "Float64Type()", type_, "builtin::Float64")
	testTypeMangledName_Exact(t, "Float64Type()", type_, "_Af3")
	testTypePadding(t, "Float64Type()", type_)
}

func TestType_Complex32(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "Complex32Type()", interp.Complex32Type)
	testTypeCanonicalName(t, "Complex32Type()", type_, "builtin::Complex32")
	testTypeMangledName_Exact(t, "Complex32Type()", type_, "_Ac1")
	testTypePadding(t, "Complex32Type()", type_)
}

func TestType_Complex64(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "Complex64Type()", interp.Complex64Type)
	testTypeCanonicalName(t, "Complex64Type()", type_, "builtin::Complex64")
	testTypeMangledName_Exact(t, "Complex64Type()", type_, "_Ac2")
	testTypePadding(t, "Complex64Type()", type_)
}

func TestType_Complex128(t *testing.T) {
	interp := GlobalTestInterp()
	type_ := testTypeSingleton(t, "Complex128Type()", interp.Complex128Type)
	testTypeCanonicalName(t, "Complex128Type()", type_, "builtin::Complex128")
	testTypeMangledName_Exact(t, "Complex128Type()", type_, "_Ac3")
	testTypePadding(t, "Complex128Type()", type_)
}
