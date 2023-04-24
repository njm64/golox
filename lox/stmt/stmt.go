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
func (s *Block) stmt()      {}
func (s *If) stmt()         {}
func (s *While) stmt()      {}
func (e *Function) stmt()   {}
func (e *Return) stmt()     {}

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

type Block struct {
	Statements []Stmt
}

type If struct {
	Condition  expr.Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type While struct {
	Condition expr.Expr
	Body      Stmt
}

type Function struct {
	Name   *tok.Token
	Params []*tok.Token
	Body   []Stmt
}

type Return struct {
	Keyword *tok.Token
	Value   expr.Expr
}
