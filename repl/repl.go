package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/hungtcs/monkey-lang/lexer"
	"github.com/hungtcs/monkey-lang/token"
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
		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
	return scanner.Err()
}
