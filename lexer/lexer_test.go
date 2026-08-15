package lexer

import (
	"strings"
	"testing"

	"github.com/KyAnhVo/goviz/types"
)

// ----------------------- TEST SUITE HELPERS -----------------------

func checkRune(t *testing.T, label string, want rune, got rune) {
	t.Helper()
	if want != got {
		t.Errorf(
			"\nCharacter:\n"+
				"%s\n"+
				"\t%s != %s\n",
			label, types.FormatRune(want), types.FormatRune(got),
		)
	}
}

func checkPos(t *testing.T, label string, want types.Pos, got types.Pos) {
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
// Should be able to test gettypes.Pos, etc. since they are all incoorporated
// ino these 2.
func TestCharRead(t *testing.T) {
	type TestCase struct {
		Rune     rune
		Position types.Pos
	}

	// this builder is just to seperate lines so it look nice.
	var builder strings.Builder
	builder.WriteString("abc\n")
	builder.WriteString("de\n")
	l := NewLexer([]rune(builder.String()))

	tcs := []TestCase{
		{'a', types.Pos{Line: 1, Column: 1, Pos: 1}},
		{'b', types.Pos{Line: 1, Column: 2, Pos: 2}},
		{'c', types.Pos{Line: 1, Column: 3, Pos: 3}},
		{'\n', types.Pos{Line: 1, Column: 4, Pos: 4}},
		{'d', types.Pos{Line: 2, Column: 1, Pos: 5}},
		{'e', types.Pos{Line: 2, Column: 2, Pos: 6}},
		{'\n', types.Pos{Line: 2, Column: 3, Pos: 7}},
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
		Position types.Pos
		Function func(*Lexer)
	}

	tcs := []TestCase{
		{
			Src:      "// Inline Comment\n",
			Label:    "Inline comment",
			Rune:     '\n',
			Position: types.Pos{Line: 1, Column: 18, Pos: 18},
			Function: func(l *Lexer) { l.getLineComment() },
		},
		{
			Src:      "/* Block Comment Inline */",
			Label:    "Block comment inline",
			Rune:     ' ',
			Position: types.PosSynthetic,
			Function: func(l *Lexer) { l.getGeneralComment() },
		},
		{
			Src:      "/* Block Comment \n With Line */",
			Label:    "Block comment with line",
			Rune:     '\n',
			Position: types.PosSynthetic,
			Function: func(l *Lexer) { l.getGeneralComment() },
		},
		{
			Src:      "/* Block Comment Check After */x",
			Label:    "Check after block comment",
			Rune:     'x',
			Position: types.Pos{Line: 1, Column: 32, Pos: 32},
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
