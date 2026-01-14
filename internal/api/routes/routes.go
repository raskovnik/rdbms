package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	handlers "github.com/raskovnik/rdbms/internal/api/handler"
	"github.com/raskovnik/rdbms/internal/api/templates"
	"github.com/raskovnik/rdbms/internal/app"
)

func NewRouter(app *app.WebApp) http.Handler {
	r := chi.NewRouter()

	// routes
	r.Get("/todos", handlers.GetTodos(app))
	r.Post("/todos", handlers.CreateTodo(app))
	r.Put("/todos/{id}", handlers.UpdateTodo(app))
	r.Delete("/todos/{id}", handlers.DeleteTodo(app))

	// serve index page
	r.Get("/", templates.ServeIndex)
	return r
}
