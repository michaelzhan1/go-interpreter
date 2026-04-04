package ast

import (
	"bytes"
	"strings"

	"github.com/michaelzhan1/go-interpreter/token"
)

// Node is the underlying interface for all nodes in an AST
type Node interface {
	TokenLiteral() string
	String() string
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
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

var _ Node = &Program{}

// Identifier is an expression node that represents a token.IDENT token
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

var _ Node = &Identifier{}

// IntegerLiteral is an expression node that represents a standalone integer
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

var _ Node = &IntegerLiteral{}

// BooleanLiteral is an expression node that represents a standalone boolean
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

var _ Node = &BooleanLiteral{}

// FunctionLiteral is an expression node that represents a function
type FunctionLiteral struct {
	Token      token.Token // fn token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	params := make([]string, len(fl.Parameters))
	for i, p := range fl.Parameters {
		params[i] = p.String()
	}

	var out bytes.Buffer

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

var _ Node = &FunctionLiteral{}

// CallExpression is an expression node that represents a function call.
// It can either represent a call from an identifier or from an inlined function.
type CallExpression struct {
	Token     token.Token // '(' token since this is technically an infix expression between the Function and the Arguments
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	args := make([]string, len(ce.Arguments))
	for i, a := range ce.Arguments {
		args[i] = a.String()
	}

	var out bytes.Buffer

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

var _ Node = &CallExpression{}

// IfExpression is an expression node that represents an if-else statement. Else is optional
type IfExpression struct {
	Token       token.Token // if token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString((ie.Alternative.String()))
	}

	return out.String()
}

var _ Node = &IfExpression{}

// PrefixExpression is a prefix expression such as "-5" or "!function(a)"
type PrefixExpression struct {
	Token    token.Token // prefix token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

var _ Node = &PrefixExpression{}

// InfixExpression is an infix expression such as "5+5"
type InfixExpression struct {
	Token    token.Token // operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

var _ Node = &InfixExpression{}

// LetStatement is a statement node that represents a token.LET token
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.Token.Literal + " ")   // let
	out.WriteString(ls.Name.String() + " = ") // variable name

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")

	return out.String()
}

var _ Node = &LetStatement{}

// ReturnStatement is a statement node thate represents a token.RETURN token
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.Token.Literal) // return

	if rs.ReturnValue != nil {
		out.WriteString(" " + rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

var _ Node = &ReturnStatement{}

// ExpressionStatement is a statement that wraps around a single expression. Allows for do-nothing statements such as `x + 5;`
type ExpressionStatement struct {
	Token      token.Token // first token in the expression, TODO: is this even used? or only used for tokenliteral interface
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	var out bytes.Buffer

	if es.Expression != nil {
		out.WriteString(es.Expression.String())
	}

	// does not print semicolon, despite being a statement. Could add it here, though

	return out.String()
}

var _ Node = &ExpressionStatement{}

// BlockStatement represents a nested statement in {}, such as in an if-else expression.
// Essentially, it's a nested program within a part of a program.
type BlockStatement struct {
	Token      token.Token // left brace { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}
