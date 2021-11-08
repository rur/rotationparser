package rotationparser

import (
	"strings"
	"testing"
)

func mustTokenize(input string) (out []Lexeme) {
	_, items := Lex("test", input)
	for item := range items {
		if item.Type == ItemEOF {
			return
		}
		out = append(out, item)
	}
	return
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name   string
		tokens []Lexeme
		want   []string
	}{
		{
			name:   "Basic",
			tokens: mustTokenize("3+4"),
			want: []string{
				`─> "+"`,
				`    ├── "4"`,
				`    └── "3"`,
			},
		},
		{
			name:   "Compound expression same precadence",
			tokens: mustTokenize("3 + 4 - 5"),
			want: []string{
				`─> "-"`,
				`    ├── "5"`,
				`    └── "+"`,
				`         ├── "4"`,
				`         └── "3"`,
			},
		},
		{
			name:   "Compound expression varying precadence",
			tokens: mustTokenize(`3 + 4 * 5`),
			want: []string{
				`─> "+"`,
				`    ├── "*"`,
				`    │    ├── "5"`,
				`    │    └── "4"`,
				`    └── "3"`,
			},
		},
		{
			name:   "apply precedence recursively",
			tokens: mustTokenize(`3 + 4 - 5 + 6 - 7`),
			want: []string{
				`─> "-"`,
				`    ├── "7"`,
				`    └── "+"`,
				`         ├── "6"`,
				`         └── "-"`,
				`              ├── "5"`,
				`              └── "+"`,
				`                   ├── "4"`,
				`                   └── "3"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expecting := strings.Join(tt.want, "\n")
			gotOut := ParseBinaryExpression(tt.tokens)
			gotPrinted := SprintNodeTree(gotOut)
			if gotPrinted != expecting {
				t.Errorf("ParseExpression() =\n%v\n----\nwant =\n%v", gotPrinted, expecting)
			}
		})
	}
}
