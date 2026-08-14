package lexer

import "testing"

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
