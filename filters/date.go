package filters

import (
	"fmt"
	"github.com/AdamShannag/goprime/placeholder"
)

type DateAfterFilter uint8

func (DateAfterFilter) Apply(column string, _, currentIndex int, placeholder placeholder.Placeholder) string {
	return fmt.Sprintf("(%s > %s)", column, placeholder.Get(currentIndex))
}

func (DateAfterFilter) EnrichValue(*any) error { return nil }

type DateBeforeFilter uint8

func (DateBeforeFilter) Apply(column string, _, currentIndex int, placeholder placeholder.Placeholder) string {
	return fmt.Sprintf("(%s < %s)", column, placeholder.Get(currentIndex))
}

func (DateBeforeFilter) EnrichValue(*any) error { return nil }

type DateIsFilter uint8

func (DateIsFilter) Apply(column string, _, currentIndex int, placeholder placeholder.Placeholder) string {
	return fmt.Sprintf("(%s = %s)", column, placeholder.Get(currentIndex))
}

func (DateIsFilter) EnrichValue(*any) error { return nil }

type DateIsNotFilter uint8

func (DateIsNotFilter) Apply(column string, _, currentIndex int, placeholder placeholder.Placeholder) string {
	return fmt.Sprintf("(%s <> %s)", column, placeholder.Get(currentIndex))
}

func (DateIsNotFilter) EnrichValue(*any) error { return nil }
