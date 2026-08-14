package lexer

import "testing"

func TestInteger(t *testing.T) {
	type TestCase struct {
		Label       string
		Src         string
		ExpectedNum Token
	}

	tcs := []TestCase{
		{"decimal no dash", "42 ", TokenIntLit("42")},
		{"decimal dash", "4_29_350 ", TokenIntLit("4_29_350")},
		{"binary no dash", "0b01 ", TokenIntLit("0b01")},
		{"binary dash", "0B_1011_1010_1101 ", TokenIntLit("0B_1011_1010_1101")},
		{"octal no dash", "0o76770406 ", TokenIntLit("0o76770406")},
		{"octal dash", "0O_0123_4567 ", TokenIntLit("0O_0123_4567")},
	}

	for _, tc := range tcs {
		l := NewLexer([]rune(tc.Src))
		token, _, err := l.getNumericToken()
		if err != nil {
			t.Errorf(
				"\nTestInteger\n"+
					"\t%s\n"+
					"\t\terror: %s\n",
				tc.Label,
				err.Error(),
			)
		} else if token != tc.ExpectedNum {
			t.Errorf(
				"\nTestInteger\n"+
					"\t%s\n"+
					"\t\tExpected: %+v\n"+
					"\t\tGot: %+v\n",
				tc.Label,
				tc.ExpectedNum,
				token,
			)
		} else if l.peekNextChar() != ' ' {
			t.Errorf(
				"\nTestInteger\n"+
					"\t%s\n"+
					"\t\terror: did not advance over token, final char: %c\n",
				tc.Label,
				l.peekNextChar(),
			)
		}
	}
}

func TestDigits(t *testing.T) {
	type TestCase struct {
		Src         string
		Num         string
		NextChar    rune
		VerifyDigit func(rune) bool
	}

	tcs := []TestCase{
		{"911", "911", '\u0000', isDecimalDigit},
		{"0010.425", "0010", '.', isBinaryDigit},
		{"36_69_420_abcd", "36_69_420_abcd", '\u0000', isHexDigit},
		{"24_67.3579", "24_67", '.', isOctalDigit},
	}

	for _, tc := range tcs {
		l := NewLexer([]rune(tc.Src))
		s, _, err := l.getDigits(tc.VerifyDigit)
		if err != nil {
			t.Errorf("TestDigits: Error:\n%s\n", err.Error())
		} else if s != tc.Num {
			t.Errorf(
				"TestDigits: Digits not the same:\n%s != %s\n",
				tc.Num, s,
			)
		} else if l.peekNextChar() != tc.NextChar {
			t.Errorf(
				"TestDigits: Next char wrong:\n%c != %c\n",
				tc.NextChar, l.peekNextChar(),
			)
		}
	}
}
