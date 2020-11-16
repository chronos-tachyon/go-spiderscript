package memory

import (
	"sync"
)

type RWLocker interface {
	sync.Locker
	RLock()
	RUnlock()
	RLocker() sync.Locker
}

type NoOpRWLocker struct{}

func (*NoOpRWLocker) Lock()    {}
func (*NoOpRWLocker) Unlock()  {}
func (*NoOpRWLocker) RLock()   {}
func (*NoOpRWLocker) RUnlock() {}

func (*NoOpRWLocker) RLocker() sync.Locker { return (*NoOpRWLocker)(nil) }

var _ RWLocker = (*NoOpRWLocker)(nil)
var _ sync.Locker = (*NoOpRWLocker)(nil)
