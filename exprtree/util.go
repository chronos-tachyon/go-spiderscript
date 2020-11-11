package exprtree

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const moduleNameComponent = `[A-Za-z][0-9A-Za-z]*(?:_[0-9A-Za-z]+)*`

var reModuleName = regexp.MustCompile(`^(?:_|` + moduleNameComponent + `(?:::` + moduleNameComponent + `)*)$`)
var reSymbolName = regexp.MustCompile(`^[A-Za-z$_][0-9A-Za-z_]*$`)
var reReservedModuleNames = regexp.MustCompile(`^(?:_|this|main|builtin|builtin::.*)$`)

var gBufferPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

type Checker interface {
	Check()
}

type Keyer interface {
	Key() string
}

type Namer interface {
	CanonicalName() string
	MangledName() string
}

type MangledNamer interface {
	MangledName() string
}

func takeBuffer() *strings.Builder {
	buffer := gBufferPool.Get().(*strings.Builder)
	return buffer
}

func giveBuffer(buffer *strings.Builder) {
	buffer.Reset()
	gBufferPool.Put(buffer)
}

func checkBug(err error) {
	if err != nil {
		panic(fmt.Errorf("BUG: %w", err))
	}
}

func checkNotNil(name string, ptr interface{}) {
	if reflect.ValueOf(ptr).IsNil() {
		panic(fmt.Errorf("BUG: %s is nil", name))
	}
}

func checkIJ(i uint, j uint, size uint) {
	if i > j {
		panic(fmt.Errorf("BUG: i > j; i=%d, j=%d", i, j))
	}
	if j > size {
		panic(fmt.Errorf("BUG: j > size; j=%d, size=%d", j, size))
	}
}

func locked(mu sync.Locker, fn func()) {
	mu.Lock()
	defer mu.Unlock()
	fn()
}

func checkEstimatedLength(buf *strings.Builder, estimatedLen uint) string {
	result := buf.String()
	actualLen := uint(len(result))
	if actualLen != estimatedLen {
		panic(fmt.Errorf("BUG: estimatedLen != actualLen: %d vs %d", estimatedLen, actualLen))
	}
	return result
}

func lengthUint(num uint) uint {
	return lengthUint64(uint64(num))
}

func lengthUint64(num uint64) uint {
	switch {
	case num < 10:
		return 1

	case num < 100:
		return 2

	default:
		return 1 + uint(math.Floor(math.Log10(float64(num))))
	}
}

func writeUint(buf *strings.Builder, num uint) {
	writeUint64(buf, uint64(num))
}

func writeUint64(buf *strings.Builder, num uint64) {
	buf.WriteString(strconv.FormatUint(num, 10))
}

func writeInt(buf *strings.Builder, num int) {
	writeInt64(buf, int64(num))
}

func writeInt64(buf *strings.Builder, num int64) {
	buf.WriteString(strconv.FormatInt(num, 10))
}

func writeNameTo(buf *strings.Builder, name string) {
	buf.WriteByte('N')
	writeUint(buf, uint(len(name)))
	buf.WriteString(name)
}

func cloneAnys(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	copy(out, in)
	return out
}

func cloneStrings(in []string) []string {
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func cloneTypes(in []*Type) []*Type {
	out := make([]*Type, len(in))
	copy(out, in)
	return out
}

func cloneGenericParams(in []GenericParam) []GenericParam {
	out := make([]GenericParam, len(in))
	copy(out, in)
	return out
}

func cloneFunctionArgs(in []FunctionArg) []FunctionArg {
	out := make([]FunctionArg, len(in))
	copy(out, in)
	return out
}

func cloneFields(in []Field) []Field {
	out := make([]Field, len(in))
	copy(out, in)
	return out
}

func cloneValues(in []Value) []Value {
	out := make([]Value, len(in))
	copy(out, in)
	return out
}

func cloneEnumItems(in []*EnumItem) []*EnumItem {
	out := make([]*EnumItem, len(in))
	copy(out, in)
	return out
}

func cloneBitfieldItems(in []*BitfieldItem) []*BitfieldItem {
	out := make([]*BitfieldItem, len(in))
	copy(out, in)
	return out
}

func cloneStructFields(in []*StructField) []*StructField {
	out := make([]*StructField, len(in))
	copy(out, in)
	return out
}

func cloneUnionFields(in []*UnionField) []*UnionField {
	out := make([]*UnionField, len(in))
	copy(out, in)
	return out
}

func cloneStackTrace(in StackTrace) StackTrace {
	out := make(StackTrace, len(in))
	copy(out, in)
	return out
}

func cloneFunctionArgMap(in map[string]FunctionArg) map[string]FunctionArg {
	out := make(map[string]FunctionArg, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func toUInt(kind TypeKind, in interface{}) (uint64, error) {
	var maxValue uint64
	switch kind {
	case U8Kind:
		maxValue = (1 << 8) - 1
	case U16Kind:
		maxValue = (1 << 16) - 1
	case U32Kind:
		maxValue = (1 << 32) - 1
	case U64Kind:
		maxValue = (1 << 64) - 1
	default:
		panic(fmt.Errorf("BUG: TypeKind argument must be U8Kind, U16Kind, U32Kind, or U64Kind; got %v", kind))
	}

	var u64 uint64
	var s64 int64
	isSigned := false
	switch v := in.(type) {
	case uint64:
		u64 = v
	case uintptr:
		u64 = uint64(v)
	case uint:
		u64 = uint64(v)
	case uint32:
		u64 = uint64(v)
	case uint16:
		u64 = uint64(v)
	case uint8:
		u64 = uint64(v)
	case int64:
		s64 = v
		isSigned = true
	case int:
		s64 = int64(v)
		isSigned = true
	case int32:
		s64 = int64(v)
		isSigned = true
	case int16:
		s64 = int64(v)
		isSigned = true
	case int8:
		s64 = int64(v)
		isSigned = true
	default:
		return 0, fmt.Errorf("expected integer type; got %T", in)
	}

	if isSigned && s64 < 0 {
		return 0, fmt.Errorf("value %d out of range for %v", s64, kind)
	}

	if isSigned {
		u64 = uint64(s64)
	}

	if u64 > maxValue {
		return 0, fmt.Errorf("value %d out of range for %v", u64, kind)
	}

	return u64, nil
}

func toSInt(kind TypeKind, in interface{}) (int64, error) {
	var minValue, maxValue int64
	switch kind {
	case S8Kind:
		minValue = -(1 << 7)
		maxValue = (1 << 7) - 1
	case S16Kind:
		minValue = -(1 << 15)
		maxValue = (1 << 15) - 1
	case S32Kind:
		minValue = -(1 << 31)
		maxValue = (1 << 31) - 1
	case S64Kind:
		minValue = -(1 << 63)
		maxValue = (1 << 63) - 1
	default:
		panic(fmt.Errorf("BUG: TypeKind argument must be S8Kind, S16Kind, S32Kind, or S64Kind; got %v", kind))
	}

	var u64 uint64
	var s64 int64
	isUnsigned := false
	switch v := in.(type) {
	case int64:
		s64 = v
	case int:
		s64 = int64(v)
	case int32:
		s64 = int64(v)
	case int16:
		s64 = int64(v)
	case int8:
		s64 = int64(v)
	case uint64:
		u64 = v
		isUnsigned = true
	case uintptr:
		u64 = uint64(v)
		isUnsigned = true
	case uint:
		u64 = uint64(v)
		isUnsigned = true
	case uint32:
		u64 = uint64(v)
		isUnsigned = true
	case uint16:
		u64 = uint64(v)
		isUnsigned = true
	case uint8:
		u64 = uint64(v)
		isUnsigned = true
	default:
		return 0, fmt.Errorf("expected integer type; got %T", in)
	}

	if isUnsigned && u64 > uint64(maxValue) {
		return 0, fmt.Errorf("value %d out of range for %v", u64, kind)
	}

	if isUnsigned {
		s64 = int64(u64)
	}

	if s64 < minValue || s64 > maxValue {
		return 0, fmt.Errorf("value %d out of range for %v", s64, kind)
	}

	return s64, nil
}

func toFloat(in interface{}) (float64, error) {
	switch v := in.(type) {
	case float64:
		return v, nil

	case float32:
		return float64(v), nil

	default:
		return 0, fmt.Errorf("expected float type; got %T", in)
	}
}

func toComplex(in interface{}) (complex128, error) {
	switch v := in.(type) {
	case complex128:
		return v, nil

	case complex64:
		return complex128(v), nil

	case float64:
		return complex(v, 0), nil

	case float32:
		return complex(float64(v), 0), nil

	default:
		return 0, fmt.Errorf("expected complex type; got %T", in)
	}
}

func myFloat16bits(f32 float32) uint16 {
	u32 := math.Float32bits(f32)
	u16 := uint16(u32) // FIXME
	return u16
}

func myFloat16frombits(u16 uint16) float32 {
	u32 := uint32(u16) // FIXME
	return math.Float32frombits(u32)
}
