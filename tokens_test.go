package rotationparser

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		wantTokens []Token
		wantErr    bool
	}{
		{
			name:       "single digit",
			code:       "9",
			wantTokens: []Token{{Factor, "9"}},
		},
		{
			name:       "multiple digits",
			code:       "1234",
			wantTokens: []Token{{Factor, "1234"}},
		},
		{
			name:       "simple expression",
			code:       "9 + 34",
			wantTokens: []Token{{Factor, "9"}, {Plus, "+"}, {Factor, "34"}},
		},
		{
			name: "compound expression",
			code: "9 + 34\t/8 * 63",
			wantTokens: []Token{
				{Factor, "9"},
				{Plus, "+"},
				{Factor, "34"},
				{Divide, "/"},
				{Factor, "8"},
				{Multiply, "*"},
				{Factor, "63"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTokens, err := Tokenize(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTokens, tt.wantTokens) {
				t.Errorf("Tokenize() = %v, want %v", gotTokens, tt.wantTokens)
			}
		})
	}
}
