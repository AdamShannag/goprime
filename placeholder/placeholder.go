package placeholder

import "fmt"

// Placeholder defines the interface for generating SQL placeholders in queries.
// Implementations of this interface provide different styles of placeholders,
// such as simple placeholders or numbered placeholders.
//
// Methods:
//
//	Get(int) string: Returns the placeholder string for the given index.
//	Numbered() bool: Indicates if the placeholder uses a numbered format.
type Placeholder interface {
	Get(int) string // Returns the placeholder string for a given index.
	Numbered() bool // Returns true if the placeholder format is numbered.
}

// UnNumbered represents a placeholder with a simple format, such as "?"
// which is commonly used in SQL queries.
type UnNumbered string

// Numbered represents a placeholder with a numbered format, such as "$1",
// which is often used in PostgreSQL and other SQL databases.
type Numbered string

// Get returns the placeholder string for UnNumbered.
// It always returns "?" regardless of the provided index.
// Parameters:
//
//	index: The index of the placeholder (ignored for UnNumbered).
//
// Returns:
//
//	A string representing the placeholder, which is always "?".
func (d UnNumbered) Get(int) string { return string(d) }

// Numbered returns false for UnNumbered, indicating it does not use a numbered format.
func (UnNumbered) Numbered() bool { return false }

// Get returns a formatted numbered placeholder string for Numbered.
// The placeholder format is based on the provided prefix (e.g., "$") followed by the index.
// Parameters:
//
//	n: The index to be used in the placeholder format (e.g., "$1").
//
// Returns:
//
//	A string representing the numbered placeholder (e.g., "$1").
func (np Numbered) Get(n int) string { return fmt.Sprintf("%s%d", np, n) }

// Numbered returns true for Numbered, indicating that it uses a numbered format.
func (Numbered) Numbered() bool { return true }
