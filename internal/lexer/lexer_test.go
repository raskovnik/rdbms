package lexer

import (
	"fmt"
	"testing"

	"github.com/raskovnik/rdbms/internal/token"
)

func TestNextToken(t *testing.T) {
	input := `CREATE TABLE users (id INT PRIMARY KEY, name TEXT)`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.CREATE, "CREATE"},
		{token.TABLE, "TABLE"},
		{token.IDENT, "users"},
		{token.LPAREN, "("},
		{token.IDENT, "id"},
		{token.TYPE_INT, "INT"},
		{token.PRIMARY, "PRIMARY"},
		{token.KEY, "KEY"},
		{token.COMMA, ","},
		{token.IDENT, "name"},
		{token.TYPE_TEXT, "TEXT"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestStringLiterals(t *testing.T) {
	input := `INSERT INTO users VALUES (1, 'Alice')`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INSERT, "INSERT"},
		{token.INTO, "INTO"},
		{token.IDENT, "users"},
		{token.VALUES, "VALUES"},
		{token.LPAREN, "("},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.STRING, "Alice"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		fmt.Printf("Type: %s, Literal: %q\n", tok.Type, tok.Literal)

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestSelectWithWhere(t *testing.T) {
	input := `SELECT * FROM users WHERE id = 1`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.SELECT, "SELECT"},
		{token.ASTERISK, "*"},
		{token.FROM, "FROM"},
		{token.IDENT, "users"},
		{token.WHERE, "WHERE"},
		{token.IDENT, "id"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestAllCommands(t *testing.T) {
	input := `
	CREATE TABLE users (id INT PRIMARY KEY, email TEXT UNIQUE)
	INSERT INTO users VALUES (1, 'test@example.com')
	SELECT name FROM users WHERE id = 1
	UPDATE users SET name = 'Bob' WHERE id = 1
	DELETE FROM users WHERE id = 1
	SELECT a.name, b.total FROM a JOIN b ON a.id = b.id
	`

	l := New(input)

	// Just make sure we can tokenize everything without ILLEGAL tokens
	for {
		tok := l.NextToken()

		if tok.Type == token.ILLEGAL {
			t.Fatalf("Found ILLEGAL token: %q", tok.Literal)
		}

		if tok.Type == token.EOF {
			break
		}
	}
}
