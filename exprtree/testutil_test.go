package exprtree

import (
	"sync"
)

var gGlobalTestInterpOnce sync.Once
var gGlobalTestInterp *Interp

func GlobalTestInterp() *Interp {
	gGlobalTestInterpOnce.Do(func() {
		gGlobalTestInterp = NewInterp(SystemCPU(), SystemOS())
	})
	return gGlobalTestInterp
}
