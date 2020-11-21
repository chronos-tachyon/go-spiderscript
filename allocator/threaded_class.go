package allocator

import (
	"fmt"
)

const threadedNumClasses = 48

const (
	align0001 = 0
	align0002 = 1
	align0004 = 2
	align0008 = 3
	align0016 = 4
	align0032 = 5
	align0064 = 6
	align0128 = 7
	align0256 = 8
	align0512 = 9
	align1024 = 10
	align2048 = 11
	align4096 = 12
)

type threadedClassRow struct {
	alignShift   uint8
	chunkSize    uint8
	pagesToGrab  uint16
	chunksToGrab uint16
}

var threadedClassData = [threadedNumClasses]threadedClassRow{
	// Comment format is:
	//   "Class #<index>: <bytesPerChunk>/<bytesPerAlign>, <pagesToGrab*pageSize>/<bytesPerChunk>=<chunksPerGrow> chunks"

	// Class  #0: 1B/1B, 4K/1B=4096 chunks
	{align0001, 1, 1, 256},

	// Class  #1: 2B/2B, 4K/2B=2048 chunks
	{align0002, 1, 1, 256},

	// Class  #2: 4B/4B, 4K/4B=1024 chunks
	{align0004, 1, 1, 256},

	// Class  #3: 8B/8B, 8K/8B=1024 chunks
	{align0008, 1, 2, 256},

	// Class  #4: 12B/4B, 12K/12B=1024 chunks
	{align0004, 3, 3, 256},

	// Class  #5: 16B/16B, 16K/16B=1024 chunks
	{align0016, 1, 4, 256},

	// Class  #6: 24B/8B, 12K/24B=512 chunks
	{align0008, 3, 3, 128},

	// Class  #7: 32B/32B, 16K/32B=512 chunks
	{align0032, 1, 4, 128},

	// Class  #8: 40B/8B, 20K/40B=512 chunks
	{align0008, 5, 5, 128},

	// Class  #9: 48B/16B, 12K/48B=256 chunks
	{align0016, 3, 3, 64},

	// Class #10: 56B/8B, 28K/56B=512 chunks
	{align0008, 7, 7, 128},

	// Class #11: 64B/64B, 32K/64B=512 chunks
	{align0064, 1, 8, 128},

	// Class #12: 80B/16B, 20K/80B=256 chunks
	{align0016, 5, 5, 64},

	// Class #13: 96B/32B, 24K/96B=256 chunks
	{align0032, 3, 6, 64},

	// Class #14: 112B/16B, 28K/112B=256 chunks
	{align0016, 7, 7, 64},

	// Class #15: 128B/128B, 32K/128B=256 chunks
	{align0128, 1, 8, 64},

	// Class #16: 160B/32B, 20K/160B=128 chunks
	{align0032, 5, 5, 32},

	// Class #17: 192B/64B, 24K/192B=128 chunks
	{align0064, 3, 6, 32},

	// Class #18: 224B/32B, 28K/224B=128 chunks
	{align0032, 7, 7, 32},

	// Class #19: 256B/256B, 32K/256B=128 chunks
	{align0256, 1, 8, 32},

	// Class #20: 320B/64B, 20K/320B=64 chunks
	{align0064, 5, 5, 16},

	// Class #21: 384B/128B, 24K/384B=64 chunks
	{align0128, 3, 6, 16},

	// Class #22: 448B/64B, 28K/448B=64 chunks
	{align0064, 7, 7, 16},

	// Class #23: 512B/512B, 32K/512B=64 chunks
	{align0512, 1, 8, 16},

	// Class #24: 640B/128B, 20K/640B=32 chunks
	{align0128, 5, 5, 8},

	// Class #25: 768B/256B, 24K/768B=32 chunks
	{align0256, 3, 6, 8},

	// Class #26: 896B/128B, 28K/896B=32 chunks
	{align0128, 7, 7, 8},

	// Class #27: 1K/1K, 32K/1K=32 chunks
	{align1024, 1, 8, 8},

	// Class #28: 1.25K/256B, 20K/1280B=16 chunks
	{align0256, 5, 5, 4},

	// Class #29: 1.50K/512B, 24K/1536B=16 chunks
	{align0512, 3, 6, 4},

	// Class #30: 1.75K/256B, 28K/1792B=16 chunks
	{align0256, 7, 7, 4},

	// Class #31: 2K/2K, 32K/2K=16 chunks
	{align2048, 1, 8, 4},

	// Class #32: 2.5K/512B, 40K/2560B=16 chunks
	{align0512, 5, 10, 4},

	// Class #33: 3K/1K, 48K/3072B=16 chunks
	{align1024, 3, 12, 4},

	// Class #34: 3.5K/512B, 56K/3584B=16 chunks
	{align0512, 7, 14, 4},

	// Class #35: 4K/4K, 64K/4K=16 chunks
	{align4096, 1, 16, 4},

	// Class #36: 5K/1K, 40K/5K=8 chunks
	{align1024, 5, 10, 2},

	// Class #37: 6K/2K, 48K/6K=8 chunks
	{align2048, 3, 12, 2},

	// Class #38: 7K/1K, 48K/7K=8 chunks
	{align1024, 7, 14, 2},

	// Class #39: 8K/4K, 64K/8K=8 chunks
	{align4096, 2, 16, 2},

	// Class #40: 10K/2K, 40K/10K=4 chunks
	{align2048, 5, 10, 2},

	// Class #41: 12K/4K, 48K/12K=4 chunks
	{align4096, 3, 12, 2},

	// Class #42: 14K/2K, 48K/14K=4 chunks
	{align2048, 7, 14, 2},

	// Class #43: 16K/4K, 64K/16K=4 chunks
	{align4096, 4, 16, 2},

	// Class #44: 20K/4K, 80K/20K=4 chunks
	{align4096, 5, 20, 2},

	// Class #45: 24K/4K, 96K/24K=4 chunks
	{align4096, 6, 24, 2},

	// Class #46: 28K/4K, 112K/28K=4 chunks
	{align4096, 7, 28, 2},

	// Class #47: 32K/4K, 128K/32K=4 chunks
	{align4096, 8, 32, 2},
}

func computeThreadedSmallAllocClass(length uint, alignShift uint) uint {
	if alignShift > pageShift {
		panic(fmt.Errorf("BUG: alignShift=%d, max=%d", alignShift, pageShift))
	}
	if length > threadedLargeThreshold {
		panic(fmt.Errorf("BUG: length=%d, threadedLargeThreshold=%d", length, threadedLargeThreshold))
	}
	for classIndex := uint(0); classIndex < threadedNumClasses; classIndex++ {
		class := threadedClassData[classIndex]
		maxLength := uint(class.chunkSize) << class.alignShift
		if alignShift <= uint(class.alignShift) && length <= maxLength {
			return classIndex
		}
	}
	panic(fmt.Errorf("BUG: no alloc class for length=%d, alignShift=%d", length, alignShift))
}
