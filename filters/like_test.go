package filters

import (
	"github.com/AdamShannag/goprime/placeholder"
	"testing"
)

func TestPatternMatchFilter_Apply(t *testing.T) {
	tests := []struct {
		filter       *PatternMatchFilter
		column       string
		currentIndex int
		expected     string
	}{
		{NewPatternMatchFilter("LIKE", PRE), "name", 1, "(name LIKE $1)"},
		{NewPatternMatchFilter("LIKE", POST), "status", 2, "(status LIKE $2)"},
		{NewPatternMatchFilter("LIKE", AROUND), "description", 3, "(description LIKE $3)"},
	}

	mockPlaceholder := placeholder.Numbered("$")

	for _, test := range tests {
		t.Run(test.column, func(t *testing.T) {
			result := test.filter.Apply(test.column, 0, test.currentIndex, mockPlaceholder)
			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestPatternMatchFilter_EnrichValue(t *testing.T) {
	tests := []struct {
		filter      *PatternMatchFilter
		input       any
		expected    any
		expectError bool
		name        string
	}{
		{
			filter:      NewPatternMatchFilter("LIKE", PRE),
			input:       "foo",
			expected:    "%foo",
			expectError: false,
			name:        "String with PRE format",
		},
		{
			filter:      NewPatternMatchFilter("LIKE", POST),
			input:       "bar",
			expected:    "bar%",
			expectError: false,
			name:        "String with POST format",
		},
		{
			filter:      NewPatternMatchFilter("LIKE", AROUND),
			input:       "baz",
			expected:    "%baz%",
			expectError: false,
			name:        "String with AROUND format",
		},
		{
			filter:      NewPatternMatchFilter("LIKE", AROUND),
			input:       123,
			expected:    123,
			expectError: false,
			name:        "Non-string input with AROUND format",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := test.input
			err := test.filter.EnrichValue(&value)

			if (err != nil) != test.expectError {
				t.Errorf("expected error %v, got %v", test.expectError, err)
			}

			if value != test.expected {
				t.Errorf("expected %v, got %v", test.expected, value)
			}
		})
	}
}
