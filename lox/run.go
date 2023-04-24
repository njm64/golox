package lox

import (
	"bufio"
	"fmt"
	"golox/lox/ast"
	"os"
)

func interpret(statements []ast.Stmt) {
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
	statements := parser.Parse()
	if HadError {
		return
	}

	resolver := NewResolver()
	resolver.ResolveStatements(statements)
	if HadError {
		return
	}

	interpret(statements)
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
