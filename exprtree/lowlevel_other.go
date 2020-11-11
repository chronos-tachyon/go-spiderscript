// +build !linux

package exprtree

import (
	"fmt"
	"reflect"
	"unsafe"
)

func malloc(slice *[]byte, length uint, hugePages HugePagesMode, isLocked bool) {
	if slice == nil {
		panic(fmt.Errorf("BUG: *[]byte is nil"))
	}

	pageSize := hugePages.PageSize()
	newSize := mallocAlign(uintptr(length), pageSize)

	const maxInt = uintptr((^uint(0)) >> 1)
	if newSize > maxInt {
		panic(fmt.Errorf("BUG: cannot allocate a single block of %d bytes", newSize))
	}

	var newSlice []byte
	if newSize != 0 {
		newSlice = make([]byte, newSize+pageSize)
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&newSlice))
		p := hdr.Data
		q := mallocAlign(p, pageSize)
		i := (q - p)
		j := i + newSize
		newSlice = newSlice[i:j]
	}

	oldSlice := *slice
	if oldSlice != nil && newSlice != nil {
		limit := uintptr(cap(oldSlice))
		if limit > newSize {
			limit = newSize
		}
		copy(newSlice[:limit], oldSlice[:limit])
	}

	if newSlice != nil && uintptr(length) != newSize {
		newSlice = newSlice[:length]
	}

	*slice = newSlice
}

func mprotect(slice []byte, r, w, x bool) error {
	if slice == nil {
		return nil
	}

	return ErrNotImplemented
}

func mlock(slice []byte, acquire bool) error {
	if slice == nil {
		return nil
	}

	return ErrNotImplemented
}
