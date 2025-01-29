package usecases

import (
	"fmt"
	"time"

	"encoding/json"

	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/domain/interfaces"
	"github.com/sirupsen/logrus"
)

type LoanCondition interface {
	SetLoanCondition(LoanCondition entities.LoanCondition) error
	GetLoanConditions() ([]entities.LoanCondition, error)
}

type LoanCondition_usecase struct {
	LoanConditionRepository interfaces.Repository[entities.LoanCondition]
	CacheRepository         interfaces.CacheRepository
	Logger                  *logrus.Logger
}

func (l *LoanCondition_usecase) SetLoanCondition(LoanCondition entities.LoanCondition) error {
	fieldsFrom := make(map[string]interface{})
	fieldsFrom["name"] = LoanCondition.Name
	fieldsFrom["InterestRate"] = LoanCondition.InterestRate
	fieldsFrom["MaxAge"] = LoanCondition.MaxAge
	fieldsFrom["MinAge"] = LoanCondition.MinAge

	err := l.LoanConditionRepository.UpdateItemCollection(LoanCondition.Name, fieldsFrom)
	if err != nil {
		l.Logger.Error("Error updating loan condition: ", err.Error())
		return err
	}

	return nil
}

func (l *LoanCondition_usecase) GetLoanConditions() ([]entities.LoanCondition, error) {
	loanConditions := []entities.LoanCondition{}

	// Check if we have the loan conditions in cache
	val, err := l.CacheRepository.Get("loan_conditions")
	if err == nil {
		err = json.Unmarshal([]byte(val), &loanConditions)
		if err != nil {
			l.Logger.Error("Error unmarshalling loan conditions from cache: ", err.Error())
		}

		l.Logger.Error(fmt.Printf("Loan conditions from cache: %v", loanConditions))
		return loanConditions, nil
	}

	// If we don't have the loan conditions in cache, let's get it from mongoDB
	conditions, err := l.LoanConditionRepository.GetItemsCollection("loan_conditions")
	l.Logger.Info(fmt.Printf("Conditions: %v", conditions))
	if err != nil {
		l.Logger.Error("Error getting loan conditions: ", err.Error())
		return nil, fmt.Errorf("Error getting loan conditions from mongoDB: %w", err)
	}

	if len(conditions) > 0 {
		jsonConditions, err := json.Marshal(conditions)
		// Save in cache, if not, let's just log the error and continue
		if err != nil {
			l.Logger.Error("Error marshalling loan conditions: ", err.Error())
		} else {
			err = l.CacheRepository.Set("loan_conditions", jsonConditions, time.Minute*10)
			if err != nil {
				l.Logger.Error("Error setting loan conditions in cache: ", err.Error())
			}
		}
	}

	l.Logger.Info(fmt.Printf("Loan conditions from mongoDB: %v", conditions))
	return conditions, nil
}

func (l *LoanCondition_usecase) InitLoanEngineConditionsData() error {

	err := l.LoanConditionRepository.TrunkCollection()
	if err != nil {
		l.Logger.Error("Error truncating loan conditions collection: ", err.Error())
		return fmt.Errorf("Error truncating loan conditions collection: %w", err)
	}

	err = l.LoanConditionRepository.SaveItemCollection(entities.LoanCondition{
		Name:         "tier1",
		InterestRate: 5,
		MinAge:       18,
		MaxAge:       25,
		ModifiedDate: time.Now(),
	})
	if err != nil {
		l.Logger.Error("Error saving default loan condition for tier 1:", err.Error())
		return fmt.Errorf("Error saving default loan condition for tier 1: %w", err)
	}

	err = l.LoanConditionRepository.SaveItemCollection(entities.LoanCondition{
		Name:         "tier2",
		InterestRate: 3,
		MinAge:       26,
		MaxAge:       40,
		ModifiedDate: time.Now(),
	})
	if err != nil {
		l.Logger.Error("Error saving default loan condition for tier 2:", err.Error())
		return fmt.Errorf("Error saving default loan condition for tier 2: %w", err)
	}

	err = l.LoanConditionRepository.SaveItemCollection(entities.LoanCondition{
		Name:         "tier2",
		InterestRate: 2,
		MinAge:       41,
		MaxAge:       60,
		ModifiedDate: time.Now(),
	})
	if err != nil {
		l.Logger.Error("Error saving default loan condition for tier 3:", err.Error())
		return fmt.Errorf("Error saving default loan condition for tier 3: %w", err)
	}

	err = l.LoanConditionRepository.SaveItemCollection(entities.LoanCondition{
		Name:         "tier2",
		InterestRate: 4,
		MinAge:       61,
		MaxAge:       100,
		ModifiedDate: time.Now(),
	})
	if err != nil {
		l.Logger.Error("Error saving default loan condition for tier 4:", err.Error())
		return fmt.Errorf("Error saving default loan condition for tier 4: %w", err)
	}

	return nil
}
