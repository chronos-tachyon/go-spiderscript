package memory

import (
	"fmt"
	"runtime"
	"sync"
)

// Memory
// {{{

type Memory struct {
	name   string
	huge   HugePagesMode
	mu     RWLocker
	bytes  []byte
	locked bool
}

func New(name string, hugePages HugePagesMode, multiThreaded bool) *Memory {
	mem := &Memory{
		name:   name,
		huge:   hugePages,
		mu:     (*NoOpRWLocker)(nil),
		bytes:  nil,
		locked: false,
	}
	if multiThreaded {
		mem.mu = new(sync.RWMutex)
	}
	runtime.SetFinalizer(mem, func(x *Memory) { x.Reset() })
	return mem
}

func (mem *Memory) String() string {
	return fmt.Sprintf("memory %q", mem.Name())
}

func (mem *Memory) GoString() string {
	return fmt.Sprintf("Memory(%q)", mem.Name())
}

func (mem *Memory) Name() string {
	return mem.name
}

func (mem *Memory) HugePagesMode() HugePagesMode {
	return mem.huge
}

func (mem *Memory) Reset() {
	mem.SetLen(0)
}

func (mem *Memory) SetLen(length uint) {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	malloc(&mem.bytes, length, mem.huge, mem.locked)
}

func (mem *Memory) Grow(n uint) {
	if n == 0 {
		return
	}

	mem.mu.Lock()
	defer mem.mu.Unlock()

	length := uint(len(mem.bytes)) + n
	malloc(&mem.bytes, length, mem.huge, mem.locked)
}

func (mem *Memory) Shrink(n uint) {
	if n == 0 {
		return
	}

	mem.mu.Lock()
	defer mem.mu.Unlock()

	length := uint(len(mem.bytes))
	if n > length {
		panic(fmt.Errorf("cannot grow to negative size: length=%d, n=%d", length, n))
	}
	length -= n
	malloc(&mem.bytes, length, mem.huge, mem.locked)
}

func (mem *Memory) Protect(r bool, w bool, x bool) error {
	if w && x {
		panic(fmt.Errorf("BUG: illegal protection W|X"))
	}

	mem.mu.Lock()
	defer mem.mu.Unlock()

	return mprotect(mem.bytes, r, w, x)
}

func (mem *Memory) LockToRAM() error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	err := mlock(mem.bytes, true)
	if err == nil {
		mem.locked = true
	}
	return err
}

func (mem *Memory) UnlockFromRAM() error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	mem.locked = false
	return mlock(mem.bytes, false)
}

func (mem *Memory) Size() uint {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	return uint(len(mem.bytes))
}

func (mem *Memory) UInt8s() UInt8Span {
	size := mem.Size()
	return UInt8Span{mem, 0, size, 12}
}

func (mem *Memory) UInt16s() UInt16Span {
	size := mem.Size()
	return UInt16Span{mem, 0, size, 12}
}

func (mem *Memory) UInt32s() UInt32Span {
	size := mem.Size()
	return UInt32Span{mem, 0, size, 12}
}

func (mem *Memory) UInt64s() UInt64Span {
	size := mem.Size()
	return UInt64Span{mem, 0, size, 12}
}

func (mem *Memory) Pages() PageSpan {
	size := mem.Size()
	return PageSpan{mem, 0, size, 12}
}

func (mem *Memory) withWriteLockImpl(i uint, j uint, fn func([]byte) error) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	checkIJ(i, j, uint(len(mem.bytes)))
	return fn(mem.bytes[i:j])
}

func (mem *Memory) withReadLockImpl(i uint, j uint, fn func([]byte) error) error {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	checkIJ(i, j, uint(len(mem.bytes)))
	return fn(mem.bytes[i:j])
}

var _ fmt.Stringer = (*Memory)(nil)
var _ fmt.GoStringer = (*Memory)(nil)

// }}}
