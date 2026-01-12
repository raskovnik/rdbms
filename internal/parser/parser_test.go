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

func TestParseCreateTable(t *testing.T) {
	input := "CREATE TABLE users (id INT PRIMARY KEY, name TEXT, email TEXT UNIQUE)"

	l := lexer.New(input)
	p := New(l)

	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement() returned error: %v", err)
	}

	createStmt, ok := stmt.(*ast.CreateStatement)
	if !ok {
		t.Fatalf("stmt is not *CreateTableStmt. got=%T", stmt)
	}

	if createStmt.Table != "users" {
		t.Errorf("table name wrong. expected=users, got=%s", createStmt.Table)
	}

	if len(createStmt.Columns) != 3 {
		t.Fatalf("wrong number of columns. expected=3, got=%d", len(createStmt.Columns))
	}

	// Test first column (id INT PRIMARY KEY)
	if createStmt.Columns[0].Name != "id" {
		t.Errorf("column[0] name wrong. expected=id, got=%s", createStmt.Columns[0].Name)
	}
	if createStmt.Columns[0].Type != "INT" {
		t.Errorf("column[0] type wrong. expected=INT, got=%s", createStmt.Columns[0].Type)
	}
	if !createStmt.Columns[0].PrimaryKey {
		t.Errorf("column[0] should be primary key")
	}

	// Test third column (email TEXT UNIQUE)
	if createStmt.Columns[2].Name != "email" {
		t.Errorf("column[2] name wrong. expected=email, got=%s", createStmt.Columns[2].Name)
	}
	if !createStmt.Columns[2].Unique {
		t.Errorf("column[2] should be unique")
	}
}
