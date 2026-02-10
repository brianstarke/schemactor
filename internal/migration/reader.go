package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

// Migration represents a migration file pair
type Migration struct {
	Number  int
	Name    string
	UpPath  string
	DownPath string
}

// Reader reads migration files from a directory
type Reader struct {
	directory string
	separator string // detected separator: "_" or "-"
}

// NewReader creates a new migration reader
func NewReader(directory string) *Reader {
	return &Reader{
		directory: directory,
	}
}

// ReadMigrations reads all migration files in order
func (r *Reader) ReadMigrations() ([]*Migration, error) {
	files, err := os.ReadDir(r.directory)
	if err != nil {
		return nil, fmt.Errorf("reading directory: %w", err)
	}

	migrationMap := make(map[int]*Migration)

	// Pattern to match migration files: 0001_name.up.sql or 0001-name.up.sql
	pattern := regexp.MustCompile(`^(\d+)([_-])([^.]+)\.(up|down)\.sql$`)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matches := pattern.FindStringSubmatch(file.Name())
		if len(matches) < 5 {
			continue
		}

		number, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}

		// Detect separator from first matched file
		if r.separator == "" {
			r.separator = matches[2]
		}

		name := matches[3]
		direction := matches[4]

		migration, exists := migrationMap[number]
		if !exists {
			migration = &Migration{
				Number: number,
				Name:   name,
			}
			migrationMap[number] = migration
		}

		fullPath := filepath.Join(r.directory, file.Name())

		if direction == "up" {
			migration.UpPath = fullPath
		} else if direction == "down" {
			migration.DownPath = fullPath
		}
	}

	// Convert map to sorted slice
	var migrations []*Migration
	for _, migration := range migrationMap {
		// Only include migrations that have an up file
		if migration.UpPath != "" {
			migrations = append(migrations, migration)
		}
	}

	// Sort by number
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Number < migrations[j].Number
	})

	return migrations, nil
}

// Separator returns the detected separator pattern ("_" or "-")
// Defaults to "_" if no files were read
func (r *Reader) Separator() string {
	if r.separator == "" {
		return "_"
	}
	return r.separator
}
