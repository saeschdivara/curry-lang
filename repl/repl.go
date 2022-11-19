package repl

import (
	"bufio"
	"curryLang/lexer"
	"curryLang/parser"
	"fmt"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)

		p := parser.New(l)
		program := p.ParseProgram()

		for _, stmt := range program.Statements {
			fmt.Println(stmt.String())
		}
	}
}
