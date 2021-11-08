package rotationparser

import (
	"strings"
	"testing"
)

func TestSprintNodeTree(t *testing.T) {
	tests := []struct {
		name string
		v    *Node
		want []string
	}{
		{
			name: "basic",
			v:    &Node{Item: Lexeme{ItemNumber, "123"}},
			want: []string{
				`─> "123"`,
			},
		},
		{
			name: "addition expression",
			v: &Node{
				Item:  Lexeme{ItemOperator, "+"},
				Left:  &Node{Item: Lexeme{ItemNumber, "2"}},
				Right: &Node{Item: Lexeme{ItemNumber, "3"}},
			},
			want: []string{
				`─> "+"`,
				`    ├── "3"`,
				`    └── "2"`,
			},
		},
		{
			name: "nested expression, left hand side",
			v: &Node{
				Item: Lexeme{ItemOperator, "+"},
				Left: &Node{
					Item:  Lexeme{ItemOperator, "*"},
					Left:  &Node{Item: Lexeme{ItemNumber, "2"}},
					Right: &Node{Item: Lexeme{ItemNumber, "3"}},
				},
				Right: &Node{Item: Lexeme{ItemNumber, "4"}},
			},
			want: []string{
				`─> "+"`,
				`    ├── "4"`,
				`    └── "*"`,
				`         ├── "3"`,
				`         └── "2"`,
			},
		},
		{
			name: "nested expression, right hand side",
			v: &Node{
				Item: Lexeme{ItemOperator, "+"},
				Right: &Node{
					Item:  Lexeme{ItemOperator, "*"},
					Left:  &Node{Item: Lexeme{ItemNumber, "2"}},
					Right: &Node{Item: Lexeme{ItemNumber, "3"}},
				},
				Left: &Node{Item: Lexeme{ItemNumber, "4"}},
			},
			want: []string{
				`─> "+"`,
				`    ├── "*"`,
				`    │    ├── "3"`,
				`    │    └── "2"`,
				`    └── "4"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expect := strings.Join(tt.want, "\n")
			if got := SprintNodeTree(tt.v); got != expect {
				t.Errorf("SprintNodeTree() =\n%v\n--------\nwant =\n%v", got, expect)
			}
		})
	}
}
