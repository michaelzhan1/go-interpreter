package main

import (
	"fmt"
	"os"

	"github.com/michaelzhan1/go-interpreter/repl"
)

func main() {
	fmt.Println("Welcome to the Monkey Programming Language!")
	repl.Start(os.Stdin, os.Stdout)
}
