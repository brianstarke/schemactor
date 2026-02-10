package consolidator

import (
	"fmt"

	"github.com/brianstarke/schemactor/internal/migration"
	"github.com/brianstarke/schemactor/internal/parser"
	"github.com/brianstarke/schemactor/internal/state"
)

// Consolidator orchestrates the migration consolidation process
type Consolidator struct {
	inputDir  string
	outputDir string
	verbose   bool
}

// NewConsolidator creates a new consolidator
func NewConsolidator(inputDir, outputDir string, verbose bool) *Consolidator {
	return &Consolidator{
		inputDir:  inputDir,
		outputDir: outputDir,
		verbose:   verbose,
	}
}

// Consolidate runs the consolidation process
func (c *Consolidator) Consolidate(dryRun bool) error {
	// Phase 1: Read migrations
	if c.verbose {
		fmt.Println("Phase 1: Reading migrations...")
	}

	reader := migration.NewReader(c.inputDir)
	migrations, err := reader.ReadMigrations()
	if err != nil {
		return fmt.Errorf("reading migrations: %w", err)
	}

	if c.verbose {
		fmt.Printf("Found %d migrations\n", len(migrations))
	}

	// Phase 2: Build cumulative state
	if c.verbose {
		fmt.Println("\nPhase 2: Building cumulative state...")
	}

	dbState := state.NewDatabaseState()
	sqlParser := parser.NewParser()
	applier := NewApplier(dbState)

	for _, mig := range migrations {
		if c.verbose {
			fmt.Printf("  Processing migration %04d: %s\n", mig.Number, mig.Name)
		}

		// Set current migration number for tracking creation order
		applier.SetCurrentMigration(mig.Number)

		statements, err := sqlParser.ParseFile(mig.UpPath)
		if err != nil {
			return fmt.Errorf("parsing migration %s: %w", mig.Name, err)
		}

		for _, stmt := range statements {
			if err := applier.Apply(stmt); err != nil {
				return fmt.Errorf("applying statement in %s: %w", mig.Name, err)
			}
		}
	}

	if c.verbose {
		fmt.Printf("\nState summary:\n")
		fmt.Printf("  Domains: %d\n", len(dbState.Domains))
		fmt.Printf("  Enums: %d\n", len(dbState.Enums))
		fmt.Printf("  Tables: %d\n", len(dbState.Tables))
		fmt.Printf("  Views: %d\n", len(dbState.Views))
	}

	// Phase 3: Analyze enum usage
	if c.verbose {
		fmt.Println("\nPhase 3: Analyzing enum usage...")
	}

	AnalyzeEnumUsage(dbState)

	// Phase 4: Build dependency graph
	if c.verbose {
		fmt.Println("\nPhase 4: Building dependency graph...")
	}

	depGraph := BuildDependencyGraph(dbState)

	// Phase 5: Topological sort
	if c.verbose {
		fmt.Println("\nPhase 5: Performing topological sort...")
		fmt.Println("Dependency edges:")
		for from, tos := range depGraph.Edges {
			for _, to := range tos {
				fmt.Printf("  %s -> %s\n", from, to)
			}
		}
	}

	orderedObjects, err := depGraph.TopologicalSort()
	if err != nil {
		return fmt.Errorf("sorting dependencies: %w", err)
	}

	if c.verbose {
		fmt.Printf("Ordered %d objects\n", len(orderedObjects))
	}

	// Phase 6: Generate consolidated migrations
	if c.verbose {
		fmt.Println("\nPhase 6: Generating consolidated migrations...")
	}

	generator := NewGenerator(dbState, depGraph)
	consolidatedMigrations, err := generator.Generate(orderedObjects)
	if err != nil {
		return fmt.Errorf("generating migrations: %w", err)
	}

	if c.verbose {
		fmt.Printf("Generated %d consolidated migrations\n", len(consolidatedMigrations))
	}

	// Phase 7: Write output
	if c.verbose {
		fmt.Println("\nPhase 7: Writing output...")
	}

	writer := migration.NewWriter(c.outputDir, reader.Separator())

	if dryRun {
		if c.verbose {
			fmt.Println("\n*** DRY RUN MODE - No files will be written ***\n")
		}
		writer.PreviewMigrations(consolidatedMigrations)
	} else {
		if err := writer.WriteMigrations(consolidatedMigrations); err != nil {
			return fmt.Errorf("writing migrations: %w", err)
		}

		if c.verbose {
			fmt.Printf("\nSuccessfully wrote %d consolidated migrations to %s\n",
				len(consolidatedMigrations), c.outputDir)
		}
	}

	return nil
}
