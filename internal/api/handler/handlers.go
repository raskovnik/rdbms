package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/raskovnik/rdbms/internal/app"
	"github.com/raskovnik/rdbms/internal/ast"
)

// GET /todos -> list all todos
func GetTodos(app *app.WebApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stmt := &ast.SelectStatement{ // define our select statement
			Table:   "todos",
			Columns: []string{"*"},
		}

		rows, err := app.DB.Execute(stmt) // execute the statement

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}

// POST /todos -> create a new todo
func CreateTodo(app *app.WebApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Task string `json:"task"`
		}

		// obtain the json
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if req.Task == "" {
			http.Error(w, "task cannot be empty", http.StatusBadRequest)
			return
		}

		// generate id
		id := app.GetNextID()

		// create statememnt definition
		stmt := &ast.InsertStatement{
			Table: "todos",
			Values: []interface{}{
				id,
				req.Task,
				0, // not complete
				time.Now().Format(time.RFC3339),
			},
		}

		_, err := app.DB.Execute(stmt)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   id,
			"task": req.Task,
		})
	}
}

// PUT /todos/{id} -> update todo (mark complete/incomplete)
func UpdateTodo(app *app.WebApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id") // get todo id

		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
		}

		var req struct {
			Completed int `json:"completed"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		stmt := &ast.UpdateStatement{ // update statement definition
			Table: "todos",
			Updates: []ast.ColumnUpdate{
				{Column: "completed", Value: req.Completed},
			},
			Where: &ast.WhereClause{
				Column:   "id",
				Operator: "=",
				Value:    id,
			},
		}

		res, err := app.DB.Execute(stmt)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if res == 0 {
			http.Error(w, "todo not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	}
}

// DELETE /todos/{id}- delete todo
func DeleteTodo(app *app.WebApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id") // get todo id

		id, err := strconv.Atoi(idStr)

		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		stmt := &ast.DeleteStatement{ // define our delete statement
			Table: "todos",
			Where: &ast.WhereClause{
				Column:   "id",
				Operator: "=",
				Value:    id,
			},
		}

		rows, err := app.DB.Execute(stmt)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rows == 0 { // no rows affected = no todo found
			http.Error(w, "todo not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
	}
}
