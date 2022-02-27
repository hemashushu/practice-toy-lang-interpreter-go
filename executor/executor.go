package executor

import (
	"fmt"
	"interpreter/evaluator"
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"log"
	"os"
)

func Exec(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	text := string(content)

	l := lexer.New(text)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		return
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}

func printParserErrors(errors []string) {
	fmt.Println("parser errors:")
	for _, msg := range errors {
		fmt.Println("\t" + msg)
	}
}
