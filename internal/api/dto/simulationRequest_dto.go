package dto

import (
	"time"
)

type SimulationRequest_dto struct {
	LoanAmount   float64
	Installments int
	BithDate     time.Time
	Currency     string
	Email        string
}
