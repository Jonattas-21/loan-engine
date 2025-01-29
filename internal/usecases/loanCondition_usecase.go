package usecases

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"encoding/json"

	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/domain/interfaces"
)

type LoanCondition interface {
	SetLoanCondition(LoanCondition entities.LoanCondition) error
	GetLoanConditions() ([]entities.LoanCondition, error)
}

type LoanCondition_usecase struct {
	LoanConditionRepository interfaces.Repository[entities.LoanCondition]
	CacheRepository         interfaces.CacheRepository
}

func (l *LoanCondition_usecase) SetLoanCondition(LoanCondition entities.LoanCondition) error {
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
	loanConditions := []entities.LoanCondition{}
	val, err := l.CacheRepository.Get("*")
	if err == nil {
		err = json.Unmarshal([]byte(val.(string)), &loanConditions)
		if err != nil {
			log.Println("Error unmarshalling loan conditions from cache: ", err.Error())
		}
		return loanConditions, nil
	}

	conditions, err := l.LoanConditionRepository.GetItemsCollection("")
	if err != nil {
		log.Println("Error getting loan conditions: ", err.Error())
		return nil, fmt.Errorf("Error getting loan conditions: %w", err)
	}

	if len(conditions) > 0 {
		err = l.CacheRepository.Set("*", conditions, time.Second*10)
		if err != nil {
			log.Println("Error setting loan conditions in cache: ", err.Error())
			return nil, fmt.Errorf("Error setting loan conditions in cache: %w", err)
		}
	}

	return conditions, nil
}

func (l *LoanCondition_usecase) InitLoanEngineConditionsData() error {
	err := l.LoanConditionRepository.SaveItemCollection(entities.LoanCondition{
		Name:         "tier1",
		InterestRate: *big.NewFloat(5),
		MinAge:       18,
		MaxAge:       25,
	})
	if err != nil {
		log.Println("Error saving default loan condition for tier 1:", err.Error())
		return fmt.Errorf("Error saving default loan condition for tier 1: %w", err)
	}

	err = l.LoanConditionRepository.SaveItemCollection(entities.LoanCondition{
		Name:         "tier2",
		InterestRate: *big.NewFloat(3),
		MinAge:       26,
		MaxAge:       40,
	})
	if err != nil {
		log.Println("Error saving default loan condition for tier 2:", err.Error())
		return fmt.Errorf("Error saving default loan condition for tier 2: %w", err)
	}

	err = l.LoanConditionRepository.SaveItemCollection(entities.LoanCondition{
		Name:         "tier2",
		InterestRate: *big.NewFloat(2),
		MinAge:       41,
		MaxAge:       60,
	})
	if err != nil {
		log.Println("Error saving default loan condition for tier 3:", err.Error())
		return fmt.Errorf("Error saving default loan condition for tier 3: %w", err)
	}

	err = l.LoanConditionRepository.SaveItemCollection(entities.LoanCondition{
		Name:         "tier2",
		InterestRate: *big.NewFloat(4),
		MinAge:       61,
		MaxAge:       100,
	})
	if err != nil {
		log.Println("Error saving default loan condition for tier 4:", err.Error())
		return fmt.Errorf("Error saving default loan condition for tier 4: %w", err)
	}

	return nil
}
