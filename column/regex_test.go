package column

import (
	"fmt"
	"testing"
)

func TestNewRegexValidator(t *testing.T) {
	validRegex := `^[a-zA-Z]+$`
	validator, err := NewRegexValidator(validRegex)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if validator == nil {
		t.Fatal("expected validator to be non-nil")
	}

	invalidRegex := `[\`
	_, err = NewRegexValidator(invalidRegex)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestRegexValidator_Validate(t *testing.T) {
	validator, err := NewRegexValidator(`^[a-zA-Z]+$`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tests := []struct {
		column   string
		expected error
	}{
		{"John", nil},
		{"Doe", nil},
		{"John123", fmt.Errorf("column John123 does not match regex")},
		{"12345", fmt.Errorf("column 12345 does not match regex")},
	}

	for _, test := range tests {
		t.Run(test.column, func(t *testing.T) {
			err := validator.Validate(test.column)
			if err != nil && err.Error() != test.expected.Error() {
				t.Errorf("expected error %v, got %v", test.expected, err)
			}
			if err == nil && test.expected != nil {
				t.Errorf("expected error %v, got nil", test.expected)
			}
		})
	}
}

func TestRegexValidator_Name(t *testing.T) {
	validator, err := NewRegexValidator(`^[a-zA-Z]+$`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedName := "RegexValidator"
	actualName := validator.Name()

	if expectedName != actualName {
		t.Errorf("expected validator name %s, got %s", expectedName, actualName)
	}
}
