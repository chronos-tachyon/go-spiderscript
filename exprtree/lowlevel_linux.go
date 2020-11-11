// +build linux

package exprtree

import (
	"fmt"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

const (
	_MAP_HUGE_SHIFT = 26
	_MAP_HUGE_2MB   = (21 << _MAP_HUGE_SHIFT)
	_MAP_HUGE_1GB   = (30 << _MAP_HUGE_SHIFT)
	_MREMAP_MAYMOVE = 1
)

func (mode HugePagesMode) mmapFlags() uintptr {
	switch mode {
	case HugePages2M:
		return syscall.MAP_HUGETLB | _MAP_HUGE_2MB
	case HugePages1G:
		return syscall.MAP_HUGETLB | _MAP_HUGE_1GB
	default:
		return (1 << 12)
	}
}

func malloc(slice *[]byte, length uint, hugePages HugePagesMode, isLocked bool) {
	if slice == nil {
		panic(fmt.Errorf("BUG: *[]byte is nil"))
	}

	pageSize := hugePages.PageSize()
	newSize := mallocAlign(uintptr(length), pageSize)

	const maxInt = uintptr((^uint(0)) >> 1)
	if newSize > maxInt {
		panic(fmt.Errorf("BUG: cannot allocate a single block of %#x bytes", newSize))
	}

	var newAddr uintptr
	var errno syscall.Errno

	oldSlice := *slice
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&oldSlice))
	oldAddr := hdr.Data
	oldSize := uintptr(uint(hdr.Cap))

	switch {
	case oldSlice == nil && newSize == 0:
		// pass

	case oldSlice == nil:
		var (
			oldAddr uintptr = 0
			prot    uintptr = syscall.PROT_READ | syscall.PROT_WRITE
			flags0  uintptr = syscall.MAP_PRIVATE | syscall.MAP_ANONYMOUS | hugePages.mmapFlags()
			fd      uintptr = ^uintptr(0)
			offset  uintptr = 0
		)

		newAddr, _, errno = syscall.Syscall6(syscall.SYS_MMAP, oldAddr, newSize, prot, flags0, fd, offset)
		if errno != 0 {
			err := &os.SyscallError{Syscall: "mmap", Err: errno}
			panic(fmt.Errorf("mmap(NULL, %#x, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) failed: %w", newSize, err))
		}

	case newSize == 0:
		_, _, errno = syscall.Syscall(syscall.SYS_MUNMAP, oldAddr, oldSize, 0)
		if errno != 0 {
			err := &os.SyscallError{Syscall: "munmap", Err: errno}
			panic(fmt.Errorf("munmap(%#x, %#x) failed: %w", oldAddr, oldSize, err))
		}

	case newSize == oldSize:
		newAddr = oldAddr

	default:
		var flags1 uintptr = _MREMAP_MAYMOVE
		newAddr, _, errno = syscall.Syscall6(syscall.SYS_MREMAP, oldAddr, oldSize, newSize, flags1, 0, 0)
		if errno != 0 {
			err := &os.SyscallError{Syscall: "mremap", Err: errno}
			panic(fmt.Errorf("mremap(%#x, %#x, %#x, MREMAP_MAYMOVE) failed: %w", oldAddr, oldSize, newSize, err))
		}
	}

	var newSlice []byte
	hdr = (*reflect.SliceHeader)(unsafe.Pointer(&newSlice))
	hdr.Data = newAddr
	hdr.Len = int(length)
	hdr.Cap = int(newSize)

	if isLocked {
		_ = mlock(newSlice, true)
	}

	*slice = newSlice
}

func mprotect(slice []byte, r, w, x bool) error {
	if slice == nil {
		return nil
	}

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	memAddr := hdr.Data
	memSize := uintptr(uint(hdr.Cap))

	var memProt uintptr
	var strProt string
	if !r && !w && !x {
		memProt = syscall.PROT_NONE
		strProt = "|PROT_NONE"
	} else {
		if r {
			memProt |= syscall.PROT_READ
			strProt += "|PROT_READ"
		}
		if w {
			memProt |= syscall.PROT_WRITE
			strProt += "|PROT_WRITE"
		}
		if x {
			memProt |= syscall.PROT_EXEC
			strProt += "|PROT_EXEC"
		}
	}

	_, _, errno := syscall.Syscall(syscall.SYS_MPROTECT, memAddr, memSize, memProt)
	if errno != 0 {
		err := &os.SyscallError{Syscall: "mprotect", Err: errno}
		return fmt.Errorf("mprotect(%#x, %#x, %s) failed: %w", memAddr, memSize, strProt[1:], err)
	}

	return nil
}

func mlock(slice []byte, acquire bool) error {
	if slice == nil {
		return nil
	}

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	memAddr := hdr.Data
	memSize := uintptr(uint(hdr.Cap))

	var syscallName string
	var errno syscall.Errno
	if acquire {
		syscallName = "mlock"
		_, _, errno = syscall.Syscall(syscall.SYS_MLOCK, memAddr, memSize, 0)
	} else {
		syscallName = "munlock"
		_, _, errno = syscall.Syscall(syscall.SYS_MUNLOCK, memAddr, memSize, 0)
	}

	if errno != 0 {
		err := &os.SyscallError{Syscall: syscallName, Err: errno}
		return fmt.Errorf("%s(%#x, %#x) failed: %w", syscallName, memAddr, memSize, err)
	}

	var memAdvice uintptr
	var strAdvice string
	if acquire {
		memAdvice = syscall.MADV_DONTFORK
		strAdvice = "MADV_DONTFORK"
	} else {
		memAdvice = syscall.MADV_DOFORK
		strAdvice = "MADV_DOFORK"
	}

	_, _, errno = syscall.Syscall(syscall.SYS_MADVISE, memAddr, memSize, memAdvice)
	if errno != 0 {
		err := &os.SyscallError{Syscall: "madvise", Err: errno}
		return fmt.Errorf("madvise(%#x, %#x, %s) failed: %w", memAddr, memSize, strAdvice, err)
	}

	return nil
}
