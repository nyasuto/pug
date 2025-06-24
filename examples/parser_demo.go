package main

import (
	"fmt"
	"github.com/nyasuto/pug/phase1"
)

func main() {
	input := `
let x = 5;
let y = 10;
let add = fn(a, b) { a + b; };
let result = add(x, y);
if (result > 10) {
    return true;
} else {
    return false;
}
`

	l := phase1.New(input)
	p := phase1.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range p.Errors() {
			fmt.Printf("  - %s\n", err)
		}
		return
	}

	fmt.Println("Parsed successfully!")
	fmt.Printf("Program has %d statements\n", len(program.Statements))
	fmt.Println("\nAST representation:")
	fmt.Println(program.String())
}
