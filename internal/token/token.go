package token

import "strings"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (

	// keywords
	CREATE  = "CREATE"
	TABLE   = "TABLE"
	INSERT  = "INSERT"
	INTO    = "INTO"
	VALUES  = "VALUES"
	SELECT  = "SELECT"
	FROM    = "FROM"
	WHERE   = "WHERE"
	UPDATE  = "UPDATE"
	SET     = "SET"
	DELETE  = "DELETE"
	JOIN    = "JOIN"
	ON      = "ON"
	PRIMARY = "PRIMARY"
	KEY     = "KEY"
	UNIQUE  = "UNIQUE"

	// identifiers & literals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// operators & delimiters
	ASSIGN    = "="
	ASTERISK  = "*"
	COMMA     = ","
	LPAREN    = "("
	RPAREN    = ")"
	DOT       = "."
	SEMICOLON = ";"
	GT        = ">"
	LT        = "<"

	// data type keywords
	TYPE_INT  = "TYPE_INT"
	TYPE_TEXT = "TYPE_TEXT"
	TYPE_BOOL = "TYPE_BOOL"

	// special
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

var keywords = map[string]TokenType{
	"create":  CREATE,
	"table":   TABLE,
	"insert":  INSERT,
	"into":    INTO,
	"values":  VALUES,
	"select":  SELECT,
	"from":    FROM,
	"where":   WHERE,
	"update":  UPDATE,
	"set":     SET,
	"delete":  DELETE,
	"join":    JOIN,
	"on":      ON,
	"primary": PRIMARY,
	"key":     KEY,
	"unique":  UNIQUE,
	"int":     TYPE_INT,
	"text":    TYPE_TEXT,
	"bool":    TYPE_BOOL,
}

// check if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	// convert to lowercase for case-insensitive matching
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}

	return IDENT
}
