package column

import (
	"fmt"
	"regexp"
)

type RegexValidator struct {
	regex *regexp.Regexp
}

func NewRegexValidator(regex string) (*RegexValidator, error) {
	compile, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	return &RegexValidator{compile}, nil
}

func (d *RegexValidator) Validate(column string) error {
	if d.regex.MatchString(column) {
		return nil
	}
	return fmt.Errorf("column %s does not match regex", column)
}

func (d *RegexValidator) Name() string {
	return "RegexValidator"
}
