package prime

import (
	"fmt"
	"github.com/AdamShannag/goprime/column"
	"github.com/AdamShannag/goprime/filter"
	"github.com/AdamShannag/goprime/placeholder"
	"iter"
	"strings"
)

// Filter manages a collection of SQL filters and their associated placeholders.
// It stores filters by their match modes and provides methods for creating new filters
// and registering additional filters. The placeholder used in SQL conditions is
// represented by the Placeholder interface, allowing for different placeholder styles.
type Filter struct {
	filters          map[filter.MatchMode]filter.Filter // Map of match modes to their corresponding filters
	placeholder      placeholder.Placeholder            // Placeholder interface used in SQL conditions
	columnValidators column.Validators                  // Validators for column values
}

// New creates a new Filter instance with the specified placeholder.
// It initializes the filters map as an empty map and sets up column validators as an empty slice.
// Parameters:
//
//	placeholder: An implementation of the Placeholder interface for SQL conditions.
//
// Returns:
//
//	A pointer to a newly created Filter instance.
func New(placeholder placeholder.Placeholder) *Filter {
	return &Filter{
		placeholder:      placeholder,
		filters:          make(map[filter.MatchMode]filter.Filter),
		columnValidators: make(column.Validators, 0),
	}
}

// NewWithFilters creates a new Filter instance with the specified placeholder
// and initializes it with a map of filters.
// Parameters:
//
//	placeholder: An implementation of the Placeholder interface for SQL conditions.
//	filters: A map of filter.MatchMode to filter.Filter to initialize the Filter with.
//
// Returns:
//
//	A pointer to a newly created Filter instance with the specified filters.
func NewWithFilters(placeholder placeholder.Placeholder, filters map[filter.MatchMode]filter.Filter) *Filter {
	return &Filter{
		placeholder:      placeholder,
		filters:          filters,
		columnValidators: make(column.Validators, 0),
	}
}

// NewWithFiltersAndValidators creates a new Filter instance with the specified placeholder,
// filters, and column validators.
// Parameters:
//
//	placeholder: An implementation of the Placeholder interface for SQL conditions.
//	filters: A map of filter.MatchMode to filter.Filter to initialize the Filter with.
//	validators: A slice of column validators to be used for validating column values.
//
// Returns:
//
//	A pointer to a newly created Filter instance with the specified filters and validators.
func NewWithFiltersAndValidators(placeholder placeholder.Placeholder, filters map[filter.MatchMode]filter.Filter, validators column.Validators) *Filter {
	return &Filter{
		placeholder:      placeholder,
		filters:          filters,
		columnValidators: validators,
	}
}

// RegisterFilter adds a new filter for a specific match mode to the Filter.
// Parameters:
//
//	matchMode: The match mode for which to register the filter.
//	filter: The filter to be registered for the specified match mode.
func (f *Filter) RegisterFilter(matchMode filter.MatchMode, filter filter.Filter) {
	f.filters[matchMode] = filter
}

// RegisterColumnValidator adds a new column validator to the Filter's list of validators.
// Parameters:
//
//	validator: A Validator to be applied to a column.
func (f *Filter) RegisterColumnValidator(validator column.Validator) {
	f.columnValidators = append(f.columnValidators, validator)
}

// Sql generates an SQL condition string and associated values based on the provided Specs.
// It constructs a WHERE clause and prepares the corresponding values for parameterized queries.
//
// Parameters:
//
//	specs: The Specs object that provides the specifications for generating the SQL conditions.
//	       It is iterated to extract match modes and other necessary details.
//
// Returns:
//
//	vals: A slice of values that correspond to the placeholders in the SQL condition.
//	condition: The SQL WHERE clause condition string composed of the conditions derived from specs.
//	err: An error, if any, encountered during the process.
func (f *Filter) Sql(specs Specs) (vals []any, condition string, err error) {
	for _, s := range specs.iter() {
		if _, ok := f.filters[s.MatchMode]; !ok {
			err = fmt.Errorf("match mode not registered [%s]", s.MatchMode)
			return
		}
	}

	var conditions []string
	for vc, iterErr := range f.conditionsIter(specs) {
		if iterErr != nil {
			err = iterErr
			return
		}
		conditions = append(conditions, f.buildSqlCondition(len(vc.values), len(vals)+1, vc.conditions))
		vals = append(vals, vc.values...)
	}

	condition = strings.Join(conditions, " and ")
	return
}

// ValidateColumns checks if the columns specified in the Specs are valid according to the registered column validators.
// Parameters:
//
//	specs: The Specs object that provides the specifications for generating the SQL conditions.
//
// Returns:
//
//	error: Returns an error if any of the columns are invalid or if a validation error occurs. Returns nil if all columns are valid.
func (f *Filter) ValidateColumns(specs Specs) error {
	if f.columnValidators == nil {
		return nil
	}

	for col, _ := range specs.iter() {
		for validator, err := range f.columnValidators.Iter(col) {
			if err != nil {
				return fmt.Errorf("%s: %s", validator, err.Error())
			}
		}
	}

	return nil
}

func (f *Filter) buildSqlCondition(totalValsInSpec, valIndex int, conditions []filter.Condition) string {
	sb := strings.Builder{}
	sb.WriteString("(")
	i := 0
	for opr, c := range f.sqlIter(totalValsInSpec, valIndex, conditions) {
		sb.WriteString(c)
		if i < len(conditions)-1 {
			sb.WriteString(" ")
			sb.WriteString(opr)
			sb.WriteString(" ")
		}
		i++
	}
	sb.WriteString(")")
	return sb.String()
}

type valuesConditionComposite struct {
	values     []any
	conditions []filter.Condition
}

func (f *Filter) conditionsIter(ps Specs) iter.Seq2[*valuesConditionComposite, error] {
	return func(yield func(*valuesConditionComposite, error) bool) {
		for col, specs := range ps {
			vals, condition, err := f.enrichAndExtract(col, specs)
			if len(specs) == 0 || len(condition) == 0 {
				continue
			}
			if !yield(&valuesConditionComposite{vals, condition}, err) {
				return
			}
		}
	}
}

func (f *Filter) sqlIter(totalVals, lastIndex int, conditions []filter.Condition) iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		ph := f.placeholder
		for _, condition := range conditions {
			if !yield(condition.Operator, condition.Filter.Apply(condition.Column, totalVals, lastIndex, ph)) {
				return
			}
			if f.placeholder.Numbered() {
				lastIndex++
			}
		}
	}
}

func (f *Filter) enrichAndExtract(col string, specs []filter.Spec) (values []any, conditions []filter.Condition, err error) {
	for _, spec := range specs {
		if spec.Value == nil {
			continue
		}
		condition := filter.Condition{
			Column:   col,
			Filter:   f.filters[spec.MatchMode],
			Operator: spec.Operator,
		}
		conditions = append(conditions, condition)
		err = condition.Filter.EnrichValue(&spec.Value)
		if err != nil {
			return
		}
		if _, ok := spec.Value.([]any); ok {
			values = append(values, spec.Value.([]any)...)
		} else {
			values = append(values, spec.Value)
		}
	}

	return
}
