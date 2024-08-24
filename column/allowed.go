package column

import (
	"fmt"
	"slices"
)

type AllowedValidator []string

func (d AllowedValidator) Validate(column string) error {
	if slices.Contains(d, column) {
		return nil
	}
	return fmt.Errorf("column [%s] is not allowed", column)
}

func (d AllowedValidator) Name() string {
	return "AllowedValidator"
}
