package lox

import "golox/lox/tok"

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		values: make(map[string]any),
	}
}

func (e *Environment) Define(name *tok.Token, value any) {
	e.values[name.Lexeme] = value
}

func (e *Environment) Get(name *tok.Token) (any, error) {
	val, ok := e.values[name.Lexeme]
	if ok {
		return val, nil
	} else {
		return nil, &Error{Token: name, Message: "Undefined variable '" + name.Lexeme + "'"}
	}
}

func (e *Environment) Assign(name *tok.Token, value any) error {
	_, ok := e.values[name.Lexeme]
	if !ok {
		return &Error{Token: name, Message: "Undefined variable '" + name.Lexeme + "'"}
	}
	e.values[name.Lexeme] = value
	return nil
}
