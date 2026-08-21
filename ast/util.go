package ast

type QualifiedIdent struct {
	PackageName string
	Name        string
}

// --------------------- Compo Lit Key --------------------------

type CompositeLiteralKey interface {
	isCompositeLiteralKey()
}

type FieldNameKey string

func (_ *FieldNameKey) isCompositeLiteralKey()

type ExpressionKey struct{ Expr Expression }

func (_ *ExpressionKey) isCompositeLiteralKey()

type LiteralValueKey []CompositeLiteralValue

func (_ *LiteralValueKey) isCompositeLiteralKey()

// -------------------- Comp lit value -------------------------

type CompositeLiteralValue interface {
	isCompositeLiteralValue()
}

type ExpressionValue struct{ Expr Expression }

func (_ *ExpressionValue) isCompositeLiteralValue() {}

type CompositeElementsValue []CompositeLiteralElement

func (_ *CompositeElementsValue) isCompositeLiteralValue() {}

// -------------------- Comp Lit

type CompositeLiteral struct {
	Type  TypeNode
	Value []CompositeLiteralElement
}

type CompositeLiteralElement struct {
	Key CompositeLiteralKey
	Val Expression
}
