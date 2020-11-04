package exprtree

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

func MangleModuleName(name string) string {
	if !reModuleName.MatchString(name) {
		panic(fmt.Errorf("BUG: invalid module name %q", name))
	}

	pieces := strings.Split(name, "::")
	m := uint(len(pieces))
	estimatedLen := 3 + m
	for _, piece := range pieces {
		n := uint(len(piece))
		estimatedLen += n + lengthUint(n)
	}

	buf := takeBuffer()
	buf.Grow(int(estimatedLen))
	defer giveBuffer(buf)

	buf.WriteString("_A")
	for _, piece := range pieces {
		buf.WriteByte('M')
		writeUint(buf, uint(len(piece)))
		buf.WriteString(piece)
	}
	buf.WriteByte('Z')

	return checkEstimatedLength(buf, estimatedLen)
}

var _ = sort.Sort
var _ = unicode.IsUpper

func isSymbolNameAlreadyMangled(name string) bool {
	if len(name) >= 3 {
		if name[0] == '_' && name[1] == '_' && unicode.IsUpper(rune(name[2])) {
			return true
		}
	}
	return false
}

func mangleSymbolImpl(parent MangledNamer, symbolByte byte, name string) string {
	if !reSymbolName.MatchString(name) {
		panic(fmt.Errorf("BUG: invalid symbol name %q", name))
	}

	isPreMangled := isSymbolNameAlreadyMangled(name)

	outerName := parent.MangledName()
	outerName = outerName[:len(outerName)-1]
	m := uint(len(outerName))
	n := uint(len(name))
	estimatedLen := 2 + m + n
	if isPreMangled {
		estimatedLen -= 2
	} else {
		estimatedLen += 1 + lengthUint(n)
	}

	buf := takeBuffer()
	buf.Grow(int(estimatedLen))
	defer giveBuffer(buf)

	buf.WriteString(outerName)
	buf.WriteByte(symbolByte)
	if isPreMangled {
		buf.WriteString(name[2:])
	} else {
		buf.WriteByte('N')
		writeUint(buf, n)
		buf.WriteString(name)
	}
	buf.WriteByte('Z')

	return checkEstimatedLength(buf, estimatedLen)
}

func MangleGlobalSymbolName(m *Module, name string) string {
	return mangleSymbolImpl(m, 'G', name)
}

/*

func MangleTypeScopedSymbolName(t *Type, isStatic bool, name string) string {
	if !reSymbolName.MatchString(name) {
		panic(fmt.Errorf("BUG: invalid symbol name %q", name))
	}

	symbolByte := byte('I')
	if isStatic {
		symbolByte = 'G'
	}

	return mangleSymbolImpl(t, symbolByte, name)
}

func MangleFunctionScopedSymbolName(f *FunctionType, isStatic bool, name string) string {
	if !reSymbolName.MatchString(name) {
		panic(fmt.Errorf("BUG: invalid symbol name %q", name))
	}

	symbolByte := byte('L')
	if isStatic {
		symbolByte = 'G'
	}

	return mangleSymbolImpl(f, symbolByte, name)
}

*/
