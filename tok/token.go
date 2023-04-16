package tok

import "fmt"

type Token struct {
	Type    Type
	Lexeme  string
	Literal any
	Line    int
}

func NewToken(t Type, lexeme string, literal any, line int) *Token {
	return &Token{
		Type:    t,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t *Token) String() string {
	if t.Literal != nil {
		return fmt.Sprintf("%s %s %v", t.Type, t.Lexeme, t.Literal)
	} else {
		return fmt.Sprintf("%s %s", t.Type, t.Lexeme)
	}
}
