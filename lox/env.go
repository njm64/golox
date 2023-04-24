package lox

import (
	"golox/lox/ast"
	"golox/lox/tok"
)

var globalEnv = NewEnvironment(nil)
var currentEnv = globalEnv
var depthMap = make(map[ast.Expr]int)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name *tok.Token) (any, error) {
	val, ok := e.values[name.Lexeme]
	if ok {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, &Error{Token: name, Message: "Undefined variable '" + name.Lexeme + "'"}
}

func (e *Environment) Assign(name *tok.Token, value any) error {
	_, ok := e.values[name.Lexeme]
	if ok {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}

	return &Error{Token: name, Message: "Undefined variable '" + name.Lexeme + "'"}
}

func (e *Environment) GetAt(distance int, name *tok.Token) (any, error) {
	return e.ancestor(distance).values[name.Lexeme], nil
}

func (e *Environment) AssignAt(distance int, name *tok.Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}

func (e *Environment) ancestor(distance int) *Environment {
	a := e
	for i := 0; i < distance; i++ {
		a = a.enclosing
	}
	return a
}
