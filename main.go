package main

import (
	"bufio"
	"fmt"
	"os"
)

func run(s string) error {
	return nil
}

func runFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return run(string(bytes))
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("> ")
		if !scanner.Scan() {
			break
		}
		_ = run(scanner.Text())
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
