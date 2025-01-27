package usecases

import (
	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
)

type loan interface{
	SetLoanConditions(LoanConditions LoanCondition) error
}