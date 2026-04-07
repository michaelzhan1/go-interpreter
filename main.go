package main

import (
	"fmt"
	"os"

	"github.com/michaelzhan1/go-interpreter/evaluator"
	"github.com/michaelzhan1/go-interpreter/lexer"
	"github.com/michaelzhan1/go-interpreter/object"
	"github.com/michaelzhan1/go-interpreter/parser"
	"github.com/michaelzhan1/go-interpreter/repl"
)

func main() {
	if len(os.Args) > 1 {
		path := os.Args[1]
		dat, err := os.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("unable to read file at %s", path))
		}
		l := lexer.New(string(dat))
		p := parser.New(l)
		program := p.ParseProgram()

		env := object.NewEnvironment()
		evaluator.Eval(program, env)
	} else {
		fmt.Println("Welcome to the Monkey Programming Language!")
		repl.Start(os.Stdin, os.Stdout)
	}

}
