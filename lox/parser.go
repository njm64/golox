package lox

import (
	"golox/lox/expr"
	"golox/lox/stmt"
	"golox/lox/tok"
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

func (p *Parser) Parse() []stmt.Stmt {
	var statements []stmt.Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() stmt.Stmt {
	// This method can return a nil statement if parsing fails.
	// Executing a nil statement would crash, but we should never
	// attempt to execute the code because it contains parse errors.
	if p.match(tok.Var) {
		s, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil
		}
		return s
	} else {
		s, err := p.statement()
		if err != nil {
			p.synchronize()
			return nil
		}
		return s
	}
}

func (p *Parser) statement() (stmt.Stmt, error) {
	if p.match(tok.Print) {
		return p.printStatement()
	} else if p.match(tok.LeftBrace) {
		block, err := p.block()
		if err != nil {
			return nil, err
		}
		return &stmt.Block{Statements: block}, nil
	} else {
		return p.expressionStatement()
	}
}

func (p *Parser) printStatement() (stmt.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tok.Semicolon, "Expect ';' after value")
	if err != nil {
		return nil, err
	}
	return &stmt.Print{Expression: value}, nil
}

func (p *Parser) varDeclaration() (stmt.Stmt, error) {
	name, err := p.consume(tok.Identifier, "Expect variable name")
	if err != nil {
		return nil, err
	}

	var initializer expr.Expr
	if p.match(tok.Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(tok.Semicolon, "Expect ';' after variable declaration")
	return &stmt.Var{Name: name, Initializer: initializer}, nil
}

func (p *Parser) expressionStatement() (stmt.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tok.Semicolon, "Expect ';' after value")
	if err != nil {
		return nil, err
	}
	return &stmt.Expression{Expression: value}, nil
}

func (p *Parser) block() ([]stmt.Stmt, error) {
	var statements []stmt.Stmt

	for !p.check(tok.RightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	_, err := p.consume(tok.RightBrace, "Expect '}' after block")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) expression() (expr.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (expr.Expr, error) {
	e, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(tok.Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		variable, ok := e.(*expr.Variable)
		if ok {
			return &expr.Assign{Name: variable.Name, Value: value}, nil
		} else {
			return nil, &Error{Token: equals, Message: "Invalid assignment target"}
		}
	}

	return e, nil
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

func (p *Parser) primary() (expr.Expr, error) {
	if p.match(tok.False) {
		return &expr.Literal{Value: false}, nil
	} else if p.match(tok.True) {
		return &expr.Literal{Value: true}, nil
	} else if p.match(tok.Nil) {
		return &expr.Literal{Value: nil}, nil
	} else if p.match(tok.Number, tok.String) {
		return &expr.Literal{Value: p.previous().Literal}, nil
	} else if p.match(tok.Identifier) {
		return &expr.Variable{Name: p.previous()}, nil
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
	return nil, p.error(p.peek(), message)
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

func (p *Parser) error(tok *tok.Token, message string) *Error {
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
