package lox

import (
	"errors"
	"golox/expr"
	"golox/tok"
)

type Parser struct {
	tokens  []*tok.Token
	current int
}

func NewParser(tokens []*tok.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() (expr.Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (expr.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (expr.Expr, error) {
	e, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(tok.BangEqual, tok.EqualEqual) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		e = &expr.Binary{Left: e, Operator: op, Right: right}
	}
	return e, nil
}

func (p *Parser) comparison() (expr.Expr, error) {
	e, err := p.term()
	if err != nil {
		return nil, err
	}
	for p.match(tok.Greater, tok.GreaterEqual, tok.Less, tok.LessEqual) {
		op := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		e = &expr.Binary{Left: e, Operator: op, Right: right}
	}
	return e, nil
}

func (p *Parser) term() (expr.Expr, error) {
	e, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(tok.Minus, tok.Plus) {
		op := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		e = &expr.Binary{Left: e, Operator: op, Right: right}
	}
	return e, nil
}

func (p *Parser) factor() (expr.Expr, error) {
	e, err := p.unary()
	if err != nil {
		return nil, err
	}
	for p.match(tok.Slash, tok.Star) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		e = &expr.Binary{Left: e, Operator: op, Right: right}
	}
	return e, nil
}

func (p *Parser) unary() (expr.Expr, error) {
	if p.match(tok.Bang, tok.Minus) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &expr.Unary{Operator: op, Right: right}, nil
	}
	return p.primary()
}

// primary        → NUMBER | STRING | "true" | "false" | "nil"
//
//	| "(" expression ")" ;
func (p *Parser) primary() (expr.Expr, error) {
	if p.match(tok.False) {
		return &expr.Literal{Value: false}, nil
	} else if p.match(tok.True) {
		return &expr.Literal{Value: true}, nil
	} else if p.match(tok.Nil) {
		return &expr.Literal{Value: nil}, nil
	} else if p.match(tok.Number, tok.String) {
		return &expr.Literal{Value: p.previous().Literal}, nil
	} else if p.match(tok.LeftParen) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(tok.RightParen, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}
		return &expr.Grouping{Expression: e}, nil
	}

	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) match(ts ...tok.Type) bool {
	for _, t := range ts {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(t tok.Type, message string) (*tok.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return nil, errors.New(message)
}

func (p *Parser) check(t tok.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() *tok.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == tok.EOF
}

func (p *Parser) peek() *tok.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *tok.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) error(tok *tok.Token, message string) error {
	ParseError(tok, message)
	return errors.New(message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == tok.Semicolon {
			return
		}

		switch p.peek().Type {
		case tok.Class, tok.Fun, tok.Var, tok.For, tok.If, tok.While, tok.Print, tok.Return:
			return
		}

		p.advance()
	}
}
