package exprtree

// Interface: Expr
// {{{

type Expr interface {
	Interp() *Interp
	Type() *Type
	Walk(to TraversalOrder, fn func(Expr))
	Eval() (*Value, error)
}

// }}}
