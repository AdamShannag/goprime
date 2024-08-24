package filter

import "github.com/AdamShannag/goprime/placeholder"

// Filter provides methods for constructing and modifying SQL filter conditions.
type Filter interface {
	// Apply creates an SQL condition string for the WHERE clause.
	// Parameters:
	//   column: The name of the SQL column to filter.
	//   totalArguments: The total number of arguments in the query.
	//   currentIndex: The index of the current argument.
	//   placeholder: The placeholder interface for generating the placeholder string (e.g., "?").
	// Returns:
	//   A string representing the SQL condition to append to the WHERE clause.
	Apply(column string, totalArguments int, currentIndex int, placeholder placeholder.Placeholder) string

	// EnrichValue modifies the value pointed to by `*any` and returns an error if the modification fails.
	// This method is used to preprocess or adjust the value before using it in the query.
	// For example, in a StartsWithFilter, you might append '%' to the value to use it with a LIKE condition.
	// Parameters:
	//   value: A pointer to the value to be modified. The value is adjusted in place.
	// Returns:
	//   An error if the value could not be modified; otherwise, nil.
	EnrichValue(value *any) error
}
