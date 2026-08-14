package lexer

import (
	"errors"
	"fmt"
	"slices"
)

type Lexer struct {
	src []rune
	ptr int

	extraBuf           rune
	currentChar        rune
	canInsertSemicolon bool

	prevLineLen int
	pos         Pos
}

func (l *Lexer) GetNextToken() (Token, Pos, error) {
	for {
		c1 := l.peekNextChar()
		c2 := l.peekOffset(2)
		//	c3 := l.peekOffset(3)

		var token Token
		var pos Pos
		if c1 == '\u0000' {
			if l.canInsertSemicolon {
				// endline and synthetic semicolon has same position
				pos = l.getCurrentPos()
				pos.column += 1
				token = TokenSemicolon
			} else {
				token = TokenEOF
				pos = PosEOF
			}
		} else if isWhitespaceNonNewline(c1) {
			l.getNextChar()
			continue
		} else if isNewline(c1) {
			if l.canInsertSemicolon {
				pos = l.getCurrentPos()
				pos.column += 1
				token = TokenSemicolon
			} else {
				l.getNextChar()
				continue
			}
		} else if isLetter(c1) {
			token, pos = l.getIdentifierOrKeyword()
		} else if c1 == '/' && c2 == '/' {
			l.getLineComment()
			continue
		} else if c1 == '/' && c2 == '*' {
			err := l.getGeneralComment()
			if err != nil {
				return TokenErr, PosErr, fmt.Errorf(
					"Error lexing: line %d, column %d, position %d: %w",
					pos.line, pos.column, pos.pos, err,
				)
			}
			continue
		} else {
			return TokenErr, PosErr, fmt.Errorf(
				"Error lexing: line %d, column %d, position %d: %w",
				pos.line, pos.column, pos.pos,
				errors.New("Uncategorized character: "+string(c1)),
			)
		}

		l.setupSemicolonInsertNewline(token)
		return token, pos, nil
	}
}

// ---------------------------- Utility ----------------------------

func NewLexer(src []rune) *Lexer {
	return &Lexer{
		src: src,
		ptr: 0,

		extraBuf:           '\u0000',
		currentChar:        '\u0000',
		canInsertSemicolon: false,

		pos: Pos{
			line:   1,
			column: 0,
			pos:    0,
		},
	}
}

func (l *Lexer) peekNextChar() rune {
	return l.peekOffset(1)
}

func (l *Lexer) getNextChar() (rune, Pos) {
	l.currentChar = l.peekNextChar()

	var currentPos Pos
	if l.currentChar != '\u0000' {
		if l.extraBuf != '\u0000' {
			// case 1: extrabuf is nonempty.
			// Then we reset extrabuf
			currentPos = PosSynthetic
			l.extraBuf = '\u0000'
		} else {
			// case 2: extrabuf is empty.
			// Thus next char src[ptr].
			currentPos = l.adjustPos(l.currentChar)
			l.ptr += 1
		}
	} else {
		currentPos = PosEOF
	}

	return l.currentChar, currentPos
}

func (l *Lexer) getCurrentPos() Pos {
	semanticPos := l.pos
	if semanticPos.column == 0 {
		semanticPos.column = l.prevLineLen + 1
		semanticPos.line -= 1
	}
	return semanticPos
}

func (l *Lexer) getCurrentChar() (rune, Pos) {
	return l.currentChar, l.getCurrentPos()
}

func (l *Lexer) adjustPos(c rune) Pos {
	if isNewline(c) {
		l.prevLineLen = l.pos.column
		l.pos.column = 0
		l.pos.line += 1
	} else {
		l.pos.column += 1
	}
	l.pos.pos += 1

	return l.getCurrentPos()
}

func (l *Lexer) peekOffset(offset int) rune {
	if offset == 0 {
		panic("offset must be positive")
	}

	// induction base case: for next char, offset = 1,
	// thus we must subtract from offset
	// inductive step: trivial
	next := l.ptr + offset - 1
	if l.extraBuf != '\u0000' {
		next -= 1
	}
	if next < l.ptr { // if this is true, l.extraBuffer must be nonempty
		return l.extraBuf
	}
	if len(l.src) > next {
		return l.src[next]
	}
	return ('\u0000')
}

// Setup if the next newline must add a semicolon before.
// refer to: rule 1 of https://go.dev/ref/spec#Semicolons
//
// When the input is broken into tokens, a semicolon is
// automatically inserted into the token stream immediately
// after a line's final token if that token is
//
//	an identifier
//	an integer, floating-point, imaginary, rune, or string literal
//	one of the keywords break, continue, fallthrough, or return
//	one of the operators and punctuation ++, --, ), ], or }
func (l *Lexer) setupSemicolonInsertNewline(token Token) {
	insertSemicolonTokenTypes := []TokenType{
		TokenTypeIdentifier, TokenTypeFloatLiteral,
		TokenTypeImaginaryLiteral, TokenTypeIntLiteral,
		TokenTypeStringLiteral, TokenTypeRuneLiteral,
	}
	insertSemicolonTokens := []Token{
		TokenKeywordBreak, TokenKeywordReturn,
		TokenKeywordFallthrough, TokenIncrement,
		TokenDecrement, TokenRBrace,
		TokenRBracket, TokenRParen,
	}

	l.canInsertSemicolon =
		slices.Contains(insertSemicolonTokens, token) ||
			slices.Contains(insertSemicolonTokenTypes, token.Type)
}
