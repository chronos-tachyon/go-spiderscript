package allocator

import (
	"fmt"
	"sort"

	"github.com/chronos-tachyon/go-spiderscript/memory"
)

func (alloc *Threaded) largeAlloc(length uint) memory.UInt8Span {
	alloc.mu.Lock()
	defer alloc.mu.Unlock()

	var facts threadedLargeAllocFacts
	facts.Init(length)
	allocStart := alloc.largeAllocLocked(facts)
	return alloc.mem.UInt8s().Span(allocStart, allocStart+length)
}

func (alloc *Threaded) largeAllocLocked(facts threadedLargeAllocFacts) uint {
	// must hold alloc.mu

	allocStart, found := grabFromFreeList(&alloc.freePages, facts.allocCount, facts.allocBytes)
	if !found {
		alloc.mem.Grow(hugePageSize)
		allocStart = alloc.sbrk
		alloc.sbrk += hugePageSize
		if facts.growCount > facts.allocCount {
			remainStart := allocStart + facts.allocBytes
			remainPages := facts.growCount - facts.allocCount
			list := alloc.freePages
			list = append(list, freeListRun{remainStart, remainPages})
			sort.Sort(list)
			alloc.freePages = list
		}
	}
	return allocStart
}

func (facts *threadedLargeAllocFacts) Init(length uint) {
	allocCount := (length + pageMask) >> pageShift
	allocBytes := allocCount << pageShift
	const growCount = uint(1) << (hugePageShift - pageShift)
	const growBytes = growCount >> pageShift

	if allocCount > growCount {
		panic(fmt.Errorf("BUG: allocCount=%d, growCount=%d", allocCount, growCount))
	}

	*facts = threadedLargeAllocFacts{
		length:     length,
		allocCount: allocCount,
		allocBytes: allocBytes,
		growCount:  growCount,
		growBytes:  growBytes,
		external:   true,
	}
}

func (facts *threadedLargeAllocFacts) InitForSteal(allocCount uint) {
	allocBytes := allocCount << pageShift
	const growCount = uint(1) << (hugePageShift - pageShift)
	const growBytes = growCount >> pageShift

	if allocCount > growCount {
		panic(fmt.Errorf("BUG: allocCount=%d, growCount=%d", allocCount, growCount))
	}

	*facts = threadedLargeAllocFacts{
		length:     0,
		allocCount: allocCount,
		allocBytes: allocBytes,
		growCount:  growCount,
		growBytes:  growBytes,
		external:   false,
	}
}
