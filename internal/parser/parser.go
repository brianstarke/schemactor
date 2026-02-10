package parser

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Parser handles SQL DDL parsing
type Parser struct {
	patterns map[string]*regexp.Regexp
}

// NewParser creates a new SQL parser
func NewParser() *Parser {
	return &Parser{
		patterns: compilePatterns(),
	}
}

// compilePatterns compiles all regex patterns for statement matching
func compilePatterns() map[string]*regexp.Regexp {
	return map[string]*regexp.Regexp{
		"CREATE_TABLE":    regexp.MustCompile(`(?i)^\s*CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?(\w+)`),
		"ALTER_TABLE":     regexp.MustCompile(`(?i)^\s*ALTER\s+TABLE\s+(?:IF\s+EXISTS\s+)?(\w+)`),
		"DROP_TABLE":      regexp.MustCompile(`(?i)^\s*DROP\s+TABLE\s+(?:IF\s+EXISTS\s+)?(\w+)`),
		"CREATE_TYPE":     regexp.MustCompile(`(?i)^\s*CREATE\s+TYPE\s+(\w+)\s+AS\s+ENUM`),
		"ALTER_TYPE":      regexp.MustCompile(`(?i)^\s*ALTER\s+TYPE\s+(\w+)\s+ADD\s+VALUE`),
		"DROP_TYPE":       regexp.MustCompile(`(?i)^\s*DROP\s+TYPE\s+(?:IF\s+EXISTS\s+)?(\w+)`),
		"CREATE_DOMAIN":   regexp.MustCompile(`(?i)^\s*CREATE\s+DOMAIN\s+(\w+)\s+AS`),
		"DROP_DOMAIN":     regexp.MustCompile(`(?i)^\s*DROP\s+DOMAIN\s+(?:IF\s+EXISTS\s+)?(\w+)`),
		"CREATE_VIEW":     regexp.MustCompile(`(?i)^\s*CREATE\s+(?:OR\s+REPLACE\s+)?VIEW\s+(\w+)`),
		"DROP_VIEW":       regexp.MustCompile(`(?i)^\s*DROP\s+VIEW\s+(?:IF\s+EXISTS\s+)?(\w+)`),
		"CREATE_INDEX":    regexp.MustCompile(`(?i)^\s*CREATE\s+(?:UNIQUE\s+)?INDEX\s+(\w+)\s+ON\s+(\w+)`),
		"DROP_INDEX":      regexp.MustCompile(`(?i)^\s*DROP\s+INDEX\s+(?:IF\s+EXISTS\s+)?(\w+)`),
		"COMMENT_ON":      regexp.MustCompile(`(?i)^\s*COMMENT\s+ON\s+(TABLE|COLUMN|TYPE|VIEW)\s+(\S+)`),
		"DO_BLOCK":        regexp.MustCompile(`(?i)^\s*DO\s+\$\$`),

		// ALTER TABLE operations
		// Type pattern handles: word, word(params), word precision, word with time zone, word[]
		"ADD_COLUMN":      regexp.MustCompile(`(?i)ADD\s+COLUMN\s+(\w+)\s+((?:(?:double|character|timestamp|time)\s+(?:precision|varying|with(?:\s+time)?\s+zone|without(?:\s+time)?\s+zone)|\w+)(?:\([^)]+\))?(?:\[\])?)`),
		"DROP_COLUMN":     regexp.MustCompile(`(?i)DROP\s+COLUMN\s+(?:IF\s+EXISTS\s+)?(\w+)`),
		"ALTER_COLUMN":    regexp.MustCompile(`(?i)ALTER\s+COLUMN\s+(\w+)`),
		"ALTER_COL_TYPE":  regexp.MustCompile(`(?i)ALTER\s+COLUMN\s+(\w+)\s+TYPE\s+((?:(?:double|character|timestamp|time)\s+(?:precision|varying|with(?:\s+time)?\s+zone|without(?:\s+time)?\s+zone)|\w+)(?:\([^)]+\))?(?:\[\])?)(?:\s+USING\s+(.+))?`),
		"ALTER_COL_NULL":  regexp.MustCompile(`(?i)ALTER\s+COLUMN\s+(\w+)\s+(SET|DROP)\s+NOT\s+NULL`),
	}
}

// ParseFile parses a migration file and returns statements
func (p *Parser) ParseFile(filepath string) ([]*Statement, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	return p.Parse(string(content))
}

// Parse parses SQL content and returns statements
func (p *Parser) Parse(sql string) ([]*Statement, error) {
	// Strip comments
	sql = StripComments(sql)

	// Split into statements
	rawStatements := SplitStatements(sql)

	var statements []*Statement
	for _, rawStmt := range rawStatements {
		stmt, err := p.parseStatement(rawStmt)
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	return statements, nil
}

// parseStatement parses a single SQL statement
func (p *Parser) parseStatement(sql string) (*Statement, error) {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil, nil
	}

	stmt := &Statement{
		Original: sql,
	}

	// Match statement type
	switch {
	case p.patterns["CREATE_TABLE"].MatchString(sql):
		return p.parseCreateTable(sql)
	case p.patterns["ALTER_TABLE"].MatchString(sql):
		return p.parseAlterTable(sql)
	case p.patterns["DROP_TABLE"].MatchString(sql):
		return p.parseDropTable(sql)
	case p.patterns["CREATE_TYPE"].MatchString(sql):
		return p.parseCreateType(sql)
	case p.patterns["ALTER_TYPE"].MatchString(sql):
		return p.parseAlterType(sql)
	case p.patterns["DROP_TYPE"].MatchString(sql):
		return p.parseDropType(sql)
	case p.patterns["CREATE_DOMAIN"].MatchString(sql):
		return p.parseCreateDomain(sql)
	case p.patterns["DROP_DOMAIN"].MatchString(sql):
		return p.parseDropDomain(sql)
	case p.patterns["CREATE_VIEW"].MatchString(sql):
		return p.parseCreateView(sql)
	case p.patterns["DROP_VIEW"].MatchString(sql):
		return p.parseDropView(sql)
	case p.patterns["CREATE_INDEX"].MatchString(sql):
		return p.parseCreateIndex(sql)
	case p.patterns["DROP_INDEX"].MatchString(sql):
		return p.parseDropIndex(sql)
	case p.patterns["COMMENT_ON"].MatchString(sql):
		return p.parseComment(sql)
	case p.patterns["DO_BLOCK"].MatchString(sql):
		return p.parseDoBlock(sql)
	default:
		// Unknown statement type - skip silently
		return nil, nil
	}

	return stmt, nil
}

func (p *Parser) parseCreateTable(sql string) (*Statement, error) {
	matches := p.patterns["CREATE_TABLE"].FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid CREATE TABLE: %s", sql)
	}

	tableName := matches[1]

	// Extract table definition (everything within parentheses)
	definition := ExtractParenthesesContent(sql)

	return &Statement{
		Type:       CreateTable,
		Original:   sql,
		ObjectName: tableName,
		Details: &CreateTableDetails{
			TableName:  tableName,
			Definition: definition,
		},
	}, nil
}

func (p *Parser) parseAlterTable(sql string) (*Statement, error) {
	matches := p.patterns["ALTER_TABLE"].FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid ALTER TABLE: %s", sql)
	}

	tableName := matches[1]
	operations := p.parseAlterOperations(sql)

	return &Statement{
		Type:       AlterTable,
		Original:   sql,
		ObjectName: tableName,
		Details: &AlterTableDetails{
			TableName:  tableName,
			Operations: operations,
		},
	}, nil
}

func (p *Parser) parseAlterOperations(sql string) []AlterOperation {
	var operations []AlterOperation

	// Split by ALTER TABLE tablename to get operations part
	// Use (?s) flag to make . match newlines
	re := regexp.MustCompile(`(?is)ALTER\s+TABLE\s+\w+\s+(.+)`)
	matches := re.FindStringSubmatch(sql)
	if len(matches) < 2 {
		return operations
	}

	opsText := matches[1]

	// Try to match ALTER COLUMN TYPE
	if matches := p.patterns["ALTER_COL_TYPE"].FindStringSubmatch(opsText); len(matches) >= 3 {
		operations = append(operations, AlterOperation{
			Type:       AlterColumn,
			ColumnName: matches[1],
			DataType:   matches[2],
			Details:    opsText,
		})
		return operations
	}

	// Try to match ALTER COLUMN SET/DROP NOT NULL
	if matches := p.patterns["ALTER_COL_NULL"].FindStringSubmatch(opsText); len(matches) >= 3 {
		operations = append(operations, AlterOperation{
			Type:       AlterColumn,
			ColumnName: matches[1],
			Details:    opsText,
		})
		return operations
	}

	// Match ADD COLUMN (can be multiple)
	addMatches := p.patterns["ADD_COLUMN"].FindAllStringSubmatch(opsText, -1)
	for _, match := range addMatches {
		if len(match) >= 3 {
			operations = append(operations, AlterOperation{
				Type:       AddColumn,
				ColumnName: match[1],
				DataType:   match[2],
				Details:    strings.TrimSpace(match[0]),
			})
		}
	}

	// Match DROP COLUMN (can be multiple)
	dropMatches := p.patterns["DROP_COLUMN"].FindAllStringSubmatch(opsText, -1)
	for _, match := range dropMatches {
		if len(match) >= 2 {
			operations = append(operations, AlterOperation{
				Type:       DropColumn,
				ColumnName: match[1],
				Details:    strings.TrimSpace(match[0]),
			})
		}
	}

	return operations
}

func (p *Parser) parseDropTable(sql string) (*Statement, error) {
	matches := p.patterns["DROP_TABLE"].FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid DROP TABLE: %s", sql)
	}

	return &Statement{
		Type:       DropTable,
		Original:   sql,
		ObjectName: matches[1],
	}, nil
}

func (p *Parser) parseCreateType(sql string) (*Statement, error) {
	matches := p.patterns["CREATE_TYPE"].FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid CREATE TYPE: %s", sql)
	}

	typeName := matches[1]

	// Extract enum values
	content := ExtractParenthesesContent(sql)
	values := p.parseEnumValues(content)

	return &Statement{
		Type:       CreateType,
		Original:   sql,
		ObjectName: typeName,
		Details: &CreateTypeDetails{
			TypeName: typeName,
			Values:   values,
		},
	}, nil
}

func (p *Parser) parseEnumValues(content string) []string {
	var values []string

	// Split by comma, handling quotes
	var current strings.Builder
	var inQuote bool

	for _, ch := range content {
		if ch == '\'' {
			inQuote = !inQuote
			continue
		}

		if ch == ',' && !inQuote {
			val := strings.TrimSpace(current.String())
			if val != "" {
				values = append(values, val)
			}
			current.Reset()
		} else {
			current.WriteRune(ch)
		}
	}

	// Add last value
	val := strings.TrimSpace(current.String())
	if val != "" {
		values = append(values, val)
	}

	return values
}

func (p *Parser) parseAlterType(sql string) (*Statement, error) {
	matches := p.patterns["ALTER_TYPE"].FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid ALTER TYPE: %s", sql)
	}

	typeName := matches[1]

	// Extract the new value
	re := regexp.MustCompile(`(?i)ADD\s+VALUE\s+(?:IF\s+NOT\s+EXISTS\s+)?'([^']+)'`)
	valueMatches := re.FindStringSubmatch(sql)
	var newValue string
	if len(valueMatches) >= 2 {
		newValue = valueMatches[1]
	}

	return &Statement{
		Type:       AlterType,
		Original:   sql,
		ObjectName: typeName,
		Details: &AlterTypeDetails{
			TypeName: typeName,
			NewValue: newValue,
		},
	}, nil
}

func (p *Parser) parseDropType(sql string) (*Statement, error) {
	matches := p.patterns["DROP_TYPE"].FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid DROP TYPE: %s", sql)
	}

	return &Statement{
		Type:       DropType,
		Original:   sql,
		ObjectName: matches[1],
	}, nil
}

func (p *Parser) parseCreateDomain(sql string) (*Statement, error) {
	matches := p.patterns["CREATE_DOMAIN"].FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid CREATE DOMAIN: %s", sql)
	}

	domainName := matches[1]

	// Extract base type (handles multi-word types like "character varying(3)")
	typeRe := regexp.MustCompile(`(?i)AS\s+(.+?)(?:\s+DEFAULT|\s+CHECK|\s+NOT\s+NULL|\s+NULL|\s+CONSTRAINT|\s*;|\s*$)`)
	typeMatches := typeRe.FindStringSubmatch(sql)
	var baseType string
	if len(typeMatches) >= 2 {
		baseType = strings.TrimSpace(typeMatches[1])
	}

	// Extract DEFAULT if present
	defaultRe := regexp.MustCompile(`(?i)DEFAULT\s+([^\s]+(?:\([^)]+\))?)`)
	defaultMatches := defaultRe.FindStringSubmatch(sql)
	var defaultVal string
	if len(defaultMatches) >= 2 {
		defaultVal = defaultMatches[1]
	}

	// Extract CHECK constraint if present
	checkRe := regexp.MustCompile(`(?i)CHECK\s*\(([^)]+)\)`)
	checkMatches := checkRe.FindStringSubmatch(sql)
	var constraint string
	if len(checkMatches) >= 2 {
		constraint = checkMatches[1]
	}

	return &Statement{
		Type:       CreateDomain,
		Original:   sql,
		ObjectName: domainName,
		Details: &CreateDomainDetails{
			DomainName: domainName,
			BaseType:   baseType,
			Default:    defaultVal,
			Constraint: constraint,
		},
	}, nil
}

func (p *Parser) parseDropDomain(sql string) (*Statement, error) {
	matches := p.patterns["DROP_DOMAIN"].FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid DROP DOMAIN: %s", sql)
	}

	return &Statement{
		Type:       DropDomain,
		Original:   sql,
		ObjectName: matches[1],
	}, nil
}

func (p *Parser) parseCreateView(sql string) (*Statement, error) {
	matches := p.patterns["CREATE_VIEW"].FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid CREATE VIEW: %s", sql)
	}

	viewName := matches[1]

	return &Statement{
		Type:       CreateView,
		Original:   sql,
		ObjectName: viewName,
		Details: &CreateViewDetails{
			ViewName:   viewName,
			Definition: sql,
		},
	}, nil
}

func (p *Parser) parseDropView(sql string) (*Statement, error) {
	matches := p.patterns["DROP_VIEW"].FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid DROP VIEW: %s", sql)
	}

	return &Statement{
		Type:       DropView,
		Original:   sql,
		ObjectName: matches[1],
	}, nil
}

func (p *Parser) parseCreateIndex(sql string) (*Statement, error) {
	matches := p.patterns["CREATE_INDEX"].FindStringSubmatch(sql)
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid CREATE INDEX: %s", sql)
	}

	indexName := matches[1]
	tableName := matches[2]

	// Check if UNIQUE
	unique := strings.Contains(strings.ToUpper(sql), "UNIQUE")

	// Extract columns
	columnsContent := ExtractParenthesesContent(sql)
	var columns []string
	for _, col := range strings.Split(columnsContent, ",") {
		col = strings.TrimSpace(col)
		// Remove DESC/ASC if present
		col = strings.Split(col, " ")[0]
		if col != "" {
			columns = append(columns, col)
		}
	}

	// Extract WHERE clause for partial index
	whereRe := regexp.MustCompile(`(?i)WHERE\s+(.+)`)
	whereMatches := whereRe.FindStringSubmatch(sql)
	var where string
	if len(whereMatches) >= 2 {
		where = strings.TrimSpace(whereMatches[1])
	}

	return &Statement{
		Type:       CreateIndex,
		Original:   sql,
		ObjectName: indexName,
		Details: &CreateIndexDetails{
			IndexName: indexName,
			TableName: tableName,
			Columns:   columns,
			Unique:    unique,
			Where:     where,
		},
	}, nil
}

func (p *Parser) parseDropIndex(sql string) (*Statement, error) {
	matches := p.patterns["DROP_INDEX"].FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid DROP INDEX: %s", sql)
	}

	return &Statement{
		Type:       DropIndex,
		Original:   sql,
		ObjectName: matches[1],
	}, nil
}

func (p *Parser) parseComment(sql string) (*Statement, error) {
	matches := p.patterns["COMMENT_ON"].FindStringSubmatch(sql)
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid COMMENT: %s", sql)
	}

	objectType := strings.ToUpper(matches[1])
	objectName := matches[2]

	// Extract comment text
	// Handle escaped quotes ('') in the comment string
	commentRe := regexp.MustCompile(`(?i)IS\s+'((?:[^']|'')+)'`)
	commentMatches := commentRe.FindStringSubmatch(sql)
	var comment string
	if len(commentMatches) >= 2 {
		// Replace escaped quotes ('') with single quotes (')
		comment = strings.ReplaceAll(commentMatches[1], "''", "'")
	}

	return &Statement{
		Type:       Comment,
		Original:   sql,
		ObjectName: objectName,
		Details: &CommentDetails{
			ObjectType: objectType,
			ObjectName: objectName,
			Comment:    comment,
		},
	}, nil
}

func (p *Parser) parseDoBlock(sql string) (*Statement, error) {
	// Try to extract ALTER TYPE ADD VALUE from DO block
	alterTypeRe := regexp.MustCompile(`(?i)ALTER\s+TYPE\s+(\w+)\s+ADD\s+VALUE\s+'([^']+)'`)
	matches := alterTypeRe.FindStringSubmatch(sql)

	var typeName, value string
	if len(matches) >= 3 {
		typeName = matches[1]
		value = matches[2]
	}

	return &Statement{
		Type:       DoBlock,
		Original:   sql,
		ObjectName: typeName,
		Details: &DoBlockDetails{
			Content:  sql,
			TypeName: typeName,
			Value:    value,
		},
	}, nil
}
