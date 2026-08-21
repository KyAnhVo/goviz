package ast

import (
	"github.com/KyAnhVo/goviz/token"
	"github.com/KyAnhVo/goviz/util"
)

type ExpressionType int

const (
	UnaryExpr ExpressionType = iota
	BinaryExpr
)

type Expression interface {
	ExpressionType() ExpressionType
	Operator() string
	Parameters() []*Expression
}

type UnaryExpression struct {
	operator  token.Token
	parameter *Expression
}

func (expr *UnaryExpression) ExpressionType() ExpressionType { return UnaryExpr }
func (expr *UnaryExpression) Operator() string {
	switch expr.operator {
	case token.TokenAdd:
		return "positive"
	case token.TokenSub:
		return "negative"
	case token.TokenNot:
		return "not"
	case token.TokenXor:
		return "bitwiseNot"
	case token.TokenMul:
		return "ptrDeref"
	case token.TokenAnd:
		return "addrOf"
	case token.TokenChannel:
		return "channel"
	}
	return ""
}
func (expr *UnaryExpression) Parameters() []*Expression { return []*Expression{expr.parameter} }

type BinaryExpression struct {
	operator token.Token
	lOperand *Expression
	rOperand *Expression
}

func (expr *BinaryExpression) ExpressionType() ExpressionType { return UnaryExpr }
func (expr *BinaryExpression) Operator() string               { return expr.operator.Value }
func (expr *BinaryExpression) Parameters() []*Expression {
	return []*Expression{expr.lOperand, expr.rOperand}
}

var unaryOps util.Set[token.Token] = util.NewSet([]token.Token{
	token.TokenAdd, token.TokenSub, token.TokenNot, token.TokenXor,
	token.TokenMul, token.TokenAnd, token.TokenChannel,
})

var mulOps util.Set[token.Token] = util.NewSet([]token.Token{
	token.TokenMul, token.TokenDiv, token.TokenMod,
	token.TokenLShift, token.TokenRShift,
	token.TokenAnd, token.TokenAndNot,
})

var addOps util.Set[token.Token] = util.NewSet([]token.Token{
	token.TokenAdd, token.TokenSub, token.TokenOr, token.TokenXor,
})

var relOps util.Set[token.Token] = util.NewSet([]token.Token{
	token.TokenEqual, token.TokenNotEqual, token.TokenLt,
	token.TokenLte, token.TokenGt, token.TokenGte,
})

var otherBinaryOps util.Set[token.Token] = util.NewSet([]token.Token{
	token.TokenLOr, token.TokenLAnd,
})

func isUnaryOp(tok token.Token) bool { return unaryOps.Contains(tok) }
func isBinOp(tok token.Token) bool {
	return isAddOp(tok) || isRelOp(tok) || isMulOp(tok) || otherBinaryOps.Contains(tok)
}
func isMulOp(tok token.Token) bool { return mulOps.Contains(tok) }
func isAddOp(tok token.Token) bool { return addOps.Contains(tok) }
func isRelOp(tok token.Token) bool { return relOps.Contains(tok) }

func lbp(tok token.Token) int {
	if isMulOp(tok) {
		return 50
	} else if isAddOp(tok) {
		return 40
	} else if isRelOp(tok) {
		return 30
	} else if tok == token.TokenLAnd {
		return 20
	} else if tok == token.TokenLOr {
		return 10
	} else {
		return -1
	}
}

func rbp(tok token.Token) int {
	return lbp(tok) + 1
}
