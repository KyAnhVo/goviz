package lexer

import "testing"

func TestInteger(t *testing.T) {
	type TestCase struct {
		Label       string
		Src         string
		ExpectedNum Token
		Err         bool
	}

	tcs := []TestCase{
		// decimal cases
		{"decimal no dash", "42 ", TokenIntLit("42"), false},
		{"decimal dash", "4_29_350 ", TokenIntLit("4_29_350"), false},

		// binary cases
		{"binary no dash", "0b01 ", TokenIntLit("0b01"), false},
		{"binary dash", "0B_1011_1010_1101 ", TokenIntLit("0B_1011_1010_1101"), false},

		// octal cases
		{"octal no dash", "0o76770406 ", TokenIntLit("0o76770406"), false},
		{"octal dash", "0O_0123_4567 ", TokenIntLit("0O_0123_4567"), false},
		{"octal start 0 no prefix 'o'", "013467 ", TokenIntLit("013467"), false},

		// hexadecimal cases
		{"hex no dash", "0xBadFace ", TokenIntLit("0xBadFace"), false},
		{"hex dash", "0X_bAD_fACE ", TokenIntLit("0X_bAD_fACE"), false},

		// error cases
		// TODO: implement error case
	}

	for _, tc := range tcs {
		l := NewLexer([]rune(tc.Src))
		token, _, err := l.getNumericLiteral()
		if tc.Err && err != nil {
			continue
		}
		if err != nil {
			t.Errorf(
				"\nTestInteger\n"+
					"\t%s\n"+
					"\t\tError: %s",
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
		} else if l.peekNextChar() != ' ' { // test to ensure correct ptr location after number
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

func TestFloat(t *testing.T) {
	type TestCase struct {
		Label       string
		Src         string
		ExpectedNum Token
		Err         bool
	}

	tcs := []TestCase{
		// decimal cases
		{"decimal integer dot", "42. ", TokenFloatLit("42."), false},
		{"decimal dot fractional", ".4_29_350 ", TokenFloatLit(".4_29_350"), false},
		{"decimal integer dot fractional", "123.456 ", TokenFloatLit("123.456"), false},
		{"decimal integer exponent", "123e+456 ", TokenFloatLit("123e+456"), false},
		{"decimal integer dot exponent", "123.e-456 ", TokenFloatLit("123.e-456"), false},
		{"decimal integer dot fractional exponent", "123.456e789 ", TokenFloatLit("123.456e789"), false},
		{"decimal dot fractional exponent", ".123e-456 ", TokenFloatLit(".123e-456"), false},

		// hexadecimal cases
		{"hex integer ", "0X020P-245 ", TokenFloatLit("0X020P-245"), false},
		{"hex integer dot", "0x_020.p+1 ", TokenFloatLit("0x_020.p+1"), false},
		{"hex integer dot fractional", "0x_Bad.Facep3 ", TokenFloatLit("0x_Bad.Facep3"), false},
		{"hex dot fractional", "0x.Bad_Facep-3 ", TokenFloatLit("0x.Bad_Facep-3"), false},

		// error cases
		// TODO: implement error cases
	}

	for _, tc := range tcs {
		l := NewLexer([]rune(tc.Src))
		token, _, err := l.getNumericLiteral()
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
		} else if token != tc.ExpectedNum {
			t.Errorf(
				"\nTestFloat\n"+
					"\t%s\n"+
					"\t\tExpected: %s\n"+
					"\t\tGot: %s\n",
				tc.Label,
				formatToken(tc.ExpectedNum),
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

func TestImaginary(t *testing.T) {
	type TestCase struct {
		Label       string
		Src         string
		ExpectedNum Token
		Err         bool
	}

	tcs := []TestCase{
		// TODO: add test cases
	}

	for _, tc := range tcs {
		l := NewLexer([]rune(tc.Src))
		token, _, err := l.getNumericLiteral()
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
		} else if token != tc.ExpectedNum {
			t.Errorf(
				"\nTestFloat\n"+
					"\t%s\n"+
					"\t\tExpected: %s\n"+
					"\t\tGot: %s\n",
				tc.Label,
				formatToken(tc.ExpectedNum),
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
