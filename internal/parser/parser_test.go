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
		t.Fatalf("stmt is not *InsertStatement. got=%T", stmt)
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
		t.Fatalf("stmt is not *CreateStatement. got=%T", stmt)
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

func TestParseSelect(t *testing.T) {
	tests := []struct {
		input    string
		table    string
		columns  []string
		hasWhere bool
	}{
		{
			"SELECT * FROM users",
			"users",
			[]string{"*"},
			false,
		},
		{
			"SELECT name, email FROM users",
			"users",
			[]string{"name", "email"},
			false,
		},
		{
			"SELECT * FROM users WHERE id = 1",
			"users",
			[]string{"*"},
			true,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		stmt, err := p.ParseStatement()
		if err != nil {
			t.Fatalf("ParseStatement() for %q returned error: %v", tt.input, err)
		}

		selectStmt, ok := stmt.(*ast.SelectStatement)
		if !ok {
			t.Fatalf("stmt is not *SelectStatement for %q. got=%T", tt.input, stmt)
		}

		if selectStmt.Table != tt.table {
			t.Errorf("table name wrong for %q. expected=%s, got=%s",
				tt.input, tt.table, selectStmt.Table)
		}

		if len(selectStmt.Columns) != len(tt.columns) {
			t.Errorf("wrong number of columns for %q. expected=%d, got=%d",
				tt.input, len(tt.columns), len(selectStmt.Columns))
		}

		if tt.hasWhere && selectStmt.Where == nil {
			t.Errorf("expected WHERE clause for %q", tt.input)
		}
	}
}

func TestParseUpdateSingleColumn(t *testing.T) {
	input := "UPDATE users SET name = 'Bob' WHERE id = 1"

	l := lexer.New(input)
	p := New(l)

	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement() returned error: %v", err)
	}

	updateStmt, ok := stmt.(*ast.UpdateStatement)
	if !ok {
		t.Fatalf("stmt is not *UpdateStatement. got=%T", stmt)
	}

	if updateStmt.Table != "users" {
		t.Errorf("table name wrong. expected=users, got=%s", updateStmt.Table)
	}

	if len(updateStmt.Updates) != 1 {
		t.Fatalf("wrong number of updates. expected=1, got=%d", len(updateStmt.Updates))
	}

	if updateStmt.Updates[0].Column != "name" {
		t.Errorf("column wrong. expected=name, got=%s", updateStmt.Updates[0].Column)
	}

	if updateStmt.Updates[0].Value != "Bob" {
		t.Errorf("value wrong. expected=Bob, got=%v", updateStmt.Updates[0].Value)
	}

	if updateStmt.Where == nil {
		t.Fatal("expected WHERE clause")
	}

	if updateStmt.Where.Column != "id" {
		t.Errorf("WHERE column wrong. expected=id, got=%s", updateStmt.Where.Column)
	}

	if updateStmt.Where.Value != 1 {
		t.Errorf("WHERE value wrong. expected=1, got=%v", updateStmt.Where.Value)
	}
}

func TestParseUpdateMultipleColumns(t *testing.T) {
	input := "UPDATE users SET name = 'Bob', email = 'bob@example.com', age = 30 WHERE id = 1"

	l := lexer.New(input)
	p := New(l)

	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement() returned error: %v", err)
	}

	updateStmt, ok := stmt.(*ast.UpdateStatement)
	if !ok {
		t.Fatalf("stmt is not *UpdateStatement. got=%T", stmt)
	}

	if updateStmt.Table != "users" {
		t.Errorf("table name wrong. expected=users, got=%s", updateStmt.Table)
	}

	if len(updateStmt.Updates) != 3 {
		t.Fatalf("wrong number of updates. expected=3, got=%d", len(updateStmt.Updates))
	}

	// Check first update (name = 'Bob')
	if updateStmt.Updates[0].Column != "name" {
		t.Errorf("updates[0] column wrong. expected=name, got=%s", updateStmt.Updates[0].Column)
	}
	if updateStmt.Updates[0].Value != "Bob" {
		t.Errorf("updates[0] value wrong. expected=Bob, got=%v", updateStmt.Updates[0].Value)
	}

	// Check second update (email = 'bob@example.com')
	if updateStmt.Updates[1].Column != "email" {
		t.Errorf("updates[1] column wrong. expected=email, got=%s", updateStmt.Updates[1].Column)
	}
	if updateStmt.Updates[1].Value != "bob@example.com" {
		t.Errorf("updates[1] value wrong. expected=bob@example.com, got=%v", updateStmt.Updates[1].Value)
	}

	// Check third update (age = 30)
	if updateStmt.Updates[2].Column != "age" {
		t.Errorf("updates[2] column wrong. expected=age, got=%s", updateStmt.Updates[2].Column)
	}
	if updateStmt.Updates[2].Value != 30 {
		t.Errorf("updates[2] value wrong. expected=30, got=%v", updateStmt.Updates[2].Value)
	}

	if updateStmt.Where == nil {
		t.Fatal("expected WHERE clause")
	}
}

func TestParseUpdateNoWhere(t *testing.T) {
	input := "UPDATE users SET name = 'Everyone'"

	l := lexer.New(input)
	p := New(l)

	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement() returned error: %v", err)
	}

	updateStmt, ok := stmt.(*ast.UpdateStatement)
	if !ok {
		t.Fatalf("stmt is not *UpdateStatement. got=%T", stmt)
	}

	if updateStmt.Table != "users" {
		t.Errorf("table name wrong. expected=users, got=%s", updateStmt.Table)
	}

	if len(updateStmt.Updates) != 1 {
		t.Fatalf("wrong number of updates. expected=1, got=%d", len(updateStmt.Updates))
	}

	if updateStmt.Where != nil {
		t.Error("expected no WHERE clause, but got one")
	}
}

func TestParseUpdateMultipleColumnsNoWhere(t *testing.T) {
	input := "UPDATE users SET name = 'Bob', status = 'active'"

	l := lexer.New(input)
	p := New(l)

	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement() returned error: %v", err)
	}

	updateStmt, ok := stmt.(*ast.UpdateStatement)
	if !ok {
		t.Fatalf("stmt is not *UpdateStatement. got=%T", stmt)
	}

	if len(updateStmt.Updates) != 2 {
		t.Fatalf("wrong number of updates. expected=2, got=%d", len(updateStmt.Updates))
	}

	if updateStmt.Where != nil {
		t.Error("expected no WHERE clause, but got one")
	}
}

func TestParseDeleteWithWhere(t *testing.T) {
	input := "DELETE FROM users WHERE id = 1"

	l := lexer.New(input)
	p := New(l)

	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement() returned error: %v", err)
	}

	deleteStmt, ok := stmt.(*ast.DeleteStatement)
	if !ok {
		t.Fatalf("stmt is not *DeleteStatement. got=%T", stmt)
	}

	if deleteStmt.Table != "users" {
		t.Errorf("table name wrong. expected=users, got=%s", deleteStmt.Table)
	}

	if deleteStmt.Where == nil {
		t.Fatal("expected WHERE clause")
	}

	if deleteStmt.Where.Column != "id" {
		t.Errorf("WHERE column wrong. expected=id, got=%s", deleteStmt.Where.Column)
	}

	if deleteStmt.Where.Value != 1 {
		t.Errorf("WHERE value wrong. expected=1, got=%v", deleteStmt.Where.Value)
	}
}

func TestParseDeleteWithoutWhere(t *testing.T) {
	input := "DELETE FROM users"

	l := lexer.New(input)
	p := New(l)

	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement() returned error: %v", err)
	}

	deleteStmt, ok := stmt.(*ast.DeleteStatement)
	if !ok {
		t.Fatalf("stmt is not *DeleteStatement. got=%T", stmt)
	}

	if deleteStmt.Table != "users" {
		t.Errorf("table name wrong. expected=users, got=%s", deleteStmt.Table)
	}

	if deleteStmt.Where != nil {
		t.Error("expected no WHERE clause, but got one")
	}
}

func TestParseDeleteWithStringCondition(t *testing.T) {
	input := "DELETE FROM users WHERE email = 'test@example.com'"

	l := lexer.New(input)
	p := New(l)

	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement() returned error: %v", err)
	}

	deleteStmt, ok := stmt.(*ast.DeleteStatement)
	if !ok {
		t.Fatalf("stmt is not *DeleteStatement. got=%T", stmt)
	}

	if deleteStmt.Where == nil {
		t.Fatal("expected WHERE clause")
	}

	if deleteStmt.Where.Column != "email" {
		t.Errorf("WHERE column wrong. expected=email, got=%s", deleteStmt.Where.Column)
	}

	if deleteStmt.Where.Value != "test@example.com" {
		t.Errorf("WHERE value wrong. expected=test@example.com, got=%v", deleteStmt.Where.Value)
	}
}
