package parser

import (
	"strings"
	"unicode"
)

// SplitStatements splits SQL text into individual statements
// Handles semicolons within quotes and dollar-quoted strings
func SplitStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	var inSingleQuote bool
	var inDoubleQuote bool
	var inDollarQuote bool
	var dollarTag string

	runes := []rune(sql)

	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		// Handle dollar quotes ($$, $tag$, etc.)
		if ch == '$' && !inSingleQuote && !inDoubleQuote {
			// Look ahead to find matching $
			tagEnd := i + 1
			for tagEnd < len(runes) && runes[tagEnd] != '$' {
				tagEnd++
			}
			if tagEnd < len(runes) {
				tag := string(runes[i:tagEnd+1])

				if inDollarQuote {
					current.WriteRune(ch)
					if tag == dollarTag {
						// End of dollar quote
						for j := i + 1; j <= tagEnd; j++ {
							current.WriteRune(runes[j])
						}
						i = tagEnd
						inDollarQuote = false
						dollarTag = ""
						continue
					}
				} else {
					// Start of dollar quote
					inDollarQuote = true
					dollarTag = tag
					for j := i; j <= tagEnd; j++ {
						current.WriteRune(runes[j])
					}
					i = tagEnd
					continue
				}
			}
		}

		// Handle single quotes
		if ch == '\'' && !inDoubleQuote && !inDollarQuote {
			current.WriteRune(ch)
			// Check for escaped quote
			if i+1 < len(runes) && runes[i+1] == '\'' {
				current.WriteRune(runes[i+1])
				i++
			} else {
				inSingleQuote = !inSingleQuote
			}
			continue
		}

		// Handle double quotes
		if ch == '"' && !inSingleQuote && !inDollarQuote {
			current.WriteRune(ch)
			inDoubleQuote = !inDoubleQuote
			continue
		}

		// Handle statement terminator
		if ch == ';' && !inSingleQuote && !inDoubleQuote && !inDollarQuote {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
			continue
		}

		current.WriteRune(ch)
	}

	// Add final statement if any
	stmt := strings.TrimSpace(current.String())
	if stmt != "" {
		statements = append(statements, stmt)
	}

	return statements
}

// NormalizeWhitespace replaces multiple whitespace characters with single space
func NormalizeWhitespace(s string) string {
	var result strings.Builder
	var prevSpace bool

	for _, ch := range s {
		if unicode.IsSpace(ch) {
			if !prevSpace {
				result.WriteRune(' ')
				prevSpace = true
			}
		} else {
			result.WriteRune(ch)
			prevSpace = false
		}
	}

	return strings.TrimSpace(result.String())
}

// ExtractParenthesesContent extracts content within parentheses
// Handles nested parentheses
func ExtractParenthesesContent(s string) string {
	start := strings.Index(s, "(")
	if start == -1 {
		return ""
	}

	depth := 0
	for i := start; i < len(s); i++ {
		if s[i] == '(' {
			depth++
		} else if s[i] == ')' {
			depth--
			if depth == 0 {
				return s[start+1 : i]
			}
		}
	}

	return ""
}

// StripComments removes SQL comments from the input
func StripComments(sql string) string {
	var result strings.Builder
	lines := strings.Split(sql, "\n")

	for _, line := range lines {
		// Remove single-line comments (--)
		if idx := strings.Index(line, "--"); idx != -1 {
			line = line[:idx]
		}

		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result.WriteString(line)
			result.WriteRune('\n')
		}
	}

	return result.String()
}
