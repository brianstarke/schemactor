package consolidator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/brianstarke/schemactor/internal/parser"
	"github.com/brianstarke/schemactor/internal/state"
)

// Applier applies parsed statements to database state
type Applier struct {
	state            *state.DatabaseState
	currentMigration int
}

// NewApplier creates a new applier
func NewApplier(dbState *state.DatabaseState) *Applier {
	return &Applier{
		state: dbState,
	}
}

// SetCurrentMigration sets the current migration number being processed
func (a *Applier) SetCurrentMigration(migrationNumber int) {
	a.currentMigration = migrationNumber
}

// Apply applies a statement to the database state
func (a *Applier) Apply(stmt *parser.Statement) error {
	switch stmt.Type {
	case parser.CreateTable:
		return a.applyCreateTable(stmt)
	case parser.AlterTable:
		return a.applyAlterTable(stmt)
	case parser.DropTable:
		return a.applyDropTable(stmt)
	case parser.CreateType:
		return a.applyCreateType(stmt)
	case parser.AlterType:
		return a.applyAlterType(stmt)
	case parser.DropType:
		return a.applyDropType(stmt)
	case parser.CreateDomain:
		return a.applyCreateDomain(stmt)
	case parser.DropDomain:
		return a.applyDropDomain(stmt)
	case parser.CreateView:
		return a.applyCreateView(stmt)
	case parser.DropView:
		return a.applyDropView(stmt)
	case parser.CreateIndex:
		return a.applyCreateIndex(stmt)
	case parser.DropIndex:
		return a.applyDropIndex(stmt)
	case parser.Comment:
		return a.applyComment(stmt)
	case parser.DoBlock:
		return a.applyDoBlock(stmt)
	default:
		return nil
	}
}

func (a *Applier) applyCreateTable(stmt *parser.Statement) error {
	details, ok := stmt.Details.(*parser.CreateTableDetails)
	if !ok {
		return fmt.Errorf("invalid CREATE TABLE details")
	}

	table := state.NewTable(details.TableName)
	table.CreatedIn = a.currentMigration

	// Parse the table definition to extract columns, constraints, etc.
	a.parseTableDefinition(table, details.Definition)

	a.state.AddOrUpdateTable(table)

	return nil
}

func (a *Applier) parseTableDefinition(table *state.Table, definition string) {
	// Split by comma at top level (not within parentheses)
	parts := splitTopLevel(definition, ',')

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Check if it's a constraint or column
		upperPart := strings.ToUpper(part)

		if strings.HasPrefix(upperPart, "PRIMARY KEY") {
			a.parsePrimaryKey(table, part)
		} else if strings.HasPrefix(upperPart, "FOREIGN KEY") {
			a.parseForeignKey(table, part)
		} else if strings.HasPrefix(upperPart, "UNIQUE") {
			a.parseUnique(table, part)
		} else if strings.HasPrefix(upperPart, "CHECK") {
			a.parseCheck(table, part)
		} else {
			// It's a column definition
			a.parseColumnDefinition(table, part)
		}
	}
}

func (a *Applier) parseColumnDefinition(table *state.Table, def string) {
	parts := strings.Fields(def)
	if len(parts) < 2 {
		return
	}

	col := &state.Column{
		Name:     parts[0],
		Nullable: true,
	}

	// Parse type
	typeIdx := 1
	colType := parts[typeIdx]

	// Handle multi-word types (double precision, character varying, etc.)
	if typeIdx+1 < len(parts) {
		nextWord := strings.ToLower(parts[typeIdx+1])
		// Check if next word is part of a multi-word type
		if nextWord == "precision" || nextWord == "varying" {
			colType += " " + parts[typeIdx+1]
			typeIdx++
		} else if nextWord == "with" || nextWord == "without" {
			// timestamp with time zone, timestamp without time zone
			if typeIdx+2 < len(parts) && strings.ToLower(parts[typeIdx+2]) == "time" {
				if typeIdx+3 < len(parts) && strings.ToLower(parts[typeIdx+3]) == "zone" {
					colType += " " + parts[typeIdx+1] + " " + parts[typeIdx+2] + " " + parts[typeIdx+3]
					typeIdx += 3
				}
			} else if typeIdx+2 < len(parts) && strings.ToLower(parts[typeIdx+2]) == "zone" {
				// time with zone, time without zone
				colType += " " + parts[typeIdx+1] + " " + parts[typeIdx+2]
				typeIdx += 2
			}
		}
	}

	// Handle types with parameters like varchar(255)
	if typeIdx+1 < len(parts) && strings.HasPrefix(parts[typeIdx+1], "(") {
		colType += " " + parts[typeIdx+1]
		typeIdx++
	} else if strings.Contains(colType, "(") {
		// Type already includes parameters
	}
	col.Type = colType

	// Parse modifiers
	remaining := strings.Join(parts[typeIdx+1:], " ")

	// Check for NOT NULL
	if strings.Contains(strings.ToUpper(remaining), "NOT NULL") {
		col.Nullable = false
	}

	// Extract DEFAULT
	defaultRe := regexp.MustCompile(`(?i)DEFAULT\s+([^\s,]+(?:\([^)]*\))?)`)
	if matches := defaultRe.FindStringSubmatch(remaining); len(matches) >= 2 {
		col.Default = matches[1]
	}

	// Check for PRIMARY KEY inline
	if strings.Contains(strings.ToUpper(remaining), "PRIMARY KEY") {
		table.PrimaryKey = &state.PrimaryKey{
			Columns: []string{col.Name},
		}
	}

	// Check for inline REFERENCES
	referencesRe := regexp.MustCompile(`(?i)REFERENCES\s+(\w+)\s*\((\w+)\)(?:\s+ON\s+DELETE\s+(\w+(?:\s+\w+)?))?(?:\s+ON\s+UPDATE\s+(\w+(?:\s+\w+)?))?`)
	if matches := referencesRe.FindStringSubmatch(remaining); len(matches) >= 3 {
		fk := &state.ForeignKey{
			Columns:           []string{col.Name},
			ReferencedTable:   matches[1],
			ReferencedColumns: []string{matches[2]},
		}
		if len(matches) >= 4 && matches[3] != "" {
			fk.OnDelete = matches[3]
		}
		if len(matches) >= 5 && matches[4] != "" {
			fk.OnUpdate = matches[4]
		}
		table.AddForeignKey(fk)
	}

	table.AddColumn(col)
}

func (a *Applier) parsePrimaryKey(table *state.Table, def string) {
	// Extract columns from PRIMARY KEY (col1, col2, ...)
	re := regexp.MustCompile(`PRIMARY\s+KEY\s*\(([^)]+)\)`)
	matches := re.FindStringSubmatch(def)
	if len(matches) >= 2 {
		cols := strings.Split(matches[1], ",")
		var columns []string
		for _, col := range cols {
			columns = append(columns, strings.TrimSpace(col))
		}
		table.PrimaryKey = &state.PrimaryKey{
			Columns: columns,
		}
	}
}

func (a *Applier) parseForeignKey(table *state.Table, def string) {
	// FOREIGN KEY (col1, col2) REFERENCES other_table (col1, col2) ON DELETE CASCADE
	fkRe := regexp.MustCompile(`FOREIGN\s+KEY\s*\(([^)]+)\)\s+REFERENCES\s+(\w+)\s*\(([^)]+)\)`)
	matches := fkRe.FindStringSubmatch(def)
	if len(matches) >= 4 {
		cols := strings.Split(matches[1], ",")
		var columns []string
		for _, col := range cols {
			columns = append(columns, strings.TrimSpace(col))
		}

		refCols := strings.Split(matches[3], ",")
		var refColumns []string
		for _, col := range refCols {
			refColumns = append(refColumns, strings.TrimSpace(col))
		}

		fk := &state.ForeignKey{
			Columns:           columns,
			ReferencedTable:   matches[2],
			ReferencedColumns: refColumns,
		}

		// Check for ON DELETE
		onDeleteRe := regexp.MustCompile(`ON\s+DELETE\s+(\w+(?:\s+\w+)?)`)
		if onDeleteMatches := onDeleteRe.FindStringSubmatch(def); len(onDeleteMatches) >= 2 {
			fk.OnDelete = onDeleteMatches[1]
		}

		// Check for ON UPDATE
		onUpdateRe := regexp.MustCompile(`ON\s+UPDATE\s+(\w+(?:\s+\w+)?)`)
		if onUpdateMatches := onUpdateRe.FindStringSubmatch(def); len(onUpdateMatches) >= 2 {
			fk.OnUpdate = onUpdateMatches[1]
		}

		table.AddForeignKey(fk)
	}
}

func (a *Applier) parseUnique(table *state.Table, def string) {
	// UNIQUE (col1, col2)
	re := regexp.MustCompile(`UNIQUE\s*\(([^)]+)\)`)
	matches := re.FindStringSubmatch(def)
	if len(matches) >= 2 {
		cols := strings.Split(matches[1], ",")
		var columns []string
		for _, col := range cols {
			columns = append(columns, strings.TrimSpace(col))
		}
		table.AddUnique(&state.UniqueConstraint{
			Columns: columns,
		})
	}
}

func (a *Applier) parseCheck(table *state.Table, def string) {
	// CHECK (expression)
	re := regexp.MustCompile(`CHECK\s*\((.+)\)`)
	matches := re.FindStringSubmatch(def)
	if len(matches) >= 2 {
		table.AddCheck(&state.CheckConstraint{
			Expression: matches[1],
		})
	}
}

func (a *Applier) applyAlterTable(stmt *parser.Statement) error {
	details, ok := stmt.Details.(*parser.AlterTableDetails)
	if !ok {
		return fmt.Errorf("invalid ALTER TABLE details")
	}

	table, exists := a.state.GetTable(details.TableName)
	if !exists {
		// Table doesn't exist yet - create it
		table = state.NewTable(details.TableName)
		a.state.AddOrUpdateTable(table)
	}

	for _, op := range details.Operations {
		switch op.Type {
		case parser.AddColumn:
			a.applyAddColumn(table, op)
		case parser.DropColumn:
			a.applyDropColumn(table, op)
		case parser.AlterColumn:
			a.applyAlterColumn(table, op)
		}
	}

	return nil
}

func (a *Applier) applyAddColumn(table *state.Table, op parser.AlterOperation) {
	col := &state.Column{
		Name:     op.ColumnName,
		Type:     op.DataType,
		Nullable: true,
	}

	// Parse the full details for additional info
	details := op.Details

	// Check for NOT NULL
	if strings.Contains(strings.ToUpper(details), "NOT NULL") {
		col.Nullable = false
	}

	// Extract DEFAULT
	defaultRe := regexp.MustCompile(`(?i)DEFAULT\s+([^\s,]+(?:\([^)]*\))?)`)
	if matches := defaultRe.FindStringSubmatch(details); len(matches) >= 2 {
		col.Default = matches[1]
	}

	table.AddColumn(col)
}

func (a *Applier) applyDropColumn(table *state.Table, op parser.AlterOperation) {
	// Get indexes that will be removed before dropping the column
	for _, idx := range table.Indexes {
		for _, col := range idx.Columns {
			if col == op.ColumnName {
				// Remove from global index tracking
				a.state.DropIndex(idx.Name)
				break
			}
		}
	}

	table.DropColumn(op.ColumnName)
}

func (a *Applier) applyAlterColumn(table *state.Table, op parser.AlterOperation) {
	// Handle ALTER COLUMN TYPE
	if op.DataType != "" {
		table.AlterColumn(op.ColumnName, func(col *state.Column) {
			col.Type = op.DataType
		})
	}

	// Handle SET/DROP NOT NULL
	if strings.Contains(strings.ToUpper(op.Details), "SET NOT NULL") {
		table.AlterColumn(op.ColumnName, func(col *state.Column) {
			col.Nullable = false
		})
	} else if strings.Contains(strings.ToUpper(op.Details), "DROP NOT NULL") {
		table.AlterColumn(op.ColumnName, func(col *state.Column) {
			col.Nullable = true
		})
	}
}

func (a *Applier) applyDropTable(stmt *parser.Statement) error {
	a.state.DropTable(stmt.ObjectName)
	return nil
}

func (a *Applier) applyCreateType(stmt *parser.Statement) error {
	details, ok := stmt.Details.(*parser.CreateTypeDetails)
	if !ok {
		return fmt.Errorf("invalid CREATE TYPE details")
	}

	enum := state.NewEnum(details.TypeName)
	enum.CreatedIn = a.currentMigration
	for _, value := range details.Values {
		enum.AddValue(value)
	}

	a.state.AddOrUpdateEnum(enum)

	return nil
}

func (a *Applier) applyAlterType(stmt *parser.Statement) error {
	details, ok := stmt.Details.(*parser.AlterTypeDetails)
	if !ok {
		return fmt.Errorf("invalid ALTER TYPE details")
	}

	enum, exists := a.state.GetEnum(details.TypeName)
	if !exists {
		enum = state.NewEnum(details.TypeName)
		enum.CreatedIn = a.currentMigration
		a.state.AddOrUpdateEnum(enum)
	}

	if details.NewValue != "" {
		enum.AddValue(details.NewValue)
	}

	return nil
}

func (a *Applier) applyDropType(stmt *parser.Statement) error {
	a.state.DropEnum(stmt.ObjectName)
	return nil
}

func (a *Applier) applyCreateDomain(stmt *parser.Statement) error {
	details, ok := stmt.Details.(*parser.CreateDomainDetails)
	if !ok {
		return fmt.Errorf("invalid CREATE DOMAIN details")
	}

	domain := state.NewDomain(details.DomainName)
	domain.CreatedIn = a.currentMigration
	domain.BaseType = details.BaseType
	domain.Default = details.Default
	domain.Constraint = details.Constraint

	a.state.AddOrUpdateDomain(domain)

	return nil
}

func (a *Applier) applyDropDomain(stmt *parser.Statement) error {
	a.state.DropDomain(stmt.ObjectName)
	return nil
}

func (a *Applier) applyCreateView(stmt *parser.Statement) error {
	details, ok := stmt.Details.(*parser.CreateViewDetails)
	if !ok {
		return fmt.Errorf("invalid CREATE VIEW details")
	}

	view := state.NewView(details.ViewName)
	view.CreatedIn = a.currentMigration
	view.Definition = details.Definition
	view.ExtractDependencies()

	a.state.AddOrUpdateView(view)

	return nil
}

func (a *Applier) applyDropView(stmt *parser.Statement) error {
	a.state.DropView(stmt.ObjectName)
	return nil
}

func (a *Applier) applyCreateIndex(stmt *parser.Statement) error {
	details, ok := stmt.Details.(*parser.CreateIndexDetails)
	if !ok {
		return fmt.Errorf("invalid CREATE INDEX details")
	}

	idx := &state.Index{
		Name:    details.IndexName,
		Columns: details.Columns,
		Unique:  details.Unique,
		Where:   details.Where,
	}

	// Add to global index tracking
	a.state.AddIndex(idx)

	// Also add to the table
	if table, exists := a.state.GetTable(details.TableName); exists {
		table.AddIndex(idx)
	}

	return nil
}

func (a *Applier) applyDropIndex(stmt *parser.Statement) error {
	a.state.DropIndex(stmt.ObjectName)
	return nil
}

func (a *Applier) applyComment(stmt *parser.Statement) error {
	details, ok := stmt.Details.(*parser.CommentDetails)
	if !ok {
		return fmt.Errorf("invalid COMMENT details")
	}

	switch details.ObjectType {
	case "TABLE":
		if table, exists := a.state.GetTable(details.ObjectName); exists {
			table.TableComment = details.Comment
		}
	case "COLUMN":
		// Parse table.column format
		parts := strings.Split(details.ObjectName, ".")
		if len(parts) == 2 {
			tableName := parts[0]
			colName := parts[1]
			if table, exists := a.state.GetTable(tableName); exists {
				table.SetColumnComment(colName, details.Comment)
			}
		}
	case "TYPE":
		if enum, exists := a.state.GetEnum(details.ObjectName); exists {
			enum.TypeComment = details.Comment
		}
	case "VIEW":
		if view, exists := a.state.GetView(details.ObjectName); exists {
			view.Comment = details.Comment
		}
	}

	return nil
}

func (a *Applier) applyDoBlock(stmt *parser.Statement) error {
	details, ok := stmt.Details.(*parser.DoBlockDetails)
	if !ok {
		return fmt.Errorf("invalid DO BLOCK details")
	}

	// If it contains ALTER TYPE ADD VALUE, apply it
	if details.TypeName != "" && details.Value != "" {
		enum, exists := a.state.GetEnum(details.TypeName)
		if !exists {
			enum = state.NewEnum(details.TypeName)
			enum.CreatedIn = a.currentMigration
			a.state.AddOrUpdateEnum(enum)
		}
		enum.AddValue(details.Value)
	}

	return nil
}

// splitTopLevel splits a string by delimiter at top level (not within parentheses)
func splitTopLevel(s string, delim rune) []string {
	var parts []string
	var current strings.Builder
	depth := 0

	for _, ch := range s {
		if ch == '(' {
			depth++
			current.WriteRune(ch)
		} else if ch == ')' {
			depth--
			current.WriteRune(ch)
		} else if ch == delim && depth == 0 {
			parts = append(parts, current.String())
			current.Reset()
		} else {
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}
