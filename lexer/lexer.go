package lexer

import (
	"errors"
)

type Lexer struct {
	src                []rune
	extraBuf           rune
	ptr                int
	canInsertSemicolon bool
	pos                Pos
}

func NewLexer(src []rune) *Lexer {
	return &Lexer{
		src: src,
	}
}

// ---------------------------- Utility ----------------------------

func (l *Lexer) peekNextChar() rune {
	return l.peekOffset(1)
}

func (l *Lexer) getNextChar() (rune, Pos) {
	c := l.peekNextChar()

	// if not at end then push
	var currentPos Pos
	if c != '\u0000' {
		if l.extraBuf != '\u0000' {
			currentPos = SyntheticPos
			l.extraBuf = '\u0000'
		} else {
			currentPos = l.adjustPos(c)
			l.ptr += 1
		}
	} else {
		currentPos = EofPos
	}

	return c, currentPos
}

func (l *Lexer) adjustPos(c rune) Pos {
	oldColCount := l.pos.column

	if isNewline(c) {
		l.pos.column = -1
		l.pos.line += 1
	} else {
		l.pos.column += 1
	}
	l.pos.pos += 1

	semanticPos := l.pos
	if semanticPos.column == -1 {
		semanticPos.column = oldColCount + 1
		semanticPos.line -= 1
	}
	return semanticPos
}

func (l *Lexer) peekOffset(offset int) rune {
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
		l.getNextChar()
		if c1 == '\u0000' {
			return errors.New("comment not terminated")
		}
		hasNewLine = isNewline(c1)
		c1, c2 = l.peekNextChar(), l.peekOffset(2)
	}
	l.getNextChar()
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
