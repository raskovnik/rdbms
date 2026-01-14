package app

import (
	"github.com/raskovnik/rdbms/internal/ast"
	"github.com/raskovnik/rdbms/internal/engine"
)

type WebApp struct {
	DB *engine.Database
}

func NewWebApp(db *engine.Database) *WebApp {
	return &WebApp{DB: db}
}

func (app *WebApp) SetupSchema() (interface{}, error) {
	stmt := &ast.CreateStatement{
		Table: "todos",
		Columns: []ast.ColumnDef{
			{Name: "id", Type: "INT", PrimaryKey: true},
			{Name: "task", Type: "TEXT"},
			{Name: "completed", Type: "INT"},
			{Name: "created_at", Type: "TEXT"},
		},
	}

	return app.DB.Execute(stmt)
}

// handler to get next available id
func (app *WebApp) GetNextID() int {
	stmt := &ast.SelectStatement{
		Table:   "todos",
		Columns: []string{"*"},
	}

	res, _ := app.DB.Execute(stmt)
	maxID := 0

	rows, _ := res.([]engine.Row)

	for _, row := range rows {
		if id, ok := row["id"].(int); ok && id > maxID {
			maxID = id
		}
	}

	return maxID + 1
}
