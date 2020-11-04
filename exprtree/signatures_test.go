package exprtree

import (
	"testing"
)

func TestGenericSignature(t *testing.T) {
	interp := GlobalTestInterp()
	builder := interp.GenericSignatureBuilder()

	emptySignature := builder.Reset().Build()

	oneGenericSignature := builder.Reset().WithType().Build()
	twoTypesSignature := builder.Reset().WithType().WithType().Build()
	oneIntSignature := builder.Reset().WithUInt().Build()
	twoIntsSignature := builder.Reset().WithUInt().WithUInt().Build()
	typeAndIntSignature := builder.Reset().WithType().WithUInt().Build()
	intAndGenericSignature := builder.Reset().WithUInt().WithType().Build()

	type testRow struct {
		Name             string
		Input            *GenericSignature
		ExpectedString   string
		ExpectedGoString string
	}

	testData := []testRow{
		{
			Name:             "emptySignature",
			Input:            emptySignature,
			ExpectedString:   "[]",
			ExpectedGoString: "GenericSignature()",
		},
		{
			Name:             "oneGenericSignature",
			Input:            oneGenericSignature,
			ExpectedString:   "[type]",
			ExpectedGoString: "GenericSignature(type)",
		},
		{
			Name:             "twoTypesSignature",
			Input:            twoTypesSignature,
			ExpectedString:   "[type, type]",
			ExpectedGoString: "GenericSignature(type, type)",
		},
		{
			Name:             "oneIntSignature",
			Input:            oneIntSignature,
			ExpectedString:   "[uint]",
			ExpectedGoString: "GenericSignature(uint)",
		},
		{
			Name:             "twoIntsSignature",
			Input:            twoIntsSignature,
			ExpectedString:   "[uint, uint]",
			ExpectedGoString: "GenericSignature(uint, uint)",
		},
		{
			Name:             "typeAndIntSignature",
			Input:            typeAndIntSignature,
			ExpectedString:   "[type, uint]",
			ExpectedGoString: "GenericSignature(type, uint)",
		},
		{
			Name:             "intAndGenericSignature",
			Input:            intAndGenericSignature,
			ExpectedString:   "[uint, type]",
			ExpectedGoString: "GenericSignature(uint, type)",
		},
	}

	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
			if actual := row.Input.String(); actual != row.ExpectedString {
				t.Errorf("String(): expected %q, actual %q", row.ExpectedString, actual)
			}
			if actual := row.Input.GoString(); actual != row.ExpectedGoString {
				t.Errorf("GoString(): expected %q, actual %q", row.ExpectedGoString, actual)
			}
		})
	}
}

func TestFunctionSignature(t *testing.T) {
	interp := GlobalTestInterp()
	builder := interp.FunctionSignatureBuilder()

	u64 := interp.UInt64Type()

	voidSignature := builder.Reset().Build()
	allPositionalSignature := builder.Reset().WithReturn(u64).WithPositionalArg(u64).WithPositionalArg(u64).Build()

	type testRow struct {
		Name             string
		Input            *FunctionSignature
		ExpectedString   string
		ExpectedGoString string
	}

	testData := []testRow{
		{
			Name:             "voidSignature",
			Input:            voidSignature,
			ExpectedString:   "(): builtin::Void",
			ExpectedGoString: "FunctionSignature(builtin::Void)",
		},
		{
			Name:             "allPositionalSignature",
			Input:            allPositionalSignature,
			ExpectedString:   "(builtin::UInt64, builtin::UInt64): builtin::UInt64",
			ExpectedGoString: "FunctionSignature(builtin::UInt64, builtin::UInt64, builtin::UInt64)",
		},
	}

	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
			if actual := row.Input.String(); actual != row.ExpectedString {
				t.Errorf("String(): expected %q, actual %q", row.ExpectedString, actual)
			}
			if actual := row.Input.GoString(); actual != row.ExpectedGoString {
				t.Errorf("GoString(): expected %q, actual %q", row.ExpectedGoString, actual)
			}
		})
	}
}
