package main

import (
	"curryLang/evaluator"
	"curryLang/lexer"
	"curryLang/parser"
	"curryLang/repl"
	"fmt"
	"io"
	"os"
	"os/user"
)

func main() {

	args := os.Args[1:]
	if len(args) > 0 {
		fileToExecute, err := os.Open(args[0])
		if err != nil {
			panic("Failed to open file")
		}

		data, err := io.ReadAll(fileToExecute)
		if err != nil {
			panic("Failed to read file")
		}

		l := lexer.New(string(data))
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Println("Error: ", err)
			}

			os.Exit(1)
		}

		engine := evaluator.ExecutionEngine{}
		evalResult := engine.Eval(program)

		if evalResult != nil {
			fmt.Println(evalResult.Inspect())
		}
	} else {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the Curry programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}
