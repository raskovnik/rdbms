package parser

import (
	"fmt"
	"strconv"

	"github.com/raskovnik/rdbms/internal/ast"
	"github.com/raskovnik/rdbms/internal/lexer"
	"github.com/raskovnik/rdbms/internal/token"
)

type Parser struct {
	lexer     *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	// initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p

}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseStatement() (ast.Statement, error) {
	switch p.curToken.Type {
	case token.INSERT:
		return p.parseInsert()
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.curToken.Type)
	}
}

func (p *Parser) parseInsert() (*ast.InsertStatement, error) {
	stmt := &ast.InsertStatement{}

	// current token is insert
	if !p.expectPeek(token.INTO) {
		return nil, fmt.Errorf("expected INTO after INSERT")
	}

	// get table name
	if !p.expectPeek(token.IDENT) {
		return nil, fmt.Errorf("expected table name")
	}
	stmt.Table = p.curToken.Literal

	// expect values keyword
	if !p.expectPeek(token.VALUES) {
		return nil, fmt.Errorf("expected VALUES")
	}

	// expect opening parenthesis
	if !p.expectPeek(token.LPAREN) {
		return nil, fmt.Errorf("expected ( after VALUES")
	}

	// parse values
	stmt.Values = []interface{}{}
	p.nextToken() // move to the first value

	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		stmt.Values = append(stmt.Values, val)

		// check for comma or closing parenthesis
		if p.peekTokenIs(token.COMMA) {
			p.nextToken() // consume comma
			p.nextToken() // move to next value
		} else if p.peekTokenIs(token.RPAREN) {
			p.nextToken() // move to closing paren
		} else {
			return nil, fmt.Errorf("expected ',' or ')' after value")
		}
	}

	if !p.curTokenIs(token.RPAREN) {
		return nil, fmt.Errorf("expected ')' to close values")
	}

	return stmt, nil

}

func (p *Parser) parseValue() (interface{}, error) {
	switch p.curToken.Type {
	case token.INT:
		val, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			return nil, fmt.Errorf("could not parse %s as integer", p.curToken.Literal)
		}
		return val, nil
	case token.STRING:
		return p.curToken.Literal, nil
	default:
		return nil, fmt.Errorf("unexpected value type: %s", p.curToken.Type)
	}
}
