package main

import (
	"fmt"
	"golox/lox"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Printf("Usage: golox [script]\n")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		if err := lox.RunFile(os.Args[1]); err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
	} else {
		lox.RunPrompt()
	}
}
