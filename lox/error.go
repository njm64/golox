package lox

import (
	"fmt"
	"golox/tok"
)

var HadError = false

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	HadError = true
}

func ErrorWithToken(t *tok.Token, message string) {
	if t.Type == tok.EOF {
		report(t.Line, " at end", message)
	} else {
		report(t.Line, " at '"+t.Lexeme+"'", message)
	}
}
