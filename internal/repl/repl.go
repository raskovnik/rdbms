package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/raskovnik/rdbms/internal/engine"
	"github.com/raskovnik/rdbms/internal/lexer"
	"github.com/raskovnik/rdbms/internal/parser"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	db := engine.NewDB() // create DB

	for {
		fmt.Fprint(out, "db> ")
		if !scanner.Scan() {
			break // EOF
		}

		line := scanner.Text()

		lex := lexer.New(line)    // lexer
		parser := parser.New(lex) // parser

		stmt, err := parser.ParseStatement()
		if err != nil {
			fmt.Fprintln(out, "Parse error:", err)
			continue
		}

		// execute the query
		res, err := db.Execute(stmt)
		if err != nil {
			fmt.Fprintln(out, "Error:", err)
			continue
		}

		// print output
		switch v := res.(type) {
		case int:
			fmt.Fprintf(out, "rows affected: %d\n", v)
		case []engine.Row:
			printRows(v, out)
		default:
			fmt.Fprintln(out, "OK")
		}
	}
}

func printRows(rows []engine.Row, out io.Writer) {
	for _, row := range rows {
		for col, val := range row {
			fmt.Fprintf(out, "%s=%v ", col, val)
		}
		fmt.Fprintln(out)
	}
}
