package parser

// StatementType represents the type of SQL DDL statement
type StatementType int

const (
	Unknown StatementType = iota
	CreateTable
	AlterTable
	DropTable
	CreateType
	AlterType
	DropType
	CreateDomain
	DropDomain
	CreateView
	DropView
	CreateIndex
	DropIndex
	Comment
	DoBlock
)

func (st StatementType) String() string {
	switch st {
	case CreateTable:
		return "CREATE TABLE"
	case AlterTable:
		return "ALTER TABLE"
	case DropTable:
		return "DROP TABLE"
	case CreateType:
		return "CREATE TYPE"
	case AlterType:
		return "ALTER TYPE"
	case DropType:
		return "DROP TYPE"
	case CreateDomain:
		return "CREATE DOMAIN"
	case DropDomain:
		return "DROP DOMAIN"
	case CreateView:
		return "CREATE VIEW"
	case DropView:
		return "DROP VIEW"
	case CreateIndex:
		return "CREATE INDEX"
	case DropIndex:
		return "DROP INDEX"
	case Comment:
		return "COMMENT"
	case DoBlock:
		return "DO BLOCK"
	default:
		return "UNKNOWN"
	}
}

// AlterTableOperation represents specific ALTER TABLE operations
type AlterTableOperation int

const (
	AddColumn AlterTableOperation = iota
	DropColumn
	AlterColumn
	AddConstraint
	DropConstraint
)

// Statement represents a parsed SQL DDL statement
type Statement struct {
	Type       StatementType
	Original   string
	ObjectName string
	Details    interface{}
}

// CreateTableDetails contains details for CREATE TABLE statements
type CreateTableDetails struct {
	TableName  string
	Definition string // Full table definition including columns and constraints
}

// AlterTableDetails contains details for ALTER TABLE statements
type AlterTableDetails struct {
	TableName  string
	Operations []AlterOperation
}

// AlterOperation represents a single operation within an ALTER TABLE statement
type AlterOperation struct {
	Type       AlterTableOperation
	ColumnName string
	DataType   string
	Details    string // Full operation text for complex operations
}

// CreateTypeDetails contains details for CREATE TYPE (enum) statements
type CreateTypeDetails struct {
	TypeName string
	Values   []string
}

// AlterTypeDetails contains details for ALTER TYPE statements
type AlterTypeDetails struct {
	TypeName string
	NewValue string
}

// CreateDomainDetails contains details for CREATE DOMAIN statements
type CreateDomainDetails struct {
	DomainName string
	BaseType   string
	Default    string
	Constraint string
}

// CreateViewDetails contains details for CREATE VIEW statements
type CreateViewDetails struct {
	ViewName   string
	Definition string // Full view SQL
}

// CreateIndexDetails contains details for CREATE INDEX statements
type CreateIndexDetails struct {
	IndexName string
	TableName string
	Columns   []string
	Unique    bool
	Where     string // Partial index WHERE clause
}

// CommentDetails contains details for COMMENT ON statements
type CommentDetails struct {
	ObjectType string // TABLE, COLUMN, TYPE, VIEW
	ObjectName string
	Comment    string
}

// DoBlockDetails contains details for DO $$ blocks
type DoBlockDetails struct {
	Content string // Full block content
	TypeName string // If it's an ALTER TYPE, the type name
	Value    string // If it's an ALTER TYPE ADD VALUE, the value
}
