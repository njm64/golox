package lox

import (
	"errors"
	"fmt"
	"golox/lox/ast"
	"golox/lox/tok"
)

func Eval(ex ast.Expr) (any, error) {
	switch e := ex.(type) {
	case *ast.Binary:
		return evalBinary(e)
	case *ast.Grouping:
		return Eval(e.Expression)
	case *ast.Literal:
		return e.Value, nil
	case *ast.Logical:
		return evalLogical(e)
	case *ast.Unary:
		return evalUnary(e)
	case *ast.Variable:
		return currentEnv.Get(e.Name)
	case *ast.Assign:
		return evalAssign(e)
	case *ast.Call:
		return evalCall(e)
	default:
		return nil, errors.New("unhandled expression type")
	}
}

func evalUnary(e *ast.Unary) (any, error) {
	right, err := Eval(e.Right)
	if err != nil {
		return nil, err
	}
	switch e.Operator.Type {
	case tok.Minus:
		err = checkNumberOperand(e.Operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	case tok.Bang:
		return !isTruthy(right), nil
	default:
		return nil, &Error{Token: e.Operator, Message: "unhandled unary expression"}
	}
}

func evalBinary(e *ast.Binary) (any, error) {
	left, err := Eval(e.Left)
	if err != nil {
		return nil, err
	}
	right, err := Eval(e.Right)
	if err != nil {
		return nil, err
	}
	switch e.Operator.Type {
	case tok.Greater:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case tok.GreaterEqual:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case tok.Less:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case tok.LessEqual:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case tok.EqualEqual:
		return isEqual(left, right), nil
	case tok.BangEqual:
		return !isEqual(left, right), nil
	case tok.Minus:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case tok.Slash:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case tok.Percent:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return float64(int(left.(float64)) % int(right.(float64))), nil
	case tok.Star:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) + right.(float64), nil
	case tok.Plus:
		if isNumber(left) && isNumber(right) {
			return left.(float64) + right.(float64), nil
		} else if isString(left) && isString(right) {
			return left.(string) + right.(string), nil
		} else {
			return nil, &Error{
				Token:   e.Operator,
				Message: "operands should be numbers or strings"}
		}
	default:
		return nil, &Error{Token: e.Operator, Message: "unexpected token"}
	}
}

func evalAssign(e *ast.Assign) (any, error) {
	value, err := Eval(e.Value)
	if err != nil {
		return nil, err
	}
	err = currentEnv.Assign(e.Name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func evalLogical(e *ast.Logical) (any, error) {
	left, err := Eval(e.Left)
	if err != nil {
		return nil, err
	}
	if e.Operator.Type == tok.Or {
		if isTruthy(left) {
			return left, nil
		}
	} else {
		if !isTruthy(left) {
			return left, nil
		}
	}
	return Eval(e.Right)
}

func evalCall(e *ast.Call) (any, error) {
	callee, err := Eval(e.Callee)
	if err != nil {
		return nil, err
	}

	var arguments []any
	for _, a := range e.Arguments {
		arg, err := Eval(a)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, arg)
	}

	f, ok := callee.(Callable)
	if !ok {
		return nil, &Error{
			Token:   e.Paren,
			Message: "Can only call functions and classes",
		}
	}

	if len(arguments) != f.Arity() {
		return nil, &Error{
			Token: e.Paren,
			Message: fmt.Sprintf("Expected %d arguments but got %d",
				f.Arity(), len(arguments)),
		}
	}

	return f.Call(arguments)
}

func checkNumberOperand(tok *tok.Token, operand any) error {
	if !isNumber(operand) {
		return &Error{Token: tok, Message: "operand must be a number"}
	}
	return nil
}

func checkNumberOperands(tok *tok.Token, left any, right any) error {
	if !isNumber(left) || !isNumber(right) {
		return &Error{
			Token:   tok,
			Message: "operands must be numbers",
		}
	}
	return nil
}

func isNumber(value any) bool {
	_, ok := value.(float64)
	return ok
}

func isString(value any) bool {
	_, ok := value.(string)
	return ok
}

func isTruthy(value any) bool {
	return value != nil && value != false
}

func isEqual(a any, b any) bool {
	// TODO: Check this
	return a == b
}
