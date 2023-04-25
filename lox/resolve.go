package lox

import (
	"golox/lox/expr"
	"golox/lox/stmt"
	"golox/lox/tok"
)

type Scope map[string]bool

type FunctionType int

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
	FunctionTypeMethod
	FunctionTypeInitializer
)

type ClassType int

const (
	ClassTypeNone ClassType = iota
	ClassTypeClass
)

type Resolver struct {
	scopes          []Scope
	currentFunction FunctionType
	currentClass    ClassType
}

func NewResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) ResolveStatements(statements []stmt.Stmt) {
	for _, st := range statements {
		r.ResolveStatement(st)
	}
}

func (r *Resolver) ResolveStatement(st stmt.Stmt) {
	switch s := st.(type) {
	case *stmt.Function:
		r.declare(s.Name)
		r.define(s.Name)
		r.resolveFunction(s, FunctionTypeFunction)
	case *stmt.Block:
		r.beginScope()
		r.ResolveStatements(s.Statements)
		r.endScope()
	case *stmt.Var:
		r.declare(s.Name)
		if s.Initializer != nil {
			r.ResolveExpression(s.Initializer)
		}
		r.define(s.Name)
	case *stmt.Expression:
		r.ResolveExpression(s.Expression)
	case *stmt.If:
		r.ResolveExpression(s.Condition)
		r.ResolveStatement(s.ThenBranch)
		if s.ElseBranch != nil {
			r.ResolveStatement(s.ElseBranch)
		}
	case *stmt.Print:
		r.ResolveExpression(s.Expression)
	case *stmt.Return:
		r.returnStmt(s)
	case *stmt.While:
		r.ResolveExpression(s.Condition)
		r.ResolveStatement(s.Body)
	case *stmt.Class:
		r.classStmt(s)
	}
}

func (r *Resolver) returnStmt(s *stmt.Return) {
	if r.currentFunction == FunctionTypeNone {
		ReportParseError(&Error{
			Token:   s.Keyword,
			Message: "Can't return from top-level code",
		})
	}

	if s.Value != nil {
		if r.currentFunction == FunctionTypeInitializer {
			ReportParseError(&Error{
				Token:   s.Keyword,
				Message: "Can't return a value from an initializer",
			})
		}
		r.ResolveExpression(s.Value)
	}
}

func (r *Resolver) classStmt(s *stmt.Class) {
	enclosingClass := r.currentClass
	r.currentClass = ClassTypeClass

	r.declare(s.Name)
	r.define(s.Name)

	r.beginScope()
	r.peekScope()["this"] = true

	for _, m := range s.Methods {
		if m.Name.Lexeme == "init" {
			r.resolveFunction(m, FunctionTypeInitializer)
		} else {
			r.resolveFunction(m, FunctionTypeMethod)
		}
	}

	r.endScope()
	r.currentClass = enclosingClass
}

func (r *Resolver) ResolveExpression(ex expr.Expr) {
	switch e := ex.(type) {
	case *expr.Variable:
		r.variableExpr(e)
	case *expr.Assign:
		r.ResolveExpression(e.Value)
		r.resolveLocal(e, e.Name)
	case *expr.Binary:
		r.ResolveExpression(e.Left)
		r.ResolveExpression(e.Right)
	case *expr.Call:
		r.ResolveExpression(e.Callee)
		for _, arg := range e.Arguments {
			r.ResolveExpression(arg)
		}
	case *expr.Grouping:
		r.ResolveExpression(e.Expression)
	case *expr.Literal:
	case *expr.Logical:
		r.ResolveExpression(e.Left)
		r.ResolveExpression(e.Right)
	case *expr.Unary:
		r.ResolveExpression(e.Right)
	case *expr.Get:
		r.ResolveExpression(e.Object)
	case *expr.Set:
		r.ResolveExpression(e.Object)
		r.ResolveExpression(e.Value)
	case *expr.This:
		r.thisExpr(e)
	}
}

func (r *Resolver) thisExpr(e *expr.This) {
	if r.currentClass == ClassTypeNone {
		ReportParseError(&Error{
			Token:   e.Keyword,
			Message: "Can't use 'this' outside a class",
		})
		return
	}
	r.resolveLocal(e, e.Keyword)
}

func (r *Resolver) variableExpr(e *expr.Variable) {
	if len(r.scopes) > 0 {
		defined, declared := r.peekScope()[e.Name.Lexeme]
		if declared && !defined {
			ReportParseError(&Error{Token: e.Name, Message: "Can't read local variable in its own initializer"})
		}
	}
	r.resolveLocal(e, e.Name)
}

func (r *Resolver) resolveLocal(e expr.Expr, name *tok.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		_, declared := r.scopes[i][name.Lexeme]
		if declared {
			depthMap[e] = len(r.scopes) - i - 1
			return
		}
	}
}

func (r *Resolver) resolveFunction(s *stmt.Function, ft FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = ft
	r.beginScope()
	for _, param := range s.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStatements(s.Body)
	r.endScope()
	r.currentFunction = enclosingFunction
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

	scope := r.peekScope()
	_, defined := scope[name.Lexeme]
	if defined {
		ReportParseError(&Error{
			Token:   name,
			Message: "Already a variable with this name in this scope"})
	}

	scope[name.Lexeme] = false
}

func (r *Resolver) define(name *tok.Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.peekScope()[name.Lexeme] = true
}
