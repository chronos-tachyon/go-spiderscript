package value

import (
	"strconv"
	"strings"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
)

type FormatSpecification struct {
	HasTick  bool
	HasHash  bool
	HasPlus  bool
	HasMinus bool
	HasSpace bool
	HasZero  bool

	HasWidth                bool
	WidthIsExternal         bool
	WidthArgumentIsExplicit bool
	WidthArgumentIsNamed    bool

	HasPrecision                bool
	PrecisionIsExternal         bool
	PrecisionArgumentIsExplicit bool
	PrecisionArgumentIsNamed    bool

	ValueArgumentIsExplicit bool
	ValueArgumentIsNamed    bool

	Conversion byte

	PrecomputedStringLength   uint
	PrecomputedGoStringLength uint
	FixedWidth                uint
	FixedPrecision            uint
	WidthArgumentIndex        uint
	PrecisionArgumentIndex    uint
	ValueArgumentIndex        uint

	WidthArgumentName     string
	PrecisionArgumentName string
	ValueArgumentName     string
}

func (spec *FormatSpecification) String() string {
	return util.StringImpl(spec)
}

func (spec *FormatSpecification) GoString() string {
	return util.GoStringImpl(spec)
}

func (spec *FormatSpecification) EstimateStringLength() uint {
	return spec.PrecomputedStringLength
}

func (spec *FormatSpecification) EstimateGoStringLength() uint {
	return spec.PrecomputedGoStringLength
}

func (spec *FormatSpecification) WriteStringTo(out *strings.Builder) {
	out.WriteByte('%')

	if spec.HasTick {
		out.WriteByte('\'')
	}
	if spec.HasHash {
		out.WriteByte('#')
	}
	if spec.HasPlus {
		out.WriteByte('+')
	}
	if spec.HasMinus {
		out.WriteByte('-')
	}
	if spec.HasSpace {
		out.WriteByte(' ')
	}
	if spec.HasZero {
		out.WriteByte('0')
	}

	if !spec.HasWidth {
		// pass
	} else if !spec.WidthIsExternal {
		out.WriteString(util.Itoa(spec.FixedWidth))
	} else if !spec.WidthArgumentIsExplicit {
		out.WriteByte('*')
	} else if !spec.WidthArgumentIsNamed {
		out.WriteByte('[')
		out.WriteString(util.Itoa(spec.WidthArgumentIndex))
		out.WriteByte(']')
		out.WriteByte('*')
	} else {
		out.WriteByte('[')
		out.WriteString(spec.WidthArgumentName)
		out.WriteByte(']')
		out.WriteByte('*')
	}

	if !spec.HasPrecision {
		// pass
	} else if !spec.PrecisionIsExternal {
		out.WriteByte('.')
		out.WriteString(util.Itoa(spec.FixedPrecision))
	} else if !spec.PrecisionArgumentIsExplicit {
		out.WriteByte('.')
		out.WriteByte('*')
	} else if !spec.PrecisionArgumentIsNamed {
		out.WriteByte('.')
		out.WriteByte('[')
		out.WriteString(util.Itoa(spec.PrecisionArgumentIndex))
		out.WriteByte(']')
		out.WriteByte('*')
	} else {
		out.WriteByte('.')
		out.WriteByte('[')
		out.WriteString(spec.PrecisionArgumentName)
		out.WriteByte(']')
		out.WriteByte('*')
	}

	if !spec.ValueArgumentIsExplicit {
		// pass
	} else if !spec.ValueArgumentIsNamed {
		out.WriteByte('[')
		out.WriteString(util.Itoa(spec.ValueArgumentIndex))
		out.WriteByte(']')
	} else {
		out.WriteByte('[')
		out.WriteString(spec.ValueArgumentName)
		out.WriteByte(']')
	}

	out.WriteByte(spec.Conversion)
}

func (spec *FormatSpecification) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("&value.FormatSpecification{")

	if spec.HasTick {
		out.WriteString("'\\'',")
	}
	if spec.HasHash {
		out.WriteString("'#',")
	}
	if spec.HasPlus {
		out.WriteString("'+',")
	}
	if spec.HasMinus {
		out.WriteString("'-',")
	}
	if spec.HasSpace {
		out.WriteString("' ',")
	}
	if spec.HasZero {
		out.WriteString("'0',")
	}

	if !spec.HasWidth {
		out.WriteByte('_')
	} else if !spec.WidthIsExternal {
		out.WriteByte('W')
		out.WriteByte('/')
		out.WriteString(util.Itoa(spec.FixedWidth))
	} else if !spec.WidthArgumentIsExplicit {
		out.WriteByte('W')
		out.WriteByte('/')
		out.WriteByte('*')
	} else if !spec.WidthArgumentIsNamed {
		out.WriteByte('W')
		out.WriteByte('/')
		out.WriteByte('[')
		out.WriteString(util.Itoa(spec.WidthArgumentIndex))
		out.WriteByte(']')
		out.WriteByte('*')
	} else {
		out.WriteByte('W')
		out.WriteByte('/')
		out.WriteByte('[')
		out.WriteString(strconv.Quote(spec.WidthArgumentName))
		out.WriteByte(']')
		out.WriteByte('*')
	}
	out.WriteByte(',')

	if !spec.HasPrecision {
		out.WriteByte('_')
	} else if !spec.PrecisionIsExternal {
		out.WriteByte('P')
		out.WriteByte('/')
		out.WriteString(util.Itoa(spec.FixedPrecision))
	} else if !spec.PrecisionArgumentIsExplicit {
		out.WriteByte('P')
		out.WriteByte('/')
		out.WriteByte('*')
	} else if !spec.PrecisionArgumentIsNamed {
		out.WriteByte('P')
		out.WriteByte('/')
		out.WriteByte('[')
		out.WriteString(util.Itoa(spec.PrecisionArgumentIndex))
		out.WriteByte(']')
		out.WriteByte('*')
	} else {
		out.WriteByte('P')
		out.WriteByte('/')
		out.WriteByte('[')
		out.WriteString(strconv.Quote(spec.PrecisionArgumentName))
		out.WriteByte(']')
		out.WriteByte('*')
	}
	out.WriteByte(',')

	if !spec.ValueArgumentIsExplicit {
		out.WriteByte(spec.Conversion)
	} else if !spec.ValueArgumentIsNamed {
		out.WriteByte(spec.Conversion)
		out.WriteByte('/')
		out.WriteByte('[')
		out.WriteString(util.Itoa(spec.ValueArgumentIndex))
		out.WriteByte(']')
	} else {
		out.WriteByte(spec.Conversion)
		out.WriteByte('/')
		out.WriteByte('[')
		out.WriteString(strconv.Quote(spec.ValueArgumentName))
		out.WriteByte(']')
	}

	out.WriteByte('}')
}

func (spec *FormatSpecification) SetEstimatedStringLength(length uint) {
	spec.PrecomputedStringLength = length
}

func (spec *FormatSpecification) SetEstimatedGoStringLength(length uint) {
	spec.PrecomputedGoStringLength = length
}

func (spec *FormatSpecification) GetFlag(flagCh rune) bool {
	switch flagCh {
	case '\'':
		return spec.HasTick
	case '#':
		return spec.HasHash
	case '+':
		return spec.HasPlus
	case '-':
		return spec.HasMinus
	case ' ':
		return spec.HasSpace
	case '0':
		return spec.HasZero
	default:
		return false
	}
}

func (spec *FormatSpecification) SetFlag(flagCh rune, boolValue bool) {
	switch flagCh {
	case '\'':
		spec.HasTick = boolValue
	case '#':
		spec.HasHash = boolValue
	case '+':
		spec.HasPlus = boolValue
	case '-':
		spec.HasMinus = boolValue
	case ' ':
		spec.HasSpace = boolValue
	case '0':
		spec.HasZero = boolValue
	}
}

var _ Value = (*FormatSpecification)(nil)
var _ util.Estimable = (*FormatSpecification)(nil)
