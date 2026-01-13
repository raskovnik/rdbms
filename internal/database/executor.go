package database

import (
	"fmt"

	"github.com/raskovnik/rdbms/internal/ast"
)

func (db *Database) executeCreate(stmt *ast.CreateStatement) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// check if table already exists
	if _, exists := db.tables[stmt.Table]; exists {
		return fmt.Errorf("table %s already exists", stmt.Table)
	}

	// validate schema
	if len(stmt.Columns) == 0 {
		return fmt.Errorf("table must have at least one column")
	}

	// check for multiple primary keys
	pkCount := 0
	for _, col := range stmt.Columns {
		if col.PrimaryKey {
			pkCount++
		}
	}

	if pkCount > 1 {
		return fmt.Errorf("table can only have one primary key")
	}

	// create the table
	db.tables[stmt.Table] = NewTable(stmt.Table, stmt.Columns)

	return nil
}

func (db *Database) executeInsert(stmt *ast.InsertStatement) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// insert into users values(a, b, c)

	// get table
	table, exists := db.tables[stmt.Table]
	if !exists {
		return fmt.Errorf("table %s does not exist", stmt.Table)
	}

	// validate value count matches schema
	if len(stmt.Values) != len(table.Schema) {
		return fmt.Errorf("value count %d does not match column count %d", len(stmt.Values), len(table.Schema))
	}

	// create row
	row := make(Row)
	for i, col := range table.Schema {
		value := stmt.Values[i]

		// basic type validation
		if !isValidType(value, col.Type) {
			return fmt.Errorf("value %v is not valid type for %s", value, col.Type)
		}

		row[col.Name] = value
	}

	// check constraints before adding row
	rowIndex := len(table.Rows) // -> index after appending

	for colName, index := range table.Indexes {
		value := row[colName]

		// cehck if value already exists (violates pk or unique)
		if index.Exists(value) {
			return fmt.Errorf("duplicate value %v for column %s", value, colName)
		}
	}

	// add row to table
	table.Rows = append(table.Rows, row)

	// update indexes
	for colName, index := range table.Indexes {
		value := row[colName]
		index.Add(value, rowIndex)
	}
	return nil
}

func isValidType(value interface{}, expectedType string) bool {
	switch expectedType {
	case "INT":
		_, ok := value.(int)
		return ok
	case "TEXT":
		_, ok := value.(string)
		return ok
	case "BOOL":
		_, ok := value.(bool)
		return ok
	default:
		return false
	}
}
