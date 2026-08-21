package parser

import (
	"fmt"

	"github.com/KyAnhVo/goviz/ast"
	"github.com/KyAnhVo/goviz/token"
)

func (p *Parser) typeclass() (ast.TypeNode, error) {
	switch p.currToken.token {
	case token.TokenLBracket:
		if p.nextToken.token == token.TokenRBracket {
			return p.sliceType()
		} else {
			return p.arrType()
		}
	case token.TokenKeywordStruct:
		return p.structType()
	case token.TokenKeywordInterface:
		return p.interfaceType()
	case token.TokenIdentifier(p.currToken.token.Value):
		return p.typename()
	case token.TokenLParen:
		p.advance()
		t, err := p.typeclass()
		if err != nil {
			return nil, err
		}
		err = p.ExpectToken(token.TokenRParen)
		if err != nil {
			return nil, err
		}
		return t, nil
	}

	return nil, fmt.Errorf("Expected type creation, received %s, at %s",
		token.FormatToken(p.currToken.token), token.FormatPos(p.currToken.pos))
}

func (p *Parser) sliceType() (*ast.Slice, error) {
	panic("")
}

func (p *Parser) arrType() (*ast.Array, error) {
	panic("")
}

func (p *Parser) structType() (*ast.Struct, error) {
	panic("")
}

func (p *Parser) definedType() (*ast.DefinedType, error) {
	panic("")
}

func (p *Parser) typename() (*ast.TypeName, error) {
	panic("")
}

func (p *Parser) interfaceType() (*ast.Interface, error) {
	panic("")
}
