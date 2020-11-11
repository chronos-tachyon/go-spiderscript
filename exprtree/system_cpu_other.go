// +build !386,!amd64,!amd64p32

package exprtree

func SystemCPU() RuntimeCPU {
	return 0
}
