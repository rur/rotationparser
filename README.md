## Rotation Parser

_Expression parsing exploration_

Explore using tree rotations with a recursive descent parser
to enforce operator precedence when building an AST.

Just curious to find out how far this approach will take me.

### Repl command

Interactive expression parser to help while experimenting with
the expression syntax.

```bash
$ go run ./cmd/parseXpr
...

eXpr: 1 + 2
=> "+"
    ├── "2"
    └── "1"
...

```