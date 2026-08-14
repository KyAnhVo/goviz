package lexer

import (
	"strings"
	"testing"
)

// ----------------------- TEST SUITE HELPERS -----------------------

func formatRune(c rune) string {
	switch c {
	case '\n':
		return "<newline>"
	case '\t':
		return "<tab>"
	case ' ':
		return "<space>"
	case '\r':
		return "<cr>"
	case '<':
		return "<lessThan>"
	case '>':
		return "<greaterThan>"
	default:
		return string(c)
	}
}

func formatToken(token Token) string {
	var builder strings.Builder
	builder.WriteString(string(token.Type))
	builder.WriteRune('<')
	for _, c := range []rune(token.Value) {
		builder.WriteString(formatRune(c))
	}
	builder.WriteRune('>')
	return builder.String()
}

func checkRune(t *testing.T, label string, want rune, got rune) {
	t.Helper()
	if want != got {
		t.Errorf(
			"\nCharacter:\n"+
				"%s\n"+
				"\t%s != %s\n",
			label, formatRune(want), formatRune(got),
		)
	}
}

func checkPos(t *testing.T, label string, want Pos, got Pos) {
	t.Helper()
	if want != got {
		t.Errorf(
			"\nPosition:\n"+
				"%s\n"+
				"\t%+v != %+v\n",
			label, want, got,
		)
	}
}

// ----------------------- TEST SUITE -----------------------

// Test the primitive functions like getNextChar, getCurrentChar.
// Should be able to test getPos, etc. since they are all incoorporated
// ino these 2.
func TestCharRead(t *testing.T) {
	type TestCase struct {
		Rune     rune
		Position Pos
	}

	// this builder is just to seperate lines so it look nice.
	var builder strings.Builder
	builder.WriteString("abc\n")
	builder.WriteString("de\n")
	l := NewLexer([]rune(builder.String()))

	tcs := []TestCase{
		{'a', Pos{line: 1, column: 1, pos: 1}},
		{'b', Pos{line: 1, column: 2, pos: 2}},
		{'c', Pos{line: 1, column: 3, pos: 3}},
		{'\n', Pos{line: 1, column: 4, pos: 4}},
		{'d', Pos{line: 2, column: 1, pos: 5}},
		{'e', Pos{line: 2, column: 2, pos: 6}},
		{'\n', Pos{line: 2, column: 3, pos: 7}},
	}

	for _, tc := range tcs {
		cNext, pNext := l.getNextChar()
		cCurr, pCurr := l.getCurrentChar()
		checkRune(t, "getNext then getCurr", cNext, cCurr)
		checkPos(t, "getNext then getCurr", pNext, pCurr)
		checkRune(t, "correctness", cNext, tc.Rune)
		checkPos(t, "correctness", pNext, pCurr)
	}
}

// Test comments. Involving Inline comment, Block comment
// behaviours when used as newline vs as space, check the
// next character after the line
func TestComment(t *testing.T) {
	// A test case is accepting the source,
	// run the Function, then check the rune and
	// position right after the Function is complete
	type TestCase struct {
		Src      string
		Label    string
		Rune     rune
		Position Pos
		Function func(*Lexer)
	}

	tcs := []TestCase{
		{
			Src:      "// Inline Comment\n",
			Label:    "Inline comment",
			Rune:     '\n',
			Position: Pos{line: 1, column: 18, pos: 18},
			Function: func(l *Lexer) { l.getLineComment() },
		},
		{
			Src:      "/* Block Comment Inline */",
			Label:    "Block comment inline",
			Rune:     ' ',
			Position: PosSynthetic,
			Function: func(l *Lexer) { l.getGeneralComment() },
		},
		{
			Src:      "/* Block Comment \n With Line */",
			Label:    "Block comment with line",
			Rune:     '\n',
			Position: PosSynthetic,
			Function: func(l *Lexer) { l.getGeneralComment() },
		},
		{
			Src:      "/* Block Comment Check After */x",
			Label:    "Check after block comment",
			Rune:     'x',
			Position: Pos{line: 1, column: 32, pos: 32},
			Function: func(l *Lexer) {
				l.getGeneralComment()
				l.getNextChar()
			},
		},
	}

	for _, tc := range tcs {
		l := NewLexer([]rune(tc.Src))
		tc.Function(l)
		c, pos := l.getNextChar()
		checkRune(t, tc.Label, tc.Rune, c)
		checkPos(t, tc.Label, tc.Position, pos)
	}
}
