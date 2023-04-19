package lox

import (
	"errors"
	"fmt"
	"golox/lox/stmt"
)

func Exec(st stmt.Stmt) error {
	switch s := st.(type) {
	case *stmt.Print:
		val, err := Eval(s.Expression)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", val)
		return nil
	case *stmt.Expression:
		_, err := Eval(s.Expression)
		return err
	default:
		return errors.New("unhandled statement")
	}
}
