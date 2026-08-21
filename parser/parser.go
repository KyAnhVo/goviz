package parser

import (
	"fmt"

	"github.com/KyAnhVo/goviz/lexer"
	"github.com/KyAnhVo/goviz/token"
)

type Parser struct {
	l         *lexer.Lexer
	currToken struct {
		token token.Token
		pos   token.Pos
	}
	nextToken struct {
		token token.Token
		pos   token.Pos
	}
}

func New(l *lexer.Lexer) (*Parser, error) {
	curr, currPos, err := l.GetNextToken()
	if err != nil {
		return nil, fmt.Errorf("Parser create first token failed: %w", err)
	}
	next, nextPos, err := l.GetNextToken()
	if err != nil {
		return nil, fmt.Errorf("Parser create second token failed: %w", err)
	}
	return &Parser{
		l: l,
		currToken: struct {
			token token.Token
			pos   token.Pos
		}{curr, currPos},
		nextToken: struct {
			token token.Token
			pos   token.Pos
		}{next, nextPos},
	}, nil
}

// advance the token
func (p *Parser) advance() error {
	p.currToken = p.nextToken
	next, pos, err := p.l.GetNextToken()
	if err != nil {
		return fmt.Errorf("advance: %w", err)
	}
	p.nextToken = struct {
		token token.Token
		pos   token.Pos
	}{next, pos}
	return nil
}

// expect a token, if not return an error
func (p *Parser) ExpectToken(tok token.Token) error {
	if p.currToken.token != tok {
		return fmt.Errorf("Expected %s, received %s, at %s",
			token.FormatToken(tok),
			token.FormatToken(p.currToken.token),
			token.FormatPos(p.currToken.pos))
	} else {
		return nil
	}
}

// expect a token type, if not return an error
func (p *Parser) ExpectTokenType(tokenType token.TokenType) error {
	if p.currToken.token.Type != tokenType {
		return fmt.Errorf("Expected %s, received %s, at %s",
			tokenType, token.FormatToken(p.currToken.token),
			token.FormatPos(p.currToken.pos))
	} else {
		return nil
	}
}
