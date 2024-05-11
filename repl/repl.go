package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/hungtcs/monkey-lang/syntax"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, PROMPT)
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		parser := syntax.NewParser(line)
		program := parser.Parse()
		if len(parser.Errors()) > 0 {
			printParseErrors(out, parser.Errors())
			continue
		}
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
	return scanner.Err()
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
