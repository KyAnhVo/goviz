package lexer

import (
	"testing"

	types "github.com/KyAnhVo/goviz/token"
)

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
		l := createLexer(tc.Src)
		token, _, err := l.getRuneToken()
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
		} else if expectedToken := types.TokenRuneLit(tc.Expected); token != expectedToken {
			t.Errorf(
				"\nTestFloat\n"+
					"\t%s\n"+
					"\t\tExpected: %s\n"+
					"\t\tGot: %s\n",
				tc.Label,
				types.FormatToken(expectedToken),
				types.FormatToken(token),
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

func TestString(t *testing.T) {
	type TestCase struct {
		Label string
		Src   string
		F     func(l *Lexer) (types.Token, types.Pos, error)
		Err   bool
	}

	tcs := []TestCase{
		// happy test
		{
			"Interpreted string",
			"\"Hello\\u1234\\U1BADFACE\\xFA\\017\\nWorld\"",
			func(l *Lexer) (types.Token, types.Pos, error) { return l.getInterpretedStringToken() },
			false,
		},

		// error test
		// TODO: add error tests
	}
	for _, tc := range tcs {
		l := createLexer(tc.Src + " ")
		token, _, err := l.getInterpretedStringToken()
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
		} else if expectedToken := types.TokenStringLit(tc.Src); token != expectedToken {
			t.Errorf(
				"\nTestFloat\n"+
					"\t%s\n"+
					"\t\tExpected: %s\n"+
					"\t\tGot: %s\n",
				tc.Label,
				types.FormatToken(expectedToken),
				types.FormatToken(token),
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
