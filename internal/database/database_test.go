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
