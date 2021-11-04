package rotationparser

import "fmt"

type Node struct {
	Token Token

	Left  *Node
	Right *Node
}

type Token struct {
	Type int
	Code string
}

// ParseExpression([]Token{"12", "*", "5", "+ "7"})
// => {+ {* {12} {5}} {7}}
func ParseExpression(tokens []Token) (out *Node) {
	defer applyPrecadence(&out)
	var lhs *Node
	for i, token := range tokens {
		switch token.Type {
		case Factor:
			lhs = &Node{Token: token}
		case Minus, Plus, Multiply, Divide: // infix operators
			out = &Node{
				Token: token,
				Left:  lhs,
				Right: ParseExpression(tokens[i+1:]), // TODO: guard against index error
			}
			return
		default:
			panic(fmt.Sprintf("Not yet supported: %v from %v", token, tokens[i:]))
		}
	}
	out = lhs
	return
}

func applyPrecadence(result **Node) {
	if *result == nil || (*result).Right == nil {
		return
	}
	delta := Precadence((*result).Token.Type) - Precadence((*result).Right.Token.Type)
	if delta > 0 || delta == 0 && RightAssociative((*result).Token.Type) {
		// do nothing
		return
	}
	// perform a left rotation
	top := *result
	lift := top.Right
	*result = lift
	top.Right = lift.Left
	lift.Left = top
}
