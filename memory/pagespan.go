package memory

import (
	"fmt"
	"reflect"
	"unsafe"
)

// PageSpan
// {{{

type PageSpan struct {
	mem        *Memory
	i          uint
	j          uint
	alignShift uint
}

func (span PageSpan) String() string {
	return fmt.Sprintf("memory %q span [%d:%d] shift=%d", span.mem.Name(), span.i, span.j, span.alignShift)
}

func (span PageSpan) GoString() string {
	return fmt.Sprintf("PageSpan(%q, %#x, %#x, %d)", span.mem.Name(), span.i, span.j, span.alignShift)
}

func (span PageSpan) Memory() *Memory {
	return span.mem
}

func (span PageSpan) StartOffset() uint {
	return span.i
}

func (span PageSpan) EndOffset() uint {
	return span.j
}

func (span PageSpan) AlignShift() uint {
	return span.alignShift
}

func (span PageSpan) AlignBytes() uint {
	return uint(1) << span.alignShift
}

func (span PageSpan) Size() uint {
	return (span.j - span.i) >> 12
}

func (span PageSpan) Absolute() (*Memory, uint, uint) {
	return span.mem, span.i, span.j
}

func (span PageSpan) UInt8s() UInt8Span {
	return UInt8Span{span.mem, span.i, span.j, span.alignShift}
}

func (span PageSpan) UInt16s() UInt16Span {
	return UInt16Span{span.mem, span.i, span.j, span.alignShift}
}

func (span PageSpan) UInt32s() UInt32Span {
	return UInt32Span{span.mem, span.i, span.j, span.alignShift}
}

func (span PageSpan) UInt64s() UInt64Span {
	return UInt64Span{span.mem, span.i, span.j, span.alignShift}
}

func (span PageSpan) Pages() PageSpan {
	return span
}

func (span PageSpan) Span(i, j uint) PageSpan {
	checkIJ(i, j, span.Size())
	i = span.i + (i << 12)
	j = span.i + (j << 12)

	alignShift := span.alignShift
	for alignShift > 12 {
		alignSize := uint(1) << alignShift
		alignMask := alignSize - 1
		if (i & alignMask) == 0 {
			break
		}
		alignShift--
	}

	return PageSpan{span.mem, i, j, alignShift}
}

func (span PageSpan) AllWithWriteLock(fn func([][4096]byte) error) error {
	return span.WithWriteLock(0, span.Size(), fn)
}

func (span PageSpan) AllWithReadLock(fn func([][4096]byte) error) error {
	return span.WithReadLock(0, span.Size(), fn)
}

func (span PageSpan) WithWriteLock(i uint, j uint, fn func([][4096]byte) error) error {
	checkIJ(i, j, span.Size())
	i = span.i + (i << 12)
	j = span.i + (j << 12)

	return span.mem.withWriteLockImpl(i, j, func(bytes []byte) error {
		var data [][4096]byte
		byteSliceToPage(&data, bytes)
		return fn(data)
	})
}

func (span PageSpan) WithReadLock(i uint, j uint, fn func([][4096]byte) error) error {
	checkIJ(i, j, span.Size())
	i = span.i + (i << 12)
	j = span.i + (j << 12)

	return span.mem.withReadLockImpl(i, j, func(bytes []byte) error {
		var data [][4096]byte
		byteSliceToPage(&data, bytes)
		return fn(data)
	})
}

func (span PageSpan) Zero() {
	_ = span.AllWithWriteLock(func(data [][4096]byte) error {
		for i := range data {
			data[i] = [4096]byte{}
		}
		return nil
	})
}

func (span PageSpan) IsZero() bool {
	zeroPage := [4096]byte{}
	result := true
	_ = span.AllWithReadLock(func(data [][4096]byte) error {
		for i := range data {
			if data[i] != zeroPage {
				result = false
				break
			}
		}
		return nil
	})
	return result
}

var _ fmt.Stringer = PageSpan{}
var _ fmt.GoStringer = PageSpan{}

// }}}

func byteSliceToPage(out *[][4096]byte, in []byte) {
	hdrIn := (*reflect.SliceHeader)(unsafe.Pointer(&in))
	hdrOut := (*reflect.SliceHeader)(unsafe.Pointer(out))
	hdrOut.Data = hdrIn.Data
	hdrOut.Len = int(uint(hdrIn.Len) >> 12)
	hdrOut.Cap = int(uint(hdrIn.Cap) >> 12)
}
