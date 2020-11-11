package exprtree

import (
	"fmt"
	"runtime"
	"sync"
)

// Interface: MemoryView
// {{{

type MemoryView interface {
	Size() uint
	All() MemorySpan
	Span(i, j uint) MemorySpan
	Absolute() (*Memory, uint, uint)

	AllWithWriteLock(fn func([]byte) error) error
	AllWithReadLock(fn func([]byte) error) error

	SpanWithWriteLock(i uint, j uint, fn func([]byte) error) error
	SpanWithReadLock(i uint, j uint, fn func([]byte) error) error

	String() string
	GoString() string
}

// Memory
// {{{

type Memory struct {
	mu    sync.RWMutex
	name  string
	bytes []byte
	locked bool
	huge  HugePagesMode
}

func NewMemory(name string, hugePages HugePagesMode) *Memory {
	mem := new(Memory)
	mem.name = name
	mem.huge = hugePages
	runtime.SetFinalizer(mem, func(x *Memory) { x.Reset() })
	return mem
}

func (mem *Memory) Name() string {
	return mem.name
}

func (mem *Memory) HugePagesMode() HugePagesMode {
	return mem.huge
}

func (mem *Memory) Reset() {
	mem.Truncate(0)
}

func (mem *Memory) Truncate(n uint) {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	malloc(&mem.bytes, n, mem.huge, mem.locked)
}

func (mem *Memory) ProtectPages(r bool, w bool, x bool) error {
	if w && x {
		panic(fmt.Errorf("BUG: illegal protection W|X"))
	}

	mem.mu.Lock()
	defer mem.mu.Unlock()

	return mprotect(mem.bytes, r, w, x)
}

func (mem *Memory) LockPages() error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	err := mlock(mem.bytes, true)
	if err == nil {
		mem.locked = true
	}
	return err
}

func (mem *Memory) UnlockPages() error {
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

func (mem *Memory) All() MemorySpan {
	size := mem.Size()
	return MemorySpan{mem: mem, i: 0, j: size}
}

func (mem *Memory) Span(i uint, j uint) MemorySpan {
	size := mem.Size()
	checkIJ(i, j, size)
	return MemorySpan{mem: mem, i: i, j: j}
}

func (mem *Memory) Absolute() (*Memory, uint, uint) {
	size := mem.Size()
	return mem, 0, size
}

func (mem *Memory) AllWithWriteLock(fn func([]byte) error) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	return fn(mem.bytes)
}

func (mem *Memory) AllWithReadLock(fn func([]byte) error) error {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	return fn(mem.bytes)
}

func (mem *Memory) SpanWithWriteLock(i uint, j uint, fn func([]byte) error) error {
	mem.mu.Lock()
	defer mem.mu.Unlock()

	checkIJ(i, j, uint(len(mem.bytes)))
	return fn(mem.bytes[i:j])
}

func (mem *Memory) SpanWithReadLock(i uint, j uint, fn func([]byte) error) error {
	mem.mu.RLock()
	defer mem.mu.RUnlock()

	checkIJ(i, j, uint(len(mem.bytes)))
	return fn(mem.bytes[i:j])
}

func (mem *Memory) String() string {
	return fmt.Sprintf("memory %q", mem.Name())
}

func (mem *Memory) GoString() string {
	return fmt.Sprintf("Memory(%q)", mem.Name())
}

var _ MemoryView = (*Memory)(nil)
var _ fmt.Stringer = (*Memory)(nil)
var _ fmt.GoStringer = (*Memory)(nil)

// }}}

// MemorySpan
// {{{

type MemorySpan struct {
	mem *Memory
	i   uint
	j   uint
}

func (span MemorySpan) StartOffset() uint {
	return span.i
}

func (span MemorySpan) EndOffset() uint {
	return span.j
}

func (span MemorySpan) Size() uint {
	return span.j - span.i
}

func (span MemorySpan) All() MemorySpan {
	return span
}

func (span MemorySpan) Span(i, j uint) MemorySpan {
	size := span.Size()
	checkIJ(i, j, size)
	if i == 0 && j == size {
		return span
	}
	i += span.i
	j += span.i
	return MemorySpan{mem: span.mem, i: i, j: j}
}

func (span MemorySpan) Absolute() (*Memory, uint, uint) {
	return span.mem, span.i, span.j
}

func (span MemorySpan) AllWithWriteLock(fn func([]byte) error) error {
	return span.mem.SpanWithWriteLock(span.i, span.j, fn)
}

func (span MemorySpan) AllWithReadLock(fn func([]byte) error) error {
	return span.mem.SpanWithReadLock(span.i, span.j, fn)
}

func (span MemorySpan) SpanWithWriteLock(i uint, j uint, fn func([]byte) error) error {
	checkIJ(i, j, span.Size())
	i += span.i
	j += span.i
	return span.mem.SpanWithWriteLock(i, j, fn)
}

func (span MemorySpan) SpanWithReadLock(i uint, j uint, fn func([]byte) error) error {
	checkIJ(i, j, span.Size())
	i += span.i
	j += span.i
	return span.mem.SpanWithReadLock(i, j, fn)
}

func (span MemorySpan) String() string {
	return fmt.Sprintf("memory %q span [%d:%d]", span.mem.Name(), span.i, span.j)
}

func (span MemorySpan) GoString() string {
	return fmt.Sprintf("MemorySpan(%q, %#x, %#x)", span.mem.Name(), span.i, span.j)
}

var _ MemoryView = MemorySpan{}
var _ fmt.Stringer = MemorySpan{}
var _ fmt.GoStringer = MemorySpan{}

// }}}

// }}}
