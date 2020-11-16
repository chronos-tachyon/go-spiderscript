package allocator

const (
	pageShift = 12
	pageSize  = (1 << pageShift) // 4KiB
	pageMask  = (pageSize - 1)

	hugePageShift = 21
	hugePageSize  = (1 << hugePageShift) // 2MiB
	hugePageMask  = (hugePageSize - 1)
)
