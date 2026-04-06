package parser

import (
	"fmt"
	"strconv"

	"github.com/michaelzhan1/go-interpreter/ast"
	"github.com/michaelzhan1/go-interpreter/lexer"
	"github.com/michaelzhan1/go-interpreter/token"
)

// operator precedence, favoring higher values
const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
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
	token.LBRACKET: INDEX,
}

type (
	// prefixParseFn is a function that parses prefix expressions
	prefixParseFn func() ast.Expression

	// infixParseFn is a function that parses infix expressions
	infixParseFn func(ast.Expression) ast.Expression
)

// Parser is the parser for the monkey language
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New returns a new Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	// prime curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// registerPrefix adds a prefixParseFn to the parser under a given token type
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix adds a infixParseFn to the parser under a given token type
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Errors returns the errors in the parser
func (p *Parser) Errors() []string {
	return p.errors
}

// peekError adds an error associated with an unexpected token type while peeking
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// noPrefixParseFnError adds an error associated with parsing an expression without the proper parse function
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// nextToken moves the token forward
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// curPrecedence returns the precedence of the current token
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

// peekPrecedence returns the precedence of the peek token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// ParseProgram parses the tokens and constructs a Program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parses a statement starting from a token
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement parses a let statement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeekAndAdvance(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken}

	if !p.expectPeekAndAdvance(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses a return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement parses an expression statement (unassigned expression, so it's its own "statement")
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// at this point, the expression statement should be fully parsed
	// semicolons for expression statements are optional (so just typing "5+5" rather than "5+5;" is ok), so check just in case
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseBlockStatement parses a block statement
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) { // only hits EOF if program is written incorrectly
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseExpression parses the expression with the current precedence given the parser's current state and returns an expression.
// It acts upon the current state of the parser.
// The precedence is the _previous_ token's precedence EXCEPT for prefixes, which effectively combines both a prefix
// and its operand into one prefix "blob"
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// increasing precedence means there's still more work to do (guarantees an operator)
	// decreasing precedence is covered by infix() and may or may not be another operator
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

// parseIdentifier parses an identifier into an ast.Identifier expression
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken}
}

// parseIntegerLiteral parses an integer into an ast.IntegerLiteral expression
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

// parseBooleanLiteral parses a boolean into an ast.BooleanLiteral expression
func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// parseStringLiteral parses a string into an ast.StringLiteral expression
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseArrayLiteral parses an array into an ast.ArrayLiteral expression
func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: p.curToken}
	arr.Elements = p.parseExpressionList(token.RBRACKET)
	return arr
}

// parseIndexExpression parses an array index into an ast.IndexExpression
func (p *Parser) parseIndexExpression(arr ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Arr: arr}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeekAndAdvance(token.RBRACKET) {
		return nil
	}

	return exp
}

// parseFunctionLiteral parses a function into an ast.FunctionLiteral expression
func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeekAndAdvance(token.LPAREN) {
		return nil
	}

	fn.Parameters = p.parseFunctionParameters()

	if !p.expectPeekAndAdvance(token.LBRACE) {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

// parseFunctionParameters parses a function's parameters. It assumes that the curToken is on the opening token.LPAREN.
// It is conservative in checking cases to better catch syntax errors.
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	params := []*ast.Identifier{}

	// 0 param case
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()
	params = append(params, &ast.Identifier{
		Token: p.curToken,
	})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		params = append(params, &ast.Identifier{
			Token: p.curToken,
		})
	}

	// check for closing paren
	if !p.expectPeekAndAdvance(token.RPAREN) {
		return nil
	}

	return params
}

// parseCallExpression parses a function call into an ast.CallExpression
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

// parseCallExpressionArguments parses a function call's arguments. It expects the current token to be the opening paren of the call.
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(end) {
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

	if !p.expectPeekAndAdvance(end) {
		return nil
	}

	return args
}

// parseIfExpression parses an if statement with an optional else
func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeekAndAdvance(token.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeekAndAdvance(token.RPAREN) {
		return nil
	}
	if !p.expectPeekAndAdvance(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeekAndAdvance(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

// parsePrefixExpression parses a prefix expression such as negative or not
func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	exp.Right = p.parseExpression(PREFIX)

	return exp
}

// parseInfixExpression parses an infix expression
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)

	return exp
}

// parseGroupedExpression parses a nested expression that uses parens
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeekAndAdvance(token.RPAREN) {
		return nil
	}

	return exp
}

// curTokenIs checks that the parser's curToken is of a given type
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs checks that the parser's peekToken is of a given type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeekAndAdvance checks the parser's peekToken and advances
func (p *Parser) expectPeekAndAdvance(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
