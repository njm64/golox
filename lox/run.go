package lox

import (
	"bufio"
	"fmt"
	"golox/lox/stmt"
	"os"
)

func interpret(statements []stmt.Stmt) {
	for _, s := range statements {
		err := Exec(s)
		if err != nil {
			runtimeError, ok := err.(*Error)
			if ok {
				ReportRuntimeError(runtimeError)
			} else {
				fmt.Printf("Error: %s\n", err)
			}
			break
		}
	}
}

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := NewParser(tokens)
	e, err := parser.Parse()
	if err != nil || HadError {
		// Do we need the global HadError flag as well as looking at the error
		// return? I think so, since the parser will do recovery/synchronisation and
		// try to log as many errors as possible. Once this is implemented, perhaps
		// the Parse method should return an error at all?
		return
	}

	interpret(e)
}

func RunFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	run(string(bytes))
	if HadError {
		os.Exit(65)
	}
	if HadRuntimeError {
		os.Exit(70)
	}
	return nil
}

func RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("> ")
		if !scanner.Scan() {
			break
		}
		run(scanner.Text())
		HadError = false
	}
}
