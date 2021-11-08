## Rotation Parser

_Expression parsing exploration_

Explore using tree rotations with a recursive descent parser
to enforce operator precedence when building an AST.

Just curious to find out how far this approach will take me.

### Repl command

Interactive expression parser to help while experimenting with
the expression syntax.

```
$ go run ./cmd/parseXpr
...

eXpr: 1 + 2
=> "+"
    ├── "2"
    └── "1"

eXpr: -123 * 3.4 + 5 - 19 / 4.3
=> "-"
    ├── "/"
    │    ├── "4.3"
    │    └── "19"
    └── "+"
         ├── "5"
         └── "*"
              ├── "3.4"
              └── "-123"

eXpr: ...

```