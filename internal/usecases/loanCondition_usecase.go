package usecases

import (
	"github.com/Jonattas-21/loan-engine/internal/domain/interfaces"
	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
)

type LoanCondition interface {
	SetLoanConditions(LoanConditions []entities.LoanCondition) error
	GetLoanConditions() ([]entities.LoanCondition, error)
}

type LoanCondition_usecase struct {
	LoanConditionRepository interfaces.Repository
}

func (l *LoanCondition_usecase) SetLoanConditions(LoanConditions []entities.LoanCondition) error {
 	//todo
	return nil
}

func (l *LoanCondition_usecase) GetLoanConditions() ([]entities.LoanCondition, error) {
 	//todo
	return nil, nil
}
