package lox

import (
	"golox/tok"
	"strconv"
	"unicode"
)

var identifierMap = map[string]tok.Type{
	"and":    tok.And,
	"class":  tok.Class,
	"else":   tok.Else,
	"false":  tok.False,
	"for":    tok.For,
	"fun":    tok.Fun,
	"if":     tok.If,
	"nil":    tok.Nil,
	"or":     tok.Or,
	"print":  tok.Print,
	"return": tok.Return,
	"super":  tok.Super,
	"this":   tok.This,
	"true":   tok.True,
	"var":    tok.Var,
	"while":  tok.While,
}

type Scanner struct {
	source  string
	tokens  []*tok.Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		line:   1,
	}
}

func (s *Scanner) ScanTokens() []*tok.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, tok.NewToken(tok.EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(tok.LeftParen)
	case ')':
		s.addToken(tok.RightParen)
	case '{':
		s.addToken(tok.LeftBrace)
	case '}':
		s.addToken(tok.RightBrace)
	case ',':
		s.addToken(tok.Comma)
	case '.':
		s.addToken(tok.Dot)
	case '-':
		s.addToken(tok.Minus)
	case '+':
		s.addToken(tok.Plus)
	case ';':
		s.addToken(tok.Semicolon)
	case '*':
		s.addToken(tok.Star)
	case '!':
		if s.match('=') {
			s.addToken(tok.BangEqual)
		} else {
			s.addToken(tok.Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(tok.EqualEqual)
		} else {
			s.addToken(tok.Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(tok.LessEqual)
		} else {
			s.addToken(tok.Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(tok.GreaterEqual)
		} else {
			s.addToken(tok.Greater)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(tok.Slash)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace
		break
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			Error(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) identifier() {
	for isAlphanumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tokenType, ok := identifierMap[text]
	if !ok {
		tokenType = tok.Identifier
	}
	s.addToken(tokenType)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()
		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	n, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)

	s.addLiteralToken(tok.Number, n)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		Error(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes
	value := s.source[s.start+1 : s.current-1]
	s.addLiteralToken(tok.String, value)
}

func (s *Scanner) advance() rune {
	c := s.source[s.current]
	s.current++
	return rune(c)
}

func (s *Scanner) addToken(tokenType tok.Type) {
	s.addLiteralToken(tokenType, nil)
}

func (s *Scanner) addLiteralToken(tokenType tok.Type, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, tok.NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.source[s.current]) != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	} else {
		return rune(s.source[s.current])
	}
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return rune(s.source[s.current+1])
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphanumeric(c rune) bool {
	return isDigit(c) || isAlpha(c)
}
