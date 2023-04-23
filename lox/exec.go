package lox

import (
	"errors"
	"fmt"
	"golox/lox/ast"
)

func Exec(st ast.Stmt) error {
	switch s := st.(type) {
	case *ast.Print:
		val, err := Eval(s.Expression)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", val)
		return nil
	case *ast.Expression:
		_, err := Eval(s.Expression)
		return err
	case *ast.If:
		return execIf(s)
	case *ast.While:
		return execWhile(s)
	case *ast.Var:
		return execVar(s)
	case *ast.Block:
		return execBlock(s.Statements, NewEnvironment(currentEnv))
	case *ast.Function:
		return execFunction(s)
	default:
		return errors.New("unhandled statement")
	}
}

func execIf(s *ast.If) error {
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

func execWhile(s *ast.While) error {
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

func execVar(s *ast.Var) error {
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

func execBlock(statements []ast.Stmt, env *Environment) error {
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

func execFunction(f *ast.Function) error {
	currentEnv.Define(f.Name.Lexeme, &Function{Declaration: f})
	return nil
}
