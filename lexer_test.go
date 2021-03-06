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
			name:  "bitwise operator",
			input: "123& 567",
			want: []Lexeme{
				{ItemNumber, "123"},
				{ItemOperator, "&"},
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
				input:  "de??????f",
				start:  0,
				cursor: len("de??????f"),
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

func Test_matchCharset(t *testing.T) {
	type args struct {
		char rune
		mask charType
	}
	tests := []struct {
		name      string
		args      args
		wantMatch bool
	}{
		{
			name:      "basic alpha",
			args:      args{'a', charAlphaLower},
			wantMatch: true,
		},
		{
			name:      "basic uppercase",
			args:      args{'A', charAlphaUpper},
			wantMatch: true,
		},
		{
			name:      "basic digit",
			args:      args{'9', charNonZeroDigit},
			wantMatch: true,
		},
		{
			name:      "zero",
			args:      args{'0', charZero},
			wantMatch: true,
		},
		{
			name:      "digit including zero",
			args:      args{'0', charNonZeroDigit | charZero},
			wantMatch: true,
		},
		{
			name:      "case insensitive 1",
			args:      args{'a', charAlphaLower | charAlphaUpper},
			wantMatch: true,
		},
		{
			name:      "case insensitive 2",
			args:      args{'A', charAlphaLower | charAlphaUpper},
			wantMatch: true,
		},
		{
			name:      "not character",
			args:      args{'_', charAlphaLower | charAlphaUpper},
			wantMatch: false,
		},
		{
			name:      "not digit",
			args:      args{'_', charNonZeroDigit},
			wantMatch: false,
		},
		{
			name:      "symbol",
			args:      args{'+', charSymbol},
			wantMatch: true,
		},
		{
			name:      "not symbol",
			args:      args{'#', charSymbol},
			wantMatch: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMatch := matchCharset(tt.args.char, tt.args.mask); gotMatch != tt.wantMatch {
				t.Errorf("matchCharset() = %v, want %v", gotMatch, tt.wantMatch)
			}
		})
	}
}
