package verifier

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// MigrationFile represents a migration file
type MigrationFile struct {
	Number   int
	Name     string
	FilePath string
	IsUp     bool
}

// Verifier verifies consolidated migrations
type Verifier struct {
	outputDir string
	verbose   bool
}

// NewVerifier creates a new verifier
func NewVerifier(outputDir string, verbose bool) *Verifier {
	return &Verifier{
		outputDir: outputDir,
		verbose:   verbose,
	}
}

// Verify runs the verification process
func (v *Verifier) Verify(ctx context.Context) error {
	if v.verbose {
		fmt.Println("\n🐳 Starting PostgreSQL container...")
	}

	// Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		return fmt.Errorf("failed to start postgres container: %w", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			fmt.Printf("Warning: failed to terminate container: %v\n", err)
		}
	}()

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to get connection string: %w", err)
	}

	if v.verbose {
		fmt.Println("✓ PostgreSQL container started")
		fmt.Println("\n📖 Reading migration files...")
	}

	// Read migration files
	upMigrations, downMigrations, err := v.readMigrations()
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	if v.verbose {
		fmt.Printf("✓ Found %d up migrations and %d down migrations\n", len(upMigrations), len(downMigrations))
	}

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Run up migrations
	if v.verbose {
		fmt.Println("\n⬆️  Running UP migrations...")
	}
	for _, mig := range upMigrations {
		if err := v.runMigration(db, mig); err != nil {
			return fmt.Errorf("UP migration failed (%s): %w", mig.Name, err)
		}
		if v.verbose {
			fmt.Printf("  ✓ %s\n", mig.Name)
		}
	}

	if v.verbose {
		fmt.Println("✓ All UP migrations succeeded")
	}

	// Run down migrations (in reverse order)
	if v.verbose {
		fmt.Println("\n⬇️  Running DOWN migrations...")
	}
	for i := len(downMigrations) - 1; i >= 0; i-- {
		mig := downMigrations[i]
		if err := v.runMigration(db, mig); err != nil {
			return fmt.Errorf("DOWN migration failed (%s): %w", mig.Name, err)
		}
		if v.verbose {
			fmt.Printf("  ✓ %s\n", mig.Name)
		}
	}

	if v.verbose {
		fmt.Println("✓ All DOWN migrations succeeded")
	}

	return nil
}

// readMigrations reads and sorts migration files
func (v *Verifier) readMigrations() ([]MigrationFile, []MigrationFile, error) {
	upFiles, err := filepath.Glob(filepath.Join(v.outputDir, "*.up.sql"))
	if err != nil {
		return nil, nil, err
	}

	downFiles, err := filepath.Glob(filepath.Join(v.outputDir, "*.down.sql"))
	if err != nil {
		return nil, nil, err
	}

	upMigrations := make([]MigrationFile, 0, len(upFiles))
	for _, f := range upFiles {
		mig, err := parseMigrationFile(f, true)
		if err != nil {
			return nil, nil, err
		}
		upMigrations = append(upMigrations, mig)
	}

	downMigrations := make([]MigrationFile, 0, len(downFiles))
	for _, f := range downFiles {
		mig, err := parseMigrationFile(f, false)
		if err != nil {
			return nil, nil, err
		}
		downMigrations = append(downMigrations, mig)
	}

	// Sort by number
	sort.Slice(upMigrations, func(i, j int) bool {
		return upMigrations[i].Number < upMigrations[j].Number
	})

	sort.Slice(downMigrations, func(i, j int) bool {
		return downMigrations[i].Number < downMigrations[j].Number
	})

	return upMigrations, downMigrations, nil
}

// parseMigrationFile parses a migration filename
func parseMigrationFile(path string, isUp bool) (MigrationFile, error) {
	filename := filepath.Base(path)
	var number int
	var name string

	// Parse format: NNNN-name.up.sql or NNNN_name.up.sql
	// Try both separators
	var separator string
	if _, err := fmt.Sscanf(filename, "%d-", &number); err == nil {
		separator = "-"
	} else if _, err := fmt.Sscanf(filename, "%d_", &number); err == nil {
		separator = "_"
	} else {
		return MigrationFile{}, fmt.Errorf("failed to parse migration number from %s", filename)
	}

	// Extract name (between number and .up.sql or .down.sql)
	prefix := fmt.Sprintf("%04d%s", number, separator)
	if isUp {
		name = strings.TrimSuffix(strings.TrimPrefix(filename, prefix), ".up.sql")
	} else {
		name = strings.TrimSuffix(strings.TrimPrefix(filename, prefix), ".down.sql")
	}

	return MigrationFile{
		Number:   number,
		Name:     name,
		FilePath: path,
		IsUp:     isUp,
	}, nil
}

// runMigration runs a single migration file
func (v *Verifier) runMigration(db *sql.DB, mig MigrationFile) error {
	// Read SQL file
	content, err := os.ReadFile(mig.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Execute SQL
	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	return nil
}
