package usecases_test

import (
	"math/big"
	"testing"

	"encoding/json"
	"fmt"
	"github.com/Jonattas-21/loan-engine/internal/api/dto"
	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/logger"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
	internalMock "github.com/Jonattas-21/loan-engine/tests"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"time"
)

var (
	mockSimulationDatabaseRepo = new(internalMock.MockRepository[entities.LoanSimulation])
	loanSimulationUsecase      = &usecases.LoanSimulation_usecase{
		CacheRepository:          mockCacheRepo,
		LoanSimulationRepository: mockSimulationDatabaseRepo,
		Logger:                   logger.LogSetup(),
	}
)

func setupSimulation() {
	// Reset the mocks
	mockConditionDatabaseRepo = new(internalMock.MockRepository[entities.LoanCondition])
	mockCacheRepo = new(internalMock.MockCacheRepository)
	loanConditionUsecase = &usecases.LoanCondition_usecase{
		CacheRepository:         mockCacheRepo,
		LoanConditionRepository: mockConditionDatabaseRepo,
		Logger:                  logger.LogSetup(),
	}
	mockSimulationDatabaseRepo = new(internalMock.MockRepository[entities.LoanSimulation])
	loanSimulationUsecase      = &usecases.LoanSimulation_usecase{
		CacheRepository:          mockCacheRepo,
		LoanSimulationRepository: mockSimulationDatabaseRepo,
		Logger:                   logger.LogSetup(),
	}
}


func TestCalculatePower_ok(t *testing.T) {
	assert := assert.New(t)

	// Define test cases
	testCases := []struct {
		base     *big.Float
		exponent int
		expected *big.Float
	}{
		{big.NewFloat(2), 3, big.NewFloat(8)},      // 2^3 = 8
		{big.NewFloat(5), 0, big.NewFloat(1)},      // 5^0 = 1
		{big.NewFloat(3), 2, big.NewFloat(9)},      // 3^2 = 9
		{big.NewFloat(7), 1, big.NewFloat(7)},      // 7^1 = 7
		{big.NewFloat(10), 4, big.NewFloat(10000)}, // 10^4 = 10000
	}

	for _, tc := range testCases {
		// Call the function
		result := loanSimulationUsecase.CalculatePower(tc.base, tc.exponent)

		// Assert the result
		assert.Equal(0, result.Cmp(tc.expected), "Base: %v, Exponent: %v", tc.base, tc.exponent)
	}
}

func TestTruncateToTwoDecimals_ok(t *testing.T) {
	assert := assert.New(t)
	// Define test cases
	testCases := []struct {
		input    float64
		expected float64
	}{
		{123.456789, 123.45},
		{987.654321, 987.65},
		{1.999999, 1.99},
		{0.123456, 0.12},
		{100.0, 100.0},
	}

	for _, tc := range testCases {
		result := loanSimulationUsecase.TruncateToTwoDecimals(tc.input)
		assert.Equal(tc.expected, result)
	}
}

func TestCalculateLoan_ok(t *testing.T) {
	assert := assert.New(t)
	setupSimulation()

	// Define test data
	simulationRequest := dto.SimulationRequest_dto{
		Email:        "test@example.com",
		LoanAmount:   10000,
		Installments: 6,
		BithDate:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Currency:     "R$",
	}
	loanConditions := []entities.LoanCondition{
		{Name: "tier1", InterestRate: 5, MinAge: 18, MaxAge: 25},
		{Name: "tier2", InterestRate: 3, MinAge: 26, MaxAge: 40},
	}

	jsonConditions, _ := json.Marshal(loanConditions)

	// Mock expectations
	//mockConditionDatabaseRepo.On("GetLoanConditions").Return(loanConditions, nil)
	mockConditionDatabaseRepo.On("GetItemsCollection", "loan_conditions").Return(loanConditions, nil)
	mockCacheRepo.On("Get", "loan_conditions").Return("", fmt.Errorf("not found"))
	mockCacheRepo.On("Set", "loan_conditions", jsonConditions, time.Minute*10).Return(nil)

	// Call the function
	result, err := loanSimulationUsecase.CalculateLoan(simulationRequest)

	// Assertions
	assert.NoError(err)
	assert.NotNil(result)
	assert.Equal(10000.00, result.LoanAmount)
	assert.Equal(10087.68, result.AmountTobePaid)
	assert.Equal(87.68, result.AmountFeeTobePaid)
	assert.Equal(3.0, result.FeeAmountPercentage)
	assert.Equal(6, result.TotalInstallments)
	assert.Equal("R$", result.Currency)
	assert.Equal("test@example.com", result.Email)
	assert.NotEmpty(result.Installments)

	// Verify mock expectations
	mockConditionDatabaseRepo.AssertExpectations(t)
}

func TestCalculateLoan_volume(t *testing.T) {
	assert := assert.New(t)
	setupCondition()

	// Define test data
	simulationRequest := []dto.SimulationRequest_dto{}

	for i := 0; i < 1000; i++ {
		randInstallmentNumber := rand.Intn(48-12+1) + 12
		rendLoanValue := 1000 + rand.Float64()*(50000-1000)

		simulationRequest = append(simulationRequest, dto.SimulationRequest_dto{
			Email:        "teste@teste.com",
			LoanAmount:   rendLoanValue,
			Installments: randInstallmentNumber,
			BithDate:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Currency:     "R$",
		})
	}

	loanConditions := []entities.LoanCondition{
		{Name: "tier1", InterestRate: 3.0, MinAge: 30, MaxAge: 40},
	}

	// Mock expectations
	mockConditionDatabaseRepo.On("GetLoanConditions").Return(loanConditions, nil)
	mockSimulationDatabaseRepo.On("SendLoanSimulationEmailMessage").Return(nil)

	// Call the function
	for _, request := range simulationRequest {
		result, err := loanSimulationUsecase.CalculateLoan(request)
		assert.NoError(err)
		assert.NotNil(result)
		assert.NotEmpty(result.Installments)
	}
}
