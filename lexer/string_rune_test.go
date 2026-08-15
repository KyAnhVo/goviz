package lexer

import "testing"

func TestRune(t *testing.T) {
	type TestCase struct {
		Label    string
		Src      string
		Expected string
		Err      bool
	}

	tcs := []TestCase{
		// happy test
		{"Rune default", "'a' ", "'a'", false},
		{"Rune escape", "'\\n' ", "'\\n'", false},
		{"Rune hex byte", "'\\xBA' ", "'\\xBA'", false},
		{"Rune octal byte", "'\\107' ", "'\\107'", false},
		{"Rune small u unicode", "'\\uFACE' ", "'\\uFACE'", false},
		{"Rune big U unicode", "'\\U0123abcd' ", "'\\U0123abcd'", false},

		// error test
		// TODO: add error tests
	}
	for _, tc := range tcs {
		l := NewLexer([]rune(tc.Src))
		token, _, err := l.getRuneLiteralToken()
		if tc.Err && err != nil {
			continue
		}
		if err != nil {
			t.Errorf(
				"\nTestFloat\n"+
					"\t%s\n"+
					"\t\tError: %s",
				tc.Label,
				err.Error(),
			)
		} else if expectedToken := TokenRuneLit(tc.Expected); token != expectedToken {
			t.Errorf(
				"\nTestFloat\n"+
					"\t%s\n"+
					"\t\tExpected: %s\n"+
					"\t\tGot: %s\n",
				tc.Label,
				formatToken(expectedToken),
				formatToken(token),
			)
		} else if l.peekNextChar() != ' ' { // test to ensure correct ptr location after number
			t.Errorf(
				"\nTestFloat\n"+
					"\t%s\n"+
					"\t\terror: did not advance over token, final char: %c\n",
				tc.Label,
				l.peekNextChar(),
			)
		}
	}
}
