package lexer

// A token type
type TokenType int

const (
	TokenTypeErr TokenType = iota
	TokenTypeIdentifier
	TokenTypeKeyword
	TokenTypeOperatorPunctuation
	TokenTypeLiteral
	TokenTypeWhitespace
)

// A token has a type and a value (source code value essentially)
type Token struct {
	Type  TokenType
	Value []rune
}

// Several defaults
var (
	TokenErr = Token{
		Type: TokenTypeErr,
	}
	TokenSemicolon = Token{
		Type:  TokenTypeOperatorPunctuation,
		Value: []rune{';'},
	}
	TokenLBrace = Token{
		Type:  TokenTypeOperatorPunctuation,
		Value: []rune{'{'},
	}
	TokenRBrace = Token{
		Type:  TokenTypeOperatorPunctuation,
		Value: []rune{'}'},
	}
	TokenLParen = Token{
		Type:  TokenTypeOperatorPunctuation,
		Value: []rune{'('},
	}
	TokenRParen = Token{
		Type:  TokenTypeOperatorPunctuation,
		Value: []rune{')'},
	}
	TokenLBracket = Token{
		Type:  TokenTypeOperatorPunctuation,
		Value: []rune{'['},
	}
	TokenRBracket = Token{
		Type:  TokenTypeOperatorPunctuation,
		Value: []rune{']'},
	}
	TokenWhitespace = Token{
		Type: TokenTypeWhitespace,
	}
)

type Pos struct {
	line   int
	column int
	pos    int
}

var (
	SyntheticPos Pos = Pos{
		line: -1,
	}
	EofPos Pos = Pos{
		line: -2,
	}
)
