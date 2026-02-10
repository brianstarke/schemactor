package migration

import (
	"fmt"
	"os"
	"path/filepath"
)

// Writer writes consolidated migrations to disk
type Writer struct {
	outputDir string
	separator string
}

// NewWriter creates a new migration writer
func NewWriter(outputDir string, separator string) *Writer {
	return &Writer{
		outputDir: outputDir,
		separator: separator,
	}
}

// WriteMigrations writes all consolidated migrations to the output directory
func (w *Writer) WriteMigrations(migrations []*ConsolidatedMigration) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	for _, migration := range migrations {
		// Write up migration
		upPath := filepath.Join(w.outputDir,
			fmt.Sprintf("%04d%s%s.up.sql", migration.Number, w.separator, migration.Name))

		if err := os.WriteFile(upPath, []byte(migration.UpSQL), 0644); err != nil {
			return fmt.Errorf("writing up migration %s: %w", upPath, err)
		}

		// Write down migration
		downPath := filepath.Join(w.outputDir,
			fmt.Sprintf("%04d%s%s.down.sql", migration.Number, w.separator, migration.Name))

		if err := os.WriteFile(downPath, []byte(migration.DownSQL), 0644); err != nil {
			return fmt.Errorf("writing down migration %s: %w", downPath, err)
		}
	}

	return nil
}

// PreviewMigrations prints migrations to stdout for dry-run
func (w *Writer) PreviewMigrations(migrations []*ConsolidatedMigration) {
	for _, migration := range migrations {
		fmt.Printf("\n========================================\n")
		fmt.Printf("Migration %04d: %s\n", migration.Number, migration.Name)
		fmt.Printf("========================================\n\n")
		fmt.Printf("--- UP ---\n%s\n", migration.UpSQL)
		fmt.Printf("--- DOWN ---\n%s\n", migration.DownSQL)
	}
}
