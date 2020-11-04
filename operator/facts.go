package operator

type Facts struct {
	Kind         Kind
	Overloadable bool
	GoName       string
	SimpleName   string
	OperatorName string
	PythonicName string
}

func (facts Facts) IsUnary() bool {
	return facts.Kind.IsUnary()
}

func (facts Facts) IsBinary() bool {
	return facts.Kind.IsBinary()
}

func (facts Facts) IsTernary() bool {
	return facts.Kind.IsTernary()
}

func (facts Facts) IsAssignment() bool {
	return facts.Kind.IsAssignment()
}

func (facts Facts) IsMutation() bool {
	return facts.Kind.IsMutation()
}
