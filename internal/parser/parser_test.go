package parser

import (
	"testing"

	"github.com/raskovnik/rdbms/internal/ast"
	"github.com/raskovnik/rdbms/internal/lexer"
)

func TestParseInsert(t *testing.T) {
	input := "INSERT INTO users VALUES (1, 'Alice', 'alice@example.com')"

	l := lexer.New(input)
	p := New(l)

	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement() returned error: %v", err)
	}

	insertStmt, ok := stmt.(*ast.InsertStatement)
	if !ok {
		t.Fatalf("stmt is not *InsertStmt. got=%T", stmt)
	}

	if insertStmt.Table != "users" {
		t.Errorf("table name wrong. expected=users, got=%s", insertStmt.Table)
	}

	if len(insertStmt.Values) != 3 {
		t.Fatalf("wrong number of values. expected=3, got=%d", len(insertStmt.Values))
	}

	if insertStmt.Values[0] != 1 {
		t.Errorf("value[0] wrong. expected=1, got=%v", insertStmt.Values[0])
	}

	if insertStmt.Values[1] != "Alice" {
		t.Errorf("value[1] wrong. expected=Alice, got=%v", insertStmt.Values[1])
	}
}
