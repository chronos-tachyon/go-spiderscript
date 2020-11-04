package operator

import (
	"fmt"
)

type Operator uint16

const (
	InvalidOperator  Operator = 0x0000
	AddressOf        Operator = 0x0100
	DerefPointer     Operator = 0x0101
	UnaryPos         Operator = 0x0102
	UnaryNeg         Operator = 0x0103
	BitwiseNOT       Operator = 0x0104
	LogicalNOT       Operator = 0x0105
	Splat            Operator = 0x0106
	NonZero          Operator = 0x01fc
	HashCode         Operator = 0x01fd
	EnterBlock       Operator = 0x01fe
	ExitBlock        Operator = 0x01ff
	Add              Operator = 0x0200
	Sub              Operator = 0x0201
	Mul              Operator = 0x0202
	Div              Operator = 0x0203
	Mod              Operator = 0x0204
	DivMod           Operator = 0x0205
	Pow              Operator = 0x0206
	LShift           Operator = 0x0207
	RShift           Operator = 0x0208
	LRotate          Operator = 0x0209
	RRotate          Operator = 0x020a
	BitwiseAND       Operator = 0x020b
	BitwiseXOR       Operator = 0x020c
	BitwiseOR        Operator = 0x020d
	CmpCMP           Operator = 0x020e
	CmpEQ            Operator = 0x020f
	CmpNE            Operator = 0x0210
	CmpLT            Operator = 0x0211
	CmpLE            Operator = 0x0212
	CmpGT            Operator = 0x0213
	CmpGE            Operator = 0x0214
	LogicalAND       Operator = 0x0215
	LogicalXOR       Operator = 0x0216
	LogicalOR        Operator = 0x0217
	Range            Operator = 0x0218
	Slice            Operator = 0x0219
	Elvis            Operator = 0x021a
	Index            Operator = 0x02fe
	Call             Operator = 0x02ff
	Ternary          Operator = 0x0300
	IfElseExpr       Operator = 0x0301
	Assign           Operator = 0x0400
	DeclareAndAssign Operator = 0x0401
	AssignAdd        Operator = 0x0402
	AssignSub        Operator = 0x0403
	AssignMul        Operator = 0x0404
	AssignDiv        Operator = 0x0405
	AssignMod        Operator = 0x0406
	AssignPow        Operator = 0x0407
	AssignLShift     Operator = 0x0408
	AssignRShift     Operator = 0x0409
	AssignLRotate    Operator = 0x040a
	AssignRRotate    Operator = 0x040b
	AssignBitwiseAND Operator = 0x040c
	AssignBitwiseXOR Operator = 0x040d
	AssignBitwiseOR  Operator = 0x040e
	AssignLogicalAND Operator = 0x040f
	AssignLogicalXOR Operator = 0x0410
	AssignLogicalOR  Operator = 0x0411
	AssignElvis      Operator = 0x0412
	MutateBitwiseNOT Operator = 0x0500
	MutateINC        Operator = 0x0501
	MutateDEC        Operator = 0x0502
)

func (op Operator) Facts() Facts {
	if data, found := factsMap[op]; found {
		return data
	}
	str := fmt.Sprintf("Operator(%#04x)", uint(op))
	return Facts{
		Kind:       InvalidKind,
		GoName:     str,
		SimpleName: str,
	}
}

func (op Operator) GoString() string {
	return op.Facts().GoName
}

func (op Operator) String() string {
	return op.Facts().SimpleName
}

func (op Operator) Kind() Kind {
	return op.Facts().Kind
}

func (op Operator) IsUnary() bool {
	return op.Facts().Kind.IsUnary()
}

func (op Operator) IsBinary() bool {
	return op.Facts().Kind.IsBinary()
}

func (op Operator) IsTernary() bool {
	return op.Facts().Kind.IsTernary()
}

func (op Operator) IsAssignment() bool {
	return op.Facts().Kind.IsAssignment()
}

func (op Operator) IsMutation() bool {
	return op.Facts().Kind.IsMutation()
}

func (op Operator) IsOverloadable() bool {
	return op.Facts().Overloadable
}

func (op Operator) OperatorName() string {
	return op.Facts().OperatorName
}

func (op Operator) PythonicName() string {
	return op.Facts().PythonicName
}

var _ fmt.Stringer = Operator(0)
var _ fmt.GoStringer = Operator(0)

func ByOperatorName(str string) (Operator, bool) {
	op, found := operatorNameMap[str]
	return op, found
}

func ByPythonicName(str string) (Operator, bool) {
	op, found := pythonicNameMap[str]
	return op, found
}
