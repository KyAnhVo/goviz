/* Implements the Pratt parser algorithm for expression


https://go.dev/ref/spec#Operators

Expression = UnaryExpr | Expression binary_op Expression .
UnaryExpr  = PrimaryExpr | unary_op UnaryExpr .

binary_op  = "||" | "&&" | rel_op | add_op | mul_op .
rel_op     = "==" | "!=" | "<" | "<=" | ">" | ">=" .
add_op     = "+" | "-" | "|" | "^" .
mul_op     = "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" .

unary_op   = "+" | "-" | "!" | "^" | "*" | "&" | "<-" .
*/

package parser

import "github.com/KyAnhVo/goviz/ast"

func (p *Parser) expr() *ast.Expression {
	return p.parseExpr(0)
}

func (p *Parser) parseExpr(minBP int) *ast.Expression {
	panic("")
}
