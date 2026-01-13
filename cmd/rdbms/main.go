package main

import (
	"os"

	"github.com/raskovnik/rdbms/internal/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
