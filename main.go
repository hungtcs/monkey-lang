package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/hungtcs/monkey-lang/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")

	if err := repl.Start(os.Stdin, os.Stdout); err != nil {
		panic(err)
	}
}
