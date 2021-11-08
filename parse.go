package rotationparser

import (
	"fmt"
)

type Node struct {
	Item  Lexeme
	Left  *Node
	Right *Node
}

func ParseExpression(name, input string) *Node {
	_, tokenStream := Lex(name, input)
	items := make([]Lexeme, 0, len(input))
	for item := range tokenStream {
		if item.Type == ItemEOF {
			break
		}
		items = append(items, item)
	}
	return ParseBinaryExpression(items)
}

// ParseBinaryExpression returns the AST for a binary expression
func ParseBinaryExpression(items []Lexeme) (out *Node) {
	defer applyPrecadence(&out)
	var lhs *Node
	for i, item := range items {
		switch item.Type {
		case ItemNumber:
			lhs = &Node{Item: item}
		case ItemOperator: // infix operators
			out = &Node{
				Item:  item,
				Left:  lhs,
				Right: ParseBinaryExpression(items[i+1:]), // TODO: guard against index error
			}
			return
		default:
			panic(fmt.Sprintf("Not yet supported: %v from %v", item, items[i:]))
		}
	}
	out = lhs
	return
}

func applyPrecadence(result **Node) {
	if *result == nil || (*result).Right == nil {
		return
	}
	delta := precadence((*result).Item) - precadence((*result).Right.Item)
	if delta < 0 || delta == 0 && rightAssociative((*result).Item) {
		// do nothing
		return
	}
	// perform a left rotation
	top := *result
	lift := top.Right
	*result = lift
	top.Right = lift.Left
	lift.Left = top
	// apply recursively after the prev top nodes right side was changed
	applyPrecadence(&lift.Left)
}

func precadence(item Lexeme) int {
	switch item.Type {
	case ItemOperator:
		switch item.Value {
		case "&", "|":
			return 10
		case "+", "-":
			return 100
		case "*", "/":
			return 1000
		default:
			return 10000
		}
	default:
		return 999999
	}
}

func rightAssociative(item Lexeme) bool {
	// none yet
	return false
}
