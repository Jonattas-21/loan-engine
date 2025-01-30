package usecases_test

import (
	"math/big"
	"testing"

	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/logger"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
	internalMock "github.com/Jonattas-21/loan-engine/tests"
	"github.com/stretchr/testify/assert"
)

var (
	mockSimulationDatabaseRepo = new(internalMock.MockRepository[entities.LoanSimulation])
	loanSimulationUsecase      = &usecases.LoanSimulation_usecase{
		CacheRepository:          mockCacheRepo,
		LoanSimulationRepository: mockSimulationDatabaseRepo,
		Logger:                   logger.LogSetup(),
	}
)

func TestCalculatePower(t *testing.T) {
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
		assert.Equal(t, 0, result.Cmp(tc.expected), "Base: %v, Exponent: %v", tc.base, tc.exponent)
	}
}
