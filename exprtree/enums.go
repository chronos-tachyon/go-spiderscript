package exprtree

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
)

// Enum
// {{{

type Enum struct {
	list  Statements
	items []*EnumItem
	byNum map[int64]*EnumItem
	byStr map[string]*EnumItem
	first int64
	last  int64
	dense bool
	kind  TypeKind
}

func (e *Enum) Key() string {
	checkNotNil("e", e)
	return e.list.Key()
}

func (e *Enum) Statements() Statements {
	checkNotNil("e", e)
	return e.list
}

func (e *Enum) Items() []*EnumItem {
	checkNotNil("e", e)
	return cloneEnumItems(e.items)
}

func (e *Enum) ByNumber(num int64) *EnumItem {
	checkNotNil("e", e)

	if !e.dense {
		return e.byNum[num]
	}

	if num >= e.first && num <= e.last {
		index := uint64(num - e.first)
		return e.items[index]
	}

	return nil
}

func (e *Enum) ByName(str string) *EnumItem {
	checkNotNil("e", e)
	return e.byStr[str]
}

func (e *Enum) First() *EnumItem {
	checkNotNil("e", e)
	length := uint(len(e.items))
	if length == 0 {
		return nil
	}
	return e.items[0]
}

func (e *Enum) Last() *EnumItem {
	checkNotNil("e", e)
	length := uint(len(e.items))
	if length == 0 {
		return nil
	}
	return e.items[length-1]
}

func (e *Enum) FirstNumber() int64 {
	checkNotNil("e", e)
	return e.first
}

func (e *Enum) LastNumber() int64 {
	checkNotNil("e", e)
	return e.last
}

func (e *Enum) IsDense() bool {
	checkNotNil("e", e)
	return e.dense
}

func (e *Enum) Kind() TypeKind {
	checkNotNil("e", e)
	return e.kind
}

// }}}

// EnumItem
// {{{

type EnumItem struct {
	parent  *Enum
	number  int64
	name    string
	aliases []string
}

func (item *EnumItem) Parent() *Enum {
	checkNotNil("item", item)
	return item.parent
}

func (item *EnumItem) Number() int64 {
	checkNotNil("item", item)
	return item.number
}

func (item *EnumItem) Name() string {
	checkNotNil("item", item)
	return item.name
}

func (item *EnumItem) Aliases() []string {
	checkNotNil("item", item)
	return cloneStrings(item.aliases)
}

func (item *EnumItem) String() string {
	checkNotNil("item", item)
	return fmt.Sprintf("%s(%d)", item.name, item.number)
}

func (item *EnumItem) GoString() string {
	checkNotNil("item", item)

	numstr := strconv.FormatInt(item.number, 10)

	estimatedLen := 13 + uint(len(numstr)) + uint(len(item.name)) + 2*uint(len(item.aliases))
	for _, alias := range item.aliases {
		estimatedLen += uint(len(alias))
	}

	buf := takeBuffer()
	buf.Grow(int(estimatedLen))
	defer giveBuffer(buf)

	buf.WriteString("EnumItem(")
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

// enumItemsByNumber
// {{{

type enumItemsByNumber []*EnumItem

func (list enumItemsByNumber) Len() int {
	return len(list)
}

func (list enumItemsByNumber) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list enumItemsByNumber) Less(i, j int) bool {
	a, b := list[i], list[j]
	return a.number < b.number
}

var _ sort.Interface = enumItemsByNumber(nil)

// }}}

func (interp *Interp) EnumType(list Statements) (*Type, error) {
	checkNotNil("interp", interp)
	list.Check(EnumStatementContext)
	list.Sort()

	key := list.Key()

	var out *Type
	var found bool
	locked(&interp.mu, func() {
		out, found = interp.enumTypeCache[key]
		if !found {
			interp.enumTypeCache[key] = nil
			return
		}
		for out == nil {
			interp.cv.Wait()
			out = interp.enumTypeCache[key]
		}
	})
	if found {
		return out, nil
	}

	hashedBytes := sha256.Sum256([]byte(key))
	hashedHex := hex.EncodeToString(hashedBytes[:])

	var err error
	out, err = interp.createType(
		interp.BuiltinEnumModule().Symbols(),
		SymbolData{
			Kind: SimpleSymbol,
			Name: "X" + hashedHex,
		},
		func(t *Type) {
			calculateEnum(t, list)
		})

	if err != nil {
		panic(fmt.Errorf("BUG: %w", err))
	}

	locked(&interp.mu, func() {
		interp.enumTypeCache[key] = out
		interp.cv.Broadcast()
	})

	return out, err
}

func calculateEnum(t *Type, list Statements) {
	var (
		omitHash     bool
		omitCompare  bool
		omitToString bool
		omitToRepr   bool
	)

	e := &Enum{
		list:  list,
		items: make([]*EnumItem, 0, len(list)),
		byNum: make(map[int64]*EnumItem, len(list)),
		byStr: make(map[string]*EnumItem, len(list)),
		kind:  InvalidTypeKind,
	}

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

		case EnumKindStatement:
			if !enumTypeLegalKind[stmt.EnumKind] {
				panic(fmt.Errorf("BUG: enum kind is %v, expected primitive integer", stmt.EnumKind))
			}
			e.kind = stmt.EnumKind

		case EnumValueStatement:
			if seen := e.byNum[stmt.EnumNumber]; seen != nil {
				panic(fmt.Errorf("BUG: duplicate enum value %d, already assigned to %v", stmt.EnumNumber, seen))
			}

			if seen := e.byStr[stmt.EnumName]; seen != nil {
				panic(fmt.Errorf("BUG: duplicate enum name %q, already assigned to %v", stmt.EnumName, seen))
			}

			value := &EnumItem{
				parent: e,
				number: stmt.EnumNumber,
				name:   stmt.EnumName,
			}
			e.items = append(e.items, value)
			e.byStr[stmt.EnumName] = value
			e.byNum[stmt.EnumNumber] = value

		case EnumAliasStatement:
			if seen := e.byStr[stmt.EnumName]; seen != nil {
				panic(fmt.Errorf("BUG: duplicate enum name %q, already assigned to %v", stmt.EnumName, seen))
			}

			value := e.byStr[stmt.EnumAliasOf]
			if value == nil {
				panic(fmt.Errorf("BUG: enum name %q is not known", stmt.EnumAliasOf))
			}

			value.aliases = append(value.aliases, stmt.EnumName)
			e.byStr[stmt.EnumName] = value
		}
	}

	sort.Sort(enumItemsByNumber(e.items))

	length := uint(len(e.items))

	if length == 0 {
		panic(fmt.Errorf("BUG: must specify at least one enum value"))
	}

	if item := e.byNum[0]; item == nil {
		panic(fmt.Errorf("BUG: must specify an enum value for number 0"))
	}

	e.first = e.items[0].number
	e.last = e.items[length-1].number

	dense := true
	prev := e.items[0].number - 1
	for index := uint(0); index < length; index++ {
		value := e.items[index]
		if value.number != (prev + 1) {
			dense = false
			break
		}
		prev = value.number
	}
	e.dense = dense

	var backingType *Type
	switch e.kind {
	case U8Kind:
		backingType = t.interp.UInt8Type()
	case U16Kind:
		backingType = t.interp.UInt16Type()
	case U32Kind:
		backingType = t.interp.UInt32Type()
	case U64Kind:
		backingType = t.interp.UInt64Type()
	case S8Kind:
		backingType = t.interp.SInt8Type()
	case S16Kind:
		backingType = t.interp.SInt16Type()
	case S32Kind:
		backingType = t.interp.SInt32Type()
	case S64Kind:
		backingType = t.interp.SInt64Type()
	}

	t.kind = EnumKind
	t.alignShift = backingType.alignShift
	t.minSize = backingType.minSize
	t.padSize = backingType.padSize
	t.details = e

	_ = omitHash
	_ = omitCompare
	_ = omitToString
	_ = omitToRepr
}

var enumTypeLegalKind = map[TypeKind]bool{
	U8Kind:  true,
	U16Kind: true,
	U32Kind: true,
	U64Kind: true,
	S8Kind:  true,
	S16Kind: true,
	S32Kind: true,
	S64Kind: true,
}
