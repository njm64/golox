package lox

import (
	tok2 "golox/lox/tok"
	"strconv"
	"unicode"
)

var identifierMap = map[string]tok2.Type{
	"and":    tok2.And,
	"class":  tok2.Class,
	"else":   tok2.Else,
	"false":  tok2.False,
	"for":    tok2.For,
	"fun":    tok2.Fun,
	"if":     tok2.If,
	"nil":    tok2.Nil,
	"or":     tok2.Or,
	"print":  tok2.Print,
	"return": tok2.Return,
	"super":  tok2.Super,
	"this":   tok2.This,
	"true":   tok2.True,
	"var":    tok2.Var,
	"while":  tok2.While,
}

type Scanner struct {
	source  string
	tokens  []*tok2.Token
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

func (s *Scanner) ScanTokens() []*tok2.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, tok2.NewToken(tok2.EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(tok2.LeftParen)
	case ')':
		s.addToken(tok2.RightParen)
	case '{':
		s.addToken(tok2.LeftBrace)
	case '}':
		s.addToken(tok2.RightBrace)
	case ',':
		s.addToken(tok2.Comma)
	case '.':
		s.addToken(tok2.Dot)
	case '-':
		s.addToken(tok2.Minus)
	case '+':
		s.addToken(tok2.Plus)
	case ';':
		s.addToken(tok2.Semicolon)
	case '*':
		s.addToken(tok2.Star)
	case '!':
		if s.match('=') {
			s.addToken(tok2.BangEqual)
		} else {
			s.addToken(tok2.Bang)
		}
	case '=':
		if s.match('=') {
			s.addToken(tok2.EqualEqual)
		} else {
			s.addToken(tok2.Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(tok2.LessEqual)
		} else {
			s.addToken(tok2.Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(tok2.GreaterEqual)
		} else {
			s.addToken(tok2.Greater)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(tok2.Slash)
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
			ReportScanError(s.line, "Unexpected character.")
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
		tokenType = tok2.Identifier
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

	s.addLiteralToken(tok2.Number, n)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		ReportScanError(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes
	value := s.source[s.start+1 : s.current-1]
	s.addLiteralToken(tok2.String, value)
}

func (s *Scanner) advance() rune {
	c := s.source[s.current]
	s.current++
	return rune(c)
}

func (s *Scanner) addToken(tokenType tok2.Type) {
	s.addLiteralToken(tokenType, nil)
}

func (s *Scanner) addLiteralToken(tokenType tok2.Type, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, tok2.NewToken(tokenType, text, literal, s.line))
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
