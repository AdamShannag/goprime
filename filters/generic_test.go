package filters

import (
	"github.com/AdamShannag/goprime/placeholder"
	"testing"
)

func TestValueFilter_Apply(t *testing.T) {
	mockPlaceholder := placeholder.Numbered("$")

	filter := ValueFilter("=")

	tests := []struct {
		column         string
		currentIndex   int
		expectedOutput string
	}{
		{"age", 1, "(age = $1)"},
		{"price", 10, "(price = $10)"},
		{"height", 100, "(height = $100)"},
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

func TestValueFilter_EnrichValue(t *testing.T) {
	filter := ValueFilter("=")

	var testValue any = "some value"

	err := filter.EnrichValue(&testValue)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if testValue != "some value" {
		t.Errorf("expected value 'some value', got %v", testValue)
	}
}

func TestMutatedValueFilter_Apply(t *testing.T) {
	mockPlaceholder := placeholder.Numbered("$")

	filter := &MutatedValueFilter{
		Operation: ">",
		MutationFunc: func(v *any) {
			if str, ok := (*v).(string); ok {
				*v = str + "_mutated"
			}
		},
	}

	tests := []struct {
		column         string
		currentIndex   int
		expectedOutput string
	}{
		{"age", 1, "(age > $1)"},
		{"price", 10, "(price > $10)"},
		{"height", 100, "(height > $100)"},
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

func TestMutatedValueFilter_EnrichValue(t *testing.T) {
	filter := &MutatedValueFilter{
		Operation: ">",
		MutationFunc: func(v *any) {
			if str, ok := (*v).(string); ok {
				*v = str + "_mutated"
			}
		},
	}

	var testValue any = "original"

	err := filter.EnrichValue(&testValue)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if testValue != "original_mutated" {
		t.Errorf("expected mutated value 'original_mutated', got %v", testValue)
	}
}
