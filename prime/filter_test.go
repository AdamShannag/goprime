package prime

import (
	"github.com/AdamShannag/goprime/column"
	"github.com/AdamShannag/goprime/filter"
	"github.com/AdamShannag/goprime/filters"
	"github.com/AdamShannag/goprime/placeholder"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	pf := New(placeholder.UnNumbered("?"))
	if pf == nil {
		t.Fatal("New returned nil")
	}
}

func TestNewWithFilters(t *testing.T) {
	pf := NewWithFilters(placeholder.Numbered(":"), map[filter.MatchMode]filter.Filter{})
	if pf == nil {
		t.Fatal("NewWithFilters returned nil")
	}
}

func TestNewWithFiltersAndValidators(t *testing.T) {
	pf := NewWithFiltersAndValidators(placeholder.Numbered("$"), map[filter.MatchMode]filter.Filter{}, make(column.Validators, 0))
	if pf == nil {
		t.Fatal("NewWithFiltersAndValidators returned nil")
	}
}

func TestRegisterFilter(t *testing.T) {
	pf := New(placeholder.UnNumbered("?"))
	pf.RegisterFilter("equals", filters.ValueFilter("="))

	if _, ok := pf.filters["equals"]; !ok {
		t.Fatal("equals filter not registered")
	}
}

func TestRegisterColumnValidator(t *testing.T) {
	pf := New(placeholder.UnNumbered("?"))
	pf.RegisterColumnValidator(column.AllowedValidator{"name"})

	if len(pf.columnValidators) == 0 {
		t.Fatal("No column validators registered")
	}

	if pf.columnValidators[0].Name() != "AllowedValidator" {
		t.Errorf("expected column validator name 'AllowedValidator', got %s", pf.columnValidators[0].Name())
	}
}

func TestSql(t *testing.T) {
	specs := Specs{
		"name": {
			{
				Value:     "James",
				MatchMode: "startsWith",
				Operator:  "and",
			},
			{
				Value:     "Butt",
				MatchMode: "endsWith",
				Operator:  "and",
			},
		},
	}

	pf := New(placeholder.UnNumbered("?"))
	pf.RegisterFilter(filter.STARTS_WITH, filters.NewPatternMatchFilter("LIKE", filters.POST))
	pf.RegisterFilter(filter.ENDS_WITH, filters.NewPatternMatchFilter("LIKE", filters.PRE))

	vals, condition, err := pf.Sql(specs)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(vals) != 2 {
		t.Fatalf("expected 2 values to be returned, got %d", len(vals))
	}

	expectedCondition := "((name LIKE ?) and (name LIKE ?))"
	if condition != expectedCondition {
		t.Errorf("expected condition %s, got %s", expectedCondition, condition)
	}
}

func TestValidateColumns(t *testing.T) {
	specs := Specs{
		"name": {
			{
				Value:     "James",
				MatchMode: "startsWith",
				Operator:  "and",
			},
			{
				Value:     "Butt",
				MatchMode: "endsWith",
				Operator:  "and",
			},
		},
	}

	pf := New(placeholder.UnNumbered("?"))
	pf.RegisterColumnValidator(column.AllowedValidator{"name"})

	err := pf.ValidateColumns(specs)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestValidateColumnsWithInvalidData(t *testing.T) {
	specs := Specs{
		"invalidColumn": {
			{
				Value:     "James",
				MatchMode: "startsWith",
				Operator:  "and",
			},
		},
	}

	pf := New(placeholder.UnNumbered("?"))
	pf.RegisterColumnValidator(column.AllowedValidator{"name"})

	err := pf.ValidateColumns(specs)
	if err == nil {
		t.Fatal("expected an error due to invalid column, got nil")
	}

	expectedErrMsg := "AllowedValidator: column [invalidColumn] is not allowed"
	if err.Error() != expectedErrMsg {
		t.Errorf("expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestSqlWithEmptySpecs(t *testing.T) {
	specs := Specs{}
	pf := New(placeholder.UnNumbered("?"))

	vals, condition, err := pf.Sql(specs)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(vals) != 0 {
		t.Fatalf("expected 0 values, got %d", len(vals))
	}

	if condition != "" {
		t.Errorf("expected empty condition, got %s", condition)
	}
}

func TestSqlWithNumberedPlaceholder(t *testing.T) {
	specs := Specs{
		"name": {
			{
				Value:     "Alice",
				MatchMode: "startsWith",
				Operator:  "or",
			},
			{
				Value:     "Smith",
				MatchMode: "contains",
				Operator:  "or",
			},
		},
	}

	pf := NewWithFilters(placeholder.Numbered("$"), map[filter.MatchMode]filter.Filter{
		filter.STARTS_WITH: filters.NewPatternMatchFilter("LIKE", filters.POST),
		filter.CONTAINS:    filters.NewPatternMatchFilter("LIKE", filters.AROUND),
	})

	vals, condition, err := pf.Sql(specs)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(vals) != 2 {
		t.Fatalf("expected 2 values to be returned, got %d", len(vals))
	}

	expectedCondition := "((name LIKE $1) or (name LIKE $2))"
	if condition != expectedCondition {
		t.Errorf("expected condition %s, got %s", expectedCondition, condition)
	}
}

func TestSqlWithMultiValuedValue(t *testing.T) {
	specs := Specs{
		"age": {
			{
				Value:     []any{18, 23},
				MatchMode: "between",
				Operator:  "and",
			},
		},
	}

	pf := NewWithFilters(placeholder.Numbered("$"), map[filter.MatchMode]filter.Filter{
		"between": filters.BetweenFilter(0),
	})

	vals, condition, err := pf.Sql(specs)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expectedValues := []interface{}{18, 23}
	expectedCondition := "((age BETWEEN $1 AND $2))"

	if len(vals) != len(expectedValues) {
		t.Fatalf("expected %d values, got %d", len(expectedValues), len(vals))
	}

	if condition != expectedCondition {
		t.Errorf("expected condition %s, got %s", expectedCondition, condition)
	}

	for i, v := range expectedValues {
		if vals[i] != v {
			t.Errorf("expected value %v at index %d, got %v", v, i, vals[i])
		}
	}
}

func TestSqlWithMultipleSpecs(t *testing.T) {
	specs := Specs{
		"name": {
			{
				Value:     "James",
				MatchMode: "startsWith",
				Operator:  "and",
			},
			{
				Value:     "Butt",
				MatchMode: "endsWith",
				Operator:  "or",
			},
		},
		"age": {
			{
				Value:     []any{18, 23},
				MatchMode: "between",
				Operator:  "and",
			},
		},
	}

	pf := NewWithFilters(placeholder.UnNumbered("?"), map[filter.MatchMode]filter.Filter{})

	pf.RegisterFilter("startsWith", filters.NewPatternMatchFilter("LIKE", filters.POST))
	pf.RegisterFilter("endsWith", filters.NewPatternMatchFilter("LIKE", filters.PRE))
	pf.RegisterFilter("between", filters.BetweenFilter(0))

	vals, condition, err := pf.Sql(specs)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expectedValues := []interface{}{
		"James%",
		"%Butt",
		18,
		23,
	}

	expectedConditions := []string{
		"name LIKE ?",
		"name LIKE ?",
		"age BETWEEN ? AND ?",
	}

	for _, expCond := range expectedConditions {
		if !strings.Contains(condition, expCond) {
			t.Errorf("expected condition to contain '%s', got '%s'", expCond, condition)
		}
	}

	if len(vals) != len(expectedValues) {
		t.Fatalf("expected %d values, got %d", len(expectedValues), len(vals))
	}

	for _, expVal := range expectedValues {
		found := false
		for _, v := range vals {
			if v == expVal {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected value %v not found in values %v", expVal, vals)
		}
	}
}

func TestColumnValidators(t *testing.T) {
	pf := NewWithFilters(placeholder.Numbered("$"), map[filter.MatchMode]filter.Filter{})

	allowedValidator := column.AllowedValidator{"name", "age"}
	regexValidator, err := column.NewRegexValidator(`^[a-zA-Z]+$`)
	if err != nil {
		t.Fatalf("failed to create RegexValidator: %v", err)
	}

	pf.RegisterColumnValidator(allowedValidator)
	pf.RegisterColumnValidator(regexValidator)

	specsValid := Specs{
		"name": {
			{
				Value:     "James",
				MatchMode: "startsWith",
				Operator:  "and",
			},
		},
		"age": {
			{
				Value:     "25",
				MatchMode: "equals",
				Operator:  "and",
			},
		},
	}

	specsInvalid := Specs{
		"invalidColumn": {
			{
				Value:     "Some value",
				MatchMode: "startsWith",
				Operator:  "and",
			},
		},
	}

	err = pf.ValidateColumns(specsValid)
	if err != nil {
		t.Errorf("expected nil error for valid specs, got %v", err)
	}

	err = pf.ValidateColumns(specsInvalid)
	if err == nil {
		t.Error("expected error for invalid column name, got nil")
	} else {
		t.Logf("expected error for invalid column name: %v", err)
	}
}
