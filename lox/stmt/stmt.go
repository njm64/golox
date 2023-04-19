package stmt

import (
	"golox/lox/expr"
	"golox/lox/tok"
)

type Stmt interface {
	stmt()
}

func (s *Expression) stmt() {}
func (s *Print) stmt()      {}
func (s *Var) stmt()        {}

type Expression struct {
	Expression expr.Expr
}

type Print struct {
	Expression expr.Expr
}

type Var struct {
	Name        *tok.Token
	Initializer expr.Expr
}
