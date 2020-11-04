package exprtree

import (
	"fmt"
	"sort"
)

// Statements
// {{{

type Statements []Statement

func (list Statements) Check(context StatementContext) {
	counts := make(map[StatementKind]uint, len(list))
	for _, stmt := range list {
		stmt.Check(context)
		counts[stmt.Kind]++
	}
	for kind, count := range counts {
		if count > 1 && statementsNoDuplicatesAllowed[kind] {
			panic(fmt.Errorf("BUG: StatementKind %v appears %d times, but it is only allowed to appear once", kind, count))
		}
	}
}

func (list Statements) Sort() {
	for index, length := uint(0), uint(len(list)); index < length; index++ {
		list[index].originalIndex = index
	}
	sort.Sort(list)
}

func (list Statements) CheckSorted() {
	for index, length := uint(0), uint(len(list)); index < length; index++ {
		list[index].originalIndex = index
	}
	if !sort.IsSorted(list) {
		panic(fmt.Errorf("BUG: Statements list is not sorted"))
	}
}

func (list Statements) Key() string {
	list.CheckSorted()

	buf := takeBuffer()
	defer giveBuffer(buf)

	for _, stmt := range list {
		buf.WriteString(stmt.Key())
	}

	return buf.String()
}

func (list Statements) Len() int {
	return len(list)
}

func (list Statements) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list Statements) Less(i, j int) bool {
	a, b := list[i], list[j]

	if a.Kind != b.Kind {
		return a.Kind < b.Kind
	}

	return a.originalIndex < b.originalIndex
}

var _ sort.Interface = Statements(nil)

// }}}

// Statement
// {{{

type Statement struct {
	Kind StatementKind

	MinAlign uint
	MinSize  uint

	EnumName    string
	EnumAliasOf string
	EnumNumber  int64
	EnumKind    TypeKind

	TagSymbol *Symbol
	TagType   *Type
	TagItem   *EnumItem

	FieldName  string
	FieldType  *Type
	FieldValue interface{}

	MethodName      string
	MethodSignature *FunctionSignature

	originalIndex uint
}

func (stmt Statement) Check(context StatementContext) {
	if !statementsValid[context][stmt.Kind] {
		panic(fmt.Errorf("BUG: Kind %v is not permitted inside %v", stmt.Kind, context))
	}

	var (
		minAlignIsLegal   bool
		minSizeIsLegal    bool
		enumNameIsLegal   bool
		enumAliasIsLegal  bool
		enumNumberIsLegal bool
		enumKindIsLegal   bool
		fieldNameIsLegal  bool
		fieldTypeIsLegal  bool
		fieldValueIsLegal bool
		tagSymbolIsLegal  bool
		tagTypeIsLegal    bool
		tagItemIsLegal    bool
		methodIsLegal     bool
	)

	switch stmt.Kind {
	case AlignPragmaStatement:
		minAlignIsLegal = true

	case MinimumSizePragmaStatement:
		minSizeIsLegal = true

	case PreserveFieldOrderPragmaStatement:
		// pass

	case OmitNewPragmaStatement:
		// pass

	case OmitCopyPragmaStatement:
		// pass

	case OmitMovePragmaStatement:
		// pass

	case OmitHashPragmaStatement:
		// pass

	case OmitComparePragmaStatement:
		// pass

	case OmitToStringPragmaStatement:
		// pass

	case OmitToReprPragmaStatement:
		// pass

	case StaticConstantStatement, InstanceConstantStatement:
		fieldNameIsLegal = true
		fieldTypeIsLegal = true
		fieldValueIsLegal = true

	case StaticFieldStatement:
		fieldNameIsLegal = true
		fieldTypeIsLegal = true

	case EnumKindStatement, BitfieldKindStatement:
		enumKindIsLegal = true

	case EnumValueStatement, BitfieldValueStatement:
		enumNameIsLegal = true
		enumNumberIsLegal = true

	case EnumAliasStatement, BitfieldAliasStatement:
		enumNameIsLegal = true
		enumAliasIsLegal = true

	case StructFieldStatement:
		fieldNameIsLegal = true
		fieldTypeIsLegal = true

	case UnionTagStatement:
		tagSymbolIsLegal = true
		tagTypeIsLegal = true

	case UnionFieldStatement:
		tagItemIsLegal = true
		fieldNameIsLegal = true
		fieldTypeIsLegal = true

	case InterfaceFieldStatement:
		fieldNameIsLegal = true
		fieldTypeIsLegal = true

	case InterfacePropertyStatement:
		fieldNameIsLegal = true
		fieldTypeIsLegal = true

	case InterfaceMethodStatement:
		methodIsLegal = true

	default:
		panic(fmt.Errorf("BUG: Kind %v not implemented", stmt.Kind))
	}

	if !minAlignIsLegal && stmt.MinAlign != 0 {
		panic(fmt.Errorf("BUG: Kind %v does not allow MinAlign %d", stmt.Kind, stmt.MinAlign))
	}
	if stmt.MinAlign > MaxAlignShift {
		panic(fmt.Errorf("BUG: Kind %v has MinAlign %d > %d", stmt.Kind, stmt.MinAlign, MaxAlignShift))
	}

	if !minSizeIsLegal && stmt.MinSize != 0 {
		panic(fmt.Errorf("BUG: Kind %v does not allow MinSize %d", stmt.Kind, stmt.MinSize))
	}
	if stmt.MinSize > MaxStructSize {
		panic(fmt.Errorf("BUG: Kind %v has MinSize %d > %d", stmt.Kind, stmt.MinSize, MaxStructSize))
	}

	if enumNameIsLegal {
		if !reSymbolName.MatchString(stmt.EnumName) {
			panic(fmt.Errorf("BUG: Kind %v has invalid EnumName %q", stmt.Kind, stmt.EnumName))
		}
	} else {
		if stmt.EnumName != "" {
			panic(fmt.Errorf("BUG: Kind %v does not allow EnumName %q", stmt.Kind, stmt.EnumName))
		}
	}

	if enumAliasIsLegal {
		if !reSymbolName.MatchString(stmt.EnumAliasOf) {
			panic(fmt.Errorf("BUG: Kind %v has invalid EnumAliasOf %q", stmt.Kind, stmt.EnumAliasOf))
		}
	} else {
		if stmt.EnumAliasOf != "" {
			panic(fmt.Errorf("BUG: Kind %v does not allow EnumAliasOf %q", stmt.Kind, stmt.EnumAliasOf))
		}
	}

	if !enumNumberIsLegal && stmt.EnumNumber != 0 {
		panic(fmt.Errorf("BUG: Kind %v does not allow EnumNumber %d", stmt.Kind, stmt.EnumNumber))
	}

	if !enumKindIsLegal && stmt.EnumKind != 0 {
		panic(fmt.Errorf("BUG: Kind %v does not allow EnumKind %v", stmt.Kind, stmt.EnumKind))
	}

	if tagSymbolIsLegal {
		if stmt.TagSymbol == nil {
			panic(fmt.Errorf("BUG: Kind %v has nil TagSymbol", stmt.Kind))
		}
	} else {
		if stmt.TagSymbol != nil {
			panic(fmt.Errorf("BUG: Kind %v does not allow non-nil TagSymbol", stmt.Kind))
		}
	}

	if tagTypeIsLegal {
		if stmt.TagType == nil {
			panic(fmt.Errorf("BUG: Kind %v has nil TagType", stmt.Kind))
		}
	} else {
		if stmt.TagType != nil {
			panic(fmt.Errorf("BUG: Kind %v does not allow non-nil TagType", stmt.Kind))
		}
	}

	if tagItemIsLegal {
		if stmt.TagItem == nil {
			panic(fmt.Errorf("BUG: Kind %v has nil TagItem", stmt.Kind))
		}
	} else {
		if stmt.TagItem != nil {
			panic(fmt.Errorf("BUG: Kind %v does not allow non-nil TagItem", stmt.Kind))
		}
	}

	if fieldNameIsLegal {
		if !reSymbolName.MatchString(stmt.FieldName) {
			panic(fmt.Errorf("BUG: Kind %v has invalid FieldName %q", stmt.Kind, stmt.FieldName))
		}
	} else {
		if stmt.FieldName != "" {
			panic(fmt.Errorf("BUG: Kind %v does not allow FieldName %q", stmt.Kind, stmt.FieldName))
		}
	}

	if fieldTypeIsLegal {
		if stmt.FieldType == nil {
			panic(fmt.Errorf("BUG: Kind %v has nil FieldType", stmt.Kind))
		}
	} else {
		if stmt.FieldType != nil {
			panic(fmt.Errorf("BUG: Kind %v does not allow non-nil FieldType", stmt.Kind))
		}
	}

	if fieldValueIsLegal {
		if stmt.FieldValue == nil {
			panic(fmt.Errorf("BUG: Kind %v has nil FieldValue", stmt.Kind))
		}
	} else {
		if stmt.FieldValue != nil {
			panic(fmt.Errorf("BUG: Kind %v does not allow non-nil FieldValue", stmt.Kind))
		}
	}

	if methodIsLegal {
		if !reSymbolName.MatchString(stmt.MethodName) {
			panic(fmt.Errorf("BUG: Kind %v has invalid MethodName %q", stmt.Kind, stmt.MethodName))
		}
		if stmt.MethodSignature == nil {
			panic(fmt.Errorf("BUG: Kind %v has nil MethodSignature", stmt.Kind))
		}
	} else {
		if stmt.MethodName != "" {
			panic(fmt.Errorf("BUG: Kind %v does not allow MethodName %q", stmt.Kind, stmt.MethodName))
		}
		if stmt.MethodSignature != nil {
			panic(fmt.Errorf("BUG: Kind %v does not allow non-nil MethodSignature", stmt.Kind))
		}
	}
}

func (stmt Statement) Key() string {
	buf := takeBuffer()
	defer giveBuffer(buf)

	switch stmt.Kind {
	case AlignPragmaStatement:
		buf.WriteString("align,")
		writeUint(buf, stmt.MinAlign)

	case MinimumSizePragmaStatement:
		buf.WriteString("minSize,")
		writeUint(buf, stmt.MinSize)

	case PreserveFieldOrderPragmaStatement:
		buf.WriteString("preserveOrder")

	case OmitNewPragmaStatement:
		buf.WriteString("omitNew")

	case OmitCopyPragmaStatement:
		buf.WriteString("omitCopy")

	case OmitMovePragmaStatement:
		buf.WriteString("omitMove")

	case OmitHashPragmaStatement:
		buf.WriteString("omitHash")

	case OmitComparePragmaStatement:
		buf.WriteString("omitCompare")

	case OmitToStringPragmaStatement:
		buf.WriteString("omitToString")

	case OmitToReprPragmaStatement:
		buf.WriteString("omitToRepr")

	case StaticConstantStatement:
		buf.WriteString("staticConst,")
		buf.WriteString(stmt.FieldName)
		buf.WriteByte(',')
		buf.WriteString(stmt.FieldType.CanonicalName())
		buf.WriteByte(',')
		buf.WriteByte('[')
		buf.WriteString(stmt.FieldValue.(Keyer).Key())
		buf.WriteByte(']')

	case InstanceConstantStatement:
		buf.WriteString("const,")
		buf.WriteString(stmt.FieldName)
		buf.WriteByte(',')
		buf.WriteString(stmt.FieldType.CanonicalName())
		buf.WriteByte(',')
		buf.WriteByte('[')
		buf.WriteString(stmt.FieldValue.(Keyer).Key())
		buf.WriteByte(']')

	case StaticFieldStatement:
		buf.WriteString("staticField,")
		buf.WriteString(stmt.FieldName)
		buf.WriteByte(',')
		buf.WriteString(stmt.FieldType.CanonicalName())

	case EnumKindStatement:
		buf.WriteString("enumKind,")
		buf.WriteString(stmt.EnumKind.String())

	case EnumValueStatement:
		buf.WriteString("enumValue,")
		buf.WriteString(stmt.EnumName)
		buf.WriteByte(',')
		writeInt64(buf, stmt.EnumNumber)

	case EnumAliasStatement:
		buf.WriteString("enumAlias,")
		buf.WriteString(stmt.EnumName)
		buf.WriteByte(',')
		buf.WriteString(stmt.EnumAliasOf)

	case BitfieldKindStatement:
		buf.WriteString("bitKind,")
		buf.WriteString(stmt.EnumKind.String())

	case BitfieldValueStatement:
		buf.WriteString("bitValue,")
		buf.WriteString(stmt.EnumName)
		buf.WriteByte(',')
		writeInt64(buf, stmt.EnumNumber)

	case BitfieldAliasStatement:
		buf.WriteString("bitAlias,")
		buf.WriteString(stmt.EnumName)
		buf.WriteByte(',')
		buf.WriteString(stmt.EnumAliasOf)

	case StructFieldStatement:
		buf.WriteString("structField,")
		buf.WriteString(stmt.FieldName)
		buf.WriteByte(',')
		buf.WriteString(stmt.FieldType.CanonicalName())

	case UnionTagStatement:
		buf.WriteString("unionTag,")
		buf.WriteString(stmt.TagSymbol.CanonicalName())
		buf.WriteByte(',')
		buf.WriteString(stmt.TagType.CanonicalName())

	case UnionFieldStatement:
		buf.WriteString("unionField,")
		buf.WriteString(stmt.TagItem.Name())
		buf.WriteByte(',')
		buf.WriteString(stmt.FieldName)
		buf.WriteByte(',')
		buf.WriteString(stmt.FieldType.CanonicalName())

	case InterfaceFieldStatement:
		buf.WriteString("ifaceField,")
		buf.WriteString(stmt.FieldName)
		buf.WriteByte(',')
		buf.WriteString(stmt.FieldType.CanonicalName())

	case InterfacePropertyStatement:
		buf.WriteString("ifaceProperty,")
		buf.WriteString(stmt.FieldName)
		buf.WriteByte(',')
		buf.WriteString(stmt.FieldType.CanonicalName())

	case InterfaceMethodStatement:
		buf.WriteString("ifaceMethod,")
		buf.WriteString(stmt.MethodName)
		buf.WriteByte(',')
		buf.WriteString(stmt.MethodSignature.String())

	default:
		panic(fmt.Errorf("BUG: StatementKind %v not implemented", stmt.Kind))
	}

	buf.WriteByte(';')
	return buf.String()
}

// }}}

var statementsValid = map[StatementContext]map[StatementKind]bool{
	EnumStatementContext: {
		OmitHashPragmaStatement:     true,
		OmitComparePragmaStatement:  true,
		OmitToStringPragmaStatement: true,
		OmitToReprPragmaStatement:   true,
		StaticConstantStatement:     true,
		InstanceConstantStatement:   true,
		StaticFieldStatement:        true,
		EnumKindStatement:           true,
		EnumValueStatement:          true,
		EnumAliasStatement:          true,
	},
	BitfieldStatementContext: {
		OmitHashPragmaStatement:     true,
		OmitComparePragmaStatement:  true,
		OmitToStringPragmaStatement: true,
		OmitToReprPragmaStatement:   true,
		StaticConstantStatement:     true,
		InstanceConstantStatement:   true,
		StaticFieldStatement:        true,
		BitfieldKindStatement:       true,
		BitfieldValueStatement:      true,
		BitfieldAliasStatement:      true,
	},
	StructStatementContext: {
		AlignPragmaStatement:              true,
		MinimumSizePragmaStatement:        true,
		PreserveFieldOrderPragmaStatement: true,
		OmitNewPragmaStatement:            true,
		OmitCopyPragmaStatement:           true,
		OmitMovePragmaStatement:           true,
		OmitHashPragmaStatement:           true,
		OmitComparePragmaStatement:        true,
		OmitToStringPragmaStatement:       true,
		OmitToReprPragmaStatement:         true,
		StaticConstantStatement:           true,
		InstanceConstantStatement:         true,
		StaticFieldStatement:              true,
		StructFieldStatement:              true,
	},
	UnionStatementContext: {
		AlignPragmaStatement:              true,
		MinimumSizePragmaStatement:        true,
		PreserveFieldOrderPragmaStatement: true,
		OmitNewPragmaStatement:            true,
		OmitCopyPragmaStatement:           true,
		OmitMovePragmaStatement:           true,
		OmitHashPragmaStatement:           true,
		OmitComparePragmaStatement:        true,
		OmitToStringPragmaStatement:       true,
		OmitToReprPragmaStatement:         true,
		StaticConstantStatement:           true,
		InstanceConstantStatement:         true,
		StaticFieldStatement:              true,
		UnionTagStatement:                 true,
		UnionFieldStatement:               true,
	},
	InterfaceStatementContext: {
		InterfaceFieldStatement:    true,
		InterfacePropertyStatement: true,
		InterfaceMethodStatement:   true,
	},
	FunctionStatementContext: {},
}

var statementsNoDuplicatesAllowed = map[StatementKind]bool{
	AlignPragmaStatement:              true,
	MinimumSizePragmaStatement:        true,
	PreserveFieldOrderPragmaStatement: true,
	OmitNewPragmaStatement:            true,
	OmitCopyPragmaStatement:           true,
	OmitMovePragmaStatement:           true,
	OmitHashPragmaStatement:           true,
	OmitComparePragmaStatement:        true,
	OmitToStringPragmaStatement:       true,
	OmitToReprPragmaStatement:         true,
	EnumKindStatement:                 true,
	BitfieldKindStatement:             true,
	UnionTagStatement:                 true,
}
