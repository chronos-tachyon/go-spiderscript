package allocator

import (
	"fmt"
	"sync"

	"github.com/chronos-tachyon/go-spiderscript/memory"
)

const (
	threadedLargeThreshold  = pageSize << 3 // 32KiB = 8 pages
	threadedDirectThreshold = hugePageSize  // 2MiB = 1 hugepage
)

type ThreadedOptions struct {
	Name      string
	HugePages memory.HugePagesMode
}

type Threaded struct {
	mem          *memory.Memory
	threadMap    sync.Map
	directMap    sync.Map
	directNextID uint32

	mu                sync.Mutex
	sbrk              uint
	freePages         freeList
	freeChunksByClass [threadedNumClasses]freeList
}

type threadedPerThreadData struct {
	spin              Spinlock
	freeChunksByClass [threadedNumClasses]freeList
}

type threadedDirectData struct {
	id     uint32
	length uint
}

type threadedLargeAllocFacts struct {
	length     uint
	allocCount uint
	allocBytes uint
	growCount  uint
	growBytes  uint
	external   bool
}

type threadedSmallAllocFacts struct {
	length        uint
	alignShift    uint
	classIndex    uint
	allocCount    uint
	allocBytes    uint
	stealCount    uint
	stealBytes    uint
	growCount     uint
	growBytes     uint
	chunksPerGrow uint
	tid           uint
	threadLocal   *threadedPerThreadData
	myPtr         *freeList
	sharedPtr     *freeList
}

func NewThreaded(opts ThreadedOptions) *Threaded {
	alloc := new(Threaded)
	alloc.Init(opts)
	return alloc
}

func (alloc *Threaded) Init(opts ThreadedOptions) {
	*alloc = Threaded{
		mem: memory.New(opts.Name, opts.HugePages, true),
	}
}

func (alloc *Threaded) Allocate(count uint, alignShift uint) memory.UInt8Span {
	if alignShift > pageShift {
		panic(fmt.Errorf("alignment is too strict: requested=%d, max=%d", alignShift, pageShift))
	}

	length := count << alignShift
	if length >= threadedDirectThreshold {
		return alloc.directAlloc(length)
	} else if length >= threadedLargeThreshold {
		return alloc.largeAlloc(length)
	} else {
		return alloc.smallAlloc(length, alignShift)
	}
}

func (alloc *Threaded) Free(span memory.UInt8Span) {
	panic(fmt.Errorf("not implemented"))
}

func (alloc *Threaded) Trim() {
	panic(fmt.Errorf("not implemented"))
}

func (alloc *Threaded) FreeAll() {
	panic(fmt.Errorf("not implemented"))
}

func (alloc *Threaded) getPerThreadData(tid uint) *threadedPerThreadData {
	if iface, found := alloc.threadMap.Load(tid); found {
		return iface.(*threadedPerThreadData)
	}

	newData := new(threadedPerThreadData)
	if iface, found := alloc.threadMap.LoadOrStore(tid, newData); found {
		return iface.(*threadedPerThreadData)
	}

	return newData
}

var _ Allocator = (*Threaded)(nil)
