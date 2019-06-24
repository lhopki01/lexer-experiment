package lexer

import (
	"fmt"
	"unicode"

	"github.com/lhopki01/lexer-experiment/token"
)

type Lexer struct {
	input []rune // use 'rune' to handle Unicode
	start int
	end   int
	char  rune
}

func NewLexer(input []byte) *Lexer {
	l := &Lexer{input: []rune(string(input))}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.end >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.end]
	}
	l.start = l.end
	l.end += 1
}

func (l *Lexer) NewToken() token.Token {
	var tok token.Token
	skipWhitespaceAndComment(l)

	switch l.char {
	case ':':
		tok = token.NewToken(token.COLON, string(l.char))
	case '=':
		tok = token.NewToken(token.EQUAL, string(l.char))
	case ',':
		tok = token.NewToken(token.COMMA, string(l.char))
	case '{':
		tok = token.NewToken(token.LBRACE, string(l.char))
	case '}':
		tok = token.NewToken(token.RBRACE, string(l.char))
	case '[':
		tok = token.NewToken(token.LBRACKET, string(l.char))
	case ']':
		tok = token.NewToken(token.RBRACKET, string(l.char))
	case '\n':
		tok = token.NewToken(token.NEWLINE, string(l.char))
	case '+':
		tok = token.NewToken(token.PLUS, string(l.char))
	case '<':
		tok = token.NewToken(token.LTHAN, string(l.char))
	default:
		if isLibrary(l) {
			tok = token.NewToken(token.LIBRARY, string(l.input[l.start:l.end]))
		} else if isString(l) {
			tok = token.NewToken(token.STRING, string(l.input[l.start:l.end]))
		} else if isInteger(l) {
			tok = token.NewToken(token.INTEGER, string(l.input[l.start:l.end]))
		} else if l.char == rune(0) {
			tok = token.NewToken(token.EOF, "")
		} else {
			tok = token.NewToken(token.INVALID, string(l.char))
		}
	}

	fmt.Println(tok.String())
	l.readChar()
	return tok
}

func (l *Lexer) PeakToken() token.Token {
	fmt.Print("Peak: ")
	start := l.start
	end := l.end
	tok := l.NewToken()

	l.start = start
	l.end = end
	l.char = l.input[l.end-1]
	return tok
}

func (l *Lexer) PeakSecondToken() token.Token {
	fmt.Print("PeakSecond: ")
	start := l.start
	end := l.end
	l.NewToken()
	tok := l.NewToken()

	l.start = start
	l.end = end
	l.char = l.input[l.end-1]
	return tok
}

func isInteger(l *Lexer) bool {
	if !unicode.IsDigit(l.char) && string(l.char) != "-" {
		return false
	}

	l.end += 1
	l.char = l.input[l.end]

	for unicode.IsDigit(l.char) {
		l.end += 1
		l.char = l.input[l.end]
	}
	return true
}

func isLibrary(l *Lexer) bool {
	if l.char == '@' {
		for l.end < len(l.input) {
			l.end += 1
			l.char = l.input[l.end]

			if l.input[l.end] == '\n' {
				l.char = l.input[l.end]
				return true
			}
		}
	}
	return false
}

func isString(l *Lexer) bool {
	fmt.Println("in isString")
	if l.char == '"' {
		for l.end < len(l.input) {
			l.end += 1
			l.char = l.input[l.end]

			if l.input[l.end] == '"' {
				l.end += 1
				l.char = l.input[l.end]
				return true
			}
		}
	}
	if l.char == '\'' {
		for l.end < len(l.input) {
			l.end += 1
			l.char = l.input[l.end]

			if l.input[l.end] == '\'' {
				l.end += 1
				l.char = l.input[l.end]
				return true
			}
		}
	}
	if !(l.char < 'a' || l.char > 'z') || !(l.char < 'A' || l.char > 'Z') || l.char == '.' {
		for l.end < len(l.input) {
			l.end += 1
			l.char = l.input[l.end]

			c := l.input[l.end]
			if isWhitespace(c) || c == ':' || c == '=' || c == ',' {
				//l.end += 1
				l.char = l.input[l.end]
				return true
			}
		}
	}
	fmt.Println("at the end")

	return false
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\n', '\t', '\r':
		return true
	default:
		return false
	}
}

func skipWhitespaceAndComment(l *Lexer) {
	for {
		switch l.char {
		case ' ', '\n', '\t', '\r':
			l.readChar()
		case '/':
			l.readChar()
			skipComment(l)
		default:
			return
		}
	}
}

func skipComment(l *Lexer) {
	if l.char == '/' {
		for {
			switch l.char {
			case '\n':
				return
			default:
				l.readChar()
			}
		}

	}
}
