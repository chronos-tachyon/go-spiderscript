package exprtree

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
)

// Bitfield
// {{{

type Bitfield struct {
	list  Statements
	items []*BitfieldItem
	byBit map[uint64]*BitfieldItem
	byStr map[string]*BitfieldItem
	kind  TypeKind
}

func (b *Bitfield) Key() string {
	checkNotNil("b", b)
	return b.list.Key()
}

func (b *Bitfield) Statements() Statements {
	checkNotNil("b", b)
	return b.list
}

func (b *Bitfield) Items() []*BitfieldItem {
	checkNotNil("b", b)
	return cloneBitfieldItems(b.items)
}

func (b *Bitfield) ByShift(shift uint) *BitfieldItem {
	checkNotNil("b", b)
	if shift >= uint(len(b.items)) {
		return nil
	}
	return b.items[shift]
}

func (b *Bitfield) ByBit(bit uint64) *BitfieldItem {
	checkNotNil("b", b)
	return b.byBit[bit]
}

func (b *Bitfield) ByName(str string) *BitfieldItem {
	checkNotNil("b", b)
	return b.byStr[str]
}

func (b *Bitfield) Kind() TypeKind {
	checkNotNil("b", b)
	return b.kind
}

// }}}

// BitfieldItem
// {{{

type BitfieldItem struct {
	parent  *Bitfield
	shift   byte
	name    string
	aliases []string
}

func (item *BitfieldItem) Parent() *Bitfield {
	checkNotNil("item", item)
	return item.parent
}

func (item *BitfieldItem) Shift() uint {
	checkNotNil("item", item)
	return uint(item.shift)
}

func (item *BitfieldItem) Bit() uint64 {
	checkNotNil("item", item)
	return uint64(1) << item.Shift()
}

func (item *BitfieldItem) Name() string {
	checkNotNil("item", item)
	return item.name
}

func (item *BitfieldItem) Aliases() []string {
	checkNotNil("item", item)
	return cloneStrings(item.aliases)
}

func (item *BitfieldItem) String() string {
	checkNotNil("item", item)
	return fmt.Sprintf("%s(1<<%d)", item.Name(), item.Shift())
}

func (item *BitfieldItem) GoString() string {
	checkNotNil("item", item)

	numstr := strconv.FormatUint(uint64(item.shift), 10)

	estimatedLen := 16 + uint(len(numstr)) + uint(len(item.name)) + 2*uint(len(item.aliases))
	for _, alias := range item.aliases {
		estimatedLen += uint(len(alias))
	}

	buf := takeBuffer()
	buf.Grow(int(estimatedLen))
	defer giveBuffer(buf)

	buf.WriteString("BitfieldItem(")
	buf.WriteString(numstr)
	buf.WriteString(", ")
	buf.WriteString(item.name)
	for _, alias := range item.aliases {
		buf.WriteString(", ")
		buf.WriteString(alias)
	}
	buf.WriteByte(')')

	return checkEstimatedLength(buf, estimatedLen)
}

// }}}

// bitfieldItemsByShift
// {{{

type bitfieldItemsByShift []*BitfieldItem

func (list bitfieldItemsByShift) Len() int {
	return len(list)
}

func (list bitfieldItemsByShift) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list bitfieldItemsByShift) Less(i, j int) bool {
	a, b := list[i], list[j]
	return a.shift < b.shift
}

var _ sort.Interface = bitfieldItemsByShift(nil)

// }}}

func (interp *Interp) BitfieldType(list Statements) (*Type, error) {
	checkNotNil("interp", interp)
	list.Check(BitfieldStatementContext)
	list.Sort()

	key := list.Key()

	var out *Type
	var found bool
	locked(&interp.mu, func() {
		out, found = interp.bitfieldTypeCache[key]
		if !found {
			interp.bitfieldTypeCache[key] = nil
			return
		}
		for out == nil {
			interp.cv.Wait()
			out = interp.bitfieldTypeCache[key]
		}
	})
	if found {
		return out, nil
	}

	hashedBytes := sha256.Sum256([]byte(key))
	hashedHex := hex.EncodeToString(hashedBytes[:])

	var err error
	out, err = interp.createType(
		interp.BuiltinBitfieldModule().Symbols(),
		SymbolData{
			Kind: SimpleSymbol,
			Name: "X" + hashedHex,
		},
		func(t *Type) {
			calculateBitfield(t, list)
		})

	if err != nil {
		panic(fmt.Errorf("BUG: %w", err))
	}

	locked(&interp.mu, func() {
		interp.bitfieldTypeCache[key] = out
		interp.cv.Broadcast()
	})

	return out, err
}

func calculateBitfield(t *Type, list Statements) {
	var (
		omitHash     bool
		omitCompare  bool
		omitToString bool
		omitToRepr   bool
	)

	b := &Bitfield{
		list:  list,
		items: make([]*BitfieldItem, 0, 64),
		byBit: make(map[uint64]*BitfieldItem, len(list)),
		byStr: make(map[string]*BitfieldItem, len(list)),
		kind:  InvalidTypeKind,
	}

	byShift := make(map[uint]*BitfieldItem, 64)

	for _, stmt := range list {
		switch stmt.Kind {
		case OmitHashPragmaStatement:
			omitHash = true

		case OmitComparePragmaStatement:
			omitCompare = true

		case OmitToStringPragmaStatement:
			omitToString = true

		case OmitToReprPragmaStatement:
			omitToRepr = true

		case StaticConstantStatement:
			// FIXME

		case InstanceConstantStatement:
			// FIXME

		case StaticFieldStatement:
			// FIXME

		case BitfieldKindStatement:
			if !bitfieldTypeLegalKind[stmt.EnumKind] {
				panic(fmt.Errorf("BUG: bitfield kind is %v, expected primitive integer", stmt.EnumKind))
			}
			b.kind = stmt.EnumKind

		case BitfieldValueStatement:
			rawShift := stmt.EnumNumber
			if rawShift < 0 || rawShift >= 64 {
				panic(fmt.Errorf("BUG: bitfield value 1<<%d is out of range", rawShift))
			}

			shift := uint(rawShift)
			bit := uint64(1) << shift
			str := stmt.EnumName

			if seen := byShift[shift]; seen != nil {
				panic(fmt.Errorf("BUG: duplicate bitfield value 1<<%d, already assigned to %v", shift, seen))
			}

			if seen := b.byStr[str]; seen != nil {
				panic(fmt.Errorf("BUG: duplicate bitfield name %q, already assigned to %v", str, seen))
			}

			item := &BitfieldItem{
				parent: b,
				shift:  byte(shift),
				name:   str,
			}

			b.items = append(b.items, item)
			b.byBit[bit] = item
			b.byStr[str] = item
			byShift[shift] = item

		case BitfieldAliasStatement:
			if seen := b.byStr[stmt.EnumName]; seen != nil {
				panic(fmt.Errorf("BUG: duplicate bitfield name %q, already assigned to %v", stmt.EnumName, seen))
			}

			item := b.byStr[stmt.EnumAliasOf]
			if item == nil {
				panic(fmt.Errorf("BUG: bitfield name %q is not known", stmt.EnumAliasOf))
			}

			item.aliases = append(item.aliases, stmt.EnumName)
			b.byStr[stmt.EnumName] = item
		}
	}

	var limit uint
	switch b.kind {
	case U8Kind:
		limit = 8
	case U16Kind:
		limit = 16
	case U32Kind:
		limit = 32
	case U64Kind:
		limit = 64
	}
	for shift := uint(0); shift < limit; shift++ {
		if byShift[shift] != nil {
			continue
		}

		bit := uint64(1) << shift
		str := fmt.Sprintf("__reserved%d", shift)

		item := &BitfieldItem{
			parent: b,
			shift:  byte(shift),
			name:   str,
		}

		b.items = append(b.items, item)
		b.byBit[bit] = item
		b.byStr[str] = item
		byShift[shift] = item
	}

	for _, item := range b.items {
		shift := uint(item.shift)
		if shift >= limit {
			panic(fmt.Errorf("BUG: bitfield value 1<<%d is out of range for Kind %v", shift, b.kind))
		}
	}

	sort.Sort(bitfieldItemsByShift(b.items))

	var backingType *Type
	switch b.kind {
	case U8Kind:
		backingType = t.interp.UInt8Type()
	case U16Kind:
		backingType = t.interp.UInt16Type()
	case U32Kind:
		backingType = t.interp.UInt32Type()
	case U64Kind:
		backingType = t.interp.UInt64Type()
	}

	t.kind = BitfieldKind
	t.alignShift = backingType.alignShift
	t.minSize = backingType.minSize
	t.padSize = backingType.padSize
	t.details = b

	_ = omitHash
	_ = omitCompare
	_ = omitToString
	_ = omitToRepr
}

var bitfieldTypeLegalKind = map[TypeKind]bool{
	U8Kind:  true,
	U16Kind: true,
	U32Kind: true,
	U64Kind: true,
}
