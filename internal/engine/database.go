package engine

import (
	"fmt"
	"sync"

	"github.com/raskovnik/rdbms/internal/ast"
)

type Database struct {
	tables map[string]*Table
	mu     sync.RWMutex
}

func NewDB() *Database {
	return &Database{
		tables: make(map[string]*Table),
	}
}

type Table struct {
	Name     string
	Schema   []ast.ColumnDef
	Rows     []Row
	Indexes  map[string]*Index
	pkColumn string
}

type Row map[string]interface{}

func NewTable(name string, schema []ast.ColumnDef) *Table {
	table := &Table{
		Name:    name,
		Schema:  schema,
		Rows:    []Row{},
		Indexes: make(map[string]*Index),
	}

	// identify the primary key column
	for _, col := range schema {
		if col.PrimaryKey {
			table.pkColumn = col.Name
			table.Indexes[col.Name] = NewIndex(col.Name)
		}
		if col.Unique {
			table.Indexes[col.Name] = NewIndex(col.Name)
		}
	}

	return table
}

func (db *Database) Execute(stmt ast.Statement) (interface{}, error) {
	switch s := stmt.(type) {
	case *ast.CreateStatement:
		return nil, db.executeCreate(s)
	case *ast.InsertStatement:
		return nil, db.executeInsert(s)
	case *ast.SelectStatement:
		return db.executeSelect(s)
	case *ast.DeleteStatement:
		return db.executeDelete(s)
	case *ast.UpdateStatement:
		return db.executeUpdate(s)
	default:
		return nil, fmt.Errorf("unknown statement type: %T", stmt)
	}
}
