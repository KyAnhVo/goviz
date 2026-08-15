/* Golang string and rune grammar
refer to https://go.dev/ref/spec#Rune_literals, https://go.dev/ref/spec#String_literals

rune_lit               = "'" ( unicode_value | byte_value ) "'" .
string_lit             = raw_string_lit | interpreted_string_lit .
raw_string_lit         = "`" { unicode_char | newline } "`" .
interpreted_string_lit = `"` { unicode_value | byte_value } `"` .

unicode_value    = unicode_char | little_u_value | big_u_value | escaped_char .
byte_value       = octal_byte_value | hex_byte_value .
octal_byte_value = `\` octal_digit octal_digit octal_digit .
hex_byte_value   = `\` "x" hex_digit hex_digit .
little_u_value   = `\` "u" hex_digit hex_digit hex_digit hex_digit .
big_u_value      = `\` "U" hex_digit hex_digit hex_digit hex_digit
                           hex_digit hex_digit hex_digit hex_digit .
escaped_char     = `\` ( "a" | "b" | "f" | "n" | "r" | "t" | "v" | `\` | "'" | `"` ) .
*/

package lexer

import (
	"errors"
	"fmt"
	"strings"
)

type escapeMode uint

const (
	escapeModeOctalByte escapeMode = 1 << iota
	escapeModeHexByte
	escapeModeBigU
	escapeModeSmallU
	escapeModeEscapeChar

	escapeModeByte    = escapeModeOctalByte | escapeModeHexByte
	escapeModeUnicode = escapeModeBigU | escapeModeSmallU | escapeModeEscapeChar
)

func (l *Lexer) escapeChar(literalType string, notAcceptable rune, escapeFlags escapeMode) (string, error) {
	c, _ := l.getNextChar()
	switch c {
	case notAcceptable:
		return "",
			errors.New(literalType + " does not accept escaping " + string(notAcceptable))
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '\'', '"':
		if escapeFlags&escapeModeEscapeChar == 0 {
			return "",
				errors.New("escapeChar: does not accept escaping characters")
		}
		return string(c), nil
	case 'x':
		if escapeFlags&escapeModeHexByte == 0 {
			return "",
				errors.New("escapeChar: does not accept hex byte notation")
		}
		var builder strings.Builder
		builder.WriteRune('x')
		for range 2 {
			c, _ := l.getNextChar()
			if !isHexDigit(c) {
				return "",
					errors.New("escapeChar: hex byte notation requires 2 hex digits")
			}
			builder.WriteRune(c)
		}
		return builder.String(), nil
	case 'u':
		if escapeFlags&escapeModeSmallU == 0 {
			return "",
				errors.New("escapeChar: does not accept small u unicode notation")
		}
		var builder strings.Builder
		builder.WriteRune('u')
		for range 4 {
			c, _ := l.getNextChar()
			if !isHexDigit(c) {
				return "",
					errors.New("escapeChar: small u notation requires 4 hex digits")
			}
			builder.WriteRune(c)
		}
		return builder.String(), nil
	case 'U':
		if escapeFlags&escapeModeBigU == 0 {
			return "",
				errors.New("escapeChar: does not accept big U unicode notation")
		}
		var builder strings.Builder
		builder.WriteRune('U')
		for range 8 {
			c, _ := l.getNextChar()
			if !isHexDigit(c) {
				return "",
					errors.New("escapeChar: big u notation requires 8 hex digits")
			}
			builder.WriteRune(c)
		}
		return builder.String(), nil
	case '0', '1', '2', '3', '4', '5', '6', '7':
		if escapeFlags&escapeModeOctalByte == 0 {
			return "",
				errors.New("escapeChar: does not accept octal byte notation")
		}
		var builder strings.Builder
		builder.WriteRune(c)
		for range 2 {
			c, _ := l.getNextChar()
			if !isOctalDigit(c) {
				return "",
					errors.New("escapeChar: octal byte notation requires 3 octal digits")
			}
			builder.WriteRune(c)
		}
		return builder.String(), nil
	default:
		return "", errors.New("character not escapable")
	}
}

func (l *Lexer) getInterpretedStringToken() (Token, Pos, error) {
	var builder strings.Builder

	c, pos := l.getNextChar()
	if c != '"' {
		return TokenErr, PosErr, errors.New("interpretedString: must start with a double quote")
	}
	builder.WriteRune(c)

	c, _ = l.getNextChar()
	for c != '"' {
		builder.WriteRune(c)
		if isNewline(c) || c == '\u0000' {
			return TokenErr, PosErr, errors.New("interpretedString: end of line without closing string")
		}
		if c == '\\' {
			s, err := l.escapeChar("String", '\'', escapeModeUnicode|escapeModeByte)
			if err != nil {
				return TokenErr, PosErr, fmt.Errorf("interpretedString: %w", err)
			}
			builder.WriteString(s)
		}
		c, _ = l.getNextChar()
	}
	builder.WriteRune(c)

	return TokenStringLit(builder.String()), pos, nil
}

func (l *Lexer) getRawStringToken() (Token, Pos, error) {
	var builder strings.Builder

	c, pos := l.getNextChar()
	if c != '`' {
		return TokenErr, PosErr, errors.New("rawString: must start with a tick")
	}
	builder.WriteRune(c)

	c, _ = l.getNextChar()
	for c != '`' {
		builder.WriteRune(c)
		if c == '\u0000' {
			return TokenErr, PosErr, errors.New("rawString: file ended without string closing")
		}
		if c == '\\' {
			s, err := l.escapeChar("String", '\'', escapeModeUnicode)
			if err != nil {
				return TokenErr, PosErr, fmt.Errorf("rawString: %w", err)
			}
			builder.WriteString(s)
		}
		c, _ = l.getNextChar()
	}
	builder.WriteRune(c)

	return TokenStringLit(builder.String()), pos, nil
}

func (l *Lexer) getRuneToken() (Token, Pos, error) {
	open, pos := l.getNextChar()
	if open != '\'' {
		return TokenErr, PosErr, errors.New("rune: must start with a single quote")
	}

	c, _ := l.getNextChar()
	extra := ""
	var err error
	if c == '\\' {
		extra, err = l.escapeChar("Rune", '"', escapeModeUnicode|escapeModeByte)
		if err != nil {
			return TokenErr, PosErr, fmt.Errorf("rune: %w", err)
		}
	}

	close, _ := l.getNextChar()
	if close != '\'' {
		return TokenErr, PosErr, errors.New("rune: must end with a single quote")
	}

	return TokenRuneLit("'" + string(c) + extra + "'"), pos, nil
}
