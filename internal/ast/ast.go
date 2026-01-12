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

type ColumnUpdate struct {
	Column string
	Value  interface{}
}

// UPDATE table SET col = val, col = val, WHERE condition
type UpdateStatement struct {
	Table   string
	Updates []ColumnUpdate
	Where   *WhereClause
}

func (us *UpdateStatement) statementNode() {}
func (us *UpdateStatement) String() string {
	return "UPDATE " + us.Table
}

// DELETE FROM table WHERE condition
type DeleteStatement struct {
	Table string
	Where *WhereClause
}

func (ds *DeleteStatement) statementNode() {}
func (ds *DeleteStatement) String() string {
	return "DELETE FROM " + ds.Table
}

// SELECT left.col, right.col FROM left JOIN right ON left.id = right.id
type JoinStatement struct {
	LeftTable  string
	RightTable string
	LeftCols   []string // cols to select from left table
	RightCols  []string // cols to select from right table
	OnLeft     string   // left side of ON condition
	OnRight    string   // right side of ON condition
}

func (js *JoinStatement) statementNode() {}
func (js *JoinStatement) String() string {
	return "JOIN " + js.LeftTable + " " + js.RightTable
}
