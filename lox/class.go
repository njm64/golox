package lox

type Class struct {
	Name string
}

func (c *Class) Arity() int {
	return 0
}

func (c *Class) Call(args []any) (any, error) {
	return NewInstance(c), nil
}

func (c *Class) String() string {
	return c.Name
}
