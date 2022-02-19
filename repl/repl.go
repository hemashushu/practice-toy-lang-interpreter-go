// original from https://interpreterbook.com/

package repl

import (
	"bufio"
	"fmt"
	"interpreter/lexer"
	"interpreter/token"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return // exit for
		}

		line := scanner.Text()
		lx := lexer.New(line)

		for tk := lx.NextToken(); tk.Type != token.EOF; tk = lx.NextToken() {
			fmt.Printf("%+v\n", tk) // %+v 比 %v 多显示结构体的字段名称
		}
	}
}
