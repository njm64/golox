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
func (e *Call) expr()     {}
func (e *Get) expr()      {}
func (e *Set) expr()      {}
func (e *This) expr()     {}
func (e *Super) expr()    {}

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

type Call struct {
	Callee    Expr
	Paren     *tok.Token
	Arguments []Expr
}

type Get struct {
	Object Expr
	Name   *tok.Token
}

type Set struct {
	Object Expr
	Name   *tok.Token
	Value  Expr
}

type This struct {
	Keyword *tok.Token
}

type Super struct {
	Keyword *tok.Token
	Method  *tok.Token
}
