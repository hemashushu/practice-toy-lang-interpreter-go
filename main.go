// original from
// https://interpreterbook.com/

package main

import (
	"fmt"
	"interpreter/executor"
	"interpreter/repl"
	"os"
)

func main() {
	args := os.Args
	count := len(args)

	if count == 1 {
		// 进入 REPL 交互模式
		fmt.Println("Toy lang REPL")
		repl.Start(os.Stdin, os.Stdout)

	} else if count == 2 {
		// 解析脚本
		executor.Exec(args[1])

	} else {
		fmt.Println(`Toy language interpreter
Usage:

1. Launch REPL mode
$ ./toy

2. Execute toy lang script source code file
$ ./toy path_to_script_file`)
	}
}
