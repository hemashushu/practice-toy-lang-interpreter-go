// original from https://interpreterbook.com/

package main

import (
	"fmt"
	"interpreter/repl"
	"os"
)

func main() {
	fmt.Println("Toy lang REPL")
	repl.Start(os.Stdin, os.Stdout)
}
