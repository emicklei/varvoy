package main

import (
	"fmt"
	"go/parser"
	"go/token"
)

func main() {
	src1 := `
	package main

	import "fmt"

	func main() {
		fmt.Println("Hello, World!")
	}
	`
	set := token.NewFileSet()
	expr, err := parser.ParseFile(set, "main.go", src1, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	fmt.Println(expr)
	{
		src2 := `
		package main
	
		import "fmt"
	
		func main() {
			fmt.Println("Hello, Varvoy!")
		}
		`
		set := token.NewFileSet()
		expr2, err := parser.ParseFile(set, "main.go", src2, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		fmt.Println(expr2)
	}
	/**
	1: strategies for finding the changed functions:

	compare functions by comparing the start and end pos of each function
	-> need to keep this data in memory
	-> also some hashcode, content may have same size but is different

	2: dont care about the changes, just reload the whole file
	-> see how REPL works

	**/
}
