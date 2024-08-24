package filters

import (
	"fmt"
	"github.com/AdamShannag/goprime/placeholder"
)

// ValueFilter represents a filter that applies a specific SQL condition to a column.
// It uses a string to define the SQL operation (e.g., "=", "<>", ">") and formats
// SQL conditions with placeholders for parameterized queries.
type ValueFilter string

// Apply creates an SQL condition string using the column name, the filter's operation,
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
//	A string representing the SQL condition, including the column, operation, and placeholder.
//
// Example:
//
//	For column = "age", currentIndex = 1, and placeholder.Get(n) returns ":1",
//	if the filter operation is "=", the result would be: "(age = :1)"
func (f ValueFilter) Apply(column string, _ int, currentIndex int, placeholder placeholder.Placeholder) string {
	return fmt.Sprintf("(%s %s %s)", column, f, placeholder.Get(currentIndex))
}

// EnrichValue is a no-op method for ValueFilter, as it does not modify the value.
// Parameters:
//
//	value: A pointer to the value to be enriched (ignored).
//
// Returns:
//
//	Always returns nil as no modification is performed.
func (ValueFilter) EnrichValue(value *any) error { return nil }

// MutatedValueFilter represents a filter that applies a SQL condition with a specific operation
// and also allows for mutation of the value before applying the filter.
// It uses an operation string (e.g., "=", "<>", ">") and a function to mutate the value.
type MutatedValueFilter struct {
	Operation    string       // SQL operation to use in the condition (e.g., "=", "<>").
	MutationFunc func(v *any) // Function to mutate the value before applying the filter.
}

// Apply creates an SQL condition string using the column name, the filter's operation,
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
//	A string representing the SQL condition, including the column, operation, and placeholder.
//
// Example:
//
//	For column = "age", currentIndex = 1, and placeholder.Get(n) returns ":1",
//	if the operation is ">", the result would be: "(age > :1)"
func (m *MutatedValueFilter) Apply(column string, _, currentIndex int, placeholder placeholder.Placeholder) string {
	return fmt.Sprintf("(%s %s %s)", column, m.Operation, placeholder.Get(currentIndex))
}

// EnrichValue applies the mutation function to the value, modifying it before using it in a filter.
// Parameters:
//
//	value: A pointer to the value to be mutated. The value is modified in place.
//
// Returns:
//
//	Always returns nil as no error is expected from the mutation function.
func (m *MutatedValueFilter) EnrichValue(value *any) error {
	m.MutationFunc(value)
	return nil
}
