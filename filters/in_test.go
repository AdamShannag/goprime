package filters

import (
	"github.com/AdamShannag/goprime/placeholder"
	"testing"
)

func TestInFilter_Apply(t *testing.T) {
	mockPlaceholder := placeholder.Numbered("$")

	filter := InFilter(0)

	tests := []struct {
		column         string
		totalArguments int
		currentIndex   int
		expectedOutput string
	}{
		{"status", 3, 1, "(status IN ($1,$2,$3))"},
		{"category", 2, 10, "(category IN ($10,$11))"},
		{"type", 5, 100, "(type IN ($100,$101,$102,$103,$104))"},
	}

	for _, test := range tests {
		t.Run(test.column, func(t *testing.T) {
			output := filter.Apply(test.column, test.totalArguments, test.currentIndex, mockPlaceholder)
			if output != test.expectedOutput {
				t.Errorf("expected %s, got %s", test.expectedOutput, output)
			}
		})
	}
}

func TestInFilter_EnrichValue(t *testing.T) {
	filter := InFilter(0)

	var testValue any = "some value"

	err := filter.EnrichValue(&testValue)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if testValue != "some value" {
		t.Errorf("expected value 'some value', got %v", testValue)
	}
}
