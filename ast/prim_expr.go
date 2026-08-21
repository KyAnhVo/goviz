package ast

type PrimaryExpressionType int

const (
	PrimaryExprOperand PrimaryExpressionType = iota
	PrimaryExprConv
	PrimaryExprMethod
)

type PrimaryExpression struct {
}
