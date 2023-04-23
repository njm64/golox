package lox

import "golox/lox/ast"

type Callable interface {
	Arity() int
	Call(arguments []any) (any, error)
}

type Function struct {
	Declaration *ast.Function
}

func (f *Function) Arity() int {
	return len(f.Declaration.Params)
}

func (f *Function) Call(arguments []any) (any, error) {
	e := NewEnvironment(globalEnv)
	for i, param := range f.Declaration.Params {
		e.Define(param.Lexeme, arguments[i])
	}
	err := execBlock(f.Declaration.Body, e)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (f *Function) String() string {
	return "<fn " + f.Declaration.Name.Lexeme + ">"
}
