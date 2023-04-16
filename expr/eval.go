package expr

import (
	"golox/tok"
)

func Eval(expr Expr) any {
	switch e := expr.(type) {
	case *Binary:
		return evalBinary(e)
	case *Grouping:
		return Eval(e.Expression)
	case *Literal:
		return e.Value
	case *Unary:
		return evalUnary(e)
	default:
		return nil
	}
}

func evalUnary(e *Unary) any {
	right := Eval(e.Right)
	switch e.Operator.Type {
	case tok.Minus:
		checkNumberOperand(e.Operator, right)
		return -right.(float64)
	case tok.Bang:
		return !isTruthy(right)
	default:
		return nil
	}
}

func evalBinary(e *Binary) any {
	left := Eval(e.Left)
	right := Eval(e.Right)
	switch e.Operator.Type {
	case tok.Greater:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) > right.(float64)
	case tok.GreaterEqual:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) >= right.(float64)
	case tok.Less:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) < right.(float64)
	case tok.LessEqual:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) <= right.(float64)
	case tok.EqualEqual:
		return isEqual(left, right)
	case tok.BangEqual:
		checkNumberOperands(e.Operator, left, right)
		return !isEqual(left, right)
	case tok.Minus:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) - right.(float64)
	case tok.Slash:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) / right.(float64)
	case tok.Star:
		checkNumberOperands(e.Operator, left, right)
		return left.(float64) + right.(float64)
	case tok.Plus:
		if isNumber(left) && isNumber(right) {
			return left.(float64) + right.(float64)
		} else if isString(left) && isString(right) {
			return left.(string) + right.(string)
		} else {
			panic(&RuntimeError{
				Token:   e.Operator,
				Message: "operands should be numbers or strings",
			})
		}
	default:
		return nil
	}
}

func checkNumberOperand(tok *tok.Token, operand any) {
	if !isNumber(operand) {
		panic(&RuntimeError{
			Token:   tok,
			Message: "operand must be a number",
		})
	}
}

func checkNumberOperands(tok *tok.Token, left any, right any) {
	if !isNumber(left) || !isNumber(right) {
		panic(&RuntimeError{
			Token:   tok,
			Message: "operands must be numbers",
		})
	}
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
