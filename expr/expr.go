package expr

import (
	"golox/tok"
)

type Expr interface {
	expr()
}

func (e *Binary) expr()   {}
func (e *Grouping) expr() {}
func (e *Literal) expr()  {}
func (e *Unary) expr()    {}

type Binary struct {
	Left     Expr
	Operator tok.Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value any
}

type Unary struct {
	Operator tok.Token
	Right    Expr
}
