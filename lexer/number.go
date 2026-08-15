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

/* Golang imaginary grammar
refer to https://go.dev/ref/spec#Imaginary_literals

imaginary_lit = (decimal_digits | int_lit | float_lit) "i" .
*/

package lexer

import (
	"errors"
	"fmt"
	"github.com/KyAnhVo/goviz/types"
	"strings"
	"unicode"
)

func (l *Lexer) getNumericLiteral() (types.Token, types.Pos, error) {
	c1 := l.peekNextChar()
	c2 := l.peekOffset(2)

	if !isDecimalDigit(c1) && !(c1 == '.' && isDecimalDigit(c2)) {
		return types.TokenErr, types.PosErr, errors.New(
			"A numeric token must start with a decimal digit or dot",
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

func (l *Lexer) getDecimalToken() (types.Token, types.Pos, error) {
	// the trick here is, we recognize that, except for the dot start case,
	// for float and non float, we all start with a <decimal_digits>. Then
	// if there is no dot or 'e' after, we safely assume that this is an int, and
	// check first char non 0.

	var builder strings.Builder

	if l.peekNextChar() == '.' {
		// suppose we start with a dot, then we use this rule:
		// decimalFloat := "." decimal_digits [ decimal_exponent ] .
		c, pos := l.getNextChar()
		builder.WriteRune(c)

		digits, _, err := l.getDigits(isDecimalDigit)
		if err != nil {
			return types.TokenErr, types.PosErr, fmt.Errorf("getDecimalToken: leading dot: fractional: %w", err)
		}
		builder.WriteString(digits)

		if c = l.peekNextChar(); c == 'E' || c == 'e' {
			exp, err := l.getExponent('e')
			if err != nil {
				return types.TokenErr, types.PosErr, fmt.Errorf("getDecimalToken: leading dot: exponent: %w", err)
			}
			builder.WriteString(exp)
		}

		if l.peekNextChar() == 'i' {
			builder.WriteRune('i')
			l.getNextChar()
			return types.TokenImaginaryLit(builder.String()), pos, nil
		}
		return types.TokenFloatLit(builder.String()), pos, nil
	}

	// decimal_digits
	digits, pos, err := l.getDigits(isDecimalDigit)
	if err != nil {
		return types.TokenErr, types.PosErr, fmt.Errorf("getDecimalToken: %w", err)
	}
	builder.WriteString(digits)

	// int:   epsilon
	// float: ("." [decimal_digits] [decimal_exponent]) | decimal_exponent
	if c := l.peekNextChar(); c != '.' && c != 'e' && c != 'E' {
		if l.peekNextChar() == 'i' {
			l.getNextChar()
			builder.WriteRune('i')
			return types.TokenImaginaryLit(builder.String()), pos, nil
		}

		// Decimal integer
		if len(digits) > 1 && digits[0] == '0' {
			for _, c := range digits {
				if !isOctalDigit(c) {
					return types.TokenErr, types.PosErr, errors.New("decimal int cannot lead with 0")
				}
			}
			return types.TokenIntLit(digits), pos, nil
		}
		return types.TokenIntLit(builder.String()), pos, nil
	} else {
		// ("." [decimal_digits] [decimal_exponent]) | (decimal_exponent)
		if c == '.' {
			builder.WriteRune(c)
			l.getNextChar()

			if isDecimalDigit(l.peekNextChar()) {
				digits, _, err := l.getDigits(isDecimalDigit)
				if err != nil {
					return types.TokenErr, types.PosErr, fmt.Errorf("getDecimalToken: %w", err)
				}
				builder.WriteString(digits)
			}
			c = l.peekNextChar()
		}
		// For the 1st branch, this acts as [decimal_exponent].
		// For the 2nd branch, this is true via the first branch check
		// is wrong and the outside check.
		if c == 'e' || c == 'E' {
			exponent, err := l.getExponent('e')
			if err != nil {
				return types.TokenErr, types.PosErr, fmt.Errorf("getDecimalToken: %w", err)
			}
			builder.WriteString(exponent)
		}

		if l.peekNextChar() == 'i' {
			l.getNextChar()
			builder.WriteRune('i')
			return types.TokenFloatLit(builder.String()), pos, nil
		}
		return types.TokenFloatLit(builder.String()), pos, nil
	}
}

func (l *Lexer) getBinaryToken() (types.Token, types.Pos, error) {
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
		return types.TokenErr, types.PosErr, fmt.Errorf("getBinaryToken: %w", err)
	}
	builder.WriteString(s)

	if l.peekNextChar() == 'i' {
		l.getNextChar()
		builder.WriteRune('i')
		return types.TokenImaginaryLit(builder.String()), pos, nil
	}
	return types.TokenIntLit(builder.String()), pos, nil
}

func (l *Lexer) getOctalToken() (types.Token, types.Pos, error) {
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
		return types.TokenErr, types.PosErr, fmt.Errorf("getOctalToken: %w", err)
	}
	builder.WriteString(s)

	if l.peekNextChar() == 'i' {
		l.getNextChar()
		builder.WriteRune('i')
		return types.TokenImaginaryLit(builder.String()), pos, nil
	}
	return types.TokenIntLit(builder.String()), pos, nil
}

func (l *Lexer) getHexToken() (types.Token, types.Pos, error) {
	var builder strings.Builder

	c1, pos := l.getNextChar()
	builder.WriteRune(c1)

	c2, _ := l.getNextChar()
	builder.WriteRune(c2)

	integer, fractional, err := l.hexMantissa()
	if err != nil {
		return types.TokenErr, types.PosErr, fmt.Errorf("getHexToken: %w", err)
	}
	builder.WriteString(integer)
	builder.WriteString(fractional)

	// If fractional part is nonempty but there is no p/P,
	// it is flagged as error due to hex float requiring
	// the exponent part.
	if c := l.peekNextChar(); c != 'p' && c != 'P' {
		if l.peekNextChar() == 'i' {
			l.getNextChar()
			builder.WriteRune('i')
			return types.TokenImaginaryLit(builder.String()), pos, nil
		}
		return types.TokenIntLit(builder.String()), pos, nil
	}

	exponent, err := l.getExponent('p')
	if err != nil {
		return types.TokenErr, types.PosErr, fmt.Errorf("getHexToken: exponent part error: %w", err)
	}
	builder.WriteString(exponent)

	if l.peekNextChar() == 'i' {
		l.getNextChar()
		builder.WriteRune('i')
		return types.TokenImaginaryLit(builder.String()), pos, nil
	}
	return types.TokenFloatLit(builder.String()), pos, nil
}

// -------------------------------- utils --------------------------------

// we gather digit greedily via the grammar
// digits := digit { ["_"] digit } .
func (l *Lexer) getDigits(verifyDigit func(rune) bool) (string, types.Pos, error) {
	var builder strings.Builder

	// first digit
	c := l.peekNextChar()
	if !verifyDigit(c) {
		return "", types.PosErr, errors.New("digts error: digits must start with a digit")
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
				return "", types.PosErr, errors.New(
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

// the mantissa part of a hex value must be nonempty, refer to
// start of file for grammar.
//
// the mantissa comprises of 2 parts: the integer and the fractional.
// both can be empty.
//
// the integer, if nonempty, has the grammar
//
//	["_"] hexDigits
//
// the fractional, if nonempty, has the grammar
//
//	"." [hexDigits]
//
// the return value is (int, frac, err).
func (l *Lexer) hexMantissa() (string, string, error) {
	var beforeDot, fromDot strings.Builder

	// Case 1: "." digits
	if c := l.peekNextChar(); c == '.' {
		fromDot.WriteRune(c)
		l.getNextChar()
		fractionalPart, _, err := l.getDigits(isHexDigit)
		if err != nil {
			return "", "", fmt.Errorf("hexMantissa: %w", err)
		}
		fromDot.WriteString(fractionalPart)
		return "", fromDot.String(), nil
	}

	// Case 2: ["_"] digits ["." [digits]]
	// ["_"] digits
	if c := l.peekNextChar(); c == '_' {
		beforeDot.WriteRune(c)
		l.getNextChar()
	}
	integerPart, _, err := l.getDigits(isHexDigit)
	if err != nil {
		return "", "", fmt.Errorf("hexMantissa: %w", err)
	}
	beforeDot.WriteString(integerPart)

	// ["." [digits]]
	if c := l.peekNextChar(); c == '.' {
		fromDot.WriteRune(c)
		l.getNextChar()
		if c = l.peekNextChar(); isHexDigit(c) {
			fractionalPart, _, err := l.getDigits(isHexDigit)
			if err != nil {
				return "", "", fmt.Errorf("hexMantissa: %w", err)
			}
			fromDot.WriteString(fractionalPart)
		}
	}

	return beforeDot.String(), fromDot.String(), nil
}
