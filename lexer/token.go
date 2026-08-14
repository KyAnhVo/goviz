package lexer

type Pos struct {
	line   int
	column int
	pos    int
}

var (
	PosSynthetic Pos = Pos{
		line: -1,
	}
	PosEOF Pos = Pos{
		line: -2,
	}
	PosErr Pos = Pos{
		line: -3,
	}
)

// A token type
type TokenType string

const (
	// Lon-language emit tokens

	TokenTypeErr        TokenType = "Err"
	TokenTypeEof        TokenType = "EOF"
	TokenTypeWhitespace TokenType = "Whitespace"

	// Language token

	TokenTypeIdentifier TokenType = "Identifier"
	TokenTypeKeyword    TokenType = "Keyword"
	TokenTypeOperator   TokenType = "Operator"

	// Langauge literals tokens

	TokenTypeStringLiteral    TokenType = "StringLiteral"
	TokenTypeRuneLiteral      TokenType = "RuneLiteral"
	TokenTypeIntLiteral       TokenType = "IntLiteral"
	TokenTypeFloatLiteral     TokenType = "FloatLiteral"
	TokenTypeImaginaryLiteral TokenType = "ImaginaryLiteral"
)

// A token has a type and a value (source code value essentially)
type Token struct {
	Type  TokenType
	Value string
}

// ---------------------- Utils to create tokens faster ------------------------

func TokenIdentifier(s string) Token   { return Token{TokenTypeIdentifier, s} }
func TokenOperator(s string) Token     { return Token{TokenTypeOperator, s} }
func TokenKeyword(s string) Token      { return Token{TokenTypeKeyword, s} }
func TokenStringLit(s string) Token    { return Token{TokenTypeStringLiteral, s} }
func TokenRuneLit(s string) Token      { return Token{TokenTypeRuneLiteral, s} }
func TokenIntLit(s string) Token       { return Token{TokenTypeIntLiteral, s} }
func TokenFloatLit(s string) Token     { return Token{TokenTypeFloatLiteral, s} }
func TokenImaginaryLit(s string) Token { return Token{TokenTypeImaginaryLiteral, s} }

// --------------------- Token default values -------------------------

// Default semantic tokens
var (
	TokenErr = Token{TokenTypeErr, ""}
	TokenEOF = Token{TokenTypeEof, ""}
)

// Keyword tokens
var (
	TokenKeywordBreak       = TokenKeyword("break")
	TokenKeywordDefault     = TokenKeyword("default")
	TokenKeywordFunc        = TokenKeyword("func")
	TokenKeywordInterface   = TokenKeyword("interface")
	TokenKeywordSelect      = TokenKeyword("select")
	TokenKeywordCase        = TokenKeyword("case")
	TokenKeywordDefer       = TokenKeyword("defer")
	TokenKeywordGo          = TokenKeyword("go")
	TokenKeywordMap         = TokenKeyword("map")
	TokenKeywordStruct      = TokenKeyword("struct")
	TokenKeywordChan        = TokenKeyword("chan")
	TokenKeywordElse        = TokenKeyword("else")
	TokenKeywordGoto        = TokenKeyword("goto")
	TokenKeywordPackage     = TokenKeyword("package")
	TokenKeywordSwitch      = TokenKeyword("switch")
	TokenKeywordConst       = TokenKeyword("const")
	TokenKeywordFallthrough = TokenKeyword("fallthrough")
	TokenKeywordIf          = TokenKeyword("if")
	TokenKeywordRange       = TokenKeyword("range")
	TokenKeywordType        = TokenKeyword("type")
	TokenKeywordContinue    = TokenKeyword("continue")
	TokenKeywordFor         = TokenKeyword("for")
	TokenKeywordImport      = TokenKeyword("import")
	TokenKeywordReturn      = TokenKeyword("return")
	TokenKeywordVar         = TokenKeyword("var")
)

// Operator tokens
var (
	TokenAdd          = TokenOperator("+")
	TokenAnd          = TokenOperator("&")
	TokenAddAssign    = TokenOperator("+=")
	TokenAndAssign    = TokenOperator("&=")
	TokenLAnd         = TokenOperator("&&")
	TokenEqual        = TokenOperator("==")
	TokenNotEqual     = TokenOperator("!=")
	TokenLParen       = TokenOperator("(")
	TokenRParen       = TokenOperator(")")
	TokenSub          = TokenOperator("-")
	TokenOr           = TokenOperator("|")
	TokenSubAssign    = TokenOperator("-=")
	TokenOrAssign     = TokenOperator("|=")
	TokenLOr          = TokenOperator("||")
	TokenLt           = TokenOperator("<")
	TokenLte          = TokenOperator("<=")
	TokenLBracket     = TokenOperator("[")
	TokenRBracket     = TokenOperator("]")
	TokenMul          = TokenOperator("*")
	TokenXor          = TokenOperator("^")
	TokenMulAssign    = TokenOperator("*=")
	TokenXorAssign    = TokenOperator("^=")
	TokenChannel      = TokenOperator("<-")
	TokenGt           = TokenOperator(">")
	TokenGte          = TokenOperator(">=")
	TokenLBrace       = TokenOperator("{")
	TokenRBrace       = TokenOperator("}")
	TokenDiv          = TokenOperator("/")
	TokenLShift       = TokenOperator("<<")
	TokenDivAssign    = TokenOperator("/=")
	TokenLShiftAssign = TokenOperator("<<=")
	TokenIncrement    = TokenOperator("++")
	TokenAssignment   = TokenOperator("=")
	TokenInitialize   = TokenOperator(":=")
	TokenComma        = TokenOperator(",")
	TokenSemicolon    = TokenOperator(";")
	TokenMod          = TokenOperator("%")
	TokenRShift       = TokenOperator(">>")
	TokenModAssign    = TokenOperator("%=")
	TokenRShiftAssign = TokenOperator(">>=")
	TokenDecrement    = TokenOperator("--")
	TokenNot          = TokenOperator("!")
	TokenEllipsis     = TokenOperator("...")
	TokenDot          = TokenOperator(".")
	TokenColon        = TokenOperator(":")
	TokenAndNot       = TokenOperator("&^")
	TokenAndNotAssign = TokenOperator("&^=")
	TokenTilde        = TokenOperator("~")
)
