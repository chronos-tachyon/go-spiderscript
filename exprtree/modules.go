package exprtree

import (
	"fmt"
	"sync"
)

type Module struct {
	mu      sync.RWMutex
	interp  *Interp
	cname   string
	mname   string
	imports map[string]*Module
	symbols SymbolTable
}

func (mod *Module) Check() {
	if mod == nil {
		panic(fmt.Errorf("BUG: *Module is nil"))
	}
	if mod.interp == nil {
		panic(fmt.Errorf("BUG: (*Module).Interp() is nil"))
	}
	if expected := MangleModuleName(mod.cname); mod.mname != expected {
		panic(fmt.Errorf("BUG: cname=%q, mname=%q, MangleModuleName(cname)=%q", mod.cname, mod.mname, expected))
	}
}

func (mod *Module) Interp() *Interp {
	return mod.interp
}

func (mod *Module) CanonicalName() string {
	return mod.cname
}

func (mod *Module) MangledName() string {
	return mod.mname
}

func (mod *Module) Symbols() *SymbolTable {
	return &mod.symbols
}

func (mod *Module) AllImports(out map[string]*Module) {
	if out == nil {
		panic(fmt.Errorf("BUG: out is nil"))
	}

	locked(mod.mu.RLocker(), func() {
		for name, imported := range mod.imports {
			out[name] = imported
		}
	})
}

func (mod *Module) Import(name string) *Module {
	var out *Module
	locked(mod.mu.RLocker(), func() {
		out = mod.imports[name]
	})
	return out
}

func (mod *Module) AddImport(name string, imported *Module) error {
	mod.Check()
	imported.Check()

	if !reModuleName.MatchString(name) {
		return fmt.Errorf("invalid module name %q", name)
	}

	var old *Module
	locked(&mod.mu, func() {
		old = mod.imports[name]
		if old == nil {
			mod.imports[name] = imported
		}
	})

	if old != nil {
		return &DuplicateModuleError{Name: name, Old: old, New: imported}
	}
	return nil
}
