# Schemactor

A tool for consolidating SQL migration files into a minimal set representing the final database state.

## WORK IN PROGRESS / LIMITATIONS

- Assumes [migrate](https://github.com/golang-migrate/migrate) style migrations.

- Assumes Postgres.

## Overview

Schemactor analyzes a directory of migration files, tracks schema changes through multiple migrations, and outputs a consolidated set of migrations. For example, if you have 20 migration files that create and modify 4 tables, Schemactor will produce 4 clean migration files representing the final state of each table.

## Features

- **Consolidates migrations**: Combines CREATE, ALTER, and DROP operations into final schema state
- **Handles dependencies**: Automatically orders migrations based on foreign keys and type dependencies
- **Splits multi-table migrations**: Separates migrations with multiple tables into individual files
- **Preserves comments**: Maintains COMMENT ON statements for tables, columns, types, and views
- **Supports PostgreSQL DDL**:
  - Tables (CREATE/ALTER/DROP)
  - Enums/Types (CREATE TYPE, ALTER TYPE ADD VALUE)
  - Domains (CREATE DOMAIN)
  - Views (CREATE VIEW)
  - Indexes (including partial indexes)
  - Constraints (PRIMARY KEY, FOREIGN KEY, UNIQUE, CHECK)

## Installation

### From Source (Latest Release)

```bash
go install github.com/brianstarke/schemactor/cmd/schemactor@latest
```

The binary will be installed to `$GOPATH/bin/schemactor` (or `$HOME/go/bin/schemactor`).

### Build Locally

Using the Makefile (recommended):

```bash
git clone https://github.com/brianstarke/schemactor
cd schemactor
make build
```

Using Go directly (basic build):

```bash
git clone https://github.com/brianstarke/schemactor
cd schemactor
go build -o schemactor ./cmd/schemactor
```

### Installation from Local Build

```bash
make install    # Builds and installs to GOPATH/bin
```

### Version Information

Show the current version:

```bash
schemactor --version
# Output: schemactor version v0.1.0-dev-f147c76
```

## Usage

Basic usage with default directories (`./sample_migrations` → `./output`):

```bash
./schemactor
```

Specify custom input and output directories:

```bash
./schemactor <input_dir> <output_dir>
```

Verify consolidated migrations with PostgreSQL (requires Docker):

```bash
./schemactor --verify
```

This will:
1. Start a PostgreSQL container using testcontainers
2. Run all UP migrations in order
3. Run all DOWN migrations in reverse order
4. Report success or failure

Show version:

```bash
./schemactor --version
```

Show help:

```bash
./schemactor --help
```

### Example

Given this migration history:
- `0001_create-users.up.sql` - Creates users table
- `0003_add-status-to-users.up.sql` - Adds status enum column
- `0006_add-user-profile-fields.up.sql` - Adds profile fields (first_name, last_name, etc.)
- `0041_remove-phone-from-users.up.sql` - Removes phone column
- `0066_add-two-factor-auth.up.sql` - Adds 2FA fields
- ...24 total alterations to users table

Schemactor produces:
- `0001-create-users.up.sql` - Final users table with all 24 modifications consolidated
- Includes the user_status enum at the top with all values (active, inactive, suspended, deleted, banned)

## Configuration

Schemactor follows these consolidation rules:

### Domains
- **Output**: Separate migration files (e.g., `0001-create-currency-domain.up.sql`)
- Domains are ordered first due to no dependencies

### Enums/Types
- **Output**: Included at the top of the first table that uses them
- Consolidates all ALTER TYPE ADD VALUE operations
- Example: `stock_exchange` enum is included in the `stonks` table migration

### Tables
- **Output**: One migration per table
- Consolidates all CREATE TABLE and ALTER TABLE operations
- Includes indexes and constraints inline
- Properly orders based on foreign key dependencies

### Views
- **Output**: Separate migration files
- Uses the latest version (if recreated multiple times)
- Ordered after all referenced tables

## Output Format

Generated files follow the pattern: `NNNN-action-object.{up|down}.sql`

Examples:
```
0001-create-currency-domain.up.sql
0001-create-currency-domain.down.sql
0002-create-users.up.sql
0002-create-users.down.sql
0003-create-products.up.sql
0003-create-products.down.sql
0004-create-orders.up.sql
0004-create-orders.down.sql
```

## Example Output

From 73 input migrations (146 files), Schemactor generates 12 consolidated migrations (24 files):

**Input**: 73 migrations with complex schema evolution
- 10 tables with extensive modifications
- 1 domain (currency with ISO validation)
- 7 enums (some with values added over time)
- 1 view (recreated/updated)

**Output**: 12 consolidated migrations representing final schema state
- 84% reduction in migration count
- All changes properly consolidated
- Dependencies correctly ordered

## Build System

### Makefile Targets

The Makefile provides convenient targets for building and releasing:

```bash
make help        # Show all available targets
make build       # Build for current platform (development version)
make build-all   # Build for all platforms (development version)
make release     # Build release binaries (requires git tag)
make install     # Install to GOPATH/bin
make clean       # Remove build artifacts
make version     # Show version information
```

### Versioning

Schemactor uses Git-based versioning:

- **Development builds**: `v0.1.0-dev-<commit-short>` (with `-dirty` suffix if uncommitted changes)
- **Release builds**: Exact version from git tag (e.g., `v0.1.0`)

### Creating a Release

Follow these steps to create a new release:

1. **Prepare the release**:
   ```bash
   # Ensure all changes are committed and working tree is clean
   git status  # Should show no changes
   
   # Update version in documentation if needed
   # Run tests to ensure everything works
   make build-all  # Test cross-platform builds
   ```

2. **Create and push the git tag**:
   ```bash
   # Create a semantic version tag (e.g., v0.1.0, v0.2.0, v1.0.0)
   git tag v0.1.0
   
   # Push the tag to GitHub
   git push origin v0.1.0
   ```

3. **Build release binaries**:
   ```bash
   # Build release binaries for all platforms
   make release
   ```

4. **Create GitHub Release**:
   - Go to your repository on GitHub
   - Click "Releases" → "Create a new release"
   - Choose the tag you just pushed (e.g., `v0.1.0`)
   - Add release notes describing the changes
   - Upload the binaries from the `dist/` directory:
     - `schemactor-v0.1.0-linux-amd64`
     - `schemactor-v0.1.0-darwin-amd64`

5. **Verify the release**:
   ```bash
   # Test the release binaries
   ./dist/schemactor-v0.1.0-linux-amd64 --version
   # Should output: schemactor version v0.1.0
   
   # Test installation from the new release
   go install github.com/brianstarke/schemactor/cmd/schemactor@v0.1.0
   ```

### Release Checklist

Before creating a release, ensure:

- [ ] All changes are committed and working tree is clean
- [ ] Version number follows semantic versioning (MAJOR.MINOR.PATCH)
- [ ] Tests pass on all target platforms
- [ ] Documentation is updated (if needed)
- [ ] CHANGELOG is updated (if you maintain one)
- [ ] Release notes are prepared

### Release Binary Naming

Release binaries follow this pattern:
```
schemactor-VERSION-PLATFORM-ARCHITECTURE
```

Examples:
- `schemactor-v0.1.0-linux-amd64`
- `schemactor-v0.1.0-darwin-amd64`

The build system automatically creates these files in the `dist/` directory when you run `make release`.

## License

MIT
