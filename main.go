package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/davidli2010/calc/parser"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	parser := parser.NewParser()

	fmt.Println("simple calculator:")
	fmt.Println("    - type 'exit' to exit")
	fmt.Println("    - support integer and floating number")
	fmt.Println("    - support '+', '-', '*', '/' and parentheses")
	fmt.Println("    - the result is double float value")
	fmt.Println("    - example: '1.5 + 2 * (-1 - 2.2) / 10', the result is 0.860000")

	for {
		fmt.Println()

		line, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Errorf("error occured when read line: %v", err))
		}

		line = strings.TrimSpace(line)
		if line == "exit" {
			break
		} else if len(line) == 0 {
			continue
		}

		value, err := parser.Parse(line)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("= %f\n", value)
		}
	}
}
