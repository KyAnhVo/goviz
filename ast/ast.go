package ast

type AST struct {
	PackageName   string
	ImportedFiles Imports
	Types         []TypeNode
}

// all the packages the file imports from
type Imports []Import

// import a "b"
// - a is name, could be empty.
// - b is src
type Import struct {
	Name string
	Src  string
}
