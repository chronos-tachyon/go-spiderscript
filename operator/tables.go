package operator

var factsMap = map[Operator]Facts{
	InvalidOperator: {
		Kind:       InvalidKind,
		GoName:     "InvalidOperator",
		SimpleName: "???",
	},
	AddressOf: {
		Kind:       UnaryPrefix,
		GoName:     "AddressOf",
		SimpleName: "&",
	},
	DerefPointer: {
		Kind:       UnaryPrefix,
		GoName:     "DerefPointer",
		SimpleName: "*",
	},
	UnaryPos: {
		Kind:       UnaryPrefix,
		GoName:     "UnaryPos",
		SimpleName: "+",

		Overloadable: true,
		OperatorName: "+",
		PythonicName: "__pos",
	},
	UnaryNeg: {
		Kind:       UnaryPrefix,
		GoName:     "UnaryNeg",
		SimpleName: "-",

		Overloadable: true,
		OperatorName: "-",
		PythonicName: "__neg",
	},
	BitwiseNOT: {
		Kind:       UnaryPrefix,
		GoName:     "BitwiseNOT",
		SimpleName: "~",

		Overloadable: true,
		OperatorName: "~",
		PythonicName: "__invert",
	},
	LogicalNOT: {
		Kind:       UnaryPrefix,
		GoName:     "LogicalNOT",
		SimpleName: "!",
	},
	Splat: {
		Kind:       UnaryPostfix,
		GoName:     "Splat",
		SimpleName: "...",
	},
	NonZero: {
		Kind:       UnaryOther,
		GoName:     "NonZero",
		SimpleName: "nonzero()",

		Overloadable: true,
		OperatorName: "!!",
		PythonicName: "__nonzero",
	},
	HashCode: {
		Kind:       UnaryOther,
		GoName:     "HashCode",
		SimpleName: "hashcode()",

		Overloadable: true,
		OperatorName: "hashcode",
		PythonicName: "__hashcode",
	},
	Add: {
		Kind:       BinaryInfix,
		GoName:     "Add",
		SimpleName: "+",

		Overloadable: true,
		OperatorName: "+",
		PythonicName: "__add",
	},
	Sub: {
		Kind:       BinaryInfix,
		GoName:     "Sub",
		SimpleName: "-",

		Overloadable: true,
		OperatorName: "-",
		PythonicName: "__sub",
	},
	Mul: {
		Kind:       BinaryInfix,
		GoName:     "Mul",
		SimpleName: "*",

		Overloadable: true,
		OperatorName: "*",
		PythonicName: "__mul",
	},
	Div: {
		Kind:       BinaryInfix,
		GoName:     "Div",
		SimpleName: "/",

		Overloadable: true,
		OperatorName: "/",
		PythonicName: "__div",
	},
	Mod: {
		Kind:       BinaryInfix,
		GoName:     "Mod",
		SimpleName: "%",

		Overloadable: true,
		OperatorName: "%",
		PythonicName: "__mod",
	},
	DivMod: {
		Kind:       BinaryInfix,
		GoName:     "DivMod",
		SimpleName: "/%",

		Overloadable: true,
		OperatorName: "/%",
		PythonicName: "__divmod",
	},
	Pow: {
		Kind:       BinaryInfix,
		GoName:     "Pow",
		SimpleName: "**",

		Overloadable: true,
		OperatorName: "**",
		PythonicName: "__pow",
	},
	LShift: {
		Kind:       BinaryInfix,
		GoName:     "LShift",
		SimpleName: "<<",

		Overloadable: true,
		OperatorName: "<<",
		PythonicName: "__lshift",
	},
	RShift: {
		Kind:       BinaryInfix,
		GoName:     "RShift",
		SimpleName: ">>",

		Overloadable: true,
		OperatorName: ">>",
		PythonicName: "__rshift",
	},
	LRotate: {
		Kind:       BinaryInfix,
		GoName:     "LRotate",
		SimpleName: "<<|",

		Overloadable: true,
		OperatorName: "<<|",
		PythonicName: "__lrotate",
	},
	RRotate: {
		Kind:       BinaryInfix,
		GoName:     "RRotate",
		SimpleName: ">>|",

		Overloadable: true,
		OperatorName: ">>|",
		PythonicName: "__rrotate",
	},
	BitwiseAND: {
		Kind:       BinaryInfix,
		GoName:     "BitwiseAND",
		SimpleName: "&",

		Overloadable: true,
		OperatorName: "&",
		PythonicName: "__and",
	},
	BitwiseXOR: {
		Kind:       BinaryInfix,
		GoName:     "BitwiseXOR",
		SimpleName: "^",

		Overloadable: true,
		OperatorName: "^",
		PythonicName: "__xor",
	},
	BitwiseOR: {
		Kind:       BinaryInfix,
		GoName:     "BitwiseOR",
		SimpleName: "|",

		Overloadable: true,
		OperatorName: "|",
		PythonicName: "__or",
	},
	CmpCMP: {
		Kind:       BinaryInfix,
		GoName:     "CmpCMP",
		SimpleName: "<=>",

		Overloadable: true,
		OperatorName: "<=>",
		PythonicName: "__cmp",
	},
	CmpEQ: {
		Kind:       BinaryInfix,
		GoName:     "CmpEQ",
		SimpleName: "==",

		Overloadable: true,
		OperatorName: "==",
		PythonicName: "__eq",
	},
	CmpNE: {
		Kind:       BinaryInfix,
		GoName:     "CmpNE",
		SimpleName: "!=",

		Overloadable: true,
		OperatorName: "!=",
		PythonicName: "__ne",
	},
	CmpLT: {
		Kind:       BinaryInfix,
		GoName:     "CmpLT",
		SimpleName: "<",

		Overloadable: true,
		OperatorName: "<",
		PythonicName: "__lt",
	},
	CmpLE: {
		Kind:       BinaryInfix,
		GoName:     "CmpLE",
		SimpleName: "<=",

		Overloadable: true,
		OperatorName: "<=",
		PythonicName: "__le",
	},
	CmpGT: {
		Kind:       BinaryInfix,
		GoName:     "CmpGT",
		SimpleName: ">",

		Overloadable: true,
		OperatorName: ">",
		PythonicName: "__gt",
	},
	CmpGE: {
		Kind:       BinaryInfix,
		GoName:     "CmpGE",
		SimpleName: ">=",

		Overloadable: true,
		OperatorName: ">=",
		PythonicName: "__ge",
	},
	LogicalAND: {
		Kind:       BinaryInfix,
		GoName:     "LogicalAND",
		SimpleName: "&&",
	},
	LogicalXOR: {
		Kind:       BinaryInfix,
		GoName:     "LogicalXOR",
		SimpleName: "^^",
	},
	LogicalOR: {
		Kind:       BinaryInfix,
		GoName:     "LogicalOR",
		SimpleName: "||",
	},
	Range: {
		Kind:       BinaryInfix,
		GoName:     "Range",
		SimpleName: "..",

		Overloadable: true,
		OperatorName: "..",
		PythonicName: "__range",
	},
	Slice: {
		Kind:       BinaryInfix,
		GoName:     "Slice",
		SimpleName: ":",

		Overloadable: true,
		OperatorName: ":",
		PythonicName: "__slice",
	},
	Elvis: {
		Kind:       BinaryInfix,
		GoName:     "Elvis",
		SimpleName: "?:",
	},
	Index: {
		Kind:       BinaryOther,
		GoName:     "Index",
		SimpleName: "[]",

		Overloadable: true,
		OperatorName: "[]",
		PythonicName: "__index",
	},
	Call: {
		Kind:       BinaryOther,
		GoName:     "Call",
		SimpleName: "()",

		Overloadable: true,
		OperatorName: "()",
		PythonicName: "__call",
	},
	Ternary: {
		Kind:       TernaryOther,
		GoName:     "Ternary",
		SimpleName: "a ? b : c",
	},
	IfElseExpr: {
		Kind:       TernaryOther,
		GoName:     "IfElseExpr",
		SimpleName: "b if a else c",
	},
	MutateBitwiseNOT: {
		Kind:       MutateStatement,
		GoName:     "MutateBitwiseNOT",
		SimpleName: "~~",

		Overloadable: true,
		OperatorName: "~~",
		PythonicName: "__iinvert",
	},
	MutateINC: {
		Kind:       MutateStatement,
		GoName:     "MutateINC",
		SimpleName: "++",

		Overloadable: true,
		OperatorName: "++",
		PythonicName: "__inc",
	},
	MutateDEC: {
		Kind:       MutateStatement,
		GoName:     "MutateDEC",
		SimpleName: "--",

		Overloadable: true,
		OperatorName: "--",
		PythonicName: "__dec",
	},
	Assign: {
		Kind:       AssignStatement,
		GoName:     "Assign",
		SimpleName: "=",
	},
	DeclareAndAssign: {
		Kind:       AssignStatement,
		GoName:     "DeclareAndAssign",
		SimpleName: ":=",
	},
	AssignAdd: {
		Kind:       AssignStatement,
		GoName:     "AssignAdd",
		SimpleName: "+=",

		Overloadable: true,
		OperatorName: "+=",
		PythonicName: "__iadd",
	},
	AssignSub: {
		Kind:       AssignStatement,
		GoName:     "AssignSub",
		SimpleName: "-=",

		Overloadable: true,
		OperatorName: "-=",
		PythonicName: "__isub",
	},
	AssignMul: {
		Kind:       AssignStatement,
		GoName:     "AssignMul",
		SimpleName: "*=",

		Overloadable: true,
		OperatorName: "*=",
		PythonicName: "__imul",
	},
	AssignDiv: {
		Kind:       AssignStatement,
		GoName:     "AssignDiv",
		SimpleName: "/=",

		Overloadable: true,
		OperatorName: "/=",
		PythonicName: "__idiv",
	},
	AssignMod: {
		Kind:       AssignStatement,
		GoName:     "AssignMod",
		SimpleName: "%=",

		Overloadable: true,
		OperatorName: "%=",
		PythonicName: "__imod",
	},
	AssignPow: {
		Kind:       AssignStatement,
		GoName:     "AssignPow",
		SimpleName: "**=",

		Overloadable: true,
		OperatorName: "**=",
		PythonicName: "__ipow",
	},
	AssignLShift: {
		Kind:       AssignStatement,
		GoName:     "AssignLShift",
		SimpleName: "<<=",

		Overloadable: true,
		OperatorName: "<<=",
		PythonicName: "__ilshift",
	},
	AssignRShift: {
		Kind:       AssignStatement,
		GoName:     "AssignRShift",
		SimpleName: ">>=",

		Overloadable: true,
		OperatorName: ">>=",
		PythonicName: "__irshift",
	},
	AssignLRotate: {
		Kind:       AssignStatement,
		GoName:     "AssignLRotate",
		SimpleName: "<<|=",

		Overloadable: true,
		OperatorName: "<<|=",
		PythonicName: "__ilrotate",
	},
	AssignRRotate: {
		Kind:       AssignStatement,
		GoName:     "AssignRRotate",
		SimpleName: ">>|=",

		Overloadable: true,
		OperatorName: ">>|=",
		PythonicName: "__irrotate",
	},
	AssignBitwiseAND: {
		Kind:       AssignStatement,
		GoName:     "AssignBitwiseAND",
		SimpleName: "&=",

		Overloadable: true,
		OperatorName: "&=",
		PythonicName: "__iand",
	},
	AssignBitwiseXOR: {
		Kind:       AssignStatement,
		GoName:     "AssignBitwiseXOR",
		SimpleName: "^=",

		Overloadable: true,
		OperatorName: "^=",
		PythonicName: "__ixor",
	},
	AssignBitwiseOR: {
		Kind:       AssignStatement,
		GoName:     "AssignBitwiseOR",
		SimpleName: "|=",

		Overloadable: true,
		OperatorName: "|=",
		PythonicName: "__ior",
	},
	AssignLogicalAND: {
		Kind:       AssignStatement,
		GoName:     "AssignLogicalAND",
		SimpleName: "&&=",
	},
	AssignLogicalXOR: {
		Kind:       AssignStatement,
		GoName:     "AssignLogicalXOR",
		SimpleName: "^^=",
	},
	AssignLogicalOR: {
		Kind:       AssignStatement,
		GoName:     "AssignLogicalOR",
		SimpleName: "||=",
	},
	AssignElvis: {
		Kind:       AssignStatement,
		GoName:     "AssignElvis",
		SimpleName: "?=",
	},
}

var operatorNameMap map[string]Operator
var pythonicNameMap map[string]Operator

func init() {
	operatorNameMap = make(map[string]Operator, len(factsMap))
	pythonicNameMap = make(map[string]Operator, len(factsMap))
	for op, data := range factsMap {
		if data.Overloadable {
			operatorNameMap[data.OperatorName] = op
			pythonicNameMap[data.PythonicName] = op
		}
	}
}
