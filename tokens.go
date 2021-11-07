package rotationparser

import "fmt"

type TokenType int

const (
	Factor TokenType = iota + 1
	Comma
	Minus
	Plus
	Multiply
	Divide
	OpenParen
	CloseParen
)

func Precadence(tt TokenType) int {
	switch tt {
	case Multiply, Divide:
		return 100
	case Minus, Plus:
		return 50
	default:
		return 9999
	}
}

func RightAssociative(tt TokenType) bool {
	switch tt {
	// we don't have right associative operators just yet
	default:
		return false
	}
}

func Tokenize(code string) (tokens []Token, err error) {
	token := Token{}
	// basic FSM style approach
	for i, char := range code {
	CHECK:
		switch token.Type {
		case 0:
			switch char {
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				token.Type = Factor
				token.Code += string([]rune{char})
			case '-':
				tokens = append(tokens, Token{Minus, "(-)"})
			case '+':
				tokens = append(tokens, Token{Plus, "+"})
			case '*':
				tokens = append(tokens, Token{Multiply, "*"})
			case '/':
				tokens = append(tokens, Token{Divide, "/"})
			case '(':
				tokens = append(tokens, Token{OpenParen, "("})
			case ')':
				tokens = append(tokens, Token{CloseParen, ")"})
			case ' ', '\n', '\r', '\t':
				// ignore whitespace
			default:
				err = fmt.Errorf("tokenize error at byte offset [%d], unexpected character '%s'", i, string(char))
				return
			}
		case Factor:
			switch char {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				token.Code += string([]rune{char})
			default:
				tokens = append(tokens, token)
				token = Token{}
				// jump avoids keeping track of 'consumed' charaters
				goto CHECK
			}
		}
	}
	if token.Type > 0 {
		tokens = append(tokens, token)
	}
	return
}

func MustTokenize(code string) []Token {
	tokens, err := Tokenize(code)
	if err != nil {
		panic(fmt.Sprintf("Failed to tokenize expression (%s), got error: %s", code, err))
	}
	return tokens
}
