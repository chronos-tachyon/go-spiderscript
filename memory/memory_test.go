package memory

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestMemory(t *testing.T) {
	modes := []HugePagesMode{HugePagesOff, HugePages2M}

	for _, mode := range modes {
		t.Run(mode.GoString(), func(t *testing.T) {
			mem := NewMemory("test", mode, false)
			mem.Grow(16)
			mem.UInt64s().Zero()

			expectedString := fmt.Sprintf("memory %q", "test")
			actualString := mem.String()
			if expectedString != actualString {
				t.Errorf("(*Memory).String(): expected %q, actual %q", expectedString, actualString)
			}

			expectedGoString := fmt.Sprintf("Memory(%q)", "test")
			actualGoString := mem.GoString()
			if expectedGoString != actualGoString {
				t.Errorf("(*Memory).GoString(): expected %q, actual %q", expectedGoString, actualGoString)
			}

			expectedSliceLen := uint(16)
			actualSliceLen := uint(len(mem.bytes))
			if expectedSliceLen != actualSliceLen {
				t.Errorf("(*Memory).bytes.len(): expected %d, actual %d", expectedSliceLen, actualSliceLen)
			}

			expectedSliceCap := uint(mode.PageSize())
			actualSliceCap := uint(cap(mem.bytes))
			if expectedSliceCap != actualSliceCap {
				t.Errorf("(*Memory).bytes.cap(): expected %d, actual %d", expectedSliceCap, actualSliceCap)
			}

			if err := mem.Protect(true, false, false); err != nil {
				t.Errorf("(*Memory).Protect(T,F,F): %w", err)
			}

			if err := mem.LockToRAM(); err != nil {
				t.Errorf("(*Memory).LockToRAM(): %w", err)
			}

			if err := mem.UnlockFromRAM(); err != nil {
				t.Errorf("(*Memory).UnlockFromRAM(): %w", err)
			}
		})
	}
}

func formatByteSlice(slice []byte) string {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	return fmt.Sprintf("{ptr=%#x, len=%#x, cap=%#x}", hdr.Data, uint(hdr.Len), uint(hdr.Cap))
}
