package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	parser "github.com/rur/rotationparser"
)

const help = `Simple a REPL to show the precedence-adjusted AST for simple expressions.

For example

     eXpr: 1 + 2
     => "+"
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
		fmt.Printf("\neXpr: ")
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
				fmt.Println("Lex error:", item.Value)
				continue Read
			}
			items = append(items, item)
		}

		func() {
			// in-leu of proper error handling, this will have to do for now
			defer func() {
				err := recover()
				if err != nil {
					fmt.Println("Parser error: ", err)
				}
			}()
			node := parser.ParseBinaryExpression(items)
			fmt.Println(parser.SprintNodeTree(node))
		}()

	}
}
