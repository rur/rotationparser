package rotationparser

import (
	"fmt"
)

type Node struct {
	Item  Lexeme
	Left  *Node
	Right *Node
}

// ParseBinaryExpression returns the AST for a binary expression
func ParseBinaryExpression(items []Lexeme) (out *Node) {
	// correct precedence as the call stack unwinds
	defer applyPrecedence(&out)
	// Start with right recursive only, EXPR = LHS Op EXPR(rest...)
	var lhs *Node
	for i, item := range items {
		switch item.Type {
		case ItemNumber:
			lhs = &Node{Item: item}
		case ItemOperator: // infix operators
			if lhs == nil {
				// TODO: add graceful error handling
				panic(fmt.Sprintf("Infix binary expression missing left operand %v", items))
			}
			out = &Node{
				Item:  item,
				Left:  lhs,
				Right: ParseBinaryExpression(items[i+1:]), // TODO: guard against index error
			}
			return
		default:
			// TODO: add graceful error handling
			panic(fmt.Sprintf("Not yet supported: %q from %v", item, items[i:]))
		}
	}
	out = lhs
	return
}

// applyPrecedence will correct left hand side precedence rules
// for a right-weighted AST subtree
func applyPrecedence(result **Node) {
	if *result == nil || (*result).Right == nil {
		return
	}
	delta := precedence((*result).Item) - precedence((*result).Right.Item)
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
	applyPrecedence(&lift.Left)
}

// precedence tabulates numbers used to compare the evaluation priority of a term
func precedence(item Lexeme) int {
	switch item.Type {
	case ItemOperator:
		switch item.Value {
		case "&", "|", "^":
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

// rightAssociative tabulates which operators are right associative (evaluating right to left)
// it will return false if the item is left associating (evaluating left to right)
func rightAssociative(item Lexeme) bool {
	// none yet
	return false
}
