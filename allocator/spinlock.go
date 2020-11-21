package allocator

import (
	"sync"
	"sync/atomic"
)

type Spinlock struct {
	v uint32
}

func (mu *Spinlock) Lock() {
	for !atomic.CompareAndSwapUint32(&mu.v, 0, 1) {
	}
}

func (mu *Spinlock) Unlock() {
	atomic.StoreUint32(&mu.v, 0)
}

var _ sync.Locker = (*Spinlock)(nil)
