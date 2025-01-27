package entities

type LoanCondition struct {
	Name string `json:"name"`
	InterestRate int64 `json:"interest_rate"`
	MinAge int64 `json:"min_age"`
	MaxAge int64 `json:"max_age"`
}