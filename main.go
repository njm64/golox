package main

import (
	"bufio"
	"fmt"
	"golox/expr"
	"golox/lox"
	"golox/tok"
	"os"
)

func run(source string) {
	scanner := lox.NewScanner(source)
	tokens := scanner.ScanTokens()
	for _, t := range tokens {
		fmt.Printf("%s\n", t)
	}
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

	e := &expr.Binary{
		Left: &expr.Unary{
			Operator: tok.Token{Type: tok.Minus, Lexeme: "-"},
			Right:    &expr.Literal{Value: 123},
		},
		Operator: tok.Token{Type: tok.Star, Lexeme: "*"},
		Right: &expr.Grouping{
			Expression: &expr.Literal{Value: 45.67},
		},
	}

	fmt.Printf("%s\n", expr.Print(e))

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
