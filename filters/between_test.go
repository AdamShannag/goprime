package filters

import (
	"github.com/AdamShannag/goprime/placeholder"
	"testing"
)

func TestBetweenFilter_Apply(t *testing.T) {
	mockPlaceholder := placeholder.Numbered("$")

	filter := BetweenFilter(0)

	tests := []struct {
		column         string
		currentIndex   int
		expectedOutput string
	}{
		{"age", 1, "(age BETWEEN $1 AND $2)"},
		{"price", 10, "(price BETWEEN $10 AND $11)"},
		{"height", 100, "(height BETWEEN $100 AND $101)"},
	}

	for _, test := range tests {
		t.Run(test.column, func(t *testing.T) {
			output := filter.Apply(test.column, 0, test.currentIndex, mockPlaceholder)
			if output != test.expectedOutput {
				t.Errorf("expected %s, got %s", test.expectedOutput, output)
			}
		})
	}
}

func TestBetweenFilter_EnrichValue(t *testing.T) {
	filter := BetweenFilter(0)

	var testValue any = "some value"

	err := filter.EnrichValue(&testValue)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if testValue != "some value" {
		t.Errorf("expected value 'some value', got %v", testValue)
	}
}
