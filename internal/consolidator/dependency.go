package consolidator

import (
	"fmt"
	"regexp"
	"strings"

	"codeberg.org/brianstarke/schemactor/internal/state"
)

// ObjectType represents the type of database object
type ObjectType int

const (
	ObjectDomain ObjectType = iota
	ObjectEnum
	ObjectTable
	ObjectView
)

// DependencyNode represents a node in the dependency graph
type DependencyNode struct {
	Type      ObjectType
	Name      string
	CreatedIn int
}

// DependencyGraph represents the dependency graph
type DependencyGraph struct {
	Nodes map[string]*DependencyNode
	Edges map[string][]string // from -> [to...]
}

// NewDependencyGraph creates a new dependency graph
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		Nodes: make(map[string]*DependencyNode),
		Edges: make(map[string][]string),
	}
}

// AddNode adds a node to the graph
func (g *DependencyGraph) AddNode(objType ObjectType, name string, createdIn int) {
	g.Nodes[name] = &DependencyNode{
		Type:      objType,
		Name:      name,
		CreatedIn: createdIn,
	}
}

// AddEdge adds an edge from -> to (from depends on to)
func (g *DependencyGraph) AddEdge(from, to string) {
	if !contains(g.Edges[to], from) {
		g.Edges[to] = append(g.Edges[to], from)
	}
}

// BuildDependencyGraph builds a dependency graph from database state
func BuildDependencyGraph(dbState *state.DatabaseState) *DependencyGraph {
	graph := NewDependencyGraph()

	// Add all domains
	for name, domain := range dbState.Domains {
		graph.AddNode(ObjectDomain, name, domain.CreatedIn)
	}

	// Add all enums
	for name, enum := range dbState.Enums {
		graph.AddNode(ObjectEnum, name, enum.CreatedIn)
	}

	// Add all tables
	for name, table := range dbState.Tables {
		graph.AddNode(ObjectTable, name, table.CreatedIn)
	}

	// Add all views
	for name, view := range dbState.Views {
		graph.AddNode(ObjectView, name, view.CreatedIn)
	}

	// Build edges for tables
	for tableName, table := range dbState.Tables {
		// Table depends on foreign key references
		for _, fk := range table.ForeignKeys {
			if fk.ReferencedTable != tableName {
				graph.AddEdge(tableName, fk.ReferencedTable)
			}
		}

		// Table depends on enums used in columns
		enumDeps := findEnumDependencies(table, dbState)
		for _, enumName := range enumDeps {
			graph.AddEdge(tableName, enumName)
			// Track enum usage
			if enum, exists := dbState.Enums[enumName]; exists {
				enum.AddUsedBy(tableName)
			}
		}

		// Table depends on domains used in columns
		domainDeps := findDomainDependencies(table, dbState)
		for _, domainName := range domainDeps {
			graph.AddEdge(tableName, domainName)
		}
	}

	// Build edges for views
	for viewName, view := range dbState.Views {
		for _, dep := range view.DependsOn {
			if dep != viewName {
				// Only add edge if the dependency actually exists in our state
				if _, exists := graph.Nodes[dep]; exists {
					graph.AddEdge(viewName, dep)
				}
			}
		}
	}

	return graph
}

// findEnumDependencies finds enums used by a table
func findEnumDependencies(table *state.Table, dbState *state.DatabaseState) []string {
	var deps []string

	// Iterate over columns in order to ensure deterministic ordering
	for _, colName := range table.ColumnOrder {
		col := table.Columns[colName]
		// Check if column type is an enum
		colType := strings.TrimSpace(col.Type)

		// Remove type modifiers like NOT NULL
		colType = strings.Split(colType, " ")[0]

		if _, exists := dbState.Enums[colType]; exists {
			if !contains(deps, colType) {
				deps = append(deps, colType)
			}
		}
	}

	return deps
}

// findDomainDependencies finds domains used by a table
func findDomainDependencies(table *state.Table, dbState *state.DatabaseState) []string {
	var deps []string

	// Iterate over columns in order to ensure deterministic ordering
	for _, colName := range table.ColumnOrder {
		col := table.Columns[colName]
		// Check if column type is a domain
		colType := strings.TrimSpace(col.Type)

		// Remove type modifiers
		colType = strings.Split(colType, " ")[0]

		if _, exists := dbState.Domains[colType]; exists {
			if !contains(deps, colType) {
				deps = append(deps, colType)
			}
		}
	}

	return deps
}

// TopologicalSort performs a topological sort on the dependency graph
// Returns objects in order: dependencies first
func (g *DependencyGraph) TopologicalSort() ([]string, error) {
	inDegree := make(map[string]int)

	// Calculate in-degrees
	for node := range g.Nodes {
		inDegree[node] = 0
	}

	for _, edges := range g.Edges {
		for _, to := range edges {
			inDegree[to]++
		}
	}

	// Queue nodes with no dependencies
	var queue []string
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	// Sort queue by priority (Domains, Enums, Tables, Views)
	sortByPriority(queue, g)

	var result []string

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)

		// Reduce in-degree of neighbors
		for _, neighbor := range g.Edges[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
				sortByPriority(queue, g)
			}
		}
	}

	if len(result) != len(g.Nodes) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	return result, nil
}

// sortByPriority sorts nodes by type priority, with creation order as tie-breaker
func sortByPriority(nodes []string, graph *DependencyGraph) {
	// Simple bubble sort by priority, then by creation order
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			nodeI := graph.Nodes[nodes[i]]
			nodeJ := graph.Nodes[nodes[j]]
			if nodeI == nil || nodeJ == nil {
				continue
			}

			priI := getPriority(nodeI)
			priJ := getPriority(nodeJ)

			// Sort by priority first
			if priI > priJ {
				nodes[i], nodes[j] = nodes[j], nodes[i]
			} else if priI == priJ {
				// If priorities are equal, sort by creation order
				if nodeI.CreatedIn > nodeJ.CreatedIn {
					nodes[i], nodes[j] = nodes[j], nodes[i]
				}
			}
		}
	}
}

// getPriority returns priority value for object type (lower = higher priority)
func getPriority(node *DependencyNode) int {
	switch node.Type {
	case ObjectDomain:
		return 0
	case ObjectEnum:
		return 1
	case ObjectTable:
		return 2
	case ObjectView:
		return 3
	default:
		return 4
	}
}

// AnalyzeEnumUsage analyzes which tables use each enum
func AnalyzeEnumUsage(dbState *state.DatabaseState) {
	for tableName, table := range dbState.Tables {
		// Iterate over columns in order to ensure deterministic enum ordering
		for _, colName := range table.ColumnOrder {
			col := table.Columns[colName]
			colType := strings.TrimSpace(col.Type)
			colType = strings.Split(colType, " ")[0]

			if enum, exists := dbState.Enums[colType]; exists {
				enum.AddUsedBy(tableName)
				table.AddRequiredEnum(colType)
			}
		}
	}
}

// ExtractTypeFromColumn extracts the base type from a column type definition
func ExtractTypeFromColumn(colType string) string {
	// Remove parentheses and everything after
	re := regexp.MustCompile(`^([^\s(]+)`)
	matches := re.FindStringSubmatch(colType)
	if len(matches) >= 2 {
		return matches[1]
	}
	return colType
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
