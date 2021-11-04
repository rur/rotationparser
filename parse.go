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
func ParseExpression(tokens []Token) *Node {
	var lhs *Node
	for i, token := range tokens {
		switch token.Type {
		case Factor:
			lhs = &Node{Token: token}
		case Minus, Plus, Multiply, Divide: // infix operators
			return &Node{
				Token: token,
				Left:  lhs,
				Right: ParseExpression(tokens[i+1:]), // TODO: guard against index error
			}
		default:
			panic(fmt.Sprintf("Not yet supported: %v from %v", token, tokens[i:]))
		}
	}
	panic("SHOULDN'T GET HERE!")
}
