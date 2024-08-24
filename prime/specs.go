package prime

import (
	"github.com/AdamShannag/goprime/filter"
	"iter"
)

type Specs map[string][]filter.Spec

func (s Specs) iter() iter.Seq2[string, filter.Spec] {
	return func(yield func(string, filter.Spec) bool) {
		for k, v := range s {
			for _, f := range v {
				if !yield(k, f) {
					return
				}
			}
		}
	}
}
