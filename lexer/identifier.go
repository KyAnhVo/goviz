/* Golang Identifier grammar
refer to https://go.dev/ref/spec#Identifiers

identifier = letter { letter | unicode_digit } .
*/

/* Golang Keywords
refer to https://go.dev/ref/spec#Keywords

Keywords¶
The following keywords are reserved and may not be used as identifiers.

break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
*/

package lexer

import (
	"fmt"
	"strings"

	types "github.com/KyAnhVo/goviz/token"
)

func (l *Lexer) getIdentifierOrKeyword() (types.Token, types.Pos, error) {
	var token types.Token
	var builder strings.Builder

	c, pos, err := l.getNextChar()
	if err != nil {
		return types.TokenErr, types.PosErr, fmt.Errorf("Identifier: %w", err)
	}
	if !isLetter(c) {
		panic("invariant broken: " + string(c) + "is not a letter")
	}
	builder.WriteRune(c)

	c = l.peekNextChar()
	for isLetter(c) || isUnicodeDigit(c) {
		l.getNextChar()
		builder.WriteRune(c)
		c = l.peekNextChar()
	}
	token.Value = builder.String()

	_, isKeyword := keywordSet[string(token.Value)]
	if isKeyword {
		token.Type = types.TokenTypeKeyword
	} else {
		token.Type = types.TokenTypeIdentifier
	}

	return token, pos, nil
}

var keywordSet map[string]struct{} = map[string]struct{}{
	"break": {}, "case": {}, "chan": {}, "const": {}, "continue": {},
	"default": {}, "defer": {}, "else": {}, "fallthrough": {}, "for": {},
	"func": {}, "go": {}, "goto": {}, "if": {}, "import": {},
	"interface": {}, "map": {}, "package": {}, "range": {}, "return": {},
	"select": {}, "struct": {}, "switch": {}, "type": {}, "var": {},
}
