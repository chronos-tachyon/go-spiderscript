package exprtree

import (
	"testing"
)

func TestSystemOSAndSystemCPU(t *testing.T) {
	t.Logf("SystemOS() = %v", SystemOS())
	t.Logf("SystemCPU() = %v", SystemCPU())
}
