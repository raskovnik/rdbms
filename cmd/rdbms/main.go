package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/raskovnik/rdbms/internal/api/routes"
	"github.com/raskovnik/rdbms/internal/app"
	"github.com/raskovnik/rdbms/internal/engine"
	"github.com/raskovnik/rdbms/internal/repl"
)

func main() {
	mode := flag.String("mode", "repl", "Mode: repl or webapp")
	port := flag.String("port", "8080", "Port for webapp mode")
	flag.Parse()

	db := engine.NewDB()

	switch *mode {
	case "repl":
		repl.Start(os.Stdin, os.Stdout)
	case "webapp":
		app := app.NewWebApp(db)

		// setup schema
		if _, err := app.SetupSchema(); err != nil {
			log.Printf("Warning: %v (table may already exist)", err)
		}

		// routes
		router := routes.NewRouter(app)

		addr := ":" + *port
		fmt.Printf("Web app running on http://localhost%s\n", addr)
		http.ListenAndServe(addr, router)
	}
}
