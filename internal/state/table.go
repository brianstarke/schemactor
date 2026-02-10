package state

// Table represents a database table with all its properties
type Table struct {
	Name           string
	Columns        map[string]*Column
	ColumnOrder    []string
	PrimaryKey     *PrimaryKey
	ForeignKeys    []*ForeignKey
	Indexes        []*Index
	Checks         []*CheckConstraint
	Uniques        []*UniqueConstraint
	TableComment   string
	ColumnComments map[string]string
	CreatedIn      int
	DependsOn      []string
	RequiredEnums  []string
}

// NewTable creates a new empty table
func NewTable(name string) *Table {
	return &Table{
		Name:           name,
		Columns:        make(map[string]*Column),
		ColumnOrder:    []string{},
		ForeignKeys:    []*ForeignKey{},
		Indexes:        []*Index{},
		Checks:         []*CheckConstraint{},
		Uniques:        []*UniqueConstraint{},
		ColumnComments: make(map[string]string),
		DependsOn:      []string{},
		RequiredEnums:  []string{},
	}
}

// AddColumn adds a column to the table
func (t *Table) AddColumn(col *Column) {
	if _, exists := t.Columns[col.Name]; !exists {
		t.ColumnOrder = append(t.ColumnOrder, col.Name)
	}
	t.Columns[col.Name] = col
}

// DropColumn removes a column from the table
func (t *Table) DropColumn(name string) {
	delete(t.Columns, name)

	// Remove from order
	for i, colName := range t.ColumnOrder {
		if colName == name {
			t.ColumnOrder = append(t.ColumnOrder[:i], t.ColumnOrder[i+1:]...)
			break
		}
	}

	// Remove column comment if exists
	delete(t.ColumnComments, name)

	// Remove indexes that reference the dropped column
	var remainingIndexes []*Index
	for _, idx := range t.Indexes {
		if !indexReferencesColumn(idx, name) {
			remainingIndexes = append(remainingIndexes, idx)
		}
	}
	t.Indexes = remainingIndexes

	// Remove unique constraints that reference the dropped column
	var remainingUniques []*UniqueConstraint
	for _, unique := range t.Uniques {
		if !containsColumn(unique.Columns, name) {
			remainingUniques = append(remainingUniques, unique)
		}
	}
	t.Uniques = remainingUniques
}

// indexReferencesColumn checks if an index references a specific column
func indexReferencesColumn(idx *Index, colName string) bool {
	for _, col := range idx.Columns {
		if col == colName {
			return true
		}
	}
	return false
}

// containsColumn checks if a column list contains a specific column
func containsColumn(columns []string, colName string) bool {
	for _, col := range columns {
		if col == colName {
			return true
		}
	}
	return false
}

// AlterColumn updates a column's properties
func (t *Table) AlterColumn(name string, updates func(*Column)) {
	if col, exists := t.Columns[name]; exists {
		updates(col)
	}
}

// AddIndex adds an index to the table
func (t *Table) AddIndex(idx *Index) {
	t.Indexes = append(t.Indexes, idx)
}

// AddForeignKey adds a foreign key constraint
func (t *Table) AddForeignKey(fk *ForeignKey) {
	t.ForeignKeys = append(t.ForeignKeys, fk)

	// Track dependency
	if !contains(t.DependsOn, fk.ReferencedTable) {
		t.DependsOn = append(t.DependsOn, fk.ReferencedTable)
	}
}

// AddCheck adds a check constraint
func (t *Table) AddCheck(check *CheckConstraint) {
	t.Checks = append(t.Checks, check)
}

// AddUnique adds a unique constraint
func (t *Table) AddUnique(unique *UniqueConstraint) {
	t.Uniques = append(t.Uniques, unique)
}

// SetColumnComment sets a comment for a column
func (t *Table) SetColumnComment(colName, comment string) {
	t.ColumnComments[colName] = comment
}

// AddRequiredEnum adds an enum to the list of required enums
func (t *Table) AddRequiredEnum(enumName string) {
	if !contains(t.RequiredEnums, enumName) {
		t.RequiredEnums = append(t.RequiredEnums, enumName)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
