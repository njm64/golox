package lox

import "golox/lox/tok"

type Instance struct {
	class  *Class
	fields map[string]any
}

func NewInstance(class *Class) *Instance {
	return &Instance{
		class:  class,
		fields: make(map[string]any),
	}
}

func (i *Instance) String() string {
	return i.class.name + " instance"
}

func (i *Instance) Get(name *tok.Token) (any, error) {
	value, ok := i.fields[name.Lexeme]
	if ok {
		return value, nil
	}

	method := i.class.FindMethod(name.Lexeme)
	if method != nil {
		return method.Bind(i), nil
	}

	return nil, &Error{
		Token:   name,
		Message: "Undefined property '" + name.Lexeme + "'",
	}
}

func (i *Instance) Set(name *tok.Token, value any) {
	i.fields[name.Lexeme] = value
}
