package entities

import (
	"math/big"
	"time"
)

type LoanCondition struct {
	Name         string    `json:"name"`
	InterestRate big.Float `json:"interest_rate"`
	MinAge       int       `json:"min_age"`
	MaxAge       int       `json:"max_age"`
	ModifiedDate time.Time `json:"modified_date"`
}
