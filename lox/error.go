package lox

import "fmt"

var HadError = false

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error %s: %s", line, where, message)
	HadError = true
}
