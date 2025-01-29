package entities

import (
	"time"
)

type LoanSimulation struct {
	LoanAmount          float64       `json:"loan_amount"`
	AmountTobePaid      float64       `json:"amount_to_be_paid"`
	AmountFeeTobePaid   float64       `json:"amount_fee_to_be_paid"`
	FeeAmountPercentage float64       `json:"fee_amount_percentage"`
	TotalInstallments   int           `json:"total_installments"`
	SimulationDate      time.Time     `json:"simulation_date"`
	Currency            string        `json:"currency"`
	Installments        []Installment `json:"installments"`
	Email               string        `json:"email"`
}

type Installment struct {
	InstallmentNumber    int     `json:"installment_number"`
	InstallmentAmount    float64 `json:"installment_amount"`
	InstallmentFeeAmount float64 `json:"installment_fee_amount"`
	Currency             string  `json:"currency"`
}
