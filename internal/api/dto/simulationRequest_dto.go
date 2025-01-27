package dto

import (
	"time"
	"math/big"
)

type SimulationRequest_dto struct {
	LoanAmount big.Float 
	Installments int 
	BithDate time.Time 
	Currency string 
	Email string 
}
