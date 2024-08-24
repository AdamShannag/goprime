package filter

type Condition struct {
	Column   string
	Filter   Filter
	Operator string
}

type Spec struct {
	Value     any       `json:"value"`
	MatchMode MatchMode `json:"matchMode"`
	Operator  string    `json:"operator,omitempty"`
}
