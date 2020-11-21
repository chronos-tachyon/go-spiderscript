package allocator

import (
	"fmt"
	"sync/atomic"

	"github.com/chronos-tachyon/go-spiderscript/memory"
)

func (alloc *Threaded) directAlloc(length uint) memory.UInt8Span {
	directID := atomic.AddUint32(&alloc.directNextID, 1) - 1
	directName := fmt.Sprintf("%s-direct-%d", alloc.mem.Name(), directID)
	directHuge := alloc.mem.HugePagesMode()
	directMem := memory.New(directName, directHuge, true)
	directMem.SetLen(length)
	alloc.directMap.Store(directMem, threadedDirectData{directID, length})
	return directMem.UInt8s()
}
