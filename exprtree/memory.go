package exprtree

import (
	"fmt"
	"sync"
)

// Memory
// {{{

type Memory struct {
	mu    sync.RWMutex
	name  string
	bytes []byte
}

func (memory *Memory) Name() string {
	return memory.name
}

func (memory *Memory) Size() uint {
	return uint(len(memory.bytes))
}

func (memory *Memory) Bytes() []byte {
	return memory.bytes
}

func (memory *Memory) Range(offset uint, length uint) []byte {
	i := offset
	j := offset + length
	if size := uint(len(memory.bytes)); j > size {
		panic(fmt.Errorf("BUG: out of range: (*Memory).Range(%d, %d) lies beyond [0,%d)", offset, length, size))
	}
	return memory.bytes[i:j]
}

func (memory *Memory) WriteLocker() sync.Locker {
	return &memory.mu
}

func (memory *Memory) ReadLocker() sync.Locker {
	return memory.mu.RLocker()
}

func (memory *Memory) WithWriteLock(offset uint, length uint, fn func([]byte) error) error {
	memory.mu.Lock()
	defer memory.mu.Unlock()
	return fn(memory.Range(offset, length))
}

func (memory *Memory) WithReadLock(offset uint, length uint, fn func([]byte) error) error {
	memory.mu.RLock()
	defer memory.mu.RUnlock()
	return fn(memory.Range(offset, length))
}

func (memory *Memory) String() string {
	return fmt.Sprintf("memory %q", memory.Name())
}

func (memory *Memory) GoString() string {
	return fmt.Sprintf("Memory(%q+%#x)", memory.Name(), memory.Size())
}

var _ fmt.Stringer = (*Memory)(nil)
var _ fmt.GoStringer = (*Memory)(nil)

// }}}
