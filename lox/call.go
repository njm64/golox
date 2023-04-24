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
	Declaration *stmt.Function
	Closure     *Environment
}

func (f *Function) Arity() int {
	return len(f.Declaration.Params)
}

func (f *Function) Call(arguments []any) (any, error) {
	e := NewEnvironment(f.Closure)
	for i, param := range f.Declaration.Params {
		e.Define(param.Lexeme, arguments[i])
	}
	err := execBlock(f.Declaration.Body, e)
	if err != nil {
		ret, ok := err.(*Return)
		if ok {
			return ret.Value, nil
		}
		return nil, err
	}
	return nil, nil
}

func (f *Function) String() string {
	return "<fn " + f.Declaration.Name.Lexeme + ">"
}
