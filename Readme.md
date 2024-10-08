# `goprime`

## Overview

The `goprime` package is a Go library designed to convert JSON filter specifications into SQL WHERE clauses. It provides a flexible and secure way to generate SQL conditions from JSON input, created for JSON that is generated by PrimeNG or any other similar JSON source. `goprime` ensures security through the use of placeholders, preventing SQL injection, and offers flexibility for custom filter implementations and column validation.

## Features

- **Generate SQL Conditions**: Converts JSON filter specifications into SQL WHERE clauses.
- **Flexible Filter Registration**: Register custom filters for various match modes.
- **Column Validation**: Register validators to enforce constraints on column names.
- **Placeholder-Based Security**: Uses placeholders to safeguard against SQL injection.

## Installation

To install the `goprime` package, use the following `go get` command:

```bash
go get -u github.com/AdamShannag/goprime
```
## JSON Filter Specifications

`goprime` processes JSON filter specifications that follow the expected format. The format typically includes fields and match modes, and it can be generated by PrimeNG table filters or any other system that conforms to this structure.

### Example JSON

```json
{
  "name": [
    {"value": "James", "matchMode": "startsWith", "operator": "and"},
    {"value": "Butt", "matchMode": "endsWith", "operator": "and"}
  ],
  "country.name": [
    {"value": null, "matchMode": "startsWith", "operator": "and"}
  ],
  "representative": [
    {"value": ["Amy Elsner", "Anna Fali", "Bernardo Dominic", "Elwin Sharvill"], "matchMode": "in", "operator": "and"}
  ],
  "date": [
    {"value": "2024-08-12T21:00:00.000Z", "matchMode": "dateAfter", "operator": "and"},
    {"value": "2024-08-20T21:00:00.000Z", "matchMode": "dateBefore", "operator": "and"}
  ],
  "activity": [
    {"value": [68, 100], "matchMode": "between", "operator": "and"}
  ]
}
```

This JSON format is versatile and can be adapted to various sources as long as the structure is maintained.

## Generating SQL Conditions

`goprime` translates the JSON filter specifications into SQL conditions. For the provided JSON example, the generated SQL WHERE conditions might look like this:

```sql
(name LIKE $1 AND name LIKE $2) AND
(country.name LIKE $3) AND
(representative IN ($4, $5, $6, $7)) AND
(date >= $8 AND date <= $9) AND
(activity BETWEEN $10 AND $11)
```

## Registering Custom Filters

You can register custom filters for different match modes to control how each filter is applied in SQL conditions. Here’s how to set up custom filters:

```go
// Initialize the Filter with a numbered placeholder
pf := prime.New(placeholder.Numbered("$"))

// Register filters
pf.RegisterFilter("startsWith", filters.NewPatternMatchFilter("LIKE", filters.POST))
pf.RegisterFilter("endsWith", filters.NewPatternMatchFilter("LIKE", filters.PRE))
pf.RegisterFilter("in", filters.InFilter(0))
pf.RegisterFilter("between", filters.BetweenFilter(0))
pf.RegisterFilter("dateAfter", filters.ValueFilter(">="))
pf.RegisterFilter("dateBefore", filters.ValueFilter("<="))
```

### Implementing Custom Filters

To create your own custom filters, implement the `Filter` interface:

```go
// CustomFilter is an example of a custom filter implementation
type CustomFilter struct {
	// Custom fields and methods here
}

func (f *CustomFilter) Apply(column string, totalArguments int, currentIndex int, placeholder placeholder.Placeholder) string {
	// Implement custom SQL condition generation
	return fmt.Sprintf("(%s CUSTOM_CONDITION %s)", column, placeholder.Get(currentIndex))
}

func (f *CustomFilter) EnrichValue(value *any) error {
	// Implement custom value enrichment
	return nil
}
```

## Column Validation

Enforce constraints on column names by registering validators. This helps ensure that only valid columns are used in queries.

```go
// Register column validators
pf.RegisterColumnValidator(column.AllowedValidator{"name", "country.name", "representative", "date", "activity"})

regexValidator, err := column.NewRegexValidator(`^[a-zA-Z0-9_]+$`)
if err != nil {
	log.Fatal(err)
}
pf.RegisterColumnValidator(regexValidator)
```

### Implementing Custom Validators

To create your own custom validators, implement the `Validator` interface:

```go
// CustomValidator is an example of a custom column validator
type CustomValidator struct {
	// Custom fields and methods here
}

func (v *CustomValidator) Validate(column string) error {
	// Implement custom column validation logic
	return nil
}

func (v *CustomValidator) Name() string {
	// Return the name of the validator
	return "CustomValidator"
}
```

## Placeholder-Based Security

To prevent SQL injection, `goprime` uses placeholders in SQL conditions. This approach ensures that user inputs are securely handled in queries.

### Implementing Custom Placeholders

If you need a custom placeholder format, implement the `Placeholder` interface:

```go
// CustomPlaceholder is an example of a custom placeholder implementation
type CustomPlaceholder struct {
	// Custom fields and methods here
}

func (p *CustomPlaceholder) Get(index int) string {
	// Implement custom placeholder format
	return fmt.Sprintf("$%d", index)
}

func (p *CustomPlaceholder) Numbered() bool {
	// Indicate if the placeholder format is numbered
	return true
}
```

## Example Usage

Here’s an example of how to use `goprime` to generate SQL conditions from a filter specification read from a JSON string:

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/AdamShannag/goprime/column"
	"github.com/AdamShannag/goprime/filter"
	"github.com/AdamShannag/goprime/filters"
	"github.com/AdamShannag/goprime/placeholder"
	"github.com/AdamShannag/goprime/prime"
	"log"
)

func main() {
	// Example JSON filter specification
	filterRequest := `{
		"name": [
			{"value": "James", "matchMode": "startsWith", "operator": "and"},
			{"value": "Butt", "matchMode": "endsWith", "operator": "and"}
		],
		"country.name": [
			{"value": null, "matchMode": "startsWith", "operator": "and"}
		],
		"representative": [
			{"value": ["Amy Elsner", "Anna Fali", "Bernardo Dominic", "Elwin Sharvill"], "matchMode": "in", "operator": "and"}
		],
		"date": [
			{"value": "2024-08-12T21:00:00.000Z", "matchMode": "dateAfter", "operator": "and"},
			{"value": "2024-08-20T21:00:00.000Z", "matchMode": "dateBefore", "operator": "and"}
		],
		"activity": [
			{"value": [68, 100], "matchMode": "between", "operator": "and"}
		]
	}`

	// Unmarshal the JSON filter specification into the Specs type
	specs := prime.Specs{}
	err := json.Unmarshal([]byte(filterRequest), &specs)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Initialize the Filter with a numbered placeholder
	pf := prime.New(placeholder.Numbered("$"))

	// Register filters
	pf.RegisterFilter("startsWith", filters.NewPatternMatchFilter("LIKE", filters.POST))
	pf.RegisterFilter(filter.ENDS_WITH, filters.NewPatternMatchFilter("LIKE", filters.PRE))
	pf.RegisterFilter("in", filters.InFilter(0))
	pf.RegisterFilter("between", filters.BetweenFilter(0))
	pf.RegisterFilter(filter.DATE_AFTER, filters.ValueFilter(">="))
	pf.RegisterFilter("dateBefore", filters.ValueFilter("<="))

	// Register column validators
	pf.RegisterColumnValidator(column.AllowedValidator{"name", "country.name", "representative", "date", "activity"})

	// Validate column names
	err = pf.ValidateColumns(specs)
	if err != nil {
		log.Fatalf("Error validating columns: %v", err)
	}

	// Generate SQL condition and values
	vals, condition, err := pf.Sql(specs)
	if err != nil {
		log.Fatalf("Error generating SQL: %v", err)
	}

	// Print the results
	fmt.Println("SQL Condition:", condition)
	// SQL Condition: ((name LIKE $1) and (name LIKE $2)) and ((representative IN ($3,$4,$5,$6))) and ((date >= $7) and (date <= $8)) and ((activity BETWEEN $9 AND $10))
	fmt.Println("Values:", vals)
	// Values: [James% %Butt Amy Elsner Anna Fali Bernardo Dominic Elwin Sharvill 2024-08-12T21:00:00.000Z 2024-08-20T21:00:00.000Z 68 100]
}

```
## Example

Explore this [example](https://github.com/AdamShannag/goprime-example) showcasing a Golang server integrated with Angular and PrimeNG
