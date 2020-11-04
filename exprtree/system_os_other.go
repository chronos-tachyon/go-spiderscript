// +build !linux

package exprtree

func SystemOS() RuntimeOS {
	// fall back on LINUX
	return LINUX
}
