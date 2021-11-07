// Lexical scanning following the outline described by R Pike here https://www.youtube.com/watch?v=HxaD_trXwRE
package rotationparser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type ItemType int

const (
	ItemError ItemType = -1
	ItemEOF   ItemType = iota
	ItemNumber
	ItemOperator
	ItemLeftParen
	ItemRightParen
)

// item is a lexeme for this lexer
type Lexeme struct {
	Type  ItemType
	Value string
}

// String implements the stringer interface
func (l Lexeme) String() string {
	switch l.Type {
	case ItemError:
		return l.Value
	case ItemEOF:
		return "EOF"
	}
	if len(l.Value) > 10 {
		return fmt.Sprintf("%.10q...", l.Value)
	}
	return fmt.Sprintf("%q", l.Value)
}

type Lexer struct {
	Name string // for error reports

	// internal
	input  string      // the full string being scanned (TODO replace with io.ReadCloser)
	start  int         // start position of this current item
	cursor int         // current cursor position in the input
	items  chan Lexeme // channel of scanned items
}

type stateFn func(*Lexer) stateFn

// run will drive the lexer state transitions
func (lx *Lexer) run() {
	for state := lexAny; state != nil; {
		state = state(lx)
	}
	lx.emit(ItemEOF)
	close(lx.items)
}

// prevWidth will return the number of bytes from the cursor to the previous unicode character
//
// Note: This will return zero if lx.cursor is <= lx.start. Preventing code from
//       attempting to scan backward from _start_ will negate the need
//       to use the io.Scanner interface when replacing the input string with an io.Reader
func (lx *Lexer) prevWidth() int {
	for j := lx.cursor - 1; j >= lx.start; j-- {
		if lx.input[j]&192 == 128 {
			// most significant bits are 10, hence part of a uft8 multibyte suffix
			continue
		}
		return lx.cursor - j
	}
	return 0
}

const eofChar rune = 0 // nullcode used by scanner

// next will advance forward one rune width
func (lx *Lexer) next() rune {
	if lx.cursor >= len(lx.input) {
		return eofChar
	}
	char, width := utf8.DecodeRuneInString(lx.input[lx.cursor:])
	lx.cursor += width
	return char
}

// peek will show the next character without advancing the cursor
func (lx *Lexer) peek() rune {
	next := lx.next()
	lx.backup()
	return next
}

// ignore will skip all input since the last emit|ignore
func (lx *Lexer) ignore() {
	lx.start = lx.cursor
}

// backup will undo next steps since
func (lx *Lexer) backup() {
	lx.cursor -= lx.prevWidth()
}

// accept consumes the next rune if it is in the valid set
func (lx *Lexer) accept(valid string) bool {
	nxt := lx.next()
	if nxt == eofChar {
		return false
	}
	if strings.ContainsRune(valid, nxt) {
		return true
	}
	lx.backup()
	return false
}

// acceptRun will consume zero or more characters in the valid rune set
func (lx *Lexer) acceptRun(valid string) {
	for {
		if !lx.accept(valid) {
			return
		}
	}
}

// emit passes a lexeme item over the channel
func (lx *Lexer) emit(t ItemType) {
	lx.items <- Lexeme{
		Type:  t,
		Value: lx.input[lx.start:lx.cursor],
	}
	lx.start = lx.cursor
}

// error emits an lexing error and continues to scan
func (lx *Lexer) errorf(template string, params ...interface{}) stateFn {
	lx.items <- Lexeme{
		Type:  ItemError,
		Value: fmt.Sprintf(template, params...),
	}
	return lexAny
}

// lexAny is the first state for parsing
func lexAny(lx *Lexer) stateFn {
	switch char := lx.next(); {
	case matchCharset(char, charNonZeroDigit):
		lx.backup()
		return lexNumber
	case char == '-':
		// this is an unary sign if it is that the start of the line or
		// if the previous character is a left-paren (\w no spaces!)
		if prev, _ := utf8.DecodeLastRuneInString(lx.input[:lx.start]); lx.cursor == 1 || prev == '(' {
			// lex the number including the sign
			lx.backup()
			return lexNumber
		}
		fallthrough
	case matchCharset(char, charSymbol):
		lx.backup()
		return lexOperator
	case matchCharset(char, charWhitespace):
		lx.ignore()
		return lexAny
	case char == '(':
		lx.emit(ItemLeftParen)
		return lexAny
	case char == ')':
		lx.emit(ItemRightParen)
		return lexAny
	case char == eofChar:
		return nil
	default:
		return lx.errorf("Unexpected character encountered")
	}
}

// lexNumber will scan a decimal number with optional negation, eg -12.345
func lexNumber(lx *Lexer) stateFn {
	lx.accept("-")
	if !lx.accept("123456789") {
		return lx.errorf("invalid number")
	}
	lx.acceptRun("0123456789")
	if lx.accept(".") {
		if !lx.accept("0123456789") {
			return lx.errorf("invalid number")
		}
		lx.acceptRun("0123456789")
	}
	lx.emit(ItemNumber)
	return lexAny
}

// lexOperator will scan for operator symbols
func lexOperator(lx *Lexer) stateFn {
	if lx.accept("+-*/^&|") {
		lx.emit(ItemOperator)
		return lexAny
	}
	panic("Lexer bug unexpected operator character")
}

// lex will concurrently scan the input string, delivering lexeme items over the channel
// as they become available
func lex(name, input string) (*Lexer, <-chan Lexeme) {
	l := &Lexer{
		Name:  name,
		input: input,
		items: make(chan Lexeme),
	}
	go l.run() // concurrently scan input, pushing items onto the channel
	return l, l.items
}

type charType uint

const (
	charAlphaLower charType = 1 << iota
	charAlphaUpper
	charUnderscore
	charNonZeroDigit
	charZero
	charSymbol
	charWhitespace
	charNewline
)

// matchCharset checks that the character match the type mask
func matchCharset(char rune, mask charType) (match bool) {
	var typ charType
	defer func() {
		match = typ&mask != 0
	}()
	switch char {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		typ = charNonZeroDigit
	case '0':
		typ = charZero
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
		typ = charAlphaLower
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
		typ = charAlphaUpper
	case '+', '-', '*', '/', '&', '|', '^':
		typ = charSymbol
	case ' ', '\t':
		typ = charWhitespace
	case '\n', '\r':
		typ = charNewline
	}
	return
}
