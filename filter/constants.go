package filter

type MatchMode string

const (
	EQUALS              MatchMode = "equals"
	NOT_EQUALS          MatchMode = "notEquals"
	CONTAINS            MatchMode = "contains"
	NOT_CONTAINS        MatchMode = "notContains"
	STARTS_WITH         MatchMode = "startsWith"
	ENDS_WITH           MatchMode = "endsWith"
	LESS_THAN           MatchMode = "lt"
	LESS_THAN_EQUALS    MatchMode = "lte"
	GREATER_THAN        MatchMode = "gt"
	GREATER_THAN_EQUALS MatchMode = "gte"
	DATE_BEFORE         MatchMode = "dateBefore"
	DATE_AFTER          MatchMode = "dateAfter"
	DATE_IS             MatchMode = "dateIs"
	DATE_IS_NOT         MatchMode = "dateIsNot"
	IN                  MatchMode = "in"
	BETWEEN             MatchMode = "between"
)
