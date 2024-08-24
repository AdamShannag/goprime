package filters

import (
	"fmt"
	"github.com/AdamShannag/goprime/placeholder"
)

const (
	// PRE is a format string for matching values that start with a specific pattern.
	// For example, applying PRE to "bar" results in "%bar".
	PRE = "%%%s"

	// POST is a format string for matching values that end with a specific pattern.
	// For example, applying POST to "bar" results in "bar%".
	POST = "%s%%"

	// AROUND is a format string for matching values that contain a specific pattern.
	// For example, applying AROUND to "bar" results in "%bar%".
	AROUND = "%%%s%%"
)

// PatternMatchFilter represents a filter for SQL pattern matching conditions.
// It applies operations like `LIKE` to a column, using a format to build the pattern for matching.
type PatternMatchFilter struct {
	Operation string // SQL operation (e.g., "LIKE") used in the condition.
	Format    string // Format string to build the pattern for matching.
}

// NewPatternMatchFilter creates a new PatternMatchFilter with the specified operation and format.
// Parameters:
//
//	operation: The SQL operation for the condition (e.g., "LIKE").
//	format: The format string to build the pattern (e.g., "%%%s%%").
//
// Returns:
//
//	A pointer to a new PatternMatchFilter instance.
func NewPatternMatchFilter(operation, format string) *PatternMatchFilter {
	return &PatternMatchFilter{Operation: operation, Format: format}
}

// Apply constructs an SQL condition string using the column name, the filter's operation,
// and the placeholder for the current index.
// Parameters:
//
//	column: The name of the SQL column to filter.
//	_: The total number of arguments (ignored).
//	currentIndex: The index of the current placeholder in the SQL query.
//	placeholder: The placeholder interface for generating the placeholder string.
//
// Returns:
//
//	A string representing the SQL condition with the column, operation, and placeholder.
//
// Example:
//
//	For column = "name", operation = "LIKE", currentIndex = 1, and placeholder.Get(n) returns ":1",
//	the result would be: "(name LIKE :1)"
func (f *PatternMatchFilter) Apply(column string, _, currentIndex int, placeholder placeholder.Placeholder) string {
	return fmt.Sprintf("(%s %s %s)", column, f.Operation, placeholder.Get(currentIndex))
}

// EnrichValue modifies the value by applying the filter's format string if the value is a string.
// Parameters:
//
//	v: A pointer to the value to be enriched. If the value is a string, it is formatted using the filter's Format.
//
// Returns:
//
//	Always returns nil, as the function does not produce an error.
//
// Behavior:
//   - If the value is a string, it formats the value using the filter's format string (e.g., "foo" becomes "%%%s%%" if Format = "%%%s%%").
//   - If the value is not a string, it remains unchanged.
func (f *PatternMatchFilter) EnrichValue(value *any) error {
	if str, ok := (*value).(string); ok {
		*value = fmt.Sprintf(f.Format, str)
	}
	return nil
}
