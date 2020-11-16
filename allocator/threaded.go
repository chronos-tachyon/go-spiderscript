package allocator

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/chronos-tachyon/go-spiderscript/memory"
)

const (
	threadedLargeThreshold  = pageSize << 3 // 32KiB = 8 pages
	threadedDirectThreshold = hugePageSize  // 2MiB = 1 hugepage

	threadedNumFixedSizeFreePageLists = 255
)

type ThreadedOptions struct {
	Name      string
	HugePages memory.HugePagesMode
}

type Threaded struct {
	mem        *memory.Memory
	threadMap  sync.Map
	directMap  sync.Map
	directNext uint32

	mu                   sync.Mutex
	fixedSizeFreeLists   [threadedNumFixedSizeFreePageLists][]uint
	variableSizeFreeList []threadedLargeFreePageRun
	sbrk                 uint
}

type threadedPerThreadData struct {
}

type threadedLargeFreePageRun struct {
	start uint
	count uint
}

type threadedLargeFreePageRunList []threadedLargeFreePageRun

func (list threadedLargeFreePageRunList) Len() int {
	return len(list)
}

func (list threadedLargeFreePageRunList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list threadedLargeFreePageRunList) Less(i, j int) bool {
	a, b := list[i], list[j]

	if a.count != b.count {
		return a.count < b.count
	}

	return a.start < b.start
}

var _ sort.Interface = threadedLargeFreePageRunList(nil)

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
		return alloc.allocateDirect(length)
	} else if length >= threadedLargeThreshold {
		return alloc.allocateLarge(length)
	} else {
		return alloc.allocateSmall(length, alignShift)
	}
}

func (alloc *Threaded) allocateDirect(length uint) memory.UInt8Span {
	memCounter := atomic.AddUint32(&alloc.directNext, 1) - 1
	memName := fmt.Sprintf("%s-direct-%d", alloc.mem.Name(), memCounter)
	memHuge := alloc.mem.HugePagesMode()
	mem := memory.New(memName, memHuge, true)
	alloc.directMap.Store(mem, struct{}{})
	mem.SetLen(length)
	return mem.UInt8s()
}

func (alloc *Threaded) allocateLarge(length uint) memory.UInt8Span {
	reqPages := (length + pageMask) & ^uint(pageMask)

	alloc.mu.Lock()
	defer alloc.mu.Unlock()

	var spanFound bool
	var spanStart uint
	var spanPages uint
	var vsflDirty bool

	for i := reqPages; i <= threadedNumFixedSizeFreePageLists; i++ {
		freeList := alloc.fixedSizeFreeLists[i]
		n := uint(len(freeList))
		if n != 0 {
			n--
			spanStart = freeList[n]
			spanPages = i
			spanFound = true
			alloc.fixedSizeFreeLists[i] = freeList[:n]
			break
		}
	}

	if !spanFound {
		freeList := alloc.variableSizeFreeList
		n := len(freeList)
		i := sort.Search(n, func(i int) bool {
			return freeList[i].count >= reqPages
		})
		if i < n {
			run := freeList[i]
			spanStart = run.start
			spanPages = run.count
			spanFound = true

			n--
			freeList[i] = freeList[n]
			freeList = freeList[:n]
			alloc.variableSizeFreeList = freeList
			vsflDirty = true
		}
	}

	if !spanFound {
		alloc.mem.Grow(hugePageSize)
		spanStart = alloc.sbrk
		spanPages = (1 << (hugePageShift - pageShift))
		spanFound = true
		alloc.sbrk += hugePageSize
	}

	if spanPages > reqPages {
		returnStart := spanStart + (reqPages << pageShift)
		returnPages := spanPages - reqPages
		if returnPages <= threadedNumFixedSizeFreePageLists {
			freeList := alloc.fixedSizeFreeLists[returnPages]
			freeList = append(freeList, returnStart)
			alloc.fixedSizeFreeLists[returnPages] = freeList
		} else {
			freeList := alloc.variableSizeFreeList
			freeList = append(freeList, threadedLargeFreePageRun{returnStart, returnPages})
			alloc.variableSizeFreeList = freeList
			vsflDirty = true
		}
	}

	if vsflDirty {
		sort.Sort(threadedLargeFreePageRunList(alloc.variableSizeFreeList))
	}

	return alloc.mem.UInt8s().Span(spanStart, spanStart+length)
}

func (alloc *Threaded) allocateSmall(length uint, alignShift uint) memory.UInt8Span {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	tid := uint(syscall.Gettid())
	iface, found := alloc.threadMap.Load(tid)

	var data *threadedPerThreadData
	if found {
		data = iface.(*threadedPerThreadData)
	} else {
		data = new(threadedPerThreadData)
		alloc.threadMap.Store(tid, data)
	}

	classIndex := computeThreadedSmallAllocClass(length, alignShift)
	classData := threadedSmallAllocClassDataTable[classIndex]
	_ = classData

	panic(fmt.Errorf("not implemented"))
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

var _ Allocator = (*Threaded)(nil)

type threadedSmallAllocClassDataRow struct {
	alignShift uint16
	chunkSize  uint16
}

var threadedSmallAllocClassDataTable = []threadedSmallAllocClassDataRow{
	{
		alignShift: 0,
		chunkSize:  1,
	},
	{
		alignShift: 1,
		chunkSize:  1,
	},
	{
		alignShift: 2,
		chunkSize:  1,
	},
	{
		alignShift: 3,
		chunkSize:  1,
	},
	{
		alignShift: 4,
		chunkSize:  1,
	},
	{
		alignShift: 5,
		chunkSize:  1,
	},
	{
		alignShift: 6,
		chunkSize:  1,
	},
	{
		alignShift: 7,
		chunkSize:  1,
	},
	{
		alignShift: 8,
		chunkSize:  1,
	},
	{
		alignShift: 9,
		chunkSize:  1,
	},
	{
		alignShift: 10,
		chunkSize:  1,
	},
	{
		alignShift: 11,
		chunkSize:  1,
	},
	{
		alignShift: 12,
		chunkSize:  1,
	},
	{
		alignShift: 12,
		chunkSize:  2,
	},
	{
		alignShift: 12,
		chunkSize:  4,
	},
	{
		alignShift: 12,
		chunkSize:  8,
	},
}

type threadedSmallAllocClassMatchRow struct {
	minLength  uint16
	maxLength  uint16
	alignShift uint16
	allocClass uint16
}

var threadedSmallAllocClassMatchTable = []threadedSmallAllocClassMatchRow{
	{
		// 0B - 1B
		minLength:  0,
		maxLength:  1 << 0,
		alignShift: 0,
		allocClass: 0,
	},
	{
		// 0B - 2B
		minLength:  0,
		maxLength:  1 << 1,
		alignShift: 1,
		allocClass: 1,
	},
	{
		// 0B - 4B
		minLength:  0,
		maxLength:  1 << 2,
		alignShift: 2,
		allocClass: 2,
	},
	{
		// 0B - 8B
		minLength:  0,
		maxLength:  1 << 3,
		alignShift: 3,
		allocClass: 3,
	},
	{
		// 0B - 16B
		minLength:  0,
		maxLength:  1 << 4,
		alignShift: 4,
		allocClass: 4,
	},
	{
		// 0B - 32B
		minLength:  0,
		maxLength:  1 << 5,
		alignShift: 5,
		allocClass: 5,
	},
	{
		// 0B - 64B
		minLength:  0,
		maxLength:  1 << 6,
		alignShift: 6,
		allocClass: 6,
	},
	{
		// 0B - 128B
		minLength:  0,
		maxLength:  1 << 7,
		alignShift: 7,
		allocClass: 7,
	},
	{
		// 0B - 256B
		minLength:  0,
		maxLength:  1 << 8,
		alignShift: 8,
		allocClass: 8,
	},
	{
		// 0B - 512B
		minLength:  0,
		maxLength:  1 << 9,
		alignShift: 9,
		allocClass: 9,
	},
	{
		// 0B - 1KiB
		minLength:  0,
		maxLength:  1 << 10,
		alignShift: 10,
		allocClass: 10,
	},
	{
		// 0B - 2KiB
		minLength:  0,
		maxLength:  1 << 11,
		alignShift: 11,
		allocClass: 11,
	},
	{
		// 0B - 4KiB
		minLength:  0,
		maxLength:  1 << 12,
		alignShift: 12,
		allocClass: 12,
	},
	{
		// 4KiB - 8KiB
		minLength:  (1 << 12) + 1,
		maxLength:  1 << 13,
		alignShift: 12,
		allocClass: 13,
	},
	{
		// 8KiB - 16KiB
		minLength:  (1 << 13) + 1,
		maxLength:  1 << 14,
		alignShift: 12,
		allocClass: 14,
	},
	{
		// 16KiB - 32KiB
		minLength:  (1 << 14) + 1,
		maxLength:  1 << 15,
		alignShift: 12,
		allocClass: 15,
	},
}

func computeThreadedSmallAllocClass(length uint, alignShift uint) uint {
	if alignShift > pageShift {
		panic(fmt.Errorf("BUG: alignShift=%d, max=%d", alignShift, pageShift))
	}
	if length >= threadedLargeThreshold {
		panic(fmt.Errorf("BUG: length=%d, threadedLargeThreshold=%d", length, threadedLargeThreshold))
	}
	lengthU16 := uint16(length) // assumes that threadedLargeThreshold fits in uint16
	alignShiftU16 := uint16(alignShift)
	for _, row := range threadedSmallAllocClassMatchTable {
		if lengthU16 >= row.minLength && lengthU16 <= row.maxLength && alignShiftU16 <= row.alignShift {
			return uint(row.allocClass)
		}
	}
	panic(fmt.Errorf("BUG: no alloc class for length=%d, alignShift=%d", length, alignShift))
}
