package ast

type Statement interface {
	statementNode()
	String() string
}

// INSERT INTO table VALUES (values)
type InsertStatement struct {
	Table  string
	Values []interface{}
}

func (is *InsertStatement) statementNode() {}
func (is *InsertStatement) String() string {
	return "INSERT INTO " + is.Table
}

type ColumnDef struct {
	Name       string
	Type       string // int, text, bool
	PrimaryKey bool
	Unique     bool
}

// CREATE TABLE table (col dt, col dt)
type CreateStatement struct {
	Table   string
	Columns []ColumnDef
}

func (cs *CreateStatement) statementNode() {}
func (cs *CreateStatement) String() string {
	return "CREATE TABLE " + cs.Table
}
