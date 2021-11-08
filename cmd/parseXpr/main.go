package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	parser "github.com/rur/rotation-parser"
)

const help = `Simple a REPL to show the precadence-adjusted AST for simple expressions.

For example

     eXpr-> 1 + 2
     ─> "+"
         ├── "2"
         └── "1"

Commands
	?      Print error message
	exit   Exit the REPLs
`

func main() {
	fmt.Println(help)
Read:
	for {
		var line string
		fmt.Printf("\neXpr-> ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		switch strings.TrimSpace(line) {
		case "?":
			fmt.Println(help)
			continue Read

		case "exit":
			fmt.Println(">>> graceful exit")
			return

		}

		_, tokenStream := parser.Lex("parseXpr", line)
		items := make([]parser.Lexeme, 0, len(line))
		for item := range tokenStream {
			if item.Type == parser.ItemEOF {
				break
			} else if item.Type == parser.ItemError {
				fmt.Println("Error:", item.Value)
				continue Read
			}
			items = append(items, item)
		}

		node := parser.ParseBinaryExpression(items)
		fmt.Println(parser.SprintNodeTree(node))
	}
}
