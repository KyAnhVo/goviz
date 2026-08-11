package lexer

import (
	"errors"
	"slices"
	"unicode"
)

type Lexer struct {
	src                []rune
	extraBuf           rune
	ptr                int
	canInsertSemicolon bool
}

func New(src []rune) *Lexer {
	return &Lexer{
		src: src,
	}
}

func (l *Lexer) peekNextChar() rune {
	return l.peekNextN(1)
}

func (l *Lexer) getNextChar() rune {
	c := l.peekNextChar()

	// if not at end then push
	if c != '\u0000' {
		if l.extraBuf != '\u0000' {
			l.extraBuf = '\u0000'
		} else {
			l.ptr += 1
		}
	}

	return c
}

func (l *Lexer) peekNextN(offset int) rune {
	if offset == 0 {
		panic("offset must be positive")
	}

	next := l.ptr + offset
	if l.extraBuf != '\u0000' {
		next -= 1
	}
	if next <= l.ptr { // if this is true, l.extraBuffer must be nonempty
		return l.extraBuf
	}
	if len(l.src) > next {
		return l.src[next]
	}
	return ('\u0000')
}

// ---------------------------- Characters ----------------------------

func isUnicodeLetter(c rune) bool {
	return unicode.In(c, unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo)
}

func isNewline(c rune) bool {
	return c == '\u000A'
}

func isWhitespaceNonNewline(c rune) bool {
	return slices.Contains([]rune{'\u0020', '\u0009', '\u000D'}, c)
}

func isWhitespace(c rune) bool {
	return isNewline(c) || isWhitespaceNonNewline(c)
}

func isUnicodeChar(c rune) bool {
	return !isNewline(c)
}

func isUnicodeDigit(c rune) bool {
	return unicode.Is(unicode.Nd, c)
}

// ---------------------------- Letters and Digits ----------------------------

func isLetter(c rune) bool {
	return isUnicodeLetter(c) || c == '\u005F'
}

func isDecimalDigit(c rune) bool {
	return '0' <= c && 'c' <= '9'
}

func isBinaryDigit(c rune) bool {
	return c == '0' || c == '1'
}

func isOctalDigit(c rune) bool {
	return '0' <= c && c <= '7'
}

func isHexDigit(c rune) bool {
	return isDecimalDigit(c) || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// ---------------------------- Lexical element ----------------------------

// Skip over line comment, then insert a semicolon as the next character.
//
// We assume that the first '/' is consumed (essentially always assume the
// first element is consumed)
func (l *Lexer) getLineComment() {
	l.getNextChar()

	c := l.getNextChar()
	for !isNewline(c) {
		if c == '\u0000' {
			l.src = append(l.src, ';')
			return
		}
		c = l.getNextChar()
	}
	l.extraBuf = ';'
}

// Skips over the comment block, and inserts a semicolon if
// the comment is multi-line.
//
// We assume that the first '/' is consumed (essentially always assume the
// first element is consumed)
func (l *Lexer) getGeneralComment() error {
	hasNewLine := false

	c1, c2 := l.getNextChar(), l.peekNextChar()
	for c1 != '*' || c2 != '/' {
		if c1 == '\u0000' {
			return errors.New("comment not terminated")
		}
		hasNewLine = isNewline(c1)
		c1, c2 = l.getNextChar(), l.peekNextChar()
	}

	// inject semicolon if multi line
	if hasNewLine {
		l.extraBuf = ';'
	}

	return nil
}
