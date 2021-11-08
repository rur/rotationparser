// Follow implementation lexical scanning as described by R Pike here https://www.youtube.com/watch?v=HxaD_trXwRE
package rotationparser

import (
	"reflect"
	"testing"
)

func Test_lex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Lexeme
	}{
		{
			name:  "Basic",
			input: "123",
			want:  []Lexeme{{ItemNumber, "123"}, {ItemEOF, ""}},
		},
		{
			name:  "Basic negation",
			input: "-123",
			want:  []Lexeme{{ItemNumber, "-123"}, {ItemEOF, ""}},
		},
		{
			name:  "expression",
			input: "1 + 2",
			want: []Lexeme{
				{ItemNumber, "1"},
				{ItemOperator, "+"},
				{ItemNumber, "2"},
				{ItemEOF, ""},
			},
		},
		{
			name:  "conjunction operator",
			input: "123&& 567",
			want: []Lexeme{
				{ItemNumber, "123"},
				{ItemOperator, "&&"},
				{ItemNumber, "567"},
				{ItemEOF, ""},
			},
		},
		{
			name:  "compound negation",
			input: "1 + (-2)",
			want: []Lexeme{
				{ItemNumber, "1"},
				{ItemOperator, "+"},
				{ItemLeftParen, "("},
				{ItemNumber, "-2"},
				{ItemRightParen, ")"},
				{ItemEOF, ""},
			},
		},
		{
			name:  "negative expression",
			input: "199-243",
			want: []Lexeme{
				{ItemNumber, "199"},
				{ItemOperator, "-"},
				{ItemNumber, "243"},
				{ItemEOF, ""},
			},
		},
		{
			name:  "compound expression",
			input: "199-243 *   6",
			want: []Lexeme{
				{ItemNumber, "199"},
				{ItemOperator, "-"},
				{ItemNumber, "243"},
				{ItemOperator, "*"},
				{ItemNumber, "6"},
				{ItemEOF, ""},
			},
		},
		{
			name:  "decimal number",
			input: "199.243",
			want: []Lexeme{
				{ItemNumber, "199.243"},
				{ItemEOF, ""},
			},
		},
		{
			name:  "paranthesis",
			input: "(19.9 + 3) & 243",
			want: []Lexeme{
				{ItemLeftParen, "("},
				{ItemNumber, "19.9"},
				{ItemOperator, "+"},
				{ItemNumber, "3"},
				{ItemRightParen, ")"},
				{ItemOperator, "&"},
				{ItemNumber, "243"},
				{ItemEOF, ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, tokens := Lex(tt.name, tt.input)
			var got []Lexeme
			for item := range tokens {
				got = append(got, item)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_prevWidth(t *testing.T) {
	tests := []struct {
		name string
		lx   *Lexer
		want []int
	}{
		{
			name: "basic ASCII",
			lx: &Lexer{
				input:  "abcdef",
				start:  0,
				cursor: 6,
			},
			want: []int{1, 1, 1, 1, 1, 1},
		},
		{
			name: "mutlibyte rune",
			lx: &Lexer{
				input:  "de世界f",
				start:  0,
				cursor: len("de世界f"),
			},
			want: []int{1, 3, 3, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int
			for width := tt.lx.prevWidth(); width > 0; width = tt.lx.prevWidth() {
				got = append(got, width)
				tt.lx.cursor -= width
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prevWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}
