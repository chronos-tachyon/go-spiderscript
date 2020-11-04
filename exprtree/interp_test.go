package exprtree

import (
	"sort"
	"testing"
)

func TestInterp_New(t *testing.T) {
	interp := GlobalTestInterp()

	dumpModule(t, interp.BuiltinModule())
	dumpModule(t, interp.BuiltinEnumModule())
	dumpModule(t, interp.BuiltinBitfieldModule())
	dumpModule(t, interp.BuiltinStructModule())
	dumpModule(t, interp.BuiltinUnionModule())
}

func dumpModule(t *testing.T, mod *Module) {
	t.Helper()

	symbols := make(map[string]*Symbol, 64)
	mod.Symbols().All(symbols)

	symbolNames := make([]string, 0, len(symbols))
	for name := range symbols {
		symbolNames = append(symbolNames, name)
	}
	sort.Strings(symbolNames)

	for _, name := range symbolNames {
		sym := symbols[name]
		t.Logf("%s[%q] = %#v", mod.CanonicalName(), name, sym)
	}
	t.Logf("%s.", mod.CanonicalName())
}
