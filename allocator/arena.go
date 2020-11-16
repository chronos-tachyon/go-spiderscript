package allocator

import (
	"fmt"
	"sync"

	"github.com/chronos-tachyon/go-spiderscript/memory"
)

type Arena struct {
	mem    *memory.Memory
	mu     sync.Locker
	offset uint
}

func NewArena(name string, hugePages memory.HugePagesMode, multiThreaded bool) *Arena {
	alloc := new(Arena)
	alloc.Init(name, hugePages, multiThreaded)
	return alloc
}

func (alloc *Arena) Init(name string, hugePages memory.HugePagesMode, multiThreaded bool) {
	*alloc = Arena{
		mem:    memory.New(name, hugePages, multiThreaded),
		mu:     (*memory.NoOpRWLocker)(nil),
		offset: 0,
	}
	if multiThreaded {
		alloc.mu = new(sync.Mutex)
	}
}

func (alloc *Arena) Allocate(count uint, alignShift uint) memory.UInt8Span {
	if alignShift > 12 {
		panic(fmt.Errorf("alignment is too strict: requested=%d, max=%d", alignShift, 12))
	}

	alignSize := uint(1) << alignShift
	alignMask := alignSize - 1

	alloc.mu.Lock()
	defer alloc.mu.Unlock()

	allocAddr := (alloc.offset + alignMask) & ^alignMask
	allocSize := count << alignShift
	allocEnd := allocAddr + allocSize

	alloc.offset = allocEnd
	alloc.mem.SetLen(allocEnd)
	return alloc.mem.UInt8s().Span(allocAddr, allocEnd)
}

func (alloc *Arena) Free(span memory.UInt8Span) {
	// no op
}

func (alloc *Arena) Trim() {
	// no op
}

func (alloc *Arena) FreeAll() {
	alloc.mu.Lock()
	defer alloc.mu.Unlock()

	alloc.offset = 0
	alloc.mem.Reset()
}

var _ Allocator = (*Arena)(nil)
