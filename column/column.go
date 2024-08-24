package column

import (
	"iter"
)

type Validators []Validator

type Validator interface {
	Validate(string) error
	Name() string
}

func (c Validators) Iter(column string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		for _, v := range c {
			if !yield(v.Name(), v.Validate(column)) {
				return
			}
		}
	}
}
