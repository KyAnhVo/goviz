# goviz
Go function dependency tracker

## Purpose
We construct the function call graph as a directed graph, where
`f -> g` denotes function `f` calls function `g`. "function" here
is also applicable for methods:
``` golang
func (f *A) f() {
  // ...
  var b B
  b.g()
  // ...
}
```
creates the edge `A.f -> B.g`.

## Secondary purpose
This project does not use golang `go/scanner` and `go/parser`, and this project is used to learn
lexing/parsing/semantic analysis from scratch. Probably will use `go/scanner` and `go/parser` for
the main version, but not the first versions.

## State:
[X] Lexer
[ ] AST type system
[ ] Parser
