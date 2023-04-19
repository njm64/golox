package lox

import (
	"golox/lox/expr"
	"golox/lox/stmt"
	tok2 "golox/lox/tok"
)

type Parser struct {
	tokens  []*tok2.Token
	current int
}

func NewParser(tokens []*tok2.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() ([]stmt.Stmt, error) {
	var statements []stmt.Stmt
	for !p.isAtEnd() {
		s, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, s)
	}
	return statements, nil
}

func (p *Parser) expression() (expr.Expr, error) {
	return p.equality()
}

func (p *Parser) statement() (stmt.Stmt, error) {
	if p.match(tok2.Print) {
		return p.printStatement()
	} else {
		return p.expressionStatement()
	}
}

func (p *Parser) printStatement() (stmt.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tok2.Semicolon, "Expect ';' after value")
	if err != nil {
		return nil, err
	}
	return &stmt.Print{Expression: value}, nil
}

func (p *Parser) expressionStatement() (stmt.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tok2.Semicolon, "Expect ';' after value")
	if err != nil {
		return nil, err
	}
	return &stmt.Expression{Expression: value}, nil
}

func (p *Parser) equality() (expr.Expr, error) {
	e, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(tok2.BangEqual, tok2.EqualEqual) {
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
	for p.match(tok2.Greater, tok2.GreaterEqual, tok2.Less, tok2.LessEqual) {
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
	for p.match(tok2.Minus, tok2.Plus) {
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
	for p.match(tok2.Slash, tok2.Star) {
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
	if p.match(tok2.Bang, tok2.Minus) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &expr.Unary{Operator: op, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (expr.Expr, error) {
	if p.match(tok2.False) {
		return &expr.Literal{Value: false}, nil
	} else if p.match(tok2.True) {
		return &expr.Literal{Value: true}, nil
	} else if p.match(tok2.Nil) {
		return &expr.Literal{Value: nil}, nil
	} else if p.match(tok2.Number, tok2.String) {
		return &expr.Literal{Value: p.previous().Literal}, nil
	} else if p.match(tok2.LeftParen) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(tok2.RightParen, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}
		return &expr.Grouping{Expression: e}, nil
	}

	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) match(ts ...tok2.Type) bool {
	for _, t := range ts {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(t tok2.Type, message string) (*tok2.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return nil, p.error(p.peek(), message)
}

func (p *Parser) check(t tok2.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() *tok2.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == tok2.EOF
}

func (p *Parser) peek() *tok2.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *tok2.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) error(tok *tok2.Token, message string) *Error {
	err := &Error{
		Token:   tok,
		Message: message,
	}
	ReportParseError(err)
	return err
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == tok2.Semicolon {
			return
		}

		switch p.peek().Type {
		case tok2.Class, tok2.Fun, tok2.Var, tok2.For, tok2.If, tok2.While, tok2.Print, tok2.Return:
			return
		}

		p.advance()
	}
}
