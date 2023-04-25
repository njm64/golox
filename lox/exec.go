package lox

import (
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
	case *stmt.Function:
		return execFunction(s)
	case *stmt.Return:
		return execReturn(s)
	case *stmt.Class:
		return execClass(s)
	default:
		return fmt.Errorf("unhandled statement %v", st)
	}
}

func execIf(s *stmt.If) error {
	condition, err := Eval(s.Condition)
	if err != nil {
		return err
	}
	if isTruthy(condition) {
		return Exec(s.ThenBranch)
	} else if s.ElseBranch != nil {
		return Exec(s.ElseBranch)
	}
	return nil
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
	currentEnv.Define(s.Name.Lexeme, value)
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

func execFunction(s *stmt.Function) error {
	currentEnv.Define(s.Name.Lexeme, NewFunction(s, currentEnv, false))
	return nil
}

func execReturn(s *stmt.Return) error {
	var result any
	var err error
	if s.Value != nil {
		result, err = Eval(s.Value)
		if err != nil {
			return err
		}
	}
	return &Return{Value: result}
}

func execClass(s *stmt.Class) error {
	currentEnv.Define(s.Name.Lexeme, nil)
	methods := make(map[string]*Function)
	for _, m := range s.Methods {
		methods[m.Name.Lexeme] = NewFunction(m, currentEnv,
			m.Name.Lexeme == "init")
	}
	class := NewClass(s.Name.Lexeme, methods)
	return currentEnv.Assign(s.Name, class)
}
