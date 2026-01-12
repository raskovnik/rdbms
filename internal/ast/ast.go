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

// WHERE column = > < value
type WhereClause struct {
	Column   string
	Operator string
	Value    interface{}
}

// SELECT * col, col FROM table WHERE condition
type SelectStatement struct {
	Columns []string // * col, col
	Table   string
	Where   *WhereClause // nil if no clause
}

func (ss *SelectStatement) statementNode() {}
func (ss *SelectStatement) String() string {
	return "SELECT FROM " + ss.Table
}
