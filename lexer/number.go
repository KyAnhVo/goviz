/* Golang integer grammar
refer to https://go.dev/ref/spec#Integer_literals

int_lit        = decimal_lit | binary_lit | octal_lit | hex_lit .
decimal_lit    = "0" | ( "1" … "9" ) [ [ "_" ] decimal_digits ] .
binary_lit     = "0" ( "b" | "B" ) [ "_" ] binary_digits .
octal_lit      = "0" [ "o" | "O" ] [ "_" ] octal_digits .
hex_lit        = "0" ( "x" | "X" ) [ "_" ] hex_digits .
decimal_digits = decimal_digit { [ "_" ] decimal_digit } .
binary_digits  = binary_digit { [ "_" ] binary_digit } .
octal_digits   = octal_digit { [ "_" ] octal_digit } .
hex_digits     = hex_digit { [ "_" ] hex_digit } .
*/

/* Golang float grammar
refer to https://go.dev/ref/spec#Floating-point_literals

float_lit         = decimal_float_lit | hex_float_lit .
decimal_float_lit = decimal_digits "." [ decimal_digits ] [ decimal_exponent ] |
                    decimal_digits decimal_exponent |
                    "." decimal_digits [ decimal_exponent ] .
decimal_exponent  = ( "e" | "E" ) [ "+" | "-" ] decimal_digits .
hex_float_lit     = "0" ( "x" | "X" ) hex_mantissa hex_exponent .
hex_mantissa      = [ "_" ] hex_digits "." [ hex_digits ] |
                    [ "_" ] hex_digits |
                    "." hex_digits .
hex_exponent      = ( "p" | "P" ) [ "+" | "-" ] decimal_digits .
*/

package lexer

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type base int

const (
	Binary base = iota
	Octal
	Decimal
	Hex
)

type numberType int

const (
	Float numberType = iota
	Int
)

func (l *Lexer) getNumericToken() (Token, Pos, error) {
	c1 := l.peekNextChar()
	c2 := l.peekOffset(2)

	if !isDecimalDigit(c1) {
		return TokenErr, PosErr, errors.New(
			"A numeric token must start with a decimal digit",
		)
	}

	if c1 == '0' && (c2 == 'b' || c2 == 'B') {
		return l.getBinaryToken()
	} else if c1 == '0' && (c2 == 'o' || c2 == 'O') {
		return l.getOctalToken()
	} else if c1 == '0' && (c2 == 'x' || c2 == 'X') {
		return l.getHexToken()
	} else {
		return l.getDecimalToken()
	}
}

func (l *Lexer) getDecimalToken() (Token, Pos, error) {
	// the trick here is, we recognize that, except for the dot start case,
	// for float and non float, we all start with a <decimal_digits>. Then
	// if there is no dot after, we safely assume that this is an int, and
	// check first char non 0.

	var builder strings.Builder

	if l.peekNextChar() == '.' {
		// suppose we start with a dot, then we use this rule:
		// decimalFloat := "." decimal_digits [ decimal_exponent ] .
		c, pos := l.getNextChar()
		builder.WriteRune(c)

		digits, _, err := l.getDigits(isDecimalDigit)
		if err != nil {
			return TokenErr, PosErr, fmt.Errorf("getDecimalToken: %w", err)
		}
		builder.WriteString(digits)

		if c = l.peekNextChar(); c == 'E' || c == 'e' {
			exp, err := l.getExponent('e')
			if err != nil {
				return TokenErr, PosErr, fmt.Errorf("getDecimalToken: %w", err)
			}
			builder.WriteString(exp)
		}
		return TokenFloatLit(builder.String()), pos, nil
	}

	// decimal_digits
	digits, pos, err := l.getDigits(isDecimalDigit)
	if err != nil {
		return TokenErr, PosErr, fmt.Errorf("getDecimalToken: %w", err)
	}
	builder.WriteString(digits)

	// int:   epsilon
	// float: ("." [decimal_digits] [decimal_exponent]) | decimal_exponent
	if c := l.peekNextChar(); c != '.' && c != 'e' && c != 'E' {
		// Decimal integer
		if len(digits) > 1 && digits[0] == '0' {
			return TokenErr, PosErr, errors.New(
				"getDecimalToken: non-zero ints cannot start with a 0",
			)
		}
		return TokenIntLit(builder.String()), pos, nil
	} else {
		// ("." [decimal_digits] [decimal_exponent]) | (decimal_exponent)

		if c == '.' {
			builder.WriteRune(c)
			l.getNextChar()

			if isDecimalDigit(l.peekNextChar()) {
				digits, _, err := l.getDigits(isDecimalDigit)
				if err != nil {
					return TokenErr, PosErr, fmt.Errorf("getDecimalToken: %w", err)
				}
				builder.WriteString(digits)
			}
		}
		// For the 1st branch, this acts as [decimal_exponent].
		// For the 2nd branch, this is true via the first branch check
		// is wrong and the outside check.
		if c == 'e' || c == 'E' {
			exponent, err := l.getExponent('e')
			if err != nil {
				return TokenErr, PosErr, fmt.Errorf("getDecimalToken: %w", err)
			}
			builder.WriteString(exponent)
		}
		return TokenFloatLit(builder.String()), pos, nil
	}
}

func (l *Lexer) getBinaryToken() (Token, Pos, error) {
	var builder strings.Builder

	c1, pos := l.getNextChar()
	builder.WriteRune(c1)

	c2, _ := l.getNextChar()
	builder.WriteRune(c2)

	if c := l.peekNextChar(); c == '_' {
		builder.WriteRune(c)
		l.getNextChar()
	}

	s, _, err := l.getDigits(isBinaryDigit)
	if err != nil {
		return TokenErr, PosErr, fmt.Errorf("getBinaryToken: %w", err)
	}
	builder.WriteString(s)

	return TokenIntLit(builder.String()), pos, nil
}

func (l *Lexer) getOctalToken() (Token, Pos, error) {
	var builder strings.Builder

	c1, pos := l.getNextChar()
	builder.WriteRune(c1)

	c2, _ := l.getNextChar()
	builder.WriteRune(c2)

	if c := l.peekNextChar(); c == '_' {
		builder.WriteRune(c)
		l.getNextChar()
	}

	s, _, err := l.getDigits(isOctalDigit)
	if err != nil {
		return TokenErr, PosErr, fmt.Errorf("getOctalToken: %w", err)
	}
	builder.WriteString(s)

	return TokenIntLit(builder.String()), pos, nil
}

func (l *Lexer) getHexToken() (Token, Pos, error) {
	panic("")
}

// -------------------------------- utils --------------------------------

// we gather digit greedily via the grammar
// digits := digit { ["_"] digit } .
func (l *Lexer) getDigits(verifyDigit func(rune) bool) (string, Pos, error) {
	var builder strings.Builder

	// first digit
	c := l.peekNextChar()
	if !verifyDigit(c) {
		return "", PosErr, errors.New("digts error: digits must start with a digit")
	}
	builder.WriteRune(c)
	_, pos := l.getNextChar()

	// { ["_"] digit }
	c = l.peekNextChar()
	for c == '_' || verifyDigit(c) {
		l.getNextChar()
		if c == '_' {
			builder.WriteRune(c)
			c = l.peekNextChar()
			if !verifyDigit(c) {
				return "", PosErr, errors.New(
					"digits error: must have digit after underline.",
				)
			}
			l.getNextChar()
		}
		builder.WriteRune(c)
		c = l.peekNextChar()
	}

	return builder.String(), pos, nil
}

// we use the following grammars:
//
//	decimal_exponent  = ("e" | "E") ["+" | "-"] decimal_digits .
//
// and
//
//	hex_exponent      = ("p" | "P") ["+" | "-"] decimal_digits .
//
// which generalizes to
//
//	exponent = (p.lowercase() | p.uppercase()) ["+" | "-"] decimal_digits .
//
// where p is prefix.
func (l *Lexer) getExponent(prefix rune) (string, error) {
	var c rune
	var res string

	c = l.peekNextChar()
	if unicode.ToLower(c) != prefix {
		return "", errors.New(
			"exponent: exponent component must start with prefix " + string(prefix),
		)
	}
	res += string(c)
	l.getNextChar()

	c = l.peekNextChar()
	if c == '+' || c == '-' {
		res += string(c)
		l.getNextChar()
	}

	digits, _, err := l.getDigits(isDecimalDigit)
	if err != nil {
		return "", fmt.Errorf("exponent: %w", err)
	}
	res += digits

	return res, nil
}
