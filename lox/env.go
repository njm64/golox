package lox

import "golox/lox/tok"

var globalEnv = NewEnvironment(nil)
var currentEnv = globalEnv

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
