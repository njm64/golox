package lox

import (
	"errors"
	"golox/lox/expr"
	tok2 "golox/lox/tok"
)

func Eval(ex expr.Expr) (any, error) {
	switch e := ex.(type) {
	case *expr.Binary:
		return evalBinary(e)
	case *expr.Grouping:
		return Eval(e.Expression)
	case *expr.Literal:
		return e.Value, nil
	case *expr.Unary:
		return evalUnary(e)
	default:
		return nil, errors.New("unhandled expression type")
	}
}

func evalUnary(e *expr.Unary) (any, error) {
	right, err := Eval(e.Right)
	if err != nil {
		return nil, err
	}
	switch e.Operator.Type {
	case tok2.Minus:
		err = checkNumberOperand(e.Operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	case tok2.Bang:
		return !isTruthy(right), nil
	default:
		return nil, &Error{Token: e.Operator, Message: "unhandled unary expression"}
	}
}

func evalBinary(e *expr.Binary) (any, error) {
	left, err := Eval(e.Left)
	if err != nil {
		return nil, err
	}
	right, err := Eval(e.Right)
	if err != nil {
		return nil, err
	}
	switch e.Operator.Type {
	case tok2.Greater:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case tok2.GreaterEqual:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case tok2.Less:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case tok2.LessEqual:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case tok2.EqualEqual:
		return isEqual(left, right), nil
	case tok2.BangEqual:
		return !isEqual(left, right), nil
	case tok2.Minus:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case tok2.Slash:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case tok2.Star:
		err = checkNumberOperands(e.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) + right.(float64), nil
	case tok2.Plus:
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

func checkNumberOperand(tok *tok2.Token, operand any) error {
	if !isNumber(operand) {
		return &Error{Token: tok, Message: "operand must be a number"}
	}
	return nil
}

func checkNumberOperands(tok *tok2.Token, left any, right any) error {
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
