package parser

import (
	"fmt"

	"github.com/lhopki01/lexer-experiment/ast"
	"github.com/lhopki01/lexer-experiment/lexer"
	"github.com/lhopki01/lexer-experiment/token"
)

type Parser struct {
	Lexer *lexer.Lexer
}

func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{Lexer: l}
}

func (p *Parser) Parse() ast.Json {
	tok := p.Lexer.NewToken()
	fmt.Println(tok.String())
	switch tok.Type {
	case token.STRING:
		return &ast.String{string(tok.Lit)}
	case token.INTEGER:
		return &ast.Integer{string(tok.Lit)}
	case token.LBRACE:
		return parseNewlineObject(p)
	case token.LBRACKET:
		//return parseArray(p)
		return parseArrayOrObject(p)
	case token.EOF:
		return nil
	}
	return nil
}

func parseArrayOrObject(p *Parser) ast.Json {
	array := []ast.Json{}
	object := map[string]ast.Json{}

	tok := p.Lexer.PeakToken()
	secondTok := p.Lexer.PeakSecondToken()

	if tok.Type == token.RBRACKET {
		return &ast.Array{array}
	} else if secondTok.Type != token.COLON && secondTok.Type != token.EQUAL {
		fmt.Println("===array===")
		for {
			array = append(array, p.Parse())
			tok = p.Lexer.NewToken()
			if tok.Type == token.RBRACKET {
				return &ast.Array{array}
			}

			if tok.Type != token.COMMA {
				panic(fmt.Sprintf("was expecting ',' got %s in array parse", string(tok.Lit)))
			}
		}
	} else {
		fmt.Println("===object===")
		for {
			key := string(p.Lexer.NewToken().Lit)
			tok = p.Lexer.NewToken() // ':'
			if tok.Type != token.COLON {
				panic(fmt.Sprintf("was expecting ':' got %s", string(tok.Lit)))
			}
			object[key] = p.Parse()
			tok = p.Lexer.NewToken() // ','

			if tok.Type == token.RBRACKET {
				return &ast.Object{object}
			}

			peakTok := p.Lexer.PeakToken()
			if peakTok.Type == token.RBRACKET {
				p.Lexer.NewToken()
				return &ast.Object{object}
			}

			if tok.Type != token.COMMA {
				panic(fmt.Sprintf("was expecting ',' got %s", string(tok.Lit)))
			}
		}
	}

	return &ast.Array{array}
}
func parseNewlineObject(p *Parser) ast.Json {
	array := []ast.Json{}
	object := map[string]ast.Json{}

	tok := p.Lexer.PeakToken()

	if tok.Type == token.RBRACE {
		return &ast.Array{array}
	} else {
		fmt.Println("===newlineobject===")
		for {
			key := string(p.Lexer.NewToken().Lit)
			tok = p.Lexer.NewToken() // ':'
			if tok.Type != token.EQUAL {
				panic(fmt.Sprintf("was expecting '=' got %s", string(tok.Lit)))
			}
			object[key] = p.Parse()
			tok = p.Lexer.PeakToken() // ','

			if tok.Type == token.RBRACE {
				return &ast.Object{object}
			}
		}
	}

	return &ast.Array{array}
}

func parseArray(p *Parser) ast.Json {
	array := []ast.Json{}
	tok := p.Lexer.PeakToken()

	if tok.Type == token.RBRACKET {
		return &ast.Array{array}
		//} else {
		//	array = append(array, p.Parse())
		//	tok = p.Lexer.NewToken()
		//	if tok.Type == token.RBRACKET {
		//		return &ast.Array{array}
		//	}
	}

	for {
		array = append(array, p.Parse())
		tok = p.Lexer.NewToken()
		if tok.Type == token.RBRACKET {
			return &ast.Array{array}
			break
		}

		if tok.Type != token.COMMA {
			panic(fmt.Sprintf("was expecting ',' got %s in array parse", string(tok.Lit)))
		}
	}

	return &ast.Array{array}
}

func parseObject(p *Parser) ast.Json {
	object := map[string]ast.Json{}
	//tok := p.Lexer.NewToken()
	tok := p.Lexer.PeakToken()

	if tok.Type == token.RBRACE { // nothing inside
		return &ast.Object{object}
		//} else {
		//	key := string(tok.Lit)
		//	p.Lexer.NewToken() // ':'
		//	object[key] = p.Parse()
		//	tok = p.Lexer.NewToken()
		//	if tok.Type == token.RBRACE {
		//		return &ast.Object{object}
		//	}
	}

	for {
		key := string(p.Lexer.NewToken().Lit)
		tok = p.Lexer.NewToken() // ':'
		if !(tok.Type == token.COLON || tok.Type == token.EQUAL) {
			panic(fmt.Sprintf("was expecting ':' or '=' got %s", string(tok.Lit)))
		}
		object[key] = p.Parse()
		tok = p.Lexer.NewToken() // ','

		if tok.Type == token.RBRACE {
			return &ast.Object{object}
			break
		}

		if tok.Type != token.COMMA {
			panic(fmt.Sprintf("was expecting ',' got %s", string(tok.Lit)))
		}
	}
	return &ast.Object{object}
}
