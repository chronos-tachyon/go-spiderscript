package memory

import (
	"fmt"
	"reflect"
	"unsafe"
)

// UInt64Span
// {{{

type UInt64Span struct {
	mem        *Memory
	i          uint
	j          uint
	alignShift uint
}

func (span UInt64Span) String() string {
	return fmt.Sprintf("memory %q span [%d:%d] shift=%d", span.mem.Name(), span.i, span.j, span.alignShift)
}

func (span UInt64Span) GoString() string {
	return fmt.Sprintf("UInt64Span(%q, %#x, %#x, %d)", span.mem.Name(), span.i, span.j, span.alignShift)
}

func (span UInt64Span) Memory() *Memory {
	return span.mem
}

func (span UInt64Span) StartOffset() uint {
	return span.i
}

func (span UInt64Span) EndOffset() uint {
	return span.j
}

func (span UInt64Span) AlignShift() uint {
	return span.alignShift
}

func (span UInt64Span) AlignBytes() uint {
	return uint(1) << span.alignShift
}

func (span UInt64Span) Size() uint {
	return (span.j - span.i) >> 3
}

func (span UInt64Span) Absolute() (*Memory, uint, uint) {
	return span.mem, span.i, span.j
}

func (span UInt64Span) UInt8s() UInt8Span {
	return UInt8Span{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt64Span) UInt16s() UInt16Span {
	return UInt16Span{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt64Span) UInt32s() UInt32Span {
	return UInt32Span{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt64Span) UInt64s() UInt64Span {
	return span
}

func (span UInt64Span) Pages() PageSpan {
	checkCast("UInt64Span", "PageSpan", 12, span.alignShift)
	return PageSpan{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt64Span) Span(i, j uint) UInt64Span {
	checkIJ(i, j, span.Size())
	i = span.i + (i << 3)
	j = span.i + (j << 3)

	alignShift := span.alignShift
	for alignShift > 3 {
		alignSize := uint(1) << alignShift
		alignMask := alignSize - 1
		if (i & alignMask) == 0 {
			break
		}
		alignShift--
	}

	return UInt64Span{span.mem, i, j, alignShift}
}

func (span UInt64Span) AllWithWriteLock(fn func([]uint64) error) error {
	return span.WithWriteLock(0, span.Size(), fn)
}

func (span UInt64Span) AllWithReadLock(fn func([]uint64) error) error {
	return span.WithReadLock(0, span.Size(), fn)
}

func (span UInt64Span) WithWriteLock(i uint, j uint, fn func([]uint64) error) error {
	checkIJ(i, j, span.Size())
	i = span.i + (i << 3)
	j = span.i + (j << 3)

	return span.mem.withWriteLockImpl(i, j, func(bytes []byte) error {
		var data []uint64
		byteSliceToUInt64(&data, bytes)
		return fn(data)
	})
}

func (span UInt64Span) WithReadLock(i uint, j uint, fn func([]uint64) error) error {
	checkIJ(i, j, span.Size())
	i = span.i + (i << 3)
	j = span.i + (j << 3)

	return span.mem.withReadLockImpl(i, j, func(bytes []byte) error {
		var data []uint64
		byteSliceToUInt64(&data, bytes)
		return fn(data)
	})
}

func (span UInt64Span) Zero() {
	_ = span.AllWithWriteLock(func(data []uint64) error {
		for i := range data {
			data[i] = 0
		}
		return nil
	})
}

func (span UInt64Span) IsZero() bool {
	result := true
	_ = span.AllWithReadLock(func(data []uint64) error {
		for i := range data {
			if data[i] != 0 {
				result = false
				break
			}
		}
		return nil
	})
	return result
}

var _ fmt.Stringer = UInt64Span{}
var _ fmt.GoStringer = UInt64Span{}

// }}}

func byteSliceToUInt64(out *[]uint64, in []byte) {
	hdrIn := (*reflect.SliceHeader)(unsafe.Pointer(&in))
	hdrOut := (*reflect.SliceHeader)(unsafe.Pointer(out))
	hdrOut.Data = hdrIn.Data
	hdrOut.Len = int(uint(hdrIn.Len) >> 3)
	hdrOut.Cap = int(uint(hdrIn.Cap) >> 3)
}
