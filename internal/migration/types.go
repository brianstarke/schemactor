package migration

// ConsolidatedMigration represents a generated migration
type ConsolidatedMigration struct {
	Number  int
	Name    string
	UpSQL   string
	DownSQL string
}
