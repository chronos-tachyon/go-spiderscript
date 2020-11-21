// +build linux

package allocator

import "syscall"

func getThreadID() uint {
	return uint(syscall.Gettid())
}
