package parser

type TokenType int

const (
	EOF TokenType = iota
	OPEN_PAREN
	CLOSE_PAREN
	OPERATION
	CONSTANT
)

type Token struct {
	typ   TokenType
	value string
}

func NewToken(typ TokenType, value string) Token {
	return Token{typ, value}
}

func (t Token) String() string {
	return t.value
}
