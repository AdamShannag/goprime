package column

import (
	"fmt"
	"testing"
)

func TestAllowedValidator_Validate(t *testing.T) {
	allowedColumns := AllowedValidator{"name", "age", "address"}

	tests := []struct {
		column   string
		expected error
	}{
		{"name", nil},
		{"age", nil},
		{"address", nil},
		{"email", fmt.Errorf("column [email] is not allowed")},
		{"phone", fmt.Errorf("column [phone] is not allowed")},
	}

	for _, test := range tests {
		t.Run(test.column, func(t *testing.T) {
			err := allowedColumns.Validate(test.column)
			if err != nil && err.Error() != test.expected.Error() {
				t.Errorf("expected error %v, got %v", test.expected, err)
			}
			if err == nil && test.expected != nil {
				t.Errorf("expected error %v, got nil", test.expected)
			}
		})
	}
}

func TestAllowedValidator_Name(t *testing.T) {
	allowedColumns := AllowedValidator{"name", "age", "address"}

	expectedName := "AllowedValidator"
	actualName := allowedColumns.Name()

	if expectedName != actualName {
		t.Errorf("expected validator name %s, got %s", expectedName, actualName)
	}
}
