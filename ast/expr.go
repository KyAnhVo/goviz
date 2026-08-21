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
	Parameters() []Expression
}

// Unary expressions in Golang is strictly prefix.
type UnaryExpression struct {
	Op      token.Token
	Operand Expression
}

func (expr *UnaryExpression) ExpressionType() ExpressionType { return UnaryExpr }
func (expr *UnaryExpression) Operator() string {
	switch expr.Op {
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
func (expr *UnaryExpression) Parameters() []Expression { return []Expression{expr.Operand} }

// Binary expressions in golang are of the form (operand OP operand)
type BinaryExpression struct {
	Op       token.Token
	LOperand *Expression
	ROperand *Expression
}

func (expr *BinaryExpression) ExpressionType() ExpressionType { return UnaryExpr }
func (expr *BinaryExpression) Operator() string               { return expr.Op.Value }
func (expr *BinaryExpression) Parameters() []*Expression {
	return []*Expression{expr.LOperand, expr.ROperand}
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

func IsUnaryOp(tok token.Token) bool { return unaryOps.Contains(tok) }
func IsBinOp(tok token.Token) bool {
	return IsAddOp(tok) || IsRelOp(tok) || IsMulOp(tok) || otherBinaryOps.Contains(tok)
}
func IsMulOp(tok token.Token) bool { return mulOps.Contains(tok) }
func IsAddOp(tok token.Token) bool { return addOps.Contains(tok) }
func IsRelOp(tok token.Token) bool { return relOps.Contains(tok) }

func lbp(tok token.Token) int {
	if IsMulOp(tok) {
		return 50
	} else if IsAddOp(tok) {
		return 40
	} else if IsRelOp(tok) {
		return 30
	} else if tok == token.TokenLAnd {
		return 20
	} else if tok == token.TokenLOr {
		return 10
	} else {
		return -1
	}
}

func Rbp(tok token.Token) int {
	return lbp(tok) + 1
}

func Bp(tok token.Token) int {
	if IsUnaryOp(tok) {
		return 100
	} else {
		return -1
	}
}
