package token

type statePredicate interface {
	MatchState(state State) bool
}

// anyState
// {{{

type anyStatePredicate struct{}

func (pred anyStatePredicate) MatchState(state State) bool {
	return true
}

func anyState() statePredicate {
	return anyStatePredicate{}
}

// }}}

// exactState
// {{{

type exactStatePredicate struct {
	State State
}

func (pred exactStatePredicate) MatchState(state State) bool {
	return state == pred.State
}

func exactState(state State) statePredicate {
	return exactStatePredicate{state}
}

// }}}

// someState
// {{{

type someStatePredicate struct {
	List []State
}

func (pred someStatePredicate) MatchState(state State) bool {
	for _, target := range pred.List {
		if state == target {
			return true
		}
	}
	return false
}

func someState(list ...State) statePredicate {
	return someStatePredicate{list}
}

// }}}
