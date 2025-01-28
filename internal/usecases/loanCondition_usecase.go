package usecases

import (
	"log"

	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/domain/interfaces"
)

type LoanCondition interface {
	SetLoanCondition(LoanCondition entities.LoanCondition) error
	GetLoanConditions() ([]entities.LoanCondition, error)
}

type LoanCondition_usecase struct {
	LoanConditionRepository interfaces.Repository[entities.LoanCondition]
}

func (l *LoanCondition_usecase) SetLoanConditions(LoanCondition entities.LoanCondition) error {
	fieldsFrom := make(map[string]interface{})
	fieldsFrom["name"] = LoanCondition.Name
	fieldsFrom["InterestRate"] = LoanCondition.InterestRate
	fieldsFrom["MaxAge"] = LoanCondition.MaxAge
	fieldsFrom["MinAge"] = LoanCondition.MinAge

	err := l.LoanConditionRepository.UpdateItemCollection(LoanCondition.Name, fieldsFrom)
	if err != nil {
		log.Println("Error updating loan condition: ", err.Error())
		return err
	}

	return nil
}

func (l *LoanCondition_usecase) GetLoanConditions() ([]entities.LoanCondition, error) {
	conditions, err := l.LoanConditionRepository.GetItemsCollection("")
	if err != nil {
		log.Println("Error getting loan conditions: ", err.Error())
		return nil, err
	}

	return conditions, nil
}
