package exprtree

import (
	"fmt"
	"math"
)

// Value
// {{{

type Value struct {
	field  *Field
	memory *Memory
}

func (value *Value) Interp() *Interp {
	return value.Field().Interp()
}

func (value *Value) Symbol() *Symbol {
	return value.Field().Symbol()
}

func (value *Value) Field() *Field {
	return value.field
}

func (value *Value) Memory() *Memory {
	return value.memory
}

func (value *Value) CanonicalName() string {
	return value.Field().CanonicalName()
}

func (value *Value) MangledName() string {
	return value.Field().MangledName()
}

func (value *Value) Type() *Type {
	return value.Field().Type()
}

func (value *Value) Offset() uint {
	return value.Field().Offset()
}

func (value *Value) Length() uint {
	return value.Field().Length()
}

func (value *Value) Bytes() []byte {
	return value.Memory().Range(value.Offset(), value.Length())
}

func (value *Value) WithReadLock(fn func(bytes []byte) error) error {
	return value.Memory().WithReadLock(value.Offset(), value.Length(), fn)
}

func (value *Value) WithWriteLock(fn func(bytes []byte) error) error {
	return value.Memory().WithWriteLock(value.Offset(), value.Length(), fn)
}

func (value *Value) ResetToZero() {
	_ = value.WithWriteLock(func(bytes []byte) error {
		for i := range bytes {
			bytes[i] = 0
		}
		return nil
	})
}

func (value *Value) Walk(to TraversalOrder, fn func(Expr)) {
	fn(value)
}

func (value *Value) Eval() (*Value, error) {
	return value, nil
}

func (value *Value) Construct() {
	value.ResetToZero()
}

func (value *Value) Destruct() {
	value.ResetToZero()
}

func (value *Value) Get() interface{} {
	bo := value.Interp().ByteOrder()

	chased := value.Type().Chase()
	kind := chased.Kind()

	var out interface{}
	value.WithReadLock(func(bytes []byte) error {
		switch kind {
		case ReflectedTypeKind:
			t, _ := value.Interp().TypeByID(TypeID(bo.Uint32(bytes)))
			out = t
		case U8Kind:
			out = bytes[0]
		case U16Kind:
			out = bo.Uint16(bytes)
		case U32Kind:
			out = bo.Uint32(bytes)
		case U64Kind:
			out = bo.Uint64(bytes)
		case S8Kind:
			out = int8(bytes[0])
		case S16Kind:
			out = int16(bo.Uint16(bytes))
		case S32Kind:
			out = int32(bo.Uint32(bytes))
		case S64Kind:
			out = int64(bo.Uint64(bytes))
		case F16Kind:
			out = myFloat16frombits(bo.Uint16(bytes))
		case F32Kind:
			out = math.Float32frombits(bo.Uint32(bytes))
		case F64Kind:
			out = math.Float64frombits(bo.Uint64(bytes))
		case C32Kind:
			re := myFloat16frombits(bo.Uint16(bytes[0:2]))
			im := myFloat16frombits(bo.Uint16(bytes[2:4]))
			out = complex(re, im)
		case C64Kind:
			re := math.Float32frombits(bo.Uint32(bytes[0:4]))
			im := math.Float32frombits(bo.Uint32(bytes[4:8]))
			out = complex(re, im)
		case C128Kind:
			re := math.Float64frombits(bo.Uint64(bytes[0:8]))
			im := math.Float64frombits(bo.Uint64(bytes[8:16]))
			out = complex(re, im)

		case StringKind:
			bufID := BufferID(bo.Uint32(bytes[0:4]))
			offset := uint(bo.Uint32(bytes[4:8]))
			length := uint(bo.Uint32(bytes[8:12]))
			if buf, found := value.Interp().BufferByID(bufID); found {
				out = String{buf, offset, length}
			} else {
				out = String{}
			}

		case ErrorKind:
			err, _ := value.Interp().ErrorByID(ErrorID(bo.Uint32(bytes)))
			out = err

		case EnumKind:
			e := chased.Details().(*Enum)

			var s64 int64
			switch e.Kind() {
			case U8Kind:
				s64 = int64(uint64(bytes[0]))
			case U16Kind:
				s64 = int64(uint64(bo.Uint16(bytes)))
			case U32Kind:
				s64 = int64(uint64(bo.Uint32(bytes)))
			case U64Kind:
				s64 = int64(bo.Uint64(bytes))
			case S8Kind:
				s64 = int64(int8(bytes[0]))
			case S16Kind:
				s64 = int64(int16(bo.Uint16(bytes)))
			case S32Kind:
				s64 = int64(int32(bo.Uint32(bytes)))
			case S64Kind:
				s64 = int64(bo.Uint64(bytes))
			default:
				panic(fmt.Errorf("BUG: (*Value).Get(): unknown (*Enum).Kind() %v", e.Kind()))
			}

			out = e.ByNumber(s64)

		case BitfieldKind:
			b := chased.Details().(*Bitfield)

			var u64 uint64
			switch b.Kind() {
			case U8Kind:
				u64 = uint64(bytes[0])
			case U16Kind:
				u64 = uint64(bo.Uint16(bytes))
			case U32Kind:
				u64 = uint64(bo.Uint32(bytes))
			case U64Kind:
				u64 = bo.Uint64(bytes)
			default:
				panic(fmt.Errorf("BUG: (*Value).Get(): unknown (*Bitfield).Kind() %v", b.Kind()))
			}

			items := b.Items()
			bitset := make(map[*BitfieldItem]struct{}, len(items))
			for _, item := range items {
				bit := item.Bit()
				if (u64 & bit) != 0 {
					bitset[item] = struct{}{}
				}
			}

			out = bitset

		default:
			panic(fmt.Errorf("BUG: (*Value).Get(): Kind %v not implemented", kind))
		}
		return nil
	})
	return out
}

func (value *Value) Set(in interface{}) error {
	bo := value.Interp().ByteOrder()

	chased := value.Type().Chase()
	kind := chased.Kind()

	return value.WithWriteLock(func(bytes []byte) error {
		switch kind {
		case ReflectedTypeKind:
			bo.PutUint32(bytes, uint32(in.(*Type).ID()))
		case U8Kind:
			bytes[0] = in.(uint8)
		case U16Kind:
			bo.PutUint16(bytes, in.(uint16))
		case U32Kind:
			bo.PutUint32(bytes, in.(uint32))
		case U64Kind:
			bo.PutUint64(bytes, in.(uint64))
		case S8Kind:
			bytes[0] = uint8(in.(int8))
		case S16Kind:
			bo.PutUint16(bytes, uint16(in.(int16)))
		case S32Kind:
			bo.PutUint32(bytes, uint32(in.(int32)))
		case S64Kind:
			bo.PutUint64(bytes, uint64(in.(int64)))
		case F16Kind:
			bo.PutUint16(bytes, myFloat16bits(in.(float32)))
		case F32Kind:
			bo.PutUint32(bytes, math.Float32bits(in.(float32)))
		case F64Kind:
			bo.PutUint64(bytes, math.Float64bits(in.(float64)))
		case C32Kind:
			x := in.(complex64)
			bo.PutUint16(bytes[0:2], myFloat16bits(real(x)))
			bo.PutUint16(bytes[2:4], myFloat16bits(imag(x)))
		case C64Kind:
			x := in.(complex64)
			bo.PutUint32(bytes[0:4], math.Float32bits(real(x)))
			bo.PutUint32(bytes[4:8], math.Float32bits(imag(x)))
		case C128Kind:
			x := in.(complex128)
			bo.PutUint64(bytes[0:8], math.Float64bits(real(x)))
			bo.PutUint64(bytes[8:16], math.Float64bits(imag(x)))

		case StringKind:
			str := in.(*String)
			var bufID BufferID
			var offset uint
			var length uint
			if str != nil {
				bufID = str.Buffer.ID()
				offset = str.Offset
				length = str.Length
			}
			bo.PutUint32(bytes[0:4], uint32(bufID))
			bo.PutUint32(bytes[4:8], uint32(offset))
			bo.PutUint32(bytes[8:12], uint32(length))

		case ErrorKind:
			bo.PutUint32(bytes, uint32(in.(*Error).ID()))

		case EnumKind:
			e := chased.Details().(*Enum)

			var s64 int64
			switch x := in.(type) {
			case nil:
				// pass

			case int64:
				s64 = x

			case uint64:
				s64 = int64(x)

			case *EnumItem:
				checkNotNil("*EnumItem", x)
				if parent := x.Parent(); parent != e {
					panic(fmt.Errorf("BUG: (*Value).Set(): (*EnumItem).Parent() was %p, expected %p", parent, e))
				}
				s64 = x.Number()

			default:
				panic(fmt.Errorf("BUG: (*Value).Set(): wrong type for argument: expected *EnumItem, got %T", in))
			}

			switch e.Kind() {
			case U8Kind, S8Kind:
				bytes[0] = uint8(s64)
			case U16Kind, S16Kind:
				bo.PutUint16(bytes, uint16(s64))
			case U32Kind, S32Kind:
				bo.PutUint32(bytes, uint32(s64))
			case U64Kind, S64Kind:
				bo.PutUint64(bytes, uint64(s64))
			default:
				panic(fmt.Errorf("BUG: (*Value).Set(): unknown (*Enum).Kind() %v", e.Kind()))
			}

		case BitfieldKind:
			b := chased.Details().(*Bitfield)

			var u64 uint64
			switch x := in.(type) {
			case nil:
				// pass

			case uint64:
				u64 = x

			case map[*BitfieldItem]struct{}:
				for item := range x {
					checkNotNil("item", item)
					if parent := item.Parent(); parent != b {
						panic(fmt.Errorf("BUG: (*Value).Set(): (*BitfieldItem).Parent() was %p, expected %p", parent, b))
					}
					u64 |= item.Bit()
				}

			case []*BitfieldItem:
				for _, item := range x {
					checkNotNil("item", item)
					if parent := item.Parent(); parent != b {
						panic(fmt.Errorf("BUG: (*Value).Set(): (*BitfieldItem).Parent() was %p, expected %p", parent, b))
					}
					u64 |= item.Bit()
				}

			default:
				panic(fmt.Errorf("BUG: (*Value).Set(): wrong type for argument: expected map[*BitfieldItem]struct{}, got %T", in))
			}

			switch b.Kind() {
			case U8Kind:
				bytes[0] = uint8(u64)
			case U16Kind:
				bo.PutUint16(bytes, uint16(u64))
			case U32Kind:
				bo.PutUint32(bytes, uint32(u64))
			case U64Kind:
				bo.PutUint64(bytes, u64)
			default:
				panic(fmt.Errorf("BUG: (*Value).Set(): unknown (*Bitfield).Kind() %v", b.Kind()))
			}

		default:
			panic(fmt.Errorf("BUG: (*Value).Set(): Kind %v not implemented", kind))
		}
		return nil
	})
}

// }}}
