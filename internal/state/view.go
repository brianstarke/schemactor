package state

import (
	"regexp"
	"strings"
)

// View represents a database view
type View struct {
	Name       string
	Definition string
	DependsOn  []string
	Comment    string
	CreatedIn  int
	Version    int
}

// NewView creates a new view
func NewView(name string) *View {
	return &View{
		Name:      name,
		DependsOn: []string{},
	}
}

// ExtractDependencies analyzes the view definition to find table/view dependencies
func (v *View) ExtractDependencies() {
	v.DependsOn = []string{}

	// Look for FROM and JOIN clauses
	fromRe := regexp.MustCompile(`(?i)\bFROM\s+(\w+)`)
	joinRe := regexp.MustCompile(`(?i)\bJOIN\s+(\w+)`)

	// Find all FROM matches
	fromMatches := fromRe.FindAllStringSubmatch(v.Definition, -1)
	for _, match := range fromMatches {
		if len(match) >= 2 {
			tableName := match[1]
			if !contains(v.DependsOn, tableName) {
				v.DependsOn = append(v.DependsOn, tableName)
			}
		}
	}

	// Find all JOIN matches
	joinMatches := joinRe.FindAllStringSubmatch(v.Definition, -1)
	for _, match := range joinMatches {
		if len(match) >= 2 {
			tableName := match[1]
			if !contains(v.DependsOn, tableName) {
				v.DependsOn = append(v.DependsOn, tableName)
			}
		}
	}
}

// SetColumnComment sets a comment for a view column
func (v *View) SetColumnComment(colName, comment string) {
	// For views, we can store this in a structured way if needed
	// For now, we'll just track that we've seen it
}

// NormalizeDefinition normalizes the view definition
func (v *View) NormalizeDefinition() string {
	// Remove extra whitespace while preserving structure
	lines := strings.Split(v.Definition, "\n")
	var normalized []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}
	return strings.Join(normalized, "\n")
}
