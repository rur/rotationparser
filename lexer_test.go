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
		want  []item
	}{
		{
			name:  "Basic",
			input: "123",
			want:  []item{{itemNumber, "123"}, {itemEOF, ""}},
		},
		{
			name:  "Basic negation",
			input: "-123",
			want:  []item{{itemNumber, "-123"}, {itemEOF, ""}},
		},
		{
			name:  "expression",
			input: "1 + 2",
			want: []item{
				{itemNumber, "1"},
				{itemOperator, "+"},
				{itemNumber, "2"},
				{itemEOF, ""},
			},
		},
		{
			name:  "conjunction operator",
			input: "123&& 567",
			want: []item{
				{itemNumber, "123"},
				{itemOperator, "&&"},
				{itemNumber, "567"},
				{itemEOF, ""},
			},
		},
		{
			name:  "compound negation",
			input: "1 + (-2)",
			want: []item{
				{itemNumber, "1"},
				{itemOperator, "+"},
				{itemLeftParen, "("},
				{itemNumber, "-2"},
				{itemRightParen, ")"},
				{itemEOF, ""},
			},
		},
		{
			name:  "negative expression",
			input: "199-243",
			want: []item{
				{itemNumber, "199"},
				{itemOperator, "-"},
				{itemNumber, "243"},
				{itemEOF, ""},
			},
		},
		{
			name:  "compound expression",
			input: "199-243 *   6",
			want: []item{
				{itemNumber, "199"},
				{itemOperator, "-"},
				{itemNumber, "243"},
				{itemOperator, "*"},
				{itemNumber, "6"},
				{itemEOF, ""},
			},
		},
		{
			name:  "decimal number",
			input: "199.243",
			want: []item{
				{itemNumber, "199.243"},
				{itemEOF, ""},
			},
		},
		{
			name:  "paranthesis",
			input: "(19.9 + 3) & 243",
			want: []item{
				{itemLeftParen, "("},
				{itemNumber, "19.9"},
				{itemOperator, "+"},
				{itemNumber, "3"},
				{itemRightParen, ")"},
				{itemOperator, "&"},
				{itemNumber, "243"},
				{itemEOF, ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, tokens := lex(tt.name, tt.input)
			var got []item
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
		lx   *lexer
		want []int
	}{
		{
			name: "basic ASCII",
			lx: &lexer{
				input:  "abcdef",
				start:  0,
				cursor: 6,
			},
			want: []int{1, 1, 1, 1, 1, 1},
		},
		{
			name: "mutlibyte rune",
			lx: &lexer{
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
