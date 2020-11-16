package allocator

import (
	"github.com/chronos-tachyon/go-spiderscript/memory"
)

type Allocator interface {
	Allocate(count uint, alignShift uint) memory.UInt8Span
	Free(span memory.UInt8Span)
	Trim()
	FreeAll()
}
