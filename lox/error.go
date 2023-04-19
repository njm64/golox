package lox

import (
	"fmt"
	tok2 "golox/lox/tok"
)

type Error struct {
	Token   *tok2.Token
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

var HadError = false
var HadRuntimeError = false

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	HadError = true
}

func ReportScanError(line int, message string) {
	report(line, "", message)
}

func ReportParseError(err *Error) {
	if err.Token.Type == tok2.EOF {
		report(err.Token.Line, " at end", err.Message)
	} else {
		report(err.Token.Line, " at '"+err.Token.Lexeme+"'", err.Message)
	}
}

func ReportRuntimeError(err *Error) {
	fmt.Printf("%s\n[line %d]\n", err.Message, err.Token.Line)
	HadRuntimeError = true
}
