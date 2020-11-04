package exprtree

import (
	"testing"
)

func TestRegexpModuleName(t *testing.T) {
	type testRow struct {
		Input string
		Expected bool
	}

	var testData = []testRow{
		{"builtin", true},
		{"snake_case", true},
		{"CamelCase", true},
		{"typical::name", true},
		{"has::three::components", true},
		{"four::components::like::this", true},
		{"Camel::Case", true},

		{"one:colon", false},
		{"three:::colons", false},

		{"two__underscores", false},
		{"something::two__underscores", false},

		{"", false},

		{"_", true},
		{"_a", false},
		{"_leading_underscore", false},
		{"_LeadingUnderscore", false},
		{"something::_leading_underscore", false},

		{"__", false},
		{"__a", false},
		{"__two_leading_underscores", false},
		{"__TwoLeadingUnderscores", false},
		{"something::__two_leading_underscores", false},

		{"$", false},
		{"$a", false},
		{"$leading_dollar", false},
		{"something::$leading_dollar", false},
	}

	for _, row := range testData {
		if actual := reModuleName.MatchString(row.Input); actual != row.Expected {
			t.Errorf("reModuleName.MatchString(%q): expected %t, actual %t", row.Input, row.Expected, actual)
		}
	}
}

func TestRegexpSymbolName(t *testing.T) {
	type testRow struct {
		Input string
		Expected bool
	}

	var testData = []testRow{
		{"simple", true},
		{"snake_case", true},
		{"CamelCase", true},
		{"two__underscores", true},

		{"_", true},
		{"_leading_underscore", true},
		{"__two_leading_underscores", true},
		{"_leading_and_trailing_underscores_", true},
		{"__two_leading_and_two_trailing_underscores__", true},

		{"$", true},
		{"$leading_dollar", true},
		{"$_dollar_underscore", true},
		{"$__dollar_two_underscores", true},
		{"$CamelCase", true},
		{"$_CamelCase", true},
		{"$__CamelCase", true},

		{"", false},

		{"mid$dollar", false},
		{"trailing_dollar$", false},
		{"$$two_leading_dollars", false},
	}

	for _, row := range testData {
		if actual := reSymbolName.MatchString(row.Input); actual != row.Expected {
			t.Errorf("reSymbolName.MatchString(%q): expected %t, actual %t", row.Input, row.Expected, actual)
		}
	}
}

func TestLengthUint64(t *testing.T) {
	type testRow struct {
		Input    uint64
		Expected uint
	}

	var testData = []testRow{
		{0, 1},
		{1, 1},
		{9, 1},
		{10, 2},
		{99, 2},
		{100, 3},
		{999, 3},
		{1000, 4},
		{9999, 4},
		{10000, 5},
		{99999, 5},
		{100000, 6},
		{999999, 6},
		{1000000, 7},
	}

	for _, row := range testData {
		if actual := lengthUint64(row.Input); actual != row.Expected {
			t.Errorf("lengthUint64(%d): expected %d, actual %d", row.Input, row.Expected, actual)
		}
	}
}
