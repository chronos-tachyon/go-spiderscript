package exprtree

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestMemory(t *testing.T) {
	modes := make([]HugePagesMode, 0, 3)
	modes = append(modes, HugePagesOff)

	for _, mode := range modes {
		t.Run(mode.GoString(), func(t *testing.T) {
			mem := NewMemory("test", mode)
			mem.Truncate(17)

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

			expectedSliceLen := 17
			actualSliceLen := len(mem.bytes)
			if expectedSliceLen != actualSliceLen {
				t.Errorf("(*Memory).bytes.len(): expected %d, actual %d", expectedSliceLen, actualSliceLen)
			}

			expectedSliceCap := 4096
			actualSliceCap := cap(mem.bytes)
			if expectedSliceCap != actualSliceCap {
				t.Errorf("(*Memory).bytes.cap(): expected %d, actual %d", expectedSliceCap, actualSliceCap)
			}

			if err := mem.ProtectPages(true, false, false); err != nil {
				t.Errorf("(*Memory).ProtectPages(T,F,F): %w", err)
			}

			if err := mem.LockPages(); err != nil {
				t.Errorf("(*Memory).LockPages(): %w", err)
			}

			if err := mem.UnlockPages(); err != nil {
				t.Errorf("(*Memory).UnlockPages(): %w", err)
			}
		})
	}
}

func formatByteSlice(slice []byte) string {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	return fmt.Sprintf("{ptr=%#x, len=%#x, cap=%#x}", hdr.Data, uint(hdr.Len), uint(hdr.Cap))
}
