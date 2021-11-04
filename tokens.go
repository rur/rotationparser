package rotationparser

const (
	Factor = iota
	Comma
	Minus
	Plus
	Multiply
	Divide
)

func Precadence(tokenType int) int {
	switch tokenType {
	case Multiply, Divide:
		return 100
	case Minus, Plus:
		return 50
	default:
		return -1
	}
}

func RightAssociative(tokenType int) bool {
	switch tokenType {
	// we don't have right associative operators just yet
	default:
		return false
	}
}
