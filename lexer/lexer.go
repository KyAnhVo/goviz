package lexer

import (
	"errors"
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

// ---------------------------- Utility ----------------------------

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

// ---------------------------- Comments ----------------------------
/* From https://go.dev/ref/spec#Comments
Comments serve as program documentation. There are two forms:
	Line comments start with the character sequence // and stop at the end of the line.
  General comments start with the character sequence /* and stop with the first subsequent character sequence *\/.
A general comment containing no newlines acts like a space. Any other comment acts like a newline.
*/

// Skip over line comment, then insert a semicolon as the next character.
func (l *Lexer) getLineComment() {
	l.getNextChar()
	l.getNextChar()

	c := l.getNextChar()
	for !isNewline(c) {
		if c == '\u0000' {
			l.extraBuf = '\n'
			return
		}
		c = l.getNextChar()
	}
	// invariant
	if l.extraBuf != '\u0000' {
		panic("null extrabuf invariant not satisfied")
	}

	// the idea of this newline is that we have consumed the newline
	// without applying any newline logic here (semicolon injection,
	// line count, etc.). Thus we add one newline so that our lexer
	// lexes the next newline content.
	l.extraBuf = '\n'
}

// Skips over the comment block, and inserts a semicolon if
// the comment is multi-line.
func (l *Lexer) getGeneralComment() error {
	hasNewLine := false

	// skip over the ['/', '*']
	l.getNextChar()
	l.getNextChar()

	c1, c2 := l.getNextChar(), l.peekNextChar()
	for c1 != '*' || c2 != '/' {
		if c1 == '\u0000' {
			return errors.New("comment not terminated")
		}
		hasNewLine = isNewline(c1)
		c1, c2 = l.getNextChar(), l.peekNextChar()
	}

	// inject corresponding whitespace char.
	// We note that there is no way that l.extraBuf is occcupied here due to the first 2 getNextChar's.
	if l.extraBuf != '\u0000' {
		panic("null extrabuf invariant not satisfied")
	}

	// refer to the spec for how block comment works
	if hasNewLine {
		l.extraBuf = '\n'
	} else {
		l.extraBuf = ' '
	}

	return nil
}
