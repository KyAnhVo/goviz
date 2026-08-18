package lexer

import (
	"errors"
	"fmt"
)

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

	// the reason for peek-check-get is that we save the newline for later.
	c := l.peekNextChar()
	for !isNewline(c) {
		l.getNextChar()
		if c == '\u0000' {
			l.extraBuf = '\n'
			return
		}
		c = l.peekNextChar()
	}

	// invariant
	if l.extraBuf != '\u0000' {
		panic("null extrabuf invariant not satisfied")
	}
}

// Skips over the comment block, and inserts a semicolon if
// the comment is multi-line.
func (l *Lexer) getGeneralComment() error {
	hasNewLine := false

	// skip over the ['/', '*']
	l.getNextChar()
	l.getNextChar()

	c1, c2 := l.peekNextChar(), l.peekOffset(2)
	for c1 != '*' || c2 != '/' {
		if c1 == '\u0000' {
			return errors.New("comment not terminated")
		}
		hasNewLine = hasNewLine || isNewline(c1)

		var err error
		c1, _, err = l.getNextChar()
		if err != nil {
			return fmt.Errorf(
				"GetGeneralComment: line %d, column %d, position %d: %w",
				l.pos.Line, l.pos.Column, l.pos.Pos, err,
			)
		}
		c2 = l.peekNextChar()
	}
	l.getNextChar()

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
