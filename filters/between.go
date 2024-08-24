package filters

import (
	"fmt"
	"github.com/AdamShannag/goprime/placeholder"
)

// BetweenFilter represents a filter for SQL BETWEEN conditions.
// It constructs a condition that checks if a column's value falls within a range defined by two placeholders.
type BetweenFilter uint8

// Apply creates an SQL BETWEEN condition string for a column using two placeholders.
// Parameters:
//
//	column: The name of the SQL column to filter.
//	_: The total number of arguments (ignored).
//	currentIndex: The index for the starting placeholder in the SQL query.
//	placeholder: The placeholder interface for generating the placeholder strings.
//
// Returns:
//
//	A string representing the SQL BETWEEN condition, with placeholders for the range values.
//
// Example:
//
//	If column = "age", currentIndex = 1, and placeholder.Get(n) returns ":1" and ":2",
//	the result would be: "(age BETWEEN :1 AND :2)"
func (BetweenFilter) Apply(column string, _, currentIndex int, placeholder placeholder.Placeholder) string {
	return fmt.Sprintf("(%s BETWEEN %s AND %s)", column, placeholder.Get(currentIndex), placeholder.Get(currentIndex+1))
}

// EnrichValue is a no-op method for BetweenFilter, as it does not modify the value.
// Parameters:
//
//	value: A pointer to the value to be enriched (ignored).
//
// Returns:
//
//	An error, which is always nil for this filter as it does not perform any modifications.
func (BetweenFilter) EnrichValue(*any) error { return nil }
