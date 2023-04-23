package ast

import (
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

type Expression struct {
	Expression Expr
}

type Print struct {
	Expression Expr
}

type Var struct {
	Name        *tok.Token
	Initializer Expr
}

type Block struct {
	Statements []Stmt
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type While struct {
	Condition Expr
	Body      Stmt
}

type Function struct {
	Name   *tok.Token
	Params []*tok.Token
	Body   []Stmt
}
