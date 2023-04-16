package main

import (
	"bufio"
	"fmt"
	"golox/expr"
	"golox/lox"
	"os"
)

func interpret(e expr.Expr) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(*expr.RuntimeError)
			lox.RuntimeError(err)
		}
	}()
	value := expr.Eval(e)
	fmt.Printf("%v\n", value)
}

func run(source string) {
	scanner := lox.NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := lox.NewParser(tokens)
	e, err := parser.Parse()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	if lox.HadError {
		return
	}

	interpret(e)
}

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	run(string(bytes))
	if lox.HadError {
		os.Exit(65)
	}
	if lox.HadRuntimeError {
		os.Exit(70)
	}
	return nil
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("> ")
		if !scanner.Scan() {
			break
		}
		run(scanner.Text())
		lox.HadError = false
	}
}

func main() {
	if len(os.Args) > 2 {
		fmt.Printf("Usage: golox [script]\n")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		if err := runFile(os.Args[1]); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	} else {
		runPrompt()
	}
}
