package util

func TrimLeadingZeroes(str []byte) []byte {
	if len(str) == 0 {
		return str
	}

	var index uint
	for index < uint(len(str)) && str[index] == '0' {
		index++
	}
	if index >= uint(len(str)) {
		index--
	}
	return str[index:]
}

func TrimTrailingZeroes(str []byte) []byte {
	if len(str) == 0 {
		return str
	}

	index := uint(len(str))
	for index > 0 && str[index-1] == '0' {
		index--
	}
	if index == 0 {
		index++
	}
	return str[:index]
}
