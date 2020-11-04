package exprtree

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
)

// Struct
// {{{

type Struct struct {
	list       Statements
	fields     []*StructField
	alignShift uint8
	minSize    uint16
}

func (s *Struct) Key() string {
	return s.list.Key()
}

func (s *Struct) Statements() Statements {
	return s.list
}

func (s *Struct) Fields() []*StructField {
	return cloneStructFields(s.fields)
}

func (s *Struct) FieldByName(name string) *StructField {
	for _, field := range s.fields {
		if field.fname == name {
			return field
		}
	}
	return nil
}

func (s *Struct) AlignShift() uint {
	return uint(s.alignShift)
}

func (s *Struct) MinimumSize() uint {
	return uint(s.minSize)
}

// }}}

// StructField
// {{{

type StructField struct {
	parent *Struct
	fname  string
	ftype  *Type
	offset uint16
	length uint16

	originalIndex uint
}

func (field *StructField) Name() string {
	return field.fname
}

func (field *StructField) Type() *Type {
	return field.ftype
}

func (field *StructField) Offset() uint {
	return uint(field.offset)
}

func (field *StructField) Length() uint {
	return uint(field.length)
}

func (field StructField) String() string {
	typeName := field.ftype.CanonicalName()
	offset := fmt.Sprintf("+%d", field.offset)
	length := fmt.Sprintf("+%d", field.length)

	estimatedLen := 4 + uint(len(field.fname)) + uint(len(typeName)) + uint(len(offset)) + uint(len(length))
	buf := takeBuffer()
	defer giveBuffer(buf)

	buf.WriteString(field.fname)
	buf.WriteByte('(')
	buf.WriteString(offset)
	buf.WriteByte(',')
	buf.WriteString(length)
	buf.WriteByte(',')
	buf.WriteString(typeName)
	buf.WriteByte(')')

	return checkEstimatedLength(buf, estimatedLen)
}

func (field StructField) GoString() string {
	typeName := field.ftype.CanonicalName()
	offset := strconv.FormatUint(uint64(field.offset), 10)
	length := strconv.FormatUint(uint64(field.length), 10)

	estimatedLen := 21 + uint(len(field.fname)) + uint(len(typeName)) + uint(len(offset)) + uint(len(length))

	buf := takeBuffer()
	buf.Grow(int(estimatedLen))
	defer giveBuffer(buf)

	buf.WriteString("StructField(")
	buf.WriteByte('"')
	buf.WriteString(field.fname)
	buf.WriteByte('"')
	buf.WriteByte(',')
	buf.WriteByte(' ')
	buf.WriteString(typeName)
	buf.WriteByte(',')
	buf.WriteByte(' ')
	buf.WriteString(offset)
	buf.WriteByte(',')
	buf.WriteByte(' ')
	buf.WriteString(length)
	buf.WriteByte(')')

	return checkEstimatedLength(buf, estimatedLen)
}

// }}}

// structFieldsByAlignAndSize
// {{{

type structFieldsByAlignAndSize []*StructField

func (list structFieldsByAlignAndSize) Len() int {
	return len(list)
}

func (list structFieldsByAlignAndSize) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list structFieldsByAlignAndSize) Less(i, j int) bool {
	a, b := list[i], list[j]

	// Stricter alignment requirements come first
	aAlign := a.ftype.AlignShift()
	bAlign := b.ftype.AlignShift()
	if aAlign != bAlign {
		return aAlign > bAlign
	}

	// For same alignment requirement, larger size comes first
	aSize := a.ftype.MinimumBytes()
	bSize := b.ftype.MinimumBytes()
	if aSize != bSize {
		return aSize > bSize
	}

	// Fall back on the original ordering
	aIndex := a.originalIndex
	bIndex := b.originalIndex
	return aIndex < bIndex
}

var _ sort.Interface = structFieldsByAlignAndSize(nil)

// }}}

func (interp *Interp) StructType(list Statements) (*Type, error) {
	checkNotNil("interp", interp)
	list.Check(StructStatementContext)
	list.Sort()

	key := list.Key()

	var out *Type
	var found bool
	locked(&interp.mu, func() {
		out, found = interp.structTypeCache[key]
		if !found {
			interp.structTypeCache[key] = nil
			return
		}
		for out == nil {
			interp.cv.Wait()
			out = interp.structTypeCache[key]
		}
	})
	if found {
		return out, nil
	}

	hashedBytes := sha256.Sum256([]byte(key))
	hashedHex := hex.EncodeToString(hashedBytes[:])

	var err error
	out, err = interp.createType(
		interp.BuiltinStructModule().Symbols(),
		SymbolData{
			Kind: SimpleSymbol,
			Name: "X" + hashedHex,
		},
		func(t *Type) {
			calculateStruct(t, list)
		})

	if err != nil {
		panic(fmt.Errorf("BUG: %w", err))
	}

	locked(&interp.mu, func() {
		interp.structTypeCache[key] = out
		interp.cv.Broadcast()
	})

	return out, err
}

func calculateStruct(t *Type, list Statements) {
	var (
		worstCaseScenario  uint
		computedAlign      uint
		explicitAlign      uint
		explicitMinSize    uint
		hasExplicitAlign   bool
		hasExplicitMinSize bool
		isStrictOrder      bool
		omitNew            bool
		omitCopy           bool
		omitMove           bool
		omitHash           bool
		omitCompare        bool
		omitToString       bool
		omitToRepr         bool
	)

	s := &Struct{
		list:       list,
		fields:     make([]*StructField, 0, len(list)),
		alignShift: ^uint8(0),
		minSize:    ^uint16(0),
	}

	for index, length := uint(0), uint(len(list)); index < length; index++ {
		stmt := list[index]

		switch stmt.Kind {
		case AlignPragmaStatement:
			explicitAlign = stmt.MinAlign
			hasExplicitAlign = true

		case MinimumSizePragmaStatement:
			explicitMinSize = stmt.MinSize
			hasExplicitMinSize = true

		case PreserveFieldOrderPragmaStatement:
			isStrictOrder = true

		case OmitNewPragmaStatement:
			omitNew = true

		case OmitCopyPragmaStatement:
			omitCopy = true

		case OmitMovePragmaStatement:
			omitCopy = true
			omitMove = true

		case OmitHashPragmaStatement:
			omitHash = true

		case OmitComparePragmaStatement:
			omitCompare = true

		case OmitToStringPragmaStatement:
			omitToString = true

		case OmitToReprPragmaStatement:
			omitToRepr = true

		case StructFieldStatement:
			s.fields = append(s.fields, &StructField{
				parent:        s,
				fname:         stmt.FieldName,
				ftype:         stmt.FieldType,
				offset:        ^uint16(0),
				length:        ^uint16(0),
				originalIndex: stmt.originalIndex,
			})

			worstCaseScenario += stmt.FieldType.PaddedBytes()
		}
	}

	if worstCaseScenario > MaxStructSize {
		panic(fmt.Errorf("BUG: struct is too large: %d bytes > %d bytes maximum", worstCaseScenario, MaxStructSize))
	}

	var used []bool
	var bytesTotal uint

	isAvailable := func(start uint, length uint) bool {
		end := start + length
		if end > bytesTotal {
			end = bytesTotal
		}
		for i := start; i < end; i++ {
			if used[i] {
				return false
			}
		}
		return true
	}

	grow := func(newBytesTotal uint) {
		if newBytesTotal <= bytesTotal {
			return
		}

		newUsed := make([]bool, newBytesTotal)
		for i := uint(0); i < bytesTotal; i++ {
			newUsed[i] = used[i]
		}

		used = newUsed
		bytesTotal = newBytesTotal
	}

	if isStrictOrder {
		for _, field := range s.fields {
			minBytes := field.ftype.MinimumBytes()
			alignShift := field.ftype.AlignShift()
			alignMask := uint(1<<alignShift) - 1

			if computedAlign < alignShift {
				computedAlign = alignShift
			}

			// Round up to next allowed offset
			bytesTotal = (bytesTotal + alignMask) & ^alignMask
			start := bytesTotal
			end := start + minBytes

			field.offset = uint16(start)
			field.length = uint16(minBytes)

			grow(end)
			for i := start; i < end; i++ {
				used[i] = true
			}
		}
	} else {
		sort.Sort(structFieldsByAlignAndSize(s.fields))

		for _, field := range s.fields {
			minBytes := field.ftype.MinimumBytes()
			alignShift := field.ftype.AlignShift()
			alignBytes := uint(1 << alignShift)

			if computedAlign < alignShift {
				computedAlign = alignShift
			}

			var start uint
			for {
				if isAvailable(start, minBytes) {
					end := start + minBytes

					field.offset = uint16(start)
					field.length = uint16(minBytes)

					grow(end)
					for i := start; i < end; i++ {
						used[i] = true
					}
					break
				}

				start += alignBytes
			}
		}
	}

	actualAlign := computedAlign
	if hasExplicitAlign && actualAlign < explicitAlign {
		actualAlign = explicitAlign
	}

	if hasExplicitMinSize && bytesTotal < explicitMinSize {
		grow(explicitMinSize)
	}

	s.alignShift = uint8(actualAlign)
	s.minSize = uint16(bytesTotal)

	alignBytes := uint(1) << actualAlign
	padSize := alignBytes
	for padSize < bytesTotal {
		padSize += alignBytes
	}

	t.kind = StructKind
	t.alignShift = s.alignShift
	t.minSize = s.minSize
	t.padSize = uint16(padSize)
	t.details = s

	_ = omitNew
	_ = omitCopy
	_ = omitMove
	_ = omitHash
	_ = omitCompare
	_ = omitToString
	_ = omitToRepr
}
