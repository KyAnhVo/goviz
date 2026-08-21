package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"unicode/utf8"

	types "github.com/KyAnhVo/goviz/token"
	"github.com/KyAnhVo/goviz/util"
)

// this is essentially a full program invariant. Panics if not followed.
const MAX_PEEK_LENGTH = 3

type Lexer struct {
	src      *bufio.Scanner
	peekBuf  *util.Queue[rune]
	extraBuf rune

	currentChar        rune
	canInsertSemicolon bool
	isEof              bool

	prevLineLen int
	pos         types.Pos
}

var singleCharOperators = util.NewSet([]rune{
	'+', '-', '*', '/', '%',
	'^', '&', '|', '<', '>',
	'=', '!', '~', '(', ')',
	'{', '}', '[', ']', ',',
	';', '.', ':',
})

var doubleCharOperators = util.NewSet([]string{
	"<<", ">>", "&^", "+=", "-=",
	"*=", "/=", "%=", "&=", "|=",
	"^=", "&&", "||", "<-", "++",
	"--", "==", "!=", "<=", ">=",
	":=",
})

var tripleCharOperators = util.NewSet([]string{
	"<<=", ">>=", "&^=", "...",
})

func (l *Lexer) GetNextToken() (types.Token, types.Pos, error) {
	for {
		c1 := l.peekNextChar()
		c2 := l.peekOffset(2)
		c3 := l.peekOffset(3)

		twoChar := string(c1) + string(c2)
		threeChars := twoChar + string(c3)

		var token types.Token
		var pos types.Pos
		if c1 == '\u0000' { // Null and whitespace starts here ---------------------------------
			if l.canInsertSemicolon {
				// endline and synthetic semicolon has same position
				pos = l.getCurrentPos()
				pos.Column += 1
				token = types.TokenSemicolon
			} else {
				token = types.TokenEOF
				pos = types.PosEOF
			}
		} else if isWhitespaceNonNewline(c1) {
			l.getNextChar()
			continue
		} else if isNewline(c1) {
			if l.canInsertSemicolon {
				pos = l.getCurrentPos()
				pos.Column += 1
				token = types.TokenSemicolon
			} else {
				l.getNextChar()
				continue
			}
		} else if isLetter(c1) { // Variable length starts here -----------------------------------
			var err error
			token, pos, err = l.getIdentifierOrKeyword()
			if err != nil {
				pos = l.pos
				return types.TokenErr, types.PosErr, fmt.Errorf(
					"Error lexing: line %d, column %d, position %d: %w",
					pos.Line, pos.Column, pos.Pos, err,
				)
			}
		} else if twoChar == "//" {
			l.getLineComment()
			continue
		} else if twoChar == "/*" {
			err := l.getGeneralComment()
			if err != nil {
				pos = l.pos
				return types.TokenErr, types.PosErr, fmt.Errorf(
					"Error lexing: line %d, column %d, position %d: %w",
					pos.Line, pos.Column, pos.Pos, err,
				)
			}
			continue
		} else if isDecimalDigit(c1) || c1 == '.' && isDecimalDigit(c2) {
			newToken, newPos, err := l.getNumericLiteral()
			if err != nil {
				pos = l.pos
				return types.TokenErr, types.PosErr, fmt.Errorf(
					"Lexer.GetNextToken: line %d, column %d, position %d: %w",
					pos.Line, pos.Column, pos.Pos, err,
				)
			}
			token = newToken
			pos = newPos

		} else if tripleCharOperators.Contains(threeChars) { // 3 char start here ------------------
			for range 3 {
				_, _, err := l.getNextChar()
				if err != nil {
					return types.TokenErr, types.PosErr, fmt.Errorf(
						"Lexer.GetNextToken: line %d, column %d, position %d: %w",
						l.pos.Line, l.pos.Column, l.pos.Pos, err,
					)
				}
			}
			token = types.TokenOperator(threeChars)
		} else if doubleCharOperators.Contains(twoChar) { // 2 char start here ------------------
			for range 2 {
				_, _, err := l.getNextChar()
				if err != nil {
					return types.TokenErr, types.PosErr, fmt.Errorf(
						"Lexer.GetNextToken: line %d, column %d, position %d: %w",
						l.pos.Line, l.pos.Column, l.pos.Pos, err,
					)
				}
			}
			token = types.TokenOperator(twoChar)
		} else if singleCharOperators.Contains(c1) { // 1 char start here ---------------------
			_, _, err := l.getNextChar()
			if err != nil {
				return types.TokenErr, types.PosErr, fmt.Errorf(
					"Lexer.GetNextToken: line %d, column %d, position %d: %w",
					l.pos.Line, l.pos.Column, l.pos.Pos, err,
				)
			}
			token = types.TokenOperator(string(c1))
		} else if c1 == '"' {
			newToken, newPos, err := l.getInterpretedStringToken()
			if err != nil {
				pos = l.pos
				return types.TokenErr, types.PosErr, fmt.Errorf(
					"Lexer.GetNextToken: line %d, column %d, position %d: %w",
					pos.Line, pos.Column, pos.Pos, err,
				)
			}
			token, pos = newToken, newPos
		} else if c1 == '\'' {
			newToken, newPos, err := l.getRuneToken()
			if err != nil {
				pos = l.pos
				return types.TokenErr, types.PosErr, fmt.Errorf(
					"Lexer.GetNextToken: line %d, column %d, position %d: %w",
					pos.Line, pos.Column, pos.Pos, err,
				)
			}
			token, pos = newToken, newPos
		} else if c1 == '`' {
			newToken, newPos, err := l.getRawStringToken()
			if err != nil {
				pos = l.pos
				return types.TokenErr, types.PosErr, fmt.Errorf(
					"Lexer.GetNextToken: line %d, column %d, position %d: %w",
					pos.Line, pos.Column, pos.Pos, err,
				)
			}
			token, pos = newToken, newPos
		} else {
			pos = l.pos
			return types.TokenErr, types.PosErr, fmt.Errorf(
				"Lexer.GetNextToken: line %d, column %d, position %d: %w",
				pos.Line, pos.Column, pos.Pos,
				errors.New("Uncategorized character: "+string(c1)),
			)
		}

		l.setupSemicolonInsertNewline(token)
		return token, pos, nil
	}
}

// ---------------------------- Utility ----------------------------

func NewLexer(src *bufio.Scanner) (*Lexer, error) {
	src.Split(bufio.ScanRunes)

	q := util.NewQueue[rune](MAX_PEEK_LENGTH)

	for range MAX_PEEK_LENGTH {
		hasScan := src.Scan()
		if !hasScan {
			return nil, errors.New(
				"File too short (file starts with \"package\" min 7 words)")
		}
		c, _ := utf8.DecodeRune(src.Bytes())
		q.Enqueue(c)
	}

	return &Lexer{
		src:      src,
		peekBuf:  q,
		extraBuf: '\u0000',

		currentChar:        '\u0000',
		canInsertSemicolon: false,
		isEof:              false,

		pos: types.Pos{
			Line:   1,
			Column: 0,
			Pos:    0,
		},
	}, nil
}

func (l *Lexer) peekNextChar() rune {
	return l.peekOffset(1)
}

func (l *Lexer) getNextChar() (rune, types.Pos, error) {
	l.currentChar = l.peekNextChar()

	var currentPos types.Pos
	if l.currentChar != '\u0000' {
		if l.extraBuf != '\u0000' {
			// case 1: extrabuf is nonempty.
			// Then we reset extrabuf
			currentPos = types.PosSynthetic
			l.extraBuf = '\u0000'
		} else {
			// case 2: extrabuf is empty.
			// Thus next char is from peekBuf,
			// and we advance the scanner
			currentPos = l.adjustPos(l.currentChar)
			l.peekBuf.Dequeue()

			// enqueue the buffer
			if l.isEof {
				l.peekBuf.Enqueue('\u0000')
			} else {
				scanned := l.src.Scan()
				if !scanned {
					if l.src.Err() == nil {
						l.isEof = true
					} else {
						return utf8.RuneError, types.PosErr, fmt.Errorf("File error: %w", l.src.Err())
					}
				} else {
					c, _ := utf8.DecodeRune(l.src.Bytes())
					if c == utf8.RuneError {
						return utf8.RuneError, types.PosErr, errors.New("Invalid utf8 in file")
					}
					l.peekBuf.Enqueue(c)
				}
			}
		}
	} else {
		currentPos = types.PosEOF
	}

	return l.currentChar, currentPos, nil
}

func (l *Lexer) getCurrentPos() types.Pos {
	semanticPos := l.pos
	if semanticPos.Column == 0 {
		semanticPos.Column = l.prevLineLen + 1
		semanticPos.Line -= 1
	}
	return semanticPos
}

func (l *Lexer) getCurrentChar() (rune, types.Pos) {
	return l.currentChar, l.getCurrentPos()
}

func (l *Lexer) adjustPos(c rune) types.Pos {
	if isNewline(c) {
		l.prevLineLen = l.pos.Column
		l.pos.Column = 0
		l.pos.Line += 1
	} else {
		l.pos.Column += 1
	}
	l.pos.Pos += 1

	return l.getCurrentPos()
}

func (l *Lexer) peekOffset(offset int) rune {
	if offset == 0 {
		panic("peekOffset: offset must be positive")
	} else if offset > MAX_PEEK_LENGTH {
		panic("peekOffset: offset > MAX_PEEK_LENGTH")
	}

	// induction base case: for next char, offset = 1,
	// thus we must subtract from offset
	// inductive step: trivial
	next := offset - 1
	if l.extraBuf != '\u0000' {
		next -= 1
	}

	if next < 0 {
		return l.extraBuf
	} else {
		c, _ := l.peekBuf.At(next)
		return c
	}
}

var insertSemicolonTokenTypes = util.NewSet([]types.TokenType{
	types.TokenTypeIdentifier, types.TokenTypeFloatLiteral,
	types.TokenTypeImaginaryLiteral, types.TokenTypeIntLiteral,
	types.TokenTypeStringLiteral, types.TokenTypeRuneLiteral,
})
var insertSemicolonTokens = util.NewSet([]types.Token{
	types.TokenKeywordBreak, types.TokenKeywordReturn,
	types.TokenKeywordFallthrough, types.TokenIncrement,
	types.TokenDecrement, types.TokenRBrace,
	types.TokenRBracket, types.TokenRParen,
})

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
func (l *Lexer) setupSemicolonInsertNewline(token types.Token) {
	l.canInsertSemicolon =
		insertSemicolonTokens.Contains(token) || insertSemicolonTokenTypes.Contains(token.Type)
}
