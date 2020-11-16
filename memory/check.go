package memory

import (
	"fmt"
)

func checkIJ(i uint, j uint, size uint) {
	if i > j {
		panic(fmt.Errorf("BUG: i > j; i=%d, j=%d", i, j))
	}
	if j > size {
		panic(fmt.Errorf("BUG: j > size; j=%d, size=%d", j, size))
	}
}

func checkCast(fromName, toName string, minAlignShift, actualAlignShift uint) {
	if actualAlignShift < minAlignShift {
		panic(fmt.Errorf("%s is not aligned strongly enough to be used as %s: minimum alignShift %d, actual alignShift %d", fromName, toName, minAlignShift, actualAlignShift))
	}
}
