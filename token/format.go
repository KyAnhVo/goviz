package token

import (
	"fmt"
	"strings"
)

func FormatRune(c rune) string {
	switch c {
	case '\n':
		return "<newline>"
	case '\t':
		return "<tab>"
	case ' ':
		return "<space>"
	case '\r':
		return "<cr>"
	case '<':
		return "<lessThan>"
	case '>':
		return "<greaterThan>"
	default:
		return string(c)
	}
}

func FormatToken(token Token) string {
	var builder strings.Builder
	builder.WriteString(string(token.Type))
	builder.WriteString("< ")
	for _, c := range []rune(token.Value) {
		builder.WriteString(FormatRune(c))
	}
	builder.WriteString(" >")
	return builder.String()
}

func FormatPos(pos Pos) string {
	return fmt.Sprintf("line: %d, column: %d, pos: %d", pos.Line, pos.Column, pos.Pos)
}
