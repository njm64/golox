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
	var s stmt.Stmt
	var err error

	if p.match(tok.Class) {
		s, err = p.classDeclaration()
	} else if p.match(tok.Fun) {
		s, err = p.function("function")
	} else if p.match(tok.Var) {
		s, err = p.varDeclaration()
	} else {
		s, err = p.statement()
	}

	if err != nil {
		p.synchronize()
		// Return a nil statement if parsing fails. Executing a nil statement
		// would crash, but we should never attempt to execute the code
		// because it contains parse errors.
		return nil
	}

	return s
}

func (p *Parser) classDeclaration() (stmt.Stmt, error) {
	name, err := p.consume(tok.Identifier, "Expect class name")
	if err != nil {
		return nil, err
	}

	var superclass *expr.Variable
	if p.match(tok.Less) {
		_, err = p.consume(tok.Identifier, "Expect superclass name")
		if err != nil {
			return nil, err
		}
		superclass = &expr.Variable{Name: p.previous()}
	}

	_, err = p.consume(tok.LeftBrace, "Expect '{' after class name")
	if err != nil {
		return nil, err
	}

	var methods []*stmt.Function
	for !p.check(tok.RightBrace) && !p.isAtEnd() {
		f, err := p.function("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, f)
	}

	_, err = p.consume(tok.RightBrace, "Expect '}' after class body")
	if err != nil {
		return nil, err
	}

	return &stmt.Class{
		Name:       name,
		Superclass: superclass,
		Methods:    methods,
	}, nil
}

func (p *Parser) statement() (stmt.Stmt, error) {
	if p.match(tok.If) {
		return p.ifStatement()
	} else if p.match(tok.While) {
		return p.whileStatement()
	} else if p.match(tok.For) {
		return p.forStatement()
	} else if p.match(tok.Print) {
		return p.printStatement()
	} else if p.match(tok.Return) {
		return p.returnStatement()
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

func (p *Parser) ifStatement() (stmt.Stmt, error) {
	_, err := p.consume(tok.LeftParen, "Expect '(' after 'if'")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tok.RightParen, "Expect ')' after if condition")
	if err != nil {
		return nil, err
	}
	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch stmt.Stmt
	if p.match(tok.Else) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return &stmt.If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch}, nil
}

func (p *Parser) whileStatement() (stmt.Stmt, error) {
	_, err := p.consume(tok.LeftParen, "Expect '(' after 'while'")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tok.RightParen, "Expect ')' after while condition")
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return &stmt.While{Condition: condition, Body: body}, nil
}

func (p *Parser) forStatement() (stmt.Stmt, error) {
	_, err := p.consume(tok.LeftParen, "Expect '(' after 'for'")
	if err != nil {
		return nil, err
	}

	var initializer stmt.Stmt
	if p.match(tok.Semicolon) {
		initializer = nil
	} else if p.match(tok.Var) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition expr.Expr
	if !p.check(tok.Semicolon) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(tok.Semicolon, "Expect ';' after loop condition")
	if err != nil {
		return nil, err
	}

	var increment expr.Expr
	if !p.check(tok.RightParen) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(tok.RightParen, "Expect ')' after for clauses")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &stmt.Block{
			Statements: []stmt.Stmt{
				body,
				&stmt.Expression{Expression: increment},
			},
		}
	}

	if condition == nil {
		condition = &expr.Literal{Value: true}
	}

	body = &stmt.While{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &stmt.Block{
			Statements: []stmt.Stmt{
				initializer,
				body,
			},
		}
	}

	return body, nil
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

func (p *Parser) returnStatement() (stmt.Stmt, error) {
	keyword := p.previous()

	var value expr.Expr
	if !p.check(tok.Semicolon) {
		var err error
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err := p.consume(tok.Semicolon, "Expect ';' after return value")
	if err != nil {
		return nil, err
	}

	return &stmt.Return{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *Parser) function(kind string) (*stmt.Function, error) {
	name, err := p.consume(tok.Identifier, "Expect "+kind+" name")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tok.LeftParen, "Expect '(' after "+kind+" name")
	if err != nil {
		return nil, err
	}
	var params []*tok.Token
	if !p.check(tok.RightParen) {
		for {
			if len(params) >= 255 {
				return nil, p.error(p.peek(), "Can't have more than 255 parameters")
			}

			param, err := p.consume(tok.Identifier, "Expect parameter name")
			if err != nil {
				return nil, err
			}
			params = append(params, param)
			if !p.match(tok.Comma) {
				break
			}
		}
	}
	_, err = p.consume(tok.RightParen, "Expect ')' after parameters")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(tok.LeftBrace, "Expect '{' before "+kind+" body")
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &stmt.Function{
		Name:   name,
		Params: params,
		Body:   body,
	}, nil
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
	e, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(tok.Equal) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		variableExpr, ok := e.(*expr.Variable)
		if ok {
			return &expr.Assign{
				Name:  variableExpr.Name,
				Value: value}, nil
		}

		getExpr, ok := e.(*expr.Get)
		if ok {
			return &expr.Set{
				Object: getExpr.Object,
				Name:   getExpr.Name,
				Value:  value,
			}, nil
		}

		return nil, &Error{Token: equals, Message: "Invalid assignment target"}
	}

	return e, nil
}

func (p *Parser) or() (expr.Expr, error) {
	e, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(tok.Or) {
		op := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		e = &expr.Logical{Left: e, Operator: op, Right: right}
	}
	return e, nil
}

func (p *Parser) and() (expr.Expr, error) {
	e, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.match(tok.And) {
		op := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		e = &expr.Logical{Left: e, Operator: op, Right: right}
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
	for p.match(tok.Slash, tok.Star, tok.Percent) {
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
	return p.call()
}

func (p *Parser) call() (expr.Expr, error) {
	e, err := p.primary()
	if err != nil {
		return nil, err
	}
	for true {
		if p.match(tok.LeftParen) {
			e, err = p.finishCall(e)
			if err != nil {
				return nil, err
			}
		} else if p.match(tok.Dot) {
			name, err := p.consume(tok.Identifier, "Expect property name after '.'")
			if err != nil {
				return nil, err
			}
			e = &expr.Get{
				Object: e,
				Name:   name,
			}
		} else {
			break
		}
	}
	return e, nil
}

func (p *Parser) finishCall(callee expr.Expr) (expr.Expr, error) {
	var arguments []expr.Expr
	if !p.check(tok.RightParen) {
		for {
			if len(arguments) >= 255 {
				_ = p.error(p.peek(), "Can't have more than 255 arguments")
			}
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
			if !p.match(tok.Comma) {
				break
			}
		}
	}

	paren, err := p.consume(tok.RightParen, "Expect ')' after arguments")
	if err != nil {
		return nil, err
	}

	return &expr.Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}, nil
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
	} else if p.match(tok.Super) {
		keyword := p.previous()
		_, err := p.consume(tok.Dot, "Expect '.' after 'super'")
		if err != nil {
			return nil, err
		}
		method, err := p.consume(tok.Identifier, "Expect superclass method name")
		if err != nil {
			return nil, err
		}
		return &expr.Super{Keyword: keyword, Method: method}, nil
	} else if p.match(tok.This) {
		return &expr.This{Keyword: p.previous()}, nil
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
