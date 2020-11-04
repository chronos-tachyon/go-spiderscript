// +build !386,!amd64,!amd64p32

package exprtree

func SystemCPU() RuntimeCPU {
	// fall back on X86_64
	return X86_64
}
