package state

// Column represents a table column with its metadata
type Column struct {
	Name      string
	Type      string
	Nullable  bool
	Default   string
	Generated string
	Comment   string
}

// PrimaryKey represents a primary key constraint
type PrimaryKey struct {
	Columns []string
	Name    string
}

// ForeignKey represents a foreign key constraint
type ForeignKey struct {
	Name               string
	Columns            []string
	ReferencedTable    string
	ReferencedColumns  []string
	OnDelete           string
	OnUpdate           string
}

// UniqueConstraint represents a unique constraint
type UniqueConstraint struct {
	Name    string
	Columns []string
}

// CheckConstraint represents a check constraint
type CheckConstraint struct {
	Name       string
	Expression string
}

// Index represents a table index
type Index struct {
	Name    string
	Columns []string
	Unique  bool
	Where   string
	Method  string
}
