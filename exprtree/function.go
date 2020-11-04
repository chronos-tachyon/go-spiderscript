package exprtree

import (
	"sync"
)

type FunctionType struct {
	Signature *FunctionSignature
}

type Function struct {
	mu  sync.Mutex
	sym *Symbol
	sig *FunctionSignature
}

func (f *Function) Symbol() *Symbol {
	return f.sym
}

func (f *Function) CanonicalName() string {
	return f.sym.CanonicalName()
}

func (f *Function) MangledName() string {
	return f.sym.MangledName()
}

func (f *Function) Signature() *FunctionSignature {
	return f.sig
}
