package consolidator

import (
	"fmt"
	"strings"

	"github.com/brianstarke/schemactor/internal/migration"
	"github.com/brianstarke/schemactor/internal/state"
)

// Generator generates SQL from database state
type Generator struct {
	state     *state.DatabaseState
	graph     *DependencyGraph
	enumsUsed map[string]bool
}

// NewGenerator creates a new SQL generator
func NewGenerator(dbState *state.DatabaseState, graph *DependencyGraph) *Generator {
	return &Generator{
		state:     dbState,
		graph:     graph,
		enumsUsed: make(map[string]bool),
	}
}

// Generate generates consolidated migrations
func (g *Generator) Generate(orderedObjects []string) ([]*migration.ConsolidatedMigration, error) {
	var migrations []*migration.ConsolidatedMigration
	migrationNum := 1

	for _, objName := range orderedObjects {
		node, exists := g.graph.Nodes[objName]
		if !exists {
			continue
		}

		switch node.Type {
		case ObjectDomain:
			domain, exists := g.state.Domains[objName]
			if !exists {
				continue
			}
			migrations = append(migrations, &migration.ConsolidatedMigration{
				Number:  migrationNum,
				Name:    fmt.Sprintf("create-%s-domain", domain.Name),
				UpSQL:   g.GenerateDomainSQL(domain),
				DownSQL: g.GenerateDomainDownSQL(domain),
			})
			migrationNum++

		case ObjectEnum:
			// Enums are included in their first table, not as separate migrations
			continue

		case ObjectTable:
			table, exists := g.state.Tables[objName]
			if !exists {
				continue
			}

			upSQL, downSQL := g.GenerateTableMigration(table)

			migrations = append(migrations, &migration.ConsolidatedMigration{
				Number:  migrationNum,
				Name:    fmt.Sprintf("create-%s", table.Name),
				UpSQL:   upSQL,
				DownSQL: downSQL,
			})
			migrationNum++

		case ObjectView:
			view, exists := g.state.Views[objName]
			if !exists {
				continue
			}
			migrations = append(migrations, &migration.ConsolidatedMigration{
				Number:  migrationNum,
				Name:    fmt.Sprintf("create-%s-view", view.Name),
				UpSQL:   g.GenerateViewSQL(view),
				DownSQL: g.GenerateViewDownSQL(view),
			})
			migrationNum++
		}
	}

	return migrations, nil
}

// GenerateTableMigration generates both up and down SQL for a table
func (g *Generator) GenerateTableMigration(table *state.Table) (string, string) {
	var upSQL strings.Builder
	var downSQL strings.Builder

	// Generate enums first (only if not already used)
	enumSQL := g.GenerateRequiredEnums(table)
	if enumSQL != "" {
		upSQL.WriteString(enumSQL)
		upSQL.WriteString("\n")
	}

	// Generate table SQL
	upSQL.WriteString(g.GenerateTableSQL(table))

	// Generate down SQL
	downSQL.WriteString(g.GenerateTableDownSQL(table))

	return upSQL.String(), downSQL.String()
}

// GenerateRequiredEnums generates CREATE TYPE statements for enums used by table
func (g *Generator) GenerateRequiredEnums(table *state.Table) string {
	var sql strings.Builder

	for _, enumName := range table.RequiredEnums {
		if g.enumsUsed[enumName] {
			continue
		}

		enum, exists := g.state.Enums[enumName]
		if !exists {
			continue
		}

		sql.WriteString(g.GenerateEnumSQL(enum))
		sql.WriteString("\n\n")

		g.enumsUsed[enumName] = true
	}

	return strings.TrimSpace(sql.String())
}

// GenerateEnumSQL generates CREATE TYPE SQL for an enum
func (g *Generator) GenerateEnumSQL(enum *state.Enum) string {
	var sql strings.Builder

	sql.WriteString(fmt.Sprintf("DROP TYPE IF EXISTS %s;\n", enum.Name))
	sql.WriteString(fmt.Sprintf("CREATE TYPE %s AS ENUM (\n", enum.Name))

	for i, value := range enum.Values {
		sql.WriteString(fmt.Sprintf("    '%s'", value))
		if i < len(enum.Values)-1 {
			sql.WriteString(",")
		}
		sql.WriteString("\n")
	}

	sql.WriteString(");\n")

	// Add comment if exists
	if enum.TypeComment != "" {
		sql.WriteString(fmt.Sprintf("\nCOMMENT ON TYPE %s IS '%s';\n",
			enum.Name, escapeComment(enum.TypeComment)))
	}

	return sql.String()
}

// GenerateTableSQL generates CREATE TABLE SQL
func (g *Generator) GenerateTableSQL(table *state.Table) string {
	var sql strings.Builder

	sql.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", table.Name))

	// Generate column definitions
	for i, colName := range table.ColumnOrder {
		col := table.Columns[colName]
		sql.WriteString("    ")
		sql.WriteString(g.GenerateColumnDef(col))

		needsComma := i < len(table.ColumnOrder)-1 ||
			table.PrimaryKey != nil ||
			len(table.Checks) > 0 ||
			len(table.Uniques) > 0 ||
			len(table.ForeignKeys) > 0

		if needsComma {
			sql.WriteString(",")
		}
		sql.WriteString("\n")
	}

	// Add primary key
	if table.PrimaryKey != nil {
		sql.WriteString(fmt.Sprintf("    PRIMARY KEY (%s)",
			strings.Join(table.PrimaryKey.Columns, ", ")))

		needsComma := len(table.Checks) > 0 ||
			len(table.Uniques) > 0 ||
			len(table.ForeignKeys) > 0

		if needsComma {
			sql.WriteString(",")
		}
		sql.WriteString("\n")
	}

	// Add unique constraints
	for i, unique := range table.Uniques {
		sql.WriteString(fmt.Sprintf("    UNIQUE (%s)",
			strings.Join(unique.Columns, ", ")))

		needsComma := i < len(table.Uniques)-1 ||
			len(table.Checks) > 0 ||
			len(table.ForeignKeys) > 0

		if needsComma {
			sql.WriteString(",")
		}
		sql.WriteString("\n")
	}

	// Add check constraints
	for i, check := range table.Checks {
		sql.WriteString(fmt.Sprintf("    CHECK (%s)", check.Expression))

		needsComma := i < len(table.Checks)-1 || len(table.ForeignKeys) > 0

		if needsComma {
			sql.WriteString(",")
		}
		sql.WriteString("\n")
	}

	// Add foreign keys
	for i, fk := range table.ForeignKeys {
		sql.WriteString(fmt.Sprintf("    FOREIGN KEY (%s) REFERENCES %s (%s)",
			strings.Join(fk.Columns, ", "),
			fk.ReferencedTable,
			strings.Join(fk.ReferencedColumns, ", ")))

		if fk.OnDelete != "" {
			sql.WriteString(fmt.Sprintf(" ON DELETE %s", fk.OnDelete))
		}
		if fk.OnUpdate != "" {
			sql.WriteString(fmt.Sprintf(" ON UPDATE %s", fk.OnUpdate))
		}

		if i < len(table.ForeignKeys)-1 {
			sql.WriteString(",")
		}
		sql.WriteString("\n")
	}

	sql.WriteString(");\n")

	// Add indexes
	for _, idx := range table.Indexes {
		sql.WriteString("\n")
		sql.WriteString(g.GenerateIndexSQL(idx, table.Name))
	}

	// Add table comment
	if table.TableComment != "" {
		sql.WriteString(fmt.Sprintf("\nCOMMENT ON TABLE %s IS '%s';\n",
			table.Name, escapeComment(table.TableComment)))
	}

	// Add column comments
	for colName, comment := range table.ColumnComments {
		sql.WriteString(fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s';\n",
			table.Name, colName, escapeComment(comment)))
	}

	return sql.String()
}

// GenerateColumnDef generates a column definition
func (g *Generator) GenerateColumnDef(col *state.Column) string {
	var def strings.Builder

	def.WriteString(col.Name)
	def.WriteString(" ")
	def.WriteString(col.Type)

	if col.Default != "" {
		def.WriteString(" DEFAULT ")
		def.WriteString(col.Default)
	}

	if !col.Nullable {
		def.WriteString(" NOT NULL")
	}

	return def.String()
}

// GenerateIndexSQL generates CREATE INDEX SQL
func (g *Generator) GenerateIndexSQL(idx *state.Index, tableName string) string {
	var sql strings.Builder

	if idx.Unique {
		sql.WriteString("CREATE UNIQUE INDEX ")
	} else {
		sql.WriteString("CREATE INDEX ")
	}

	sql.WriteString(fmt.Sprintf("%s ON %s (%s)",
		idx.Name, tableName, strings.Join(idx.Columns, ", ")))

	if idx.Where != "" {
		sql.WriteString(fmt.Sprintf(" WHERE %s", idx.Where))
	}

	sql.WriteString(";\n")

	return sql.String()
}

// GenerateTableDownSQL generates DROP TABLE SQL
func (g *Generator) GenerateTableDownSQL(table *state.Table) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;\n", table.Name)
}

// GenerateDomainSQL generates CREATE DOMAIN SQL
func (g *Generator) GenerateDomainSQL(domain *state.Domain) string {
	var sql strings.Builder

	sql.WriteString(fmt.Sprintf("DROP DOMAIN IF EXISTS %s;\n", domain.Name))
	sql.WriteString(fmt.Sprintf("CREATE DOMAIN %s AS %s",
		domain.Name, domain.BaseType))

	if domain.Default != "" {
		sql.WriteString(fmt.Sprintf(" DEFAULT %s", domain.Default))
	}

	if domain.Constraint != "" {
		sql.WriteString(fmt.Sprintf(" CHECK (%s)", domain.Constraint))
	}

	sql.WriteString(";\n")

	if domain.Comment != "" {
		sql.WriteString(fmt.Sprintf("\nCOMMENT ON DOMAIN %s IS '%s';\n",
			domain.Name, escapeComment(domain.Comment)))
	}

	return sql.String()
}

// GenerateDomainDownSQL generates DROP DOMAIN SQL
func (g *Generator) GenerateDomainDownSQL(domain *state.Domain) string {
	return fmt.Sprintf("DROP DOMAIN IF EXISTS %s;\n", domain.Name)
}

// GenerateViewSQL generates CREATE VIEW SQL
func (g *Generator) GenerateViewSQL(view *state.View) string {
	var sql strings.Builder

	sql.WriteString(view.Definition)

	// Add semicolon if not present (it's stripped during parsing)
	if !strings.HasSuffix(strings.TrimSpace(view.Definition), ";") {
		sql.WriteString(";")
	}
	sql.WriteString("\n")

	if view.Comment != "" {
		// Extract view name from definition
		sql.WriteString(fmt.Sprintf("\nCOMMENT ON VIEW %s IS '%s';\n",
			view.Name, escapeComment(view.Comment)))
	}

	return sql.String()
}

// GenerateViewDownSQL generates DROP VIEW SQL
func (g *Generator) GenerateViewDownSQL(view *state.View) string {
	return fmt.Sprintf("DROP VIEW IF EXISTS %s CASCADE;\n", view.Name)
}

// escapeComment escapes single quotes in comments
func escapeComment(comment string) string {
	return strings.ReplaceAll(comment, "'", "''")
}
