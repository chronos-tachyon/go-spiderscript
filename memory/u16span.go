package memory

import (
	"fmt"
	"reflect"
	"unsafe"
)

// UInt16Span
// {{{

type UInt16Span struct {
	mem        *Memory
	i          uint
	j          uint
	alignShift uint
}

func (span UInt16Span) String() string {
	return fmt.Sprintf("memory %q span [%d:%d] shift=%d", span.mem.Name(), span.i, span.j, span.alignShift)
}

func (span UInt16Span) GoString() string {
	return fmt.Sprintf("UInt16Span(%q, %#x, %#x, %d)", span.mem.Name(), span.i, span.j, span.alignShift)
}

func (span UInt16Span) Memory() *Memory {
	return span.mem
}

func (span UInt16Span) StartOffset() uint {
	return span.i
}

func (span UInt16Span) EndOffset() uint {
	return span.j
}

func (span UInt16Span) AlignShift() uint {
	return span.alignShift
}

func (span UInt16Span) AlignBytes() uint {
	return uint(1) << span.alignShift
}

func (span UInt16Span) Size() uint {
	return (span.j - span.i) >> 1
}

func (span UInt16Span) Absolute() (*Memory, uint, uint) {
	return span.mem, span.i, span.j
}

func (span UInt16Span) UInt8s() UInt8Span {
	return UInt8Span{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt16Span) UInt16s() UInt16Span {
	return span
}

func (span UInt16Span) UInt32s() UInt32Span {
	checkCast("UInt16Span", "UInt32Span", 2, span.alignShift)
	return UInt32Span{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt16Span) UInt64s() UInt64Span {
	checkCast("UInt16Span", "UInt64Span", 3, span.alignShift)
	return UInt64Span{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt16Span) Pages() PageSpan {
	checkCast("UInt16Span", "PageSpan", 12, span.alignShift)
	return PageSpan{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt16Span) Span(i, j uint) UInt16Span {
	checkIJ(i, j, span.Size())
	i = span.i + (i << 1)
	j = span.i + (j << 1)

	alignShift := span.alignShift
	for alignShift > 1 {
		alignSize := uint(1) << alignShift
		alignMask := alignSize - 1
		if (i & alignMask) == 0 {
			break
		}
		alignShift--
	}

	return UInt16Span{span.mem, i, j, alignShift}
}

func (span UInt16Span) AllWithWriteLock(fn func([]uint16) error) error {
	return span.WithWriteLock(0, span.Size(), fn)
}

func (span UInt16Span) AllWithReadLock(fn func([]uint16) error) error {
	return span.WithReadLock(0, span.Size(), fn)
}

func (span UInt16Span) WithWriteLock(i uint, j uint, fn func([]uint16) error) error {
	checkIJ(i, j, span.Size())
	i = span.i + (i << 1)
	j = span.i + (j << 1)

	return span.mem.withWriteLockImpl(i, j, func(bytes []byte) error {
		var data []uint16
		byteSliceToUInt16(&data, bytes)
		return fn(data)
	})
}

func (span UInt16Span) WithReadLock(i uint, j uint, fn func([]uint16) error) error {
	checkIJ(i, j, span.Size())
	i = span.i + (i << 1)
	j = span.i + (j << 1)

	return span.mem.withReadLockImpl(i, j, func(bytes []byte) error {
		var data []uint16
		byteSliceToUInt16(&data, bytes)
		return fn(data)
	})
}

func (span UInt16Span) Zero() {
	_ = span.AllWithWriteLock(func(data []uint16) error {
		for i := range data {
			data[i] = 0
		}
		return nil
	})
}

func (span UInt16Span) IsZero() bool {
	result := true
	_ = span.AllWithReadLock(func(data []uint16) error {
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

var _ fmt.Stringer = UInt16Span{}
var _ fmt.GoStringer = UInt16Span{}

// }}}

func byteSliceToUInt16(out *[]uint16, in []byte) {
	hdrIn := (*reflect.SliceHeader)(unsafe.Pointer(&in))
	hdrOut := (*reflect.SliceHeader)(unsafe.Pointer(out))
	hdrOut.Data = hdrIn.Data
	hdrOut.Len = int(uint(hdrIn.Len) >> 1)
	hdrOut.Cap = int(uint(hdrIn.Cap) >> 1)
}
