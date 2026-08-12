package lexer

import (
	"slices"
	"unicode"
)

// ---------------------------- Characters ----------------------------

/* from https://go.dev/ref/spec#Characters
 The following terms are used to denote specific Unicode character categories:

newline        = the Unicode code point U+000A .
unicode_char   = an arbitrary Unicode code point except newline.
unicode_letter = a Unicode code point categorized as "Letter"  .
unicode_digit  = a Unicode code point categorized as "Number, decimal digit"  .

In The Unicode Standard 8.0, Section 4.5 "General Category" defines a set of character categories.
Go treats all characters in any of the Letter categories Lu, Ll, Lt, Lm, or Lo as Unicode letters,
and those in the Number category Nd as Unicode digits.
*/

func isUnicodeLetter(c rune) bool {
	return unicode.In(c, unicode.Lu, unicode.Ll, unicode.Lt, unicode.Lm, unicode.Lo)
}

func isNewline(c rune) bool {
	return c == '\u000A'
}

func isWhitespaceNonNewline(c rune) bool {
	return slices.Contains([]rune{'\u0020', '\u0009', '\u000D'}, c)
}

func isWhitespace(c rune) bool {
	return isNewline(c) || isWhitespaceNonNewline(c)
}

func isUnicodeChar(c rune) bool {
	return !isNewline(c)
}

func isUnicodeDigit(c rune) bool {
	return unicode.Is(unicode.Nd, c)
}

// ---------------------------- Letters and Digits ----------------------------

/* from https://go.dev/ref/spec#Letters_and_digits
The underscore character _ (U+005F) is considered a lowercase letter.

letter        = unicode_letter | "_" .
decimal_digit = "0" … "9" .
binary_digit  = "0" | "1" .
octal_digit   = "0" … "7" .
hex_digit     = "0" … "9" | "A" … "F" | "a" … "f" .
*/

func isLetter(c rune) bool {
	return isUnicodeLetter(c) || c == '_'
}

func isDecimalDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func isBinaryDigit(c rune) bool {
	return c == '0' || c == '1'
}

func isOctalDigit(c rune) bool {
	return '0' <= c && c <= '7'
}

func isHexDigit(c rune) bool {
	return isDecimalDigit(c) || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}
