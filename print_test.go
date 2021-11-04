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
			v:    &Node{Token: Token{Factor, "123"}},
			want: []string{
				"- 123",
			},
		},
		{
			name: "addition expression",
			v: &Node{
				Token: Token{Multiply, "+"},
				Left:  &Node{Token: Token{Factor, "2"}},
				Right: &Node{Token: Token{Factor, "3"}},
			},
			want: []string{
				"- +",
				"  |- 2",
				"  '- 3",
			},
		},
		{
			name: "nested expression, left hand side",
			v: &Node{
				Token: Token{Multiply, "+"},
				Left: &Node{
					Token: Token{Multiply, "*"},
					Left:  &Node{Token: Token{Factor, "2"}},
					Right: &Node{Token: Token{Factor, "3"}},
				},
				Right: &Node{Token: Token{Factor, "4"}},
			},
			want: []string{
				"- +",
				"  |- *",
				"  |  |- 2",
				"  |  '- 3",
				"  '- 4",
			},
		},
		{
			name: "nested expression, right hand side",
			v: &Node{
				Token: Token{Multiply, "+"},
				Right: &Node{
					Token: Token{Multiply, "*"},
					Left:  &Node{Token: Token{Factor, "2"}},
					Right: &Node{Token: Token{Factor, "3"}},
				},
				Left: &Node{Token: Token{Factor, "4"}},
			},
			want: []string{
				"- +",
				"  |- 4",
				"  '- *",
				"     |- 2",
				"     '- 3",
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
