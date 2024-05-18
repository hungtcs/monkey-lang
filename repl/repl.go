package repl

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/chzyer/readline"
	"github.com/hungtcs/monkey-lang/monkey"
	"github.com/hungtcs/monkey-lang/syntax"
)

const PROMPT = "\033[90m>>\033[0m "
const CONTINUE_PROMPT = "\033[90m..\033[0m "

var interrupted = make(chan os.Signal, 1)

func Start() (err error) {
	signal.Notify(interrupted, os.Interrupt)
	defer signal.Stop(interrupted)

	rl, err := readline.New(PROMPT)
	if err != nil {
		return err
	}
	defer rl.Close()

	env := monkey.NewEnv(nil)

	for {
		if err := repl(rl, env); err != nil {
			if err == readline.ErrInterrupt {
				fmt.Println("(To exit, press Ctrl+D)")
				continue
			}
			break
		}
	}
	return nil
}

func repl(rl *readline.Instance, env *monkey.Env) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		select {
		case <-interrupted:
			cancel()
		case <-ctx.Done():
		}
	}()

	var eof = false
	rl.SetPrompt(PROMPT)
	readline := func() (string, error) {
		line, err := rl.Readline()
		rl.SetPrompt(CONTINUE_PROMPT)
		if err != nil {
			if err == io.EOF {
				eof = true
			}
			return "", err
		}
		return line + "\n", nil
	}

	line, err := readline()
	if err != nil {
		return err
	}

	parser := syntax.NewParser(line)
	program, err := parser.Parse()
	if err != nil {
		if eof {
			return io.EOF
		}
		printError(err)
		return nil
	}
	val, err := monkey.Eval(program, env)
	if err != nil {
		printError(err)
		return nil
	}
	if val != monkey.Null {
		fmt.Fprintln(os.Stdout, val.String())
	}

	return nil
}

func printError(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
}
