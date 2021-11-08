package main

import (
	"bufio"
	"fmt"
	"os"

	parser "github.com/rur/rotation-parser"
)

func main() {
Read:
	for {
		var line string
		fmt.Printf("\neXpr-> ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if line == "exit\n" {
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
