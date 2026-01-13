package database

import (
	"testing"

	"github.com/raskovnik/rdbms/internal/ast"
)

func TestExecuteCreateTable(t *testing.T) {
	db := NewDB()

	stmt := &ast.CreateStatement{
		Table: "users",
		Columns: []ast.ColumnDef{
			{Name: "id", Type: "INT", PrimaryKey: true},
			{Name: "name", Type: "TEXT"},
			{Name: "email", Type: "TEXT", Unique: true},
		},
	}

	err := db.executeCreate(stmt)
	if err != nil {
		t.Fatalf("executeCreate failed: %v", err)
	}

	// Verify table exists
	table, exists := db.tables["users"]
	if !exists {
		t.Fatal("table was not created")
	}

	if len(table.Schema) != 3 {
		t.Errorf("wrong number of columns. expected=3, got=%d", len(table.Schema))
	}

	if table.pkColumn != "id" {
		t.Errorf("wrong primary key. expected=id, got=%s", table.pkColumn)
	}

	// Verify indexes were created
	if _, exists := table.Indexes["id"]; !exists {
		t.Error("primary key index not created")
	}

	if _, exists := table.Indexes["email"]; !exists {
		t.Error("unique index not created")
	}
}

func TestExecuteCreateTableDuplicate(t *testing.T) {
	db := NewDB()

	stmt := &ast.CreateStatement{
		Table: "users",
		Columns: []ast.ColumnDef{
			{Name: "id", Type: "INT", PrimaryKey: true},
		},
	}

	// First creation should succeed
	err := db.executeCreate(stmt)
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	// Second creation should fail
	err = db.executeCreate(stmt)
	if err == nil {
		t.Fatal("expected error for duplicate table, got nil")
	}
}

func TestExecuteInsert(t *testing.T) {
	db := NewDB()

	// Create table first
	createStmt := &ast.CreateStatement{
		Table: "users",
		Columns: []ast.ColumnDef{
			{Name: "id", Type: "INT", PrimaryKey: true},
			{Name: "name", Type: "TEXT"},
		},
	}
	db.executeCreate(createStmt)

	// Insert row
	insertStmt := &ast.InsertStatement{
		Table:  "users",
		Values: []interface{}{1, "Alice"},
	}

	err := db.executeInsert(insertStmt)
	if err != nil {
		t.Fatalf("executeInsert failed: %v", err)
	}

	// Verify row was added
	table := db.tables["users"]
	if len(table.Rows) != 1 {
		t.Fatalf("wrong number of rows. expected=1, got=%d", len(table.Rows))
	}

	if table.Rows[0]["id"] != 1 {
		t.Errorf("wrong id. expected=1, got=%v", table.Rows[0]["id"])
	}

	if table.Rows[0]["name"] != "Alice" {
		t.Errorf("wrong name. expected=Alice, got=%v", table.Rows[0]["name"])
	}
}

func TestExecuteInsertDuplicatePrimaryKey(t *testing.T) {
	db := NewDB()

	createStmt := &ast.CreateStatement{
		Table: "users",
		Columns: []ast.ColumnDef{
			{Name: "id", Type: "INT", PrimaryKey: true},
			{Name: "name", Type: "TEXT"},
		},
	}
	db.executeCreate(createStmt)

	// Insert first row
	insertStmt := &ast.InsertStatement{
		Table:  "users",
		Values: []interface{}{1, "Alice"},
	}
	db.executeInsert(insertStmt)

	// Try to insert duplicate PK
	insertStmt2 := &ast.InsertStatement{
		Table:  "users",
		Values: []interface{}{1, "Bob"},
	}

	err := db.executeInsert(insertStmt2)
	if err == nil {
		t.Fatal("expected error for duplicate primary key, got nil")
	}
}

func TestExecuteSelectAll(t *testing.T) {
	db := setupTestDB(t)

	stmt := &ast.SelectStatement{
		Table:   "users",
		Columns: []string{"*"},
	}

	results, err := db.executeSelect(stmt)
	if err != nil {
		t.Fatalf("executeSelect failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("wrong number of results. expected=2, got=%d", len(results))
	}
}

func TestExecuteSelectWithWhere(t *testing.T) {
	db := setupTestDB(t)

	stmt := &ast.SelectStatement{
		Table:   "users",
		Columns: []string{"*"},
		Where: &ast.WhereClause{
			Column:   "id",
			Operator: "=",
			Value:    1,
		},
	}

	results, err := db.executeSelect(stmt)
	if err != nil {
		t.Fatalf("executeSelect failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("wrong number of results. expected=1, got=%d", len(results))
	}

	if results[0]["name"] != "Alice" {
		t.Errorf("wrong name. expected=Alice, got=%v", results[0]["name"])
	}
}

// Helper function
func setupTestDB(t *testing.T) *Database {
	db := NewDB()

	createStmt := &ast.CreateStatement{
		Table: "users",
		Columns: []ast.ColumnDef{
			{Name: "id", Type: "INT", PrimaryKey: true},
			{Name: "name", Type: "TEXT"},
		},
	}
	db.executeCreate(createStmt)

	db.executeInsert(&ast.InsertStatement{
		Table:  "users",
		Values: []interface{}{1, "Alice"},
	})

	db.executeInsert(&ast.InsertStatement{
		Table:  "users",
		Values: []interface{}{2, "Bob"},
	})

	return db
}
