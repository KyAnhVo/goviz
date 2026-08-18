package ast

type AST struct {
	PackageName   string
	ImportedFiles Imports
	Types         []Type
}

// all the packages the file imports from
type Imports []Import

// import a "b"
// - a is name
// - b is src
type Import struct {
	Name string
	Src  string
}

// ------------------ Types -------------------

type Type interface {
	typeNode()
}

type Struct struct{}

func (st *Struct) typeNode() {}

type Function struct{}

func (f *Function) typeNode() {}
