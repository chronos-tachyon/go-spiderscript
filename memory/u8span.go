package memory

import (
	"fmt"
)

// UInt8Span
// {{{

type UInt8Span struct {
	mem        *Memory
	i          uint
	j          uint
	alignShift uint
}

func (span UInt8Span) String() string {
	return fmt.Sprintf("memory %q span [%d:%d] shift=%d", span.mem.Name(), span.i, span.j, span.alignShift)
}

func (span UInt8Span) GoString() string {
	return fmt.Sprintf("UInt8Span(%q, %#x, %#x, %d)", span.mem.Name(), span.i, span.j, span.alignShift)
}

func (span UInt8Span) Memory() *Memory {
	return span.mem
}

func (span UInt8Span) StartOffset() uint {
	return span.i
}

func (span UInt8Span) EndOffset() uint {
	return span.j
}

func (span UInt8Span) AlignShift() uint {
	return span.alignShift
}

func (span UInt8Span) AlignBytes() uint {
	return uint(1) << span.alignShift
}

func (span UInt8Span) Size() uint {
	return span.j - span.i
}

func (span UInt8Span) Absolute() (*Memory, uint, uint) {
	return span.mem, span.i, span.j
}

func (span UInt8Span) UInt8s() UInt8Span {
	return span
}

func (span UInt8Span) UInt16s() UInt16Span {
	checkCast("UInt8Span", "UInt16Span", 1, span.alignShift)
	return UInt16Span{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt8Span) UInt32s() UInt32Span {
	checkCast("UInt8Span", "UInt32Span", 2, span.alignShift)
	return UInt32Span{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt8Span) UInt64s() UInt64Span {
	checkCast("UInt8Span", "UInt64Span", 3, span.alignShift)
	return UInt64Span{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt8Span) Pages() PageSpan {
	checkCast("UInt8Span", "PageSpan", 12, span.alignShift)
	return PageSpan{span.mem, span.i, span.j, span.alignShift}
}

func (span UInt8Span) Span(i, j uint) UInt8Span {
	size := span.Size()
	checkIJ(i, j, size)

	i += span.i
	j += span.i

	alignShift := span.alignShift
	for alignShift > 0 {
		alignSize := uint(1) << alignShift
		alignMask := alignSize - 1
		if (i & alignMask) == 0 {
			break
		}
		alignShift--
	}

	return UInt8Span{span.mem, i, j, alignShift}
}

func (span UInt8Span) AllWithWriteLock(fn func([]byte) error) error {
	return span.WithWriteLock(0, span.Size(), fn)
}

func (span UInt8Span) AllWithReadLock(fn func([]byte) error) error {
	return span.WithReadLock(0, span.Size(), fn)
}

func (span UInt8Span) WithWriteLock(i uint, j uint, fn func([]byte) error) error {
	checkIJ(i, j, span.Size())
	i += span.i
	j += span.i
	return span.mem.withWriteLockImpl(i, j, fn)
}

func (span UInt8Span) WithReadLock(i uint, j uint, fn func([]byte) error) error {
	checkIJ(i, j, span.Size())
	i += span.i
	j += span.i
	return span.mem.withReadLockImpl(i, j, fn)
}

func (span UInt8Span) Zero() {
	_ = span.AllWithWriteLock(func(data []byte) error {
		for i := range data {
			data[i] = 0
		}
		return nil
	})
}

func (span UInt8Span) IsZero() bool {
	result := true
	_ = span.AllWithReadLock(func(data []byte) error {
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

var _ fmt.Stringer = UInt8Span{}
var _ fmt.GoStringer = UInt8Span{}

// }}}
