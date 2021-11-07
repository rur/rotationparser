// Follow implementation lexical scanning as described by R Pike here https://www.youtube.com/watch?v=HxaD_trXwRE
package rotationparser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type itemType int

const (
	eofChar   rune     = 0
	itemError itemType = -1
	itemEOF   itemType = iota
	itemNumber
	itemOperator
	itemLeftParen
	itemRightParen
)

// item is a lexeme for this lexer
type item struct {
	typ   itemType
	value string
}

// String implements the stringer interface
func (i item) String() string {
	switch i.typ {
	case itemError:
		return i.value
	case itemEOF:
		return "EOF"
	}
	if len(i.value) > 10 {
		return fmt.Sprintf("%.10q...", i.value)
	}
	return fmt.Sprintf("%q", i.value)
}

type lexer struct {
	name   string    // for error reports
	input  string    // the full string being scanned (TODO replace with io.ReadCloseSeeker)
	start  int       // start position of this current item
	cursor int       // current cursor position in the input
	items  chan item // channel of scanned items
}

type stateFn func(*lexer) stateFn

// run will drive the lexer state transitions
func (lx *lexer) run() {
	for state := lexAny; state != nil; {
		state = state(lx)
	}
	lx.emit(itemEOF)
	close(lx.items)
}

// prevWidth will return the number of bytes to the last unicode character
// it will return zero if lx.cursor is <= lx.start
func (lx *lexer) prevWidth() int {
	for j := lx.cursor - 1; j >= lx.start; j-- {
		if lx.input[j]&192 == 128 {
			// most significant bits are 10, hence part of a uft8 multibyte suffix
			continue
		}
		return lx.cursor - j
	}
	return 0
}

// next will advance forward one rune width
func (lx *lexer) next() rune {
	if lx.cursor >= len(lx.input) {
		return eofChar
	}
	char, width := utf8.DecodeRuneInString(lx.input[lx.cursor:])
	lx.cursor += width
	return char
}

// peek will show the next character without advancing the cursor
func (lx *lexer) peek() rune {
	next := lx.next()
	lx.backup()
	return next
}

// ignore will skip all input since the last emit|ignore
func (lx *lexer) ignore() {
	lx.start = lx.cursor
}

// backup will undo next steps since
func (lx *lexer) backup() {
	lx.cursor -= lx.prevWidth()
}

// accept consumes the next rune if it is in the valid set
func (lx *lexer) accept(valid string) bool {
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
func (lx *lexer) acceptRun(valid string) {
	for {
		if !lx.accept(valid) {
			return
		}
	}
}

// emit passes a lexeme item over the channel
func (lx *lexer) emit(t itemType) {
	lx.items <- item{
		typ:   t,
		value: lx.input[lx.start:lx.cursor],
	}
	lx.start = lx.cursor
}

// error emits an lexing error and continues to scan
func (lx *lexer) errorf(template string, params ...interface{}) stateFn {
	lx.items <- item{
		typ:   itemError,
		value: fmt.Sprintf(template, params...),
	}
	return lexAny
}

// lexAny is the first state for parsing
func lexAny(lx *lexer) stateFn {
	switch char := lx.next(); {
	case matchCharset(char, charNonZeroDigit):
		lx.backup()
		return lexNumber
	case char == '-':
		// this is an sign if it is the start of the line or
		// the previous character is a left-paren (no spaces)
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
		lx.emit(itemLeftParen)
		return lexAny
	case char == ')':
		lx.emit(itemRightParen)
		return lexAny
	case char == eofChar:
		return nil
	default:
		return lx.errorf("Unexpected character encountered")
	}
}

// lexNumber will scan a decimal number with optional negation, eg -12.345
func lexNumber(lx *lexer) stateFn {
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
	lx.emit(itemNumber)
	return lexAny
}

// lexOperator will scan for operator symbols
func lexOperator(lx *lexer) stateFn {
	if lx.accept("+-*/^&|") {
		// one character wonder!
	} else {
		panic("Lexer bug unexpected operator character")
	}
	lx.emit(itemOperator)
	return lexAny
}

// lex will concurrently scan the input string, delivering lexeme items over the channel
// as they become available
func lex(name, input string) (*lexer, <-chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run() // concurrently scan input, pushing items onto the channel
	return l, l.items
}

type charTypes uint

const (
	charAlphaLower charTypes = 1 << iota
	charAlphaUpper
	charUnderscore
	charNonZeroDigit
	charZero
	charSymbol
	charWhitespace
	charNewline
)

func matchCharset(char rune, typeMask charTypes) bool {
	switch char {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return charNonZeroDigit&typeMask != 0
	case '0':
		return charZero&typeMask != 0
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
		return charAlphaLower&typeMask != 0
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
		return charAlphaUpper&typeMask != 0
	case '+', '-', '*', '/', '&', '|', '^':
		return charSymbol&typeMask != 0
	case ' ', '\t':
		return charWhitespace&typeMask != 0
	case '\n', '\r':
		return charNewline&typeMask != 0
	}
	return false
}
