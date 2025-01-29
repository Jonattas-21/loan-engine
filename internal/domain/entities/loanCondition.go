package entities

import (
	"time"
)

type LoanCondition struct {
	Name         string    `json:"name"`
	InterestRate float64   `json:"interest_rate"`
	MinAge       int       `json:"min_age"`
	MaxAge       int       `json:"max_age"`
	ModifiedDate time.Time `json:"modified_date"`
}
