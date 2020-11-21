// +build !linux

package allocator

func getThreadID() uint {
	return 0
}
