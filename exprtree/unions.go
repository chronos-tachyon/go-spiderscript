package exprtree

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
)

// Union
// {{{

type Union struct {
	list         Statements
	tagSym       *Symbol
	tagType      *Type
	fields       []*UnionField
	byTagAndName map[*EnumItem]map[string]*UnionField
	alignShift   uint8
	minSize      uint16
}

func (u *Union) Key() string {
	return u.list.Key()
}

func (u *Union) Statements() Statements {
	return u.list
}

func (u *Union) TagSymbol() *Symbol {
	return u.tagSym
}

func (u *Union) TagType() *Type {
	return u.tagType
}

func (u *Union) Fields() []*UnionField {
	return cloneUnionFields(u.fields)
}

func (u *Union) FieldsByTag(tag *EnumItem) []*UnionField {
	out := make([]*UnionField, 0, len(u.fields))
	for _, field := range u.fields {
		if field.tag == tag {
			out = append(out, field)
		}
	}
	return out
}

func (u *Union) FieldByTagAndName(tag *EnumItem, name string) *UnionField {
	for _, field := range u.fields {
		if field.tag == tag && field.name == name {
			return field
		}
	}
	return nil
}

func (u *Union) AlignShift() uint {
	return uint(u.alignShift)
}

func (u *Union) MinimumSize() uint {
	return uint(u.minSize)
}

// }}}

// UnionField
// {{{

type UnionField struct {
	parent *Union
	tag    *EnumItem
	name   string
	type_  *Type
	offset uint16
	length uint16

	originalIndex uint
}

func (field *UnionField) Parent() *Union {
	return field.parent
}

func (field *UnionField) Tag() *EnumItem {
	return field.tag
}

func (field *UnionField) Name() string {
	return field.name
}

func (field *UnionField) Type() *Type {
	return field.type_
}

func (field *UnionField) Offset() uint {
	return uint(field.offset)
}

func (field *UnionField) Length() uint {
	return uint(field.length)
}

func (field UnionField) String() string {
	tagName := field.Tag().Name()
	typeName := field.Type().CanonicalName()
	offset := fmt.Sprintf("+%d", field.Offset())
	length := fmt.Sprintf("+%d", field.Length())

	estimatedLen := 5 + uint(len(field.name)) + uint(len(tagName)) + uint(len(typeName)) + uint(len(offset)) + uint(len(length))
	buf := takeBuffer()
	defer giveBuffer(buf)

	buf.WriteString(field.name)
	buf.WriteByte('(')
	buf.WriteString(tagName)
	buf.WriteByte(',')
	buf.WriteString(offset)
	buf.WriteByte(',')
	buf.WriteString(length)
	buf.WriteByte(',')
	buf.WriteString(typeName)
	buf.WriteByte(')')

	return checkEstimatedLength(buf, estimatedLen)
}

func (field UnionField) GoString() string {
	tagName := field.Tag().Name()
	typeName := field.Type().CanonicalName()
	offset := strconv.FormatUint(uint64(field.Offset()), 10)
	length := strconv.FormatUint(uint64(field.Length()), 10)

	estimatedLen := 22 + uint(len(field.name)) + uint(len(tagName)) + uint(len(typeName)) + uint(len(offset)) + uint(len(length))

	buf := takeBuffer()
	buf.Grow(int(estimatedLen))
	defer giveBuffer(buf)

	buf.WriteString("UnionField(")
	buf.WriteByte('"')
	buf.WriteString(field.name)
	buf.WriteByte('"')
	buf.WriteByte(',')
	buf.WriteByte(' ')
	buf.WriteString(tagName)
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

// unionFieldsByTagAlignAndSize
// {{{

type unionFieldsByTagAlignAndSize []*UnionField

func (list unionFieldsByTagAlignAndSize) Len() int {
	return len(list)
}

func (list unionFieldsByTagAlignAndSize) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list unionFieldsByTagAlignAndSize) Less(i, j int) bool {
	a, b := list[i], list[j]

	aTag := a.Tag().Number()
	bTag := b.Tag().Number()
	if aTag != bTag {
		return aTag < bTag
	}

	aType := a.Type()
	bType := b.Type()

	// Stricter alignment requirements come first
	aAlign := aType.AlignShift()
	bAlign := bType.AlignShift()
	if aAlign != bAlign {
		return aAlign > bAlign
	}

	// For same alignment requirement, larger size comes first
	aSize := aType.MinimumBytes()
	bSize := bType.MinimumBytes()
	if aSize != bSize {
		return aSize > bSize
	}

	// Fall back on the original ordering
	aIndex := a.originalIndex
	bIndex := b.originalIndex
	return aIndex < bIndex
}

var _ sort.Interface = unionFieldsByTagAlignAndSize(nil)

// }}}

// unionFieldsByTag
// {{{

type unionFieldsByTag []*UnionField

func (list unionFieldsByTag) Len() int {
	return len(list)
}

func (list unionFieldsByTag) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list unionFieldsByTag) Less(i, j int) bool {
	a, b := list[i], list[j]

	aTag := a.Tag().Number()
	bTag := b.Tag().Number()
	if aTag != bTag {
		return aTag < bTag
	}

	aIndex := a.originalIndex
	bIndex := b.originalIndex
	return aIndex < bIndex
}

var _ sort.Interface = unionFieldsByTag(nil)

// }}}

func (interp *Interp) UnionType(list Statements) (*Type, error) {
	checkNotNil("interp", interp)
	list.Check(UnionStatementContext)
	list.Sort()

	key := list.Key()

	var out *Type
	var found bool
	locked(&interp.mu, func() {
		out, found = interp.unionTypeCache[key]
		if !found {
			interp.unionTypeCache[key] = nil
			return
		}
		for out == nil {
			interp.cv.Wait()
			out = interp.unionTypeCache[key]
		}
	})
	if found {
		return out, nil
	}

	hashedBytes := sha256.Sum256([]byte(key))
	hashedHex := hex.EncodeToString(hashedBytes[:])

	var err error
	out, err = interp.createType(
		interp.BuiltinUnionModule().Symbols(),
		SymbolData{
			Kind: SimpleSymbol,
			Name: "X" + hashedHex,
		},
		func(t *Type) {
			calculateUnion(t, list)
		})

	if err != nil {
		panic(fmt.Errorf("BUG: %w", err))
	}

	locked(&interp.mu, func() {
		interp.unionTypeCache[key] = out
		interp.cv.Broadcast()
	})

	return out, err
}

func calculateUnion(t *Type, list Statements) {
	var (
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

	u := &Union{
		list:         list,
		tagSym:       nil,
		tagType:      nil,
		fields:       make([]*UnionField, 0, len(list)),
		byTagAndName: make(map[*EnumItem]map[string]*UnionField, len(list)),
		alignShift:   ^uint8(0),
		minSize:      ^uint16(0),
	}

	worstCaseScenarioByTag := make(map[*EnumItem]uint, len(list))

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

		case UnionTagStatement:
			if kind := stmt.TagType.Chase().Kind(); kind != EnumKind {
				panic(fmt.Errorf("BUG: wrong tag type: got %v, which is Kind %v, not EnumKind", stmt.TagType.CanonicalName(), kind))
			}
			if !stmt.TagSymbol.Type().Is(stmt.TagType) {
				panic(fmt.Errorf("BUG: tag type %v is mismatched with tag symbol %v of type %v", stmt.TagType.CanonicalName(), stmt.TagSymbol.CanonicalName(), stmt.TagSymbol.Type().CanonicalName()))
			}
			u.tagSym = stmt.TagSymbol
			u.tagType = stmt.TagType

		case UnionFieldStatement:
			if u.tagType == nil {
				panic(fmt.Errorf("BUG: UnionTagStatement is required"))
			}

			if e := u.tagType.Chase().details.(*Enum); stmt.TagItem.Parent() != e {
				panic(fmt.Errorf("BUG: Enum %v does not belong to %s", stmt.TagItem, u.tagType.CanonicalName()))
			}

			byName := u.byTagAndName[stmt.TagItem]
			if byName == nil {
				byName = make(map[string]*UnionField, len(list))
				u.byTagAndName[stmt.TagItem] = byName
			}

			field := byName[stmt.FieldName]
			if field != nil {
				panic(fmt.Errorf("BUG: duplicate field name %q for tag %v", stmt.FieldName, stmt.TagItem))
			}

			field = &UnionField{
				parent:        u,
				tag:           stmt.TagItem,
				name:          stmt.FieldName,
				type_:         stmt.FieldType,
				offset:        ^uint16(0),
				length:        ^uint16(0),
				originalIndex: stmt.originalIndex,
			}
			u.fields = append(u.fields, field)
			byName[stmt.FieldName] = field
			worstCaseScenarioByTag[stmt.TagItem] += stmt.FieldType.PaddedBytes()
		}
	}

	for tag, worstCaseScenario := range worstCaseScenarioByTag {
		if worstCaseScenario > MaxStructSize {
			panic(fmt.Errorf("BUG: struct is too large: %d bytes > %d bytes maximum for tag %v", worstCaseScenario, MaxStructSize, tag))
		}
	}

	usedByTag := make(map[*EnumItem][]bool, len(worstCaseScenarioByTag))
	bytesTotalByTag := make(map[*EnumItem]uint, len(worstCaseScenarioByTag))

	isAvailable := func(tag *EnumItem, start uint, length uint) bool {
		end := start + length
		if end > bytesTotalByTag[tag] {
			end = bytesTotalByTag[tag]
		}
		used := usedByTag[tag]
		for i := start; i < end; i++ {
			if used[i] {
				return false
			}
		}
		return true
	}

	grow := func(tag *EnumItem, newBytesTotal uint) {
		used := usedByTag[tag]
		bytesTotal := bytesTotalByTag[tag]

		if newBytesTotal <= bytesTotal {
			return
		}

		newUsed := make([]bool, newBytesTotal)
		for i := uint(0); i < bytesTotal; i++ {
			newUsed[i] = used[i]
		}

		usedByTag[tag] = newUsed
		bytesTotalByTag[tag] = newBytesTotal
	}

	if isStrictOrder {
		for _, field := range u.fields {
			minBytes := field.type_.MinimumBytes()
			alignShift := field.type_.AlignShift()
			alignMask := uint(1<<alignShift) - 1

			if computedAlign < alignShift {
				computedAlign = alignShift
			}

			// Round up to next allowed offset
			bytesTotal := bytesTotalByTag[field.tag]
			bytesTotal = (bytesTotal + alignMask) & ^alignMask
			bytesTotalByTag[field.tag] = bytesTotal

			start := bytesTotal
			end := start + minBytes

			field.offset = uint16(start)
			field.length = uint16(minBytes)

			grow(field.tag, end)
			used := usedByTag[field.tag]
			for i := start; i < end; i++ {
				used[i] = true
			}
		}
	} else {
		sort.Sort(unionFieldsByTagAlignAndSize(u.fields))

		for _, field := range u.fields {
			minBytes := field.type_.MinimumBytes()
			alignShift := field.type_.AlignShift()
			alignBytes := uint(1) << alignShift

			if computedAlign < alignShift {
				computedAlign = alignShift
			}

			var start uint
			for {
				if isAvailable(field.tag, start, minBytes) {
					end := start + minBytes

					field.offset = uint16(start)
					field.length = uint16(minBytes)

					grow(field.tag, end)
					used := usedByTag[field.tag]
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

	if hasExplicitMinSize {
		for tag, num := range bytesTotalByTag {
			if num < explicitMinSize {
				grow(tag, explicitMinSize)
			}
		}
	}

	var bytesTotal uint
	for _, num := range bytesTotalByTag {
		if bytesTotal < num {
			bytesTotal = num
		}
	}

	u.alignShift = uint8(actualAlign)
	u.minSize = uint16(bytesTotal)

	alignBytes := uint(1) << actualAlign
	padSize := alignBytes
	for padSize < bytesTotal {
		padSize += alignBytes
	}

	t.kind = UnionKind
	t.alignShift = u.alignShift
	t.minSize = u.minSize
	t.padSize = uint16(padSize)
	t.details = u

	_ = omitNew
	_ = omitCopy
	_ = omitMove
	_ = omitHash
	_ = omitCompare
	_ = omitToString
	_ = omitToRepr
}
