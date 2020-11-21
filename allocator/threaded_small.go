package allocator

import (
	"fmt"
	"runtime"
	"sort"

	"github.com/chronos-tachyon/go-spiderscript/memory"
)

func (alloc *Threaded) smallAlloc(length uint, alignShift uint) memory.UInt8Span {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var facts threadedSmallAllocFacts
	facts.Init(alloc, length, alignShift)

	facts.threadLocal.spin.Lock()
	defer facts.threadLocal.spin.Unlock()

	spanStart, found := alloc.smallAllocTryLocal(facts)
	if !found {
		alloc.mu.Lock()
		defer alloc.mu.Unlock()
		spanStart = alloc.smallAllocSteal(facts)
	}
	return alloc.mem.UInt8s().Span(spanStart, spanStart+length)
}

func (alloc *Threaded) smallAllocTryLocal(facts threadedSmallAllocFacts) (uint, bool) {
	// must hold facts.threadLocal.spin

	return grabFromFreeList(facts.myPtr, facts.allocCount, facts.allocBytes)
}

func (alloc *Threaded) smallAllocSteal(facts threadedSmallAllocFacts) uint {
	// must hold facts.threadLocal.spin
	// must hold alloc.mu

	stealStart, found := grabFromFreeList(facts.sharedPtr, facts.stealCount, facts.stealBytes)
	if !found {
		stealStart = alloc.smallAllocGrow(facts)
	}

	if facts.stealCount > facts.allocCount {
		remainStart := stealStart + facts.allocBytes
		remainCount := facts.stealCount - facts.allocCount
		list := *facts.myPtr
		list = append(list, freeListRun{remainStart, remainCount})
		sort.Sort(list)
		*facts.myPtr = list
	}

	return stealStart
}

func (alloc *Threaded) smallAllocGrow(facts threadedSmallAllocFacts) uint {
	// must hold facts.threadLocal.spin
	// must hold alloc.mu

	var largeFacts threadedLargeAllocFacts
	largeFacts.InitForSteal(facts.growCount)
	growStart := alloc.largeAllocLocked(largeFacts)
	if facts.chunksPerGrow > facts.stealCount {
		remainStart := growStart + facts.stealBytes
		remainCount := facts.chunksPerGrow - facts.stealCount
		list := *facts.sharedPtr
		list = append(list, freeListRun{remainStart, remainCount})
		sort.Sort(list)
		*facts.sharedPtr = list
	}
	return growStart
}

func (facts *threadedSmallAllocFacts) Init(alloc *Threaded, length uint, alignShift uint) {
	classIndex := computeThreadedSmallAllocClass(length, alignShift)
	class := threadedClassData[classIndex]
	bytesPerChunk := uint(class.chunkSize) << class.alignShift
	allocCount := (length + bytesPerChunk - 1) / bytesPerChunk
	allocBytes := allocCount * bytesPerChunk
	stealCount := uint(class.chunksToGrab)
	stealBytes := stealCount * bytesPerChunk
	growCount := uint(class.pagesToGrab)
	growBytes := growCount << pageShift
	chunksPerGrow := growBytes / bytesPerChunk

	if allocCount > stealCount {
		panic(fmt.Errorf("BUG: allocCount=%d, stealCount=%d", allocCount, stealCount))
	}

	if stealCount > chunksPerGrow {
		panic(fmt.Errorf("BUG: chunksToGrab=%d, pagesToGrab=%d, chunksPerGrow=%d", stealCount, growCount, chunksPerGrow))
	}

	tid := getThreadID()
	threadLocal := alloc.getPerThreadData(tid)

	*facts = threadedSmallAllocFacts{
		length:        length,
		alignShift:    alignShift,
		classIndex:    classIndex,
		allocCount:    allocCount,
		allocBytes:    allocBytes,
		stealCount:    stealCount,
		stealBytes:    stealBytes,
		growCount:     growCount,
		growBytes:     growBytes,
		chunksPerGrow: chunksPerGrow,
		tid:           tid,
		threadLocal:   threadLocal,
		myPtr:         &threadLocal.freeChunksByClass[classIndex],
		sharedPtr:     &alloc.freeChunksByClass[classIndex],
	}
}
