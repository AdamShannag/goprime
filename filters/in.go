package filters

import (
	"fmt"
	"github.com/AdamShannag/goprime/placeholder"
	"strings"
)

// InFilter represents a filter for SQL IN conditions. It constructs a condition that
// checks if a column's value is within a set of values, using a variable number of placeholders.
type InFilter uint8

// Apply creates an SQL IN condition string for a column using the provided placeholders.
// Parameters:
//
//	column: The name of the SQL column to filter.
//	totalArguments: The total number of placeholders to include in the IN clause.
//	currentIndex: The starting index for the placeholders in the SQL query.
//	placeholder: The placeholder interface for generating the placeholder strings.
//
// Returns:
//
//	A string representing the SQL IN condition, with placeholders for the values.
//
// Example:
//
//	For column = "status", totalArguments = 3, currentIndex = 1, and placeholder.Get(n) returns ":1", ":2", ":3",
//	the result would be: "(status IN (:1, :2, :3))"
func (InFilter) Apply(column string, totalArguments int, currentIndex int, placeholder placeholder.Placeholder) string {
	placeholders := make([]string, totalArguments)
	for i := 0; i < totalArguments; i++ {
		placeholders[i] = placeholder.Get(currentIndex)
		currentIndex++
	}

	return fmt.Sprintf("(%s IN (%s))", column, strings.Join(placeholders, ","))
}

// EnrichValue is a no-op method for InFilter, as it does not modify the value.
// Parameters:
//
//	value: A pointer to the value to be enriched (ignored).
//
// Returns:
//
//	Always returns nil as no modification is performed.
func (InFilter) EnrichValue(*any) error { return nil }
