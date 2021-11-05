package rotationparser

import (
	"strings"
	"testing"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name   string
		tokens []Token
		want   []string
	}{
		{
			name:   "Basic",
			tokens: MustTokenize("3+4"),
			want: []string{
				"- +",
				"  |- 3",
				"  '- 4",
			},
		},
		{
			name:   "Compound expression same precadence",
			tokens: MustTokenize("3 + 4 - 5"),
			want: []string{
				"- (-)",
				"  |- +",
				"  |  |- 3",
				"  |  '- 4",
				"  '- 5",
			},
		},
		{
			name:   "Compound expression varying precadence",
			tokens: MustTokenize(`3 + 4 * 5`),
			want: []string{
				"- *",
				"  |- +",
				"  |  |- 3",
				"  |  '- 4",
				"  '- 5",
			},
		},
		{
			name:   "apply precedence recursively",
			tokens: MustTokenize(`3 + 4 - 5 + 6 - 7`),
			want: []string{
				"- (-)",
				"  |- +",
				"  |  |- (-)",
				"  |  |  |- +",
				"  |  |  |  |- 3",
				"  |  |  |  '- 4",
				"  |  |  '- 5",
				"  |  '- 6",
				"  '- 7",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expecting := strings.Join(tt.want, "\n")
			gotOut := ParseExpression(tt.tokens)
			gotPrinted := SprintNodeTree(gotOut)
			if gotPrinted != expecting {
				t.Errorf("ParseExpression() =\n%v\n----\nwant =\n%v", gotPrinted, expecting)
			}
		})
	}
}
