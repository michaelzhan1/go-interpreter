package ast

import "github.com/michaelzhan1/go-interpreter/token"

// Node is the underlying interface for all nodes in an AST
type Node interface {
	TokenLiteral() string
}

// Statement is a Node that does not produce a value
type Statement interface {
	Node
	statementNode()
}

// Expression is a Node that produces a value
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of an AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// Identifier is an expression node that represents a token.IDENT token
type Identifier struct {
	Token token.Token // token.IDENT
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) expressionNode() {}

// LetStatement is a statement node that represents a token.LET token
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) statementNode() {}

// ReturnStatement is a statement node thate represents a token.RETURN token
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) statementNode() {}
