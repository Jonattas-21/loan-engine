package usecases_test

import (
	"errors"
	"fmt"
	"testing"

	"encoding/json"
	"time"

	"github.com/Jonattas-21/loan-engine/internal/api/dto"
	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/logger"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
	internalMock "github.com/Jonattas-21/loan-engine/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockConditionDatabaseRepo = new(internalMock.MockRepository[entities.LoanCondition])
	mockCacheRepo             = new(internalMock.MockCacheRepository)
	loanConditionUsecase      = &usecases.LoanCondition_usecase{
		CacheRepository:         mockCacheRepo,
		LoanConditionRepository: mockConditionDatabaseRepo,
		Logger:                  logger.LogSetup(),
	}
)

func setupCondition() {
	// Reset the mocks
	mockConditionDatabaseRepo = new(internalMock.MockRepository[entities.LoanCondition])
	mockCacheRepo = new(internalMock.MockCacheRepository)
	loanConditionUsecase = &usecases.LoanCondition_usecase{
		CacheRepository:         mockCacheRepo,
		LoanConditionRepository: mockConditionDatabaseRepo,
		Logger:                  logger.LogSetup(),
	}
}

func TestSetLoanCondition(t *testing.T) {
	setupCondition()
	assert := assert.New(t)
	// Define the test loan condition
	loanCondition := dto.LoanConditionRequest_dto{
		Name:         "tier1",
		InterestRate: 5.0,
		MaxAge:       60,
		MinAge:       18,
	}

	// Set up expected calls and returns
	fields := map[string]interface{}{
		"interestrate": 5.0,
	}
	mockConditionDatabaseRepo.On("UpdateItemCollection", loanCondition.Name, fields).Return(nil)
	mockCacheRepo.On("Set", "loan_conditions", mock.Anything, time.Minute*10).Return(nil)

	// Call the function
	err := loanConditionUsecase.SetLoanCondition(loanCondition)

	// Assertions
	assert.NoError(err)
}

func TestGetLoanConditions_db(t *testing.T) {
	setupCondition()
	assert := assert.New(t)

	// Define test data
	loanConditions := []entities.LoanCondition{
		{Name: "tier1", InterestRate: 5.0, MaxAge: 60, MinAge: 18},
		{Name: "tier2", InterestRate: 4.5, MaxAge: 65, MinAge: 20},
	}
	jsonConditions, _ := json.Marshal(loanConditions)

	mockConditionDatabaseRepo.On("GetItemsCollection", "loan_conditions").Return(loanConditions, nil)
	mockCacheRepo.On("Get", "loan_conditions").Return("", fmt.Errorf("not found"))
	mockCacheRepo.On("Set", "loan_conditions", jsonConditions, time.Minute*10).Return(nil)

	conditions, err := loanConditionUsecase.GetLoanConditions()
	assert.NoError(err)
	assert.Equal(loanConditions, conditions)
	mockCacheRepo.AssertExpectations(t)
	mockConditionDatabaseRepo.AssertExpectations(t)
}

func TestGetLoanConditions_cache(t *testing.T) {
	setupCondition()
	assert := assert.New(t)

	// Define test data
	loanConditions := []entities.LoanCondition{
		{Name: "Condition1", InterestRate: 5.0, MaxAge: 60, MinAge: 18},
		{Name: "Condition2", InterestRate: 4.5, MaxAge: 65, MinAge: 20},
	}
	jsonConditions, _ := json.Marshal(loanConditions)

	// Test case 1: Loan conditions found in cache
	mockCacheRepo.On("Get", "loan_conditions").Return(string(jsonConditions), nil)

	conditions, err := loanConditionUsecase.GetLoanConditions()
	assert.NoError(err)
	assert.Equal(loanConditions, conditions)
	mockCacheRepo.AssertExpectations(t)
}

func TestInitLoanEngineConditionsData(t *testing.T) {
	setupCondition()
	assert := assert.New(t)

	// Set up expected calls and returns
	mockConditionDatabaseRepo.On("TrunkCollection").Return(nil)
	mockConditionDatabaseRepo.On("SaveItemCollection", mock.Anything).Return(nil)

	// Call the function
	err := loanConditionUsecase.InitLoanEngineConditionsData()

	// Assertions
	assert.Nil(err)
	mockConditionDatabaseRepo.AssertExpectations(t)
}

func TestInitLoanEngineConditionsData_Error(t *testing.T) {
	setupCondition()
	assert := assert.New(t)

	// Set up expected calls and returns
	mockConditionDatabaseRepo.On("TrunkCollection").Return(nil)
	mockConditionDatabaseRepo.On("SaveItemCollection", mock.Anything).Return(errors.New("Error saving default loan condition for tier 1"))

	// Call the function
	err := loanConditionUsecase.InitLoanEngineConditionsData()

	// Assertions
	assert.NotNil(err)
	assert.Contains(err.Error(), "Error saving default loan condition for tier 1")
	mockConditionDatabaseRepo.AssertExpectations(t)
}
