package parser

import (
	"curryLang/ast"
	"curryLang/lexer"
	"curryLang/token"
	"fmt"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	// error handling
	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfElseExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionExpression)
	p.registerPrefix(token.QUOTE, p.parseStringExpression)

	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseFunctionCall)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {

	var statement ast.Statement

	switch p.curToken.Type {
	case token.LET:
		statement = p.parseLetStatement()
	case token.RETURN:
		statement = p.parseReturnStatement()
	case token.WHILE:
		statement = p.parseWhileStatement()
	case token.IDENT:
		if p.peekTokenIs(token.ASSIGN) {
			statement = p.parseAssignmentStatement()
		} else {
			statement = p.parseExpressionStatement()
		}
	default:
		statement = p.parseExpressionStatement()
	}

	return statement
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{
		Token: p.curToken,
	}

	if !p.peekTokenIs(token.IDENT) {
		p.peekError(token.IDENT)
		return nil
	}

	p.nextToken()

	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	statement.Name = name

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		statement.Value = p.parseExpression(LOWEST)

		_, isIfElse := statement.Value.(*ast.IfElseExpression)
		_, isFunction := statement.Value.(*ast.FunctionExpression)

		if !isIfElse && !isFunction && !p.expectPeek(token.SEMICOLON) {
			return nil
		}
	} else {
		if !p.expectPeek(token.SEMICOLON) {
			return nil
		}
	}

	return statement
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	statement := &ast.AssignmentStatement{}

	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	statement.Name = name

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		statement.Value = p.parseExpression(LOWEST)

		_, isIfElse := statement.Value.(*ast.IfElseExpression)
		_, isFunction := statement.Value.(*ast.FunctionExpression)

		if !isIfElse && !isFunction && !p.expectPeek(token.SEMICOLON) {
			return nil
		}
	} else {
		p.errors = append(p.errors, "Assign operator missing")
		return nil
	}

	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: p.curToken,
	}

	p.nextToken()
	statement.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	statement := &ast.WhileStatement{
		Token: p.curToken,
	}

	p.nextToken()
	statement.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		p.errors = append(p.errors, "Missing { after while condition")
		return nil
	}

	p.nextToken()

	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt == nil {
			return nil
		}

		statement.Body = append(statement.Body, stmt)
	}

	p.nextToken()

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseFunctionCall(left ast.Expression) ast.Expression {
	expression := &ast.FunctionCallExpression{
		Token:        p.curToken,
		FunctionExpr: left,
	}

	expression.Parameters = p.parseCallArguments()

	return expression
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	booleanValue := p.curToken.Literal == "true"
	return &ast.Boolean{Token: p.curToken, Value: booleanValue}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseStringExpression() ast.Expression {
	lit := &ast.StringLiteral{Token: p.curToken}
	p.nextToken()

	val := ""

	for p.curToken.Type != token.QUOTE {
		val += p.curToken.Literal
		p.nextToken()
	}

	lit.Value = val

	return lit
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfElseExpression() ast.Expression {
	lit := &ast.IfElseExpression{Token: p.curToken}
	p.nextToken()

	lit.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	p.nextToken()

	lit.Consequence = []ast.Statement{}

	for p.curToken.Type != token.EOF && p.curToken.Type != token.RBRACE {
		statement := p.parseStatement()
		if statement != nil {
			lit.Consequence = append(lit.Consequence, statement)
		}

		p.nextToken()
	}

	if p.curToken.Type == token.EOF {
		return nil
	}

	p.nextToken()

	if p.curToken.Type == token.ELSE {

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		p.nextToken()

		lit.Alternative = []ast.Statement{}

		for p.curToken.Type != token.EOF && p.curToken.Type != token.RBRACE {
			statement := p.parseStatement()
			if statement != nil {
				lit.Alternative = append(lit.Alternative, statement)
			}

			p.nextToken()
		}

		if p.curToken.Type == token.EOF {
			return nil
		}

		p.nextToken()
	}

	return lit
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	lit := &ast.FunctionExpression{Token: p.curToken}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		funcNameExpr := p.parseIdentifier()
		lit.Name = funcNameExpr.String()
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	lit.Parameters = []ast.Parameter{}

	for p.curToken.Type != token.RPAREN {
		identifier := p.parseIdentifier()

		if identifier == nil {
			p.errors = append(p.errors, "Parameter could not be parsed")
			return nil
		}

		lit.Parameters = append(lit.Parameters, ast.Parameter{Name: identifier.String()})

		p.nextToken()

		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}

	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	p.nextToken()

	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt == nil {
			return nil
		}

		lit.Body = append(lit.Body, stmt)

		p.nextToken()
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return lit
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
