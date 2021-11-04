package rotationparser

import (
	"io"
	"strings"
)

func SprintNodeTree(v *Node) string {
	if v == nil {
		return "- nil"
	}
	str := strings.Builder{}
	str.WriteString("- ")
	str.WriteString(v.Token.Code)
	fprintViewTree(&str, []byte("  "), []*Node{v.Left, v.Right})
	return str.String()
}

// fprintViewTree delves recursively into view and sub views and writes
// a tree prepresentation of the supplied view
func fprintViewTree(w io.Writer, prefix []byte, children []*Node) {
	for i, sub := range children {
		if sub == nil {
			continue
		}
		last := i == len(children)-1
		w.Write(append([]byte{'\n'}, prefix...))
		if last {
			w.Write([]byte("'- " + sub.Token.Code))
		} else {
			w.Write([]byte("|- " + sub.Token.Code))
		}
		var subPrefix []byte
		if last {
			subPrefix = append(prefix, []byte("   ")...)
		} else {
			subPrefix = append(prefix, []byte("|  ")...)
		}
		fprintViewTree(w, subPrefix, []*Node{sub.Left, sub.Right})
	}
}
