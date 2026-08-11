package lexer

// A token type
type TokenType int

const (
	TokenIdentifier TokenType = iota
	TokenKeyword
	TokenOperatorPunctuation
	TokenLiteral
)

// A token has a type and a value (source code value essentially)
type Token struct {
	Type  TokenType
	Value []rune
}
