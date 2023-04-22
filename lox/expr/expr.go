package expr

import (
	"golox/lox/tok"
)

type Expr interface {
	expr()
}

func (e *Binary) expr()   {}
func (e *Grouping) expr() {}
func (e *Literal) expr()  {}
func (e *Unary) expr()    {}
func (e *Variable) expr() {}
func (e *Assign) expr()   {}
func (e *Logical) expr()  {}

type Binary struct {
	Left     Expr
	Operator *tok.Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value any
}

type Unary struct {
	Operator *tok.Token
	Right    Expr
}

type Variable struct {
	Name *tok.Token
}

type Assign struct {
	Name  *tok.Token
	Value Expr
}

type Logical struct {
	Left     Expr
	Operator *tok.Token
	Right    Expr
}
