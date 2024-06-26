package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/hungtcs/monkey-lang/monkey"
	"github.com/hungtcs/monkey-lang/repl"
	"github.com/hungtcs/monkey-lang/syntax"
)

func startRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")

	if err := repl.Start(); err != nil {
		panic(err)
	}
}

func main() {
	var args = os.Args[1:]
	// start repl
	if len(args) < 1 {
		startRepl()
		return
	}
	// eval a file
	if len(args) == 1 {
		var filename = args[0]
		data, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		parser := syntax.NewParser(string(data))
		program, err := parser.Parse()
		if err != nil {
			fmt.Println(err)
			return
		}

		value, err := monkey.Eval(program, monkey.NewEnv(nil))
		if err != nil {
			panic(err)
		}
		fmt.Println(value)
		return
	}
	fmt.Println("usage: monkey [file]")
}
