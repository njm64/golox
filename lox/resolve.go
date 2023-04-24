package lox

import (
	"golox/lox/ast"
	"golox/lox/tok"
)

type Scope map[string]bool

type Resolver struct {
	scopes []Scope
}

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) ResolveStatements(statements []ast.Stmt) {
	for _, st := range statements {
		r.ResolveStatement(st)
	}
}

func (r *Resolver) ResolveStatement(st ast.Stmt) {
	switch s := st.(type) {
	case *ast.Function:
		r.declare(s.Name)
		r.define(s.Name)
		r.resolveFunction(s)
	case *ast.Block:
		r.beginScope()
		r.ResolveStatements(s.Statements)
		r.endScope()
	case *ast.Var:
		r.declare(s.Name)
		if s.Initializer != nil {
			r.ResolveExpression(s.Initializer)
		}
		r.define(s.Name)
	case *ast.Expression:
		r.ResolveExpression(s.Expression)
	case *ast.If:
		r.ResolveExpression(s.Condition)
		r.ResolveStatement(s.ThenBranch)
		if s.ElseBranch != nil {
			r.ResolveStatement(s.ElseBranch)
		}
	case *ast.Print:
		r.ResolveExpression(s.Expression)
	case *ast.Return:
		if s.Value != nil {
			r.ResolveExpression(s.Value)
		}
	case *ast.While:
		r.ResolveExpression(s.Condition)
		r.ResolveStatement(s.Body)
	}
}

func (r *Resolver) ResolveExpression(ex ast.Expr) {
	switch e := ex.(type) {
	case *ast.Variable:
		r.variableExpr(e)
	case *ast.Assign:
		r.ResolveExpression(e.Value)
		r.resolveLocal(e, e.Name)
	case *ast.Binary:
		r.ResolveExpression(e.Left)
		r.ResolveExpression(e.Right)
	case *ast.Call:
		r.ResolveExpression(e.Callee)
		for _, arg := range e.Arguments {
			r.ResolveExpression(arg)
		}
	case *ast.Grouping:
		r.ResolveExpression(e.Expression)
	case *ast.Literal:
	case *ast.Logical:
		r.ResolveExpression(e.Left)
		r.ResolveExpression(e.Right)
	case *ast.Unary:
		r.ResolveExpression(e.Right)
	}
}

func (r *Resolver) variableExpr(e *ast.Variable) {
	if len(r.scopes) > 0 {
		defined, declared := r.peekScope()[e.Name.Lexeme]
		if declared && !defined {
			ReportParseError(&Error{Token: e.Name, Message: "Can't read local variable in its own initializer"})
		}
	}
	r.resolveLocal(e, e.Name)
}

func (r *Resolver) resolveLocal(e ast.Expr, name *tok.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		_, declared := r.scopes[i][name.Lexeme]
		if declared {
			depthMap[e] = len(r.scopes) - i - 1
			return
		}
	}
}

func (r *Resolver) resolveFunction(s *ast.Function) {
	r.beginScope()
	for _, param := range s.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStatements(s.Body)
	r.endScope()
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(Scope))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) peekScope() Scope {
	return r.scopes[len(r.scopes)-1]
}

func (r *Resolver) declare(name *tok.Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.peekScope()[name.Lexeme] = false
}

func (r *Resolver) define(name *tok.Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.peekScope()[name.Lexeme] = true
}
