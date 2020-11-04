package exprtree

import (
	"testing"
)

func TestTypeKind(t *testing.T) {
	type testRow struct {
		Name             string
		Input            TypeKind
		ExpectedString   string
		ExpectedGoString string
	}

	testData := []testRow{
		{
			Name:             "InvalidTypeKind",
			Input:            InvalidTypeKind,
			ExpectedString:   "InvalidTypeKind",
			ExpectedGoString: "InvalidTypeKind",
		},
		{
			Name:             "ReflectedTypeKind",
			Input:            ReflectedTypeKind,
			ExpectedString:   "ReflectedTypeKind",
			ExpectedGoString: "ReflectedTypeKind",
		},
		{
			Name:             "NamedKind",
			Input:            NamedKind,
			ExpectedString:   "NamedKind",
			ExpectedGoString: "NamedKind",
		},
		{
			Name:             "NamedKind+1",
			Input:            NamedKind + 1,
			ExpectedString:   "TypeKind(30)",
			ExpectedGoString: "TypeKind(30)",
		},
	}

	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
			if actual := row.Input.String(); actual != row.ExpectedString {
				t.Errorf("String(): expect %q, actual %q", row.ExpectedString, actual)
			}
			if actual := row.Input.GoString(); actual != row.ExpectedGoString {
				t.Errorf("GoString(): expect %q, actual %q", row.ExpectedGoString, actual)
			}
		})
	}
}

func TestGenericParamKind(t *testing.T) {
	type testRow struct {
		Name             string
		Input            GenericParamKind
		ExpectedString   string
		ExpectedGoString string
	}

	testData := []testRow{
		{
			Name:             "InvalidGenericParamKind",
			Input:            InvalidGenericParamKind,
			ExpectedString:   "InvalidGenericParamKind",
			ExpectedGoString: "InvalidGenericParamKind",
		},
		{
			Name:             "TypeGenericParam",
			Input:            TypeGenericParam,
			ExpectedString:   "TypeGenericParam",
			ExpectedGoString: "TypeGenericParam",
		},
		{
			Name:             "EnumGenericParam",
			Input:            EnumGenericParam,
			ExpectedString:   "EnumGenericParam",
			ExpectedGoString: "EnumGenericParam",
		},
		{
			Name:             "EnumGenericParam+1",
			Input:            EnumGenericParam + 1,
			ExpectedString:   "GenericParamKind(4)",
			ExpectedGoString: "GenericParamKind(4)",
		},
	}

	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
			if actual := row.Input.String(); actual != row.ExpectedString {
				t.Errorf("String(): expect %q, actual %q", row.ExpectedString, actual)
			}
			if actual := row.Input.GoString(); actual != row.ExpectedGoString {
				t.Errorf("GoString(): expect %q, actual %q", row.ExpectedGoString, actual)
			}
		})
	}
}
