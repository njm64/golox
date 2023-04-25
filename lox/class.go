package lox

type Class struct {
	name       string
	superclass *Class
	methods    map[string]*Function
}

func NewClass(name string, superclass *Class, methods map[string]*Function) *Class {
	return &Class{
		name:       name,
		superclass: superclass,
		methods:    methods,
	}
}

func (c *Class) Arity() int {
	initializer := c.FindMethod("init")
	if initializer != nil {
		return initializer.Arity()
	}
	return 0
}

func (c *Class) Call(args []any) (any, error) {
	instance := NewInstance(c)

	initializer := c.FindMethod("init")
	if initializer != nil {
		_, err := initializer.Bind(instance).Call(args)
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
}

func (c *Class) String() string {
	return c.name
}

func (c *Class) FindMethod(name string) *Function {
	method := c.methods[name]
	if method != nil {
		return method
	}

	if c.superclass != nil {
		return c.superclass.FindMethod(name)
	}

	return nil
}
