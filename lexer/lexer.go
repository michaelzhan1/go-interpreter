package lexer

import (
	"github.com/michaelzhan1/go-interpreter/token"
)

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

// peekChar returns the next char in the lexer without moving the position
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

// NextToken returns the next token, then advances the lexer in preparation
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '"':
		tok = token.NewToken(token.STRING, l.readString())
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.NewToken(token.EQ, string(ch)+string(l.ch))
		} else {
			tok = token.NewToken(token.ASSIGN, string(l.ch))
		}
	case '+':
		tok = token.NewToken(token.PLUS, string(l.ch))
	case '-':
		tok = token.NewToken(token.MINUS, string(l.ch))
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.NewToken(token.NOT_EQ, string(ch)+string(l.ch))
		} else {
			tok = token.NewToken(token.BANG, string(l.ch))
		}
	case '*':
		tok = token.NewToken(token.ASTERISK, string(l.ch))
	case '/':
		tok = token.NewToken(token.SLASH, string(l.ch))
	case '<':
		tok = token.NewToken(token.LT, string(l.ch))
	case '>':
		tok = token.NewToken(token.GT, string(l.ch))
	case '[':
		tok = token.NewToken(token.LBRACKET, string(l.ch))
	case ']':
		tok = token.NewToken(token.RBRACKET, string(l.ch))
	case '(':
		tok = token.NewToken(token.LPAREN, string(l.ch))
	case ')':
		tok = token.NewToken(token.RPAREN, string(l.ch))
	case '{':
		tok = token.NewToken(token.LBRACE, string(l.ch))
	case '}':
		tok = token.NewToken(token.RBRACE, string(l.ch))
	case ',':
		tok = token.NewToken(token.COMMA, string(l.ch))
	case ';':
		tok = token.NewToken(token.SEMICOLON, string(l.ch))
	case 0:
		tok = token.NewToken(token.EOF, "")
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			return token.NewToken(token.LookupIdentOrKeyword(ident), ident) // return early without advancing more
		} else if isDigit(l.ch) {
			return token.NewToken(token.INT, l.readNumber())
		} else {
			tok = token.NewToken(token.ILLEGAL, string(l.ch))
		}
	}

	l.readChar() // prepare next char
	return tok
}

// skipWhitespace reads past any whitespace in the lexer
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads the next keyword or identifier such as a variable name
func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

// readNumber reads the next number as a string
func (l *Lexer) readNumber() string {
	start := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.pos]
}

// readString reads the next string literal
func (l *Lexer) readString() string {
	l.readChar() // move past first '"'
	start := l.pos
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	return l.input[start:l.pos]
}

// isLetter returns if ch is a letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit returns if ch is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
