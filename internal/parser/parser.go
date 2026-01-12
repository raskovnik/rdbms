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
	case token.CREATE:
		return p.parseCreateStatement()
	case token.SELECT:
		return p.parseSelectStatement()
	case token.UPDATE:
		return p.parseUpdateStatement()
	case token.DELETE:
		return p.parseDeleteStatement()
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

func (p *Parser) parseCreateStatement() (*ast.CreateStatement, error) {
	stmt := &ast.CreateStatement{}

	// current token is CREATE
	if !p.expectPeek(token.TABLE) {
		return nil, fmt.Errorf("expected TABLE after CREATE")
	}

	// get table name
	if !p.expectPeek(token.IDENT) {
		return nil, fmt.Errorf("expected table name")
	}

	stmt.Table = p.curToken.Literal

	// expect opening paren
	if !p.expectPeek(token.LPAREN) {
		return nil, fmt.Errorf("expected '(' after table name")
	}

	// parse columns
	stmt.Columns = []ast.ColumnDef{}

	// move to first column
	p.nextToken()

	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		col, err := p.parseColumnDef()
		if err != nil {
			return nil, err
		}
		stmt.Columns = append(stmt.Columns, col)

		// check for comma -> more columnr or closing rparen -> done
		if p.peekTokenIs(token.COMMA) {
			p.nextToken() // consume comma
			p.nextToken() // move to next column name
		} else if p.peekTokenIs(token.RPAREN) {
			p.nextToken() // move to closing rparen
		} else {
			return nil, fmt.Errorf("expected ',' or ')' after column definition")
		}
	}

	if !p.curTokenIs(token.RPAREN) {
		return nil, fmt.Errorf("expected ) to close column definitions")
	}

	return stmt, nil
}

func (p *Parser) parseColumnDef() (ast.ColumnDef, error) {
	col := ast.ColumnDef{}

	// current token should be column name
	if !p.curTokenIs(token.IDENT) {
		return col, fmt.Errorf("expected column name, got %s", p.curToken.Type)
	}

	col.Name = p.curToken.Literal

	// next should be the type
	p.nextToken()
	switch p.curToken.Type {
	case token.TYPE_INT:
		col.Type = "INT"
	case token.TYPE_TEXT:
		col.Type = "TEXT"
	case token.TYPE_BOOL:
		col.Type = "BOOL"
	default:
		return col, fmt.Errorf("expected type (INT, STRING, BOOL), got %s", p.curToken.Type)
	}

	// check for primary key or unique
	if p.peekTokenIs(token.PRIMARY) {
		p.nextToken() // consume primary
		if !p.expectPeek(token.KEY) {
			return col, fmt.Errorf("expected KEY after PRIMARY")
		}

		col.PrimaryKey = true
	} else if p.peekTokenIs(token.UNIQUE) {
		p.nextToken() // consume unique
		col.Unique = true
	}

	return col, nil
}

func (p *Parser) parseSelectStatement() (*ast.SelectStatement, error) {

	stmt := &ast.SelectStatement{}

	// current token should be SELECT
	p.nextToken() // move to column name or *

	// parse column(s)
	if p.curTokenIs(token.ASTERISK) {
		stmt.Columns = []string{"*"}
		p.nextToken() // move past *
	} else {
		// parse column list
		stmt.Columns = []string{}
		for {
			if !p.curTokenIs(token.IDENT) {
				return nil, fmt.Errorf("expected column name but got %s instead", p.curToken.Type)
			}

			stmt.Columns = append(stmt.Columns, p.curToken.Literal)

			if p.peekTokenIs(token.COMMA) {
				p.nextToken() // consume comma
				p.nextToken() // move to next column
			} else {
				p.nextToken() // move past last column
				break
			}
		}

	}

	// expect FROM
	if !p.curTokenIs(token.FROM) {
		return nil, fmt.Errorf("expected FROM after * or column name")
	}

	// get table name
	if !p.expectPeek(token.IDENT) {
		return nil, fmt.Errorf("expected table name")
	}

	stmt.Table = p.curToken.Literal

	// check for WHERE clause
	if p.peekTokenIs(token.WHERE) {
		p.nextToken() // consume WHERE
		wc, err := p.parseWhereClause()
		if err != nil {
			return nil, err
		}

		stmt.Where = wc
	}
	return stmt, nil
}

func (p *Parser) parseWhereClause() (*ast.WhereClause, error) {
	where := &ast.WhereClause{}

	// get column name
	if !p.expectPeek(token.IDENT) {
		return nil, fmt.Errorf("expected column name in WHERE clause, got %s", p.curToken.Type)
	}

	where.Column = p.curToken.Literal

	// get the operator
	p.nextToken() // move to the operator
	switch p.curToken.Type {
	case token.ASSIGN:
		where.Operator = "="
	case token.GT:
		where.Operator = ">"
	case token.LT:
		where.Operator = "<"
	default:
		return nil, fmt.Errorf("expected operator (=, >, <), got %s instead", p.curToken.Type)
	}

	// get the value
	p.nextToken()
	val, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	where.Value = val

	return where, nil
}

func (p *Parser) parseUpdateStatement() (*ast.UpdateStatement, error) {
	stmt := &ast.UpdateStatement{}

	// current token should be UPDATe
	if !p.expectPeek(token.IDENT) {
		return nil, fmt.Errorf("expected table name after UPDATE, got %s", p.peekToken.Type)
	}

	stmt.Table = p.curToken.Literal

	if !p.expectPeek(token.SET) {
		return nil, fmt.Errorf("expected SET after table name")
	}

	// loop through column = value pairs
	stmt.Updates = []ast.ColumnUpdate{}
	for {
		if !p.expectPeek(token.IDENT) {
			return nil, fmt.Errorf("expected column name")
		}

		colName := p.curToken.Literal

		// parse =
		if !p.expectPeek(token.ASSIGN) {
			return nil, fmt.Errorf("expected =")
		}

		p.nextToken() // consume =
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		stmt.Updates = append(stmt.Updates, ast.ColumnUpdate{Column: colName, Value: val})

		// check for comma -> more updates or WHERE/EOF -> done
		if p.peekTokenIs(token.COMMA) {
			p.nextToken() // consume comma
			continue
		} else {
			break
		}
	}

	// WHERE clause is optional
	if p.peekTokenIs(token.WHERE) {
		p.nextToken() // consume WHERE
		wc, err := p.parseWhereClause()
		if err != nil {
			return nil, err
		}
		stmt.Where = wc
	}

	return stmt, nil
}

func (p *Parser) parseDeleteStatement() (*ast.DeleteStatement, error) {
	stmt := &ast.DeleteStatement{}

	// current token should be DELETE
	if !p.expectPeek(token.FROM) {
		return nil, fmt.Errorf("expected FROM after DELETE")
	}

	// get table name
	if !p.expectPeek(token.IDENT) {
		return nil, fmt.Errorf("expected table name after FROM")
	}

	stmt.Table = p.curToken.Literal

	// WHERE is optional
	if p.peekTokenIs(token.WHERE) {
		p.nextToken() // consume where
		wc, err := p.parseWhereClause()
		if err != nil {
			return nil, err
		}

		stmt.Where = wc
	}

	return stmt, nil
}
