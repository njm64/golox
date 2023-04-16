package lox

import (
	"fmt"
	"golox/expr"
	"golox/tok"
)

var HadError = false
var HadRuntimeError = false

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	HadError = true
}

func ParseError(t *tok.Token, message string) {
	if t.Type == tok.EOF {
		report(t.Line, " at end", message)
	} else {
		report(t.Line, " at '"+t.Lexeme+"'", message)
	}
}

func RuntimeError(err *expr.RuntimeError) {
	fmt.Printf("%s\n[line %d]\n", err.Message, err.Token.Line)
	HadRuntimeError = true
}
