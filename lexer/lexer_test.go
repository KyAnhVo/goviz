package lexer

import "testing"

func TestComment(t *testing.T) {
	type TestFunction struct {
		Run  func(*Lexer)
		Name string
	}
	type TestCase struct {
		Src         string
		ExpectedPos Pos
		F           TestFunction
	}
	inline := TestFunction{
		Run: func(l *Lexer) {
			l.getLineComment()
		},
		Name: "inline",
	}
	inlineAfter := TestFunction{
		Run: func(l *Lexer) {
			l.getLineComment()
			l.getNextChar()
		},
		Name: "inlineAfter",
	}
	block := TestFunction{
		Run: func(l *Lexer) {
			l.getGeneralComment()
		},
		Name: "block",
	}
	blockAfter := TestFunction{
		Run: func(l *Lexer) {
			l.getGeneralComment()
			l.getNextChar()
		},
		Name: "blockAfter",
	}
	testCases := []TestCase{
		{
			Src: "// inline\n",
			ExpectedPos: Pos{
				line:   0,
				column: len("// inline\n") - 1,
				pos:    len("// inline\n") - 1,
			},
			F: inline,
		},
		{
			Src: "// inline after\nx",
			ExpectedPos: Pos{
				line: 1,
				pos:  len("// inline after\nx") - 1,
			},
			F: inlineAfter,
		},
		{
			Src:         "/* block */x",
			ExpectedPos: SyntheticPos,
			F:           block,
		},
		{
			Src: "/* block after */x",
			ExpectedPos: Pos{
				line:   0,
				column: len("/* block after */x") - 1,
				pos:    len("/* block after */x") - 1,
			},
			F: blockAfter,
		},
	}

	for _, testCase := range testCases {
		l := NewLexer([]rune(testCase.Src))
		testCase.F.Run(l)
		_, pos := l.getNextChar()
		if pos != testCase.ExpectedPos {
			t.Errorf(
				"Comment type %s:\nstring:\n%s\nExpected %+v\nGot%+v\n",
				testCase.F.Name,
				testCase.Src,
				testCase.ExpectedPos,
				pos,
			)
		}
	}
}
