package runepredicate

type RunePredicate interface {
	MatchRune(ch rune) bool
}

// None
// {{{

type nonePredicate struct{}

func (pred nonePredicate) MatchRune(ch rune) bool {
	return false
}

func None() RunePredicate {
	return nonePredicate{}
}

// }}}

// Any
// {{{

type anyPredicate struct{}

func (pred anyPredicate) MatchRune(ch rune) bool {
	return true
}

func Any() RunePredicate {
	return anyPredicate{}
}

// }}}

// Exactly
// {{{

type exactlyPredicate struct {
	Rune rune
}

func (pred exactlyPredicate) MatchRune(ch rune) bool {
	return ch == pred.Rune
}

func Exactly(ch rune) RunePredicate {
	return exactlyPredicate{ch}
}

// }}}

// EitherOf
// {{{

type eitherOfPredicate struct {
	A rune
	B rune
}

func (pred eitherOfPredicate) MatchRune(ch rune) bool {
	return ch == pred.A || ch == pred.B
}

func EitherOf(a, b rune) RunePredicate {
	return eitherOfPredicate{a, b}
}

// }}}

// OneOf
// {{{

type oneOfPredicate struct {
	List []rune
}

func (pred oneOfPredicate) MatchRune(ch rune) bool {
	for _, candidate := range pred.List {
		if ch == candidate {
			return true
		}
	}
	return false
}

func OneOf(list ...rune) RunePredicate {
	if len(list) == 0 {
		return nonePredicate{}
	}
	if len(list) == 1 {
		return exactlyPredicate{list[0]}
	}
	if len(list) == 2 {
		return eitherOfPredicate{list[0], list[1]}
	}
	return oneOfPredicate{list}
}

// }}}

// Not
// {{{

type notPredicate struct {
	Child RunePredicate
}

func (pred notPredicate) MatchRune(ch rune) bool {
	return !pred.Child.MatchRune(ch)
}

func Not(child RunePredicate) RunePredicate {
	return notPredicate{child}
}

// }}}

// Func
// {{{

type funcPredicate struct {
	Func func(rune) bool
}

func (pred funcPredicate) MatchRune(ch rune) bool {
	return pred.Func(ch)
}

func Func(fn func(rune) bool) RunePredicate {
	return funcPredicate{fn}
}

// }}}
