package stmt

import (
	"golox/lox/expr"
)

type Stmt interface {
	stmt()
}

func (s *Expression) stmt() {}
func (s *Print) stmt()      {}

type Expression struct {
	Expression expr.Expr
}

type Print struct {
	Expression expr.Expr
}
