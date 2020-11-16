package memory

import (
	"fmt"
	"reflect"
	"unsafe"
)

// UInt32Span
// {{{

type UInt32Span struct {
	mem        *Memory
	i          uint
	j          uint
	alignShift uint
}

func (span UInt32Span) String() string {
	return fmt.Sprintf("memory %q span [%d:%d] shift=%d", span.mem.Name(), span.i, span.j, span.alignShift)
}

func (span UInt32Span) GoString() string {
	return fmt.Sprintf("UInt32Span(%q, %#x, %#x, %d)", span.mem.Name(), span.i, span.j, span.alignShift)
}

func (span UInt32Span) Memory() *Memory {
	return span.mem
}

func (span UInt32Span) StartOffset() uint {
	return span.i
}

func (span UInt32Span) EndOffset() uint {
	return span.j
}

func (span UInt32Span) AlignShift() uint {
	return span.alignShift
}

func (span UInt32Span) AlignBytes() uint {
	return uint(1) << span.alignShift
}

func (span UInt32Span) Size() uint {
	return (span.j - span.i) >> 2
}

func (span UInt32Span) Absolute() (*Memory, uint, uint) {
	return span.mem, span.i, span.j
}

func (span UInt32Span) UInt8s() UInt8Span {
	return UInt8Span{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt32Span) UInt16s() UInt16Span {
	return UInt16Span{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt32Span) UInt32s() UInt32Span {
	return span
}

func (span UInt32Span) UInt64s() UInt64Span {
	checkCast("UInt32Span", "UInt64Span", 3, span.alignShift)
	return UInt64Span{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt32Span) Pages() PageSpan {
	checkCast("UInt32Span", "PageSpan", 12, span.alignShift)
	return PageSpan{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt32Span) Span(i, j uint) UInt32Span {
	checkIJ(i, j, span.Size())
	i = span.i + (i << 2)
	j = span.i + (j << 2)

	alignShift := span.alignShift
	for alignShift > 2 {
		alignSize := uint(1) << alignShift
		alignMask := alignSize - 1
		if (i & alignMask) == 0 {
			break
		}
		alignShift--
	}

	return UInt32Span{span.mem, i, j, alignShift}
}

func (span UInt32Span) AllWithWriteLock(fn func([]uint32) error) error {
	return span.WithWriteLock(0, span.Size(), fn)
}

func (span UInt32Span) AllWithReadLock(fn func([]uint32) error) error {
	return span.WithReadLock(0, span.Size(), fn)
}

func (span UInt32Span) WithWriteLock(i uint, j uint, fn func([]uint32) error) error {
	checkIJ(i, j, span.Size())
	i = span.i + (i << 2)
	j = span.i + (j << 2)

	return span.mem.withWriteLockImpl(i, j, func(bytes []byte) error {
		var data []uint32
		byteSliceToUInt32(&data, bytes)
		return fn(data)
	})
}

func (span UInt32Span) WithReadLock(i uint, j uint, fn func([]uint32) error) error {
	checkIJ(i, j, span.Size())
	i = span.i + (i << 2)
	j = span.i + (j << 2)

	return span.mem.withReadLockImpl(i, j, func(bytes []byte) error {
		var data []uint32
		byteSliceToUInt32(&data, bytes)
		return fn(data)
	})
}

func (span UInt32Span) Zero() {
	_ = span.AllWithWriteLock(func(data []uint32) error {
		for i := range data {
			data[i] = 0
		}
		return nil
	})
}

func (span UInt32Span) IsZero() bool {
	result := true
	_ = span.AllWithReadLock(func(data []uint32) error {
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

var _ fmt.Stringer = UInt32Span{}
var _ fmt.GoStringer = UInt32Span{}

// }}}

func byteSliceToUInt32(out *[]uint32, in []byte) {
	hdrIn := (*reflect.SliceHeader)(unsafe.Pointer(&in))
	hdrOut := (*reflect.SliceHeader)(unsafe.Pointer(out))
	hdrOut.Data = hdrIn.Data
	hdrOut.Len = int(uint(hdrIn.Len) >> 2)
	hdrOut.Cap = int(uint(hdrIn.Cap) >> 2)
}
