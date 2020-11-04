package util

func CompareBytes(a []byte, b []byte) int {
	i := uint(0)
	j := uint(len(a))
	k := uint(len(b))
	for i < j && i < k {
		ach := a[i]
		bch := b[i]
		if ach < bch {
			return -1
		}
		if ach > bch {
			return 1
		}
		i++
	}
	if j < k {
		return -1
	}
	if j > k {
		return 1
	}
	return 0
}

func EqualBytes(a []byte, b []byte) bool {
	return CompareBytes(a, b) == 0
}

func CompareRunes(a []rune, b []rune) int {
	i := uint(0)
	j := uint(len(a))
	k := uint(len(b))
	for i < j && i < k {
		ach := a[i]
		bch := b[i]
		if ach < bch {
			return -1
		}
		if ach > bch {
			return 1
		}
		i++
	}
	if j < k {
		return -1
	}
	if j > k {
		return 1
	}
	return 0
}

func EqualRunes(a []rune, b []rune) bool {
	return CompareRunes(a, b) == 0
}

func IsAllByte(b byte, a []byte) bool {
	for _, ch := range a {
		if ch != b {
			return false
		}
	}
	return true
}
