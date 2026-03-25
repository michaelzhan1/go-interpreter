package lexer

import "github.com/michaelzhan1/go-interpreter/token"

// Lexer is the struct that holds lexing information
type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
}

// New returns a new Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // "prime" the lexer with the first ch
	return l
}

// readChar reads a character at Lexer.readPos and sets Lexer.ch, then increments Lexer.readPos
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

// NextToken returns the next token, then advances the lexer in preparation
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = token.NewToken(token.ASSIGN, l.ch)
	case '+':
		tok = token.NewToken(token.PLUS, l.ch)
	case '(':
		tok = token.NewToken(token.LPAREN, l.ch)
	case ')':
		tok = token.NewToken(token.RPAREN, l.ch)
	case '{':
		tok = token.NewToken(token.LBRACE, l.ch)
	case '}':
		tok = token.NewToken(token.RBRACE, l.ch)
	case ',':
		tok = token.NewToken(token.COMMA, l.ch)
	case ';':
		tok = token.NewToken(token.SEMICOLON, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar() // prepare next char
	return tok
}
