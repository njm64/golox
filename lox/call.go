package lox

type Callable interface {
	Arity() int
	Call(arguments []any) (any, error)
}
