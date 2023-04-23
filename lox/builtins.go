package lox

import "time"

func init() {
	globalEnv.Define("clock", &ClockFn{})
}

type ClockFn struct {
}

func (c ClockFn) Arity() int {
	return 0
}

func (c ClockFn) Call(arguments []any) (any, error) {
	return float64(time.Now().UnixMilli() * 1000), nil
}

func (c ClockFn) String() string {
	return "<native fn>"
}
