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
	case *stmt.If:
		return execIf(s)
	case *stmt.While:
		return execWhile(s)
	case *stmt.Var:
		return execVar(s)
	case *stmt.Block:
		return execBlock(s.Statements, NewEnvironment(currentEnv))
	default:
		return errors.New("unhandled statement")
	}
}

func execIf(s *stmt.If) error {
	condition, err := Eval(s.Condition)
	if err != nil {
		return err
	}
	if isTruthy(condition) {
		return Exec(s.ThenBranch)
	} else {
		return Exec(s.ElseBranch)
	}
}

func execWhile(s *stmt.While) error {
	for {
		condition, err := Eval(s.Condition)
		if err != nil {
			return err
		}
		if !isTruthy(condition) {
			return nil
		}
		if err = Exec(s.Body); err != nil {
			return err
		}
	}
}

func execVar(s *stmt.Var) error {
	var value any
	var err error
	if s.Initializer != nil {
		value, err = Eval(s.Initializer)
		if err != nil {
			return err
		}
	}
	currentEnv.Define(s.Name, value)
	return nil
}

func execBlock(statements []stmt.Stmt, env *Environment) error {
	previousEnv := currentEnv
	currentEnv = env
	for _, s := range statements {
		if err := Exec(s); err != nil {
			currentEnv = previousEnv
			return err
		}
	}
	currentEnv = previousEnv
	return nil
}
