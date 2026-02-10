# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Schemactor is a Go tool that consolidates SQL migration files into a minimal set representing the final database state. It takes a directory of migrations (e.g., 73 files with incremental changes) and outputs consolidated migrations (e.g., 12 files with final schema state).

**Target**: PostgreSQL with [golang-migrate/migrate](https://github.com/golang-migrate/migrate) style migrations.

## Build & Run Commands

```bash
# Build
go build -o schemactor ./cmd/schemactor

# Run with defaults (./sample_migrations → ./output)
./schemactor

# Run with custom directories
./schemactor <input_dir> <output_dir>

# Run with verification (requires Docker)
./schemactor --verify
```

## Architecture

The consolidation pipeline executes 7 phases in `internal/consolidator/consolidator.go`:

1. **Read migrations** (`internal/migration/reader.go`) - Parse input directory for `.up.sql` files, detect separator pattern (`_` or `-`)
2. **Build cumulative state** - Parse each migration and apply changes to state
3. **Analyze enum usage** - Track which tables use which enums
4. **Build dependency graph** - Create graph of foreign keys and type dependencies
5. **Topological sort** - Order objects: Domains → Enums → Tables → Views
6. **Generate migrations** - Produce consolidated SQL from final state
7. **Write output** - Write numbered migration files

### Key Components

- **Parser** (`internal/parser/`) - Regex-based SQL DDL parser supporting CREATE/ALTER/DROP for tables, types, domains, views, indexes
- **State** (`internal/state/`) - In-memory database state tracking (DatabaseState holds Domains, Enums, Tables, Views, Indexes)
- **Applier** (`internal/consolidator/applier.go`) - Applies parsed statements to state, handling column additions/removals, constraint changes
- **Generator** (`internal/consolidator/generator.go`) - Produces SQL from final state
- **Dependency Graph** (`internal/consolidator/dependency.go`) - Handles foreign key ordering and type dependencies

### Statement Flow

```
SQL File → Parser → Statement → Applier → DatabaseState → Generator → Output SQL
```

The `Statement` struct (`internal/parser/statement.go`) contains operation type and type-specific `Details` (CreateTableDetails, AlterTableDetails, etc.).

### Enum Handling

Enums are embedded in the first table migration that uses them, not as separate files. The `RequiredEnums` field on Table tracks which enums to include.

## Filename Pattern

Migration files follow the pattern `NNNN{sep}name.{up|down}.sql` where `{sep}` is either `_` or `-`. The output files preserve the same separator pattern as the input files.

## Verification

The `--verify` flag uses testcontainers to spin up PostgreSQL, run all UP migrations, then all DOWN migrations in reverse order. Requires Docker.
