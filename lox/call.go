package lox

import (
	"golox/lox/stmt"
)

type Callable interface {
	Arity() int
	Call(arguments []any) (any, error)
}

type Return struct {
	Value any
}

func (r *Return) Error() string {
	// Not really an error, but we use the error mechanism to unwind the
	// stack for returns, the same way the Java interpreter uses exceptions.
	return ""
}

type Function struct {
	declaration   *stmt.Function
	closure       *Environment
	isInitializer bool
}

func NewFunction(declaration *stmt.Function, closure *Environment, isInitializer bool) *Function {
	return &Function{
		declaration:   declaration,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

func (f *Function) Arity() int {
	return len(f.declaration.Params)
}

func (f *Function) Call(arguments []any) (any, error) {
	e := NewEnvironment(f.closure)
	for i, param := range f.declaration.Params {
		e.Define(param.Lexeme, arguments[i])
	}
	err := execBlock(f.declaration.Body, e)
	if err != nil {
		ret, ok := err.(*Return)
		if ok {
			if f.isInitializer {
				return f.closure.GetAt(0, "this")
			}
			return ret.Value, nil
		}
		return nil, err
	}

	if f.isInitializer {
		return f.closure.GetAt(0, "this")
	}

	return nil, nil
}

func (f *Function) Bind(instance *Instance) *Function {
	e := NewEnvironment(f.closure)
	e.Define("this", instance)
	return NewFunction(f.declaration, e, f.isInitializer)
}

func (f *Function) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}
