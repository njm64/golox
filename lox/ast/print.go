package ast

import (
	"fmt"
	"strings"
)

func PrintExpr(expr Expr) string {
	switch e := expr.(type) {
	case *Binary:
		return parenthesize(e.Operator.Lexeme, e.Left, e.Right)
	case *Grouping:
		return parenthesize("group", e.Expression)
	case *Literal:
		if e.Value == nil {
			return "nil"
		} else {
			return fmt.Sprintf("%v", e.Value)
		}
	case *Unary:
		return parenthesize(e.Operator.Lexeme, e.Right)
	default:
		return ""
	}
}

func parenthesize(name string, es ...Expr) string {
	sb := &strings.Builder{}
	sb.WriteRune('(')
	sb.WriteString(name)
	for _, e := range es {
		sb.WriteRune(' ')
		sb.WriteString(PrintExpr(e))
	}
	sb.WriteRune(')')
	return sb.String()
}
