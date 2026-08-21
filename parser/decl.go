package parser

import (
	"github.com/KyAnhVo/goviz/token"
	"github.com/KyAnhVo/goviz/util"
)

var topLevelDeclStarters util.Set[token.Token] = util.NewSet([]token.Token{
	token.TokenKeywordConst,
	token.TokenKeywordVar,
	token.TokenKeywordType,
	token.TokenKeywordFunc,
})

func (p *Parser) topLevelDecl() {
	switch p.currToken.token {
	case token.TokenKeywordConst, token.TokenKeywordVar, token.TokenKeywordType:
		p.decl()
	case token.TokenKeywordFunc:
		if p.nextToken.token == token.TokenLParen {
			p.methodDecl()
		} else if p.nextToken.token.Type == token.TokenTypeIdentifier {
			p.funcDecl()
		}
	}
}

func (p *Parser) decl() {
	switch p.currToken.token {
	case token.TokenKeywordConst:
		p.constDecl()
	case token.TokenKeywordVar:
		p.varDecl()
	case token.TokenKeywordType:
		p.typeDecl()
	}
}

func (p *Parser) constDecl() {
	// TODO: implement constDecl
}

func (p *Parser) varDecl() {
	// TODO: implement varDecl
}

func (p *Parser) typeDecl() {
	// TODO: implement typeDecl
}

func (p *Parser) funcDecl() {
	// TODO: implement funcDecl
}

func (p *Parser) methodDecl() {
	// TODO: implement methodDecl
}
