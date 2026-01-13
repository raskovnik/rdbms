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

func (db *Database) executeSelect(stmt *ast.SelectStatement) ([]Row, error) {
	// select */[]columns(string) from tablename() where (optional) condition

	db.mu.Lock()
	defer db.mu.Unlock()

	// get table
	table, exists := db.tables[stmt.Table]
	if !exists {
		return nil, fmt.Errorf("table %s does not exist", stmt.Table)
	}

	var results []Row

	// if WHERE clause exists and targets an indexed column, use index for better performance
	if stmt.Where != nil {
		if index, indexed := table.Indexes[stmt.Where.Column]; indexed {
			// use index lookup
			rowIndices := index.Lookup(stmt.Where.Value)

			for _, idx := range rowIndices {
				row := table.Rows[idx]
				if evaluateWhere(row, stmt.Where) {
					results = append(results, row)
				}
			}
		} else {
			// full table scan
			for _, row := range table.Rows {
				if evaluateWhere(row, stmt.Where) {
					results = append(results, row)
				}
			}
		}
	} else {
		// no where clause -> return all rows
		results = table.Rows
	}

	// filter columns if not SELECT *
	if len(stmt.Columns) == 1 && stmt.Columns[0] == "*" {
		return results, nil
	}

	// project specific columns
	projected := make([]Row, len(results))
	for i, row := range results {
		projectedRow := make(Row)
		for _, col := range stmt.Columns {
			if val, exists := row[col]; exists {
				projectedRow[col] = val
			} else {
				return nil, fmt.Errorf("column %s does not exist", col)
			}
		}
		projected[i] = projectedRow
	}

	return projected, nil
}

func evaluateWhere(row Row, where *ast.WhereClause) bool {
	value, exists := row[where.Column]
	if !exists {
		return false
	}

	switch where.Operator {
	case "=":
		return value == where.Value
	case ">":
		if intVal, ok := value.(int); ok {
			if intWhere, ok := where.Value.(int); ok {
				return intVal > intWhere
			}
		}
		return false
	case "<":
		if intVal, ok := value.(int); ok {
			if intWhere, ok := where.Value.(int); ok {
				return intVal < intWhere
			}
		}
		return false
	default:
		return false
	}
}
func (db *Database) executeDelete(stmt *ast.DeleteStatement) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	table, exists := db.tables[stmt.Table]
	if !exists {
		return 0, fmt.Errorf("table %s does not exist", stmt.Table)
	}

	// DELETE FROM table (no WHERE clause) - delete all
	if stmt.Where == nil {
		deletedCount := len(table.Rows)
		table.Rows = []Row{}
		// Clear all indexes
		for _, index := range table.Indexes {
			index.Data = make(map[interface{}][]int)
		}
		return deletedCount, nil
	}

	// Filter out rows that match WHERE condition
	var newRows []Row
	deletedCount := 0

	// Use index if available for optimization
	if index, indexed := table.Indexes[stmt.Where.Column]; indexed {
		// Build set of row indices to delete
		toDelete := make(map[int]bool)
		candidates := index.Lookup(stmt.Where.Value)
		for _, idx := range candidates {
			if evaluateWhere(table.Rows[idx], stmt.Where) {
				toDelete[idx] = true
			}
		}

		// Keep rows that aren't in toDelete set
		for i, row := range table.Rows {
			if toDelete[i] {
				deletedCount++
			} else {
				newRows = append(newRows, row)
			}
		}
	} else {
		// Full table scan - keep rows that don't match
		for _, row := range table.Rows {
			if evaluateWhere(row, stmt.Where) {
				deletedCount++
			} else {
				newRows = append(newRows, row)
			}
		}
	}

	// Replace old rows with filtered rows
	table.Rows = newRows

	// Rebuild indexes to reflect new row positions
	db.rebuildIndexes(table)

	return deletedCount, nil
}

func (db *Database) rebuildIndexes(table *Table) {
	// Clear existing index data
	for _, index := range table.Indexes {
		index.Data = make(map[interface{}][]int)
	}

	// Rebuild indexes from current rows
	for rowIndex, row := range table.Rows {
		for colName, index := range table.Indexes {
			if value, exists := row[colName]; exists {
				index.Add(value, rowIndex)
			}
		}
	}
}
