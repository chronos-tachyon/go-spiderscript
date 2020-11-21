package allocator

import (
	"sort"
)

type freeListRun struct {
	start uint
	count uint
}

type freeList []freeListRun

func (list freeList) Len() int {
	return len(list)
}

func (list freeList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list freeList) Less(i, j int) bool {
	a, b := list[i], list[j]

	if a.count != b.count {
		return a.count < b.count
	}

	return a.start < b.start
}

var _ sort.Interface = freeList(nil)

type freeListByStart []freeListRun

func (list freeListByStart) Len() int {
	return len(list)
}

func (list freeListByStart) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list freeListByStart) Less(i, j int) bool {
	a, b := list[i], list[j]
	return a.start < b.start
}

var _ sort.Interface = freeListByStart(nil)

func grabFromFreeList(ptr *freeList, allocCount uint, allocBytes uint) (uint, bool) {
	list := *ptr
	i := sort.Search(len(list), func(i int) bool {
		return list[i].count >= allocCount
	})
	if i >= len(list) {
		return 0, false
	}
	run := list[i]
	allocStart := run.start
	if run.count > allocCount {
		remainStart := run.start + allocBytes
		remainCount := run.count - allocCount
		list[i] = freeListRun{remainStart, remainCount}
	} else {
		j := len(list) - 1
		list[i] = list[j]
		list[j] = freeListRun{}
		list = list[:j]
	}
	sort.Sort(list)
	*ptr = list
	return allocStart, true
}
