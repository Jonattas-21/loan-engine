package entities

import (
	"math/big"
	"time"
)

type LoanSimulation struct {
	AmountTobePaid    big.Float     `json:"loan_amount"`
	AmountFeeTobePaid big.Float     `json:"fee_amount"`
	SimulationDate    time.Time     `json:"simulation_date"`
	Currency          string        `json:"currency"`
	Installments      []Installment `json:"installments"`
	Email             string        `json:"email"`
}

type Installment struct {
	InstallmentNumber    int       `json:"installment_number"`
	InstallmentAmount    big.Float `json:"installment_amount"`
	InstallmentFeeAmount big.Float `json:"installment_fee_amount"`
	Currency             string    `json:"currency"`
}
