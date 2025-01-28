package usecases

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/api/dto"
	"github.com/Jonattas-21/loan-engine/internal/domain/interfaces"
)

type LoanSimulation interface {
	GetLoanSimulation(SimulationRequests [] dto.SimulationRequest_dto) ([] entities.LoanCondition, error)
	CalculateLoan(SimulationRequest dto.SimulationRequest_dto) (entities.LoanCondition, error)
}

type LoanSimulation_usecase struct {
	LoanSimulationRepository interfaces.Repository[entities.LoanSimulation]
	LoanCondition  LoanCondition
}

func (l *LoanSimulation_usecase) GetLoanSimulation(SimulationRequests []dto.SimulationRequest_dto) ([]entities.LoanSimulation, error) {

	var simulationResponses []entities.LoanSimulation

	//loop through all simulation requests
	//todo async
	for _, SimulationRequest := range SimulationRequests {
		simulationResponse, err := l.CalculateLoan(SimulationRequest)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error calculating loan, %v", err.Error()))
		}
		simulationResponses = append(simulationResponses, simulationResponse)
	}

	return simulationResponses, nil
}

func (l *LoanSimulation_usecase) CalculateLoan(SimulationRequest dto.SimulationRequest_dto) (entities.LoanSimulation, error) {

	//get fee conditions
	conditions, err := l.LoanCondition.GetLoanConditions()
	if err != nil {
		return entities.LoanSimulation{}, errors.New(fmt.Sprintf("Error getting loan conditions, %v", err.Error()))
	}

	//calculate age
	today := time.Now()
	age := today.Year() - SimulationRequest.BithDate.Year()

	// adjust the age if the birthdate has not occurred yet this year
	if today.YearDay() < SimulationRequest.BithDate.YearDay() {
		age--
	}

	//get interest rate
	var interestRate big.Float
	for _, condition := range conditions {
		if condition.MinAge >= age && condition.MaxAge < age {
			interestRate = condition.InterestRate
		}
	}

	//calculate instalment from a given loan value
	monthlyInterestRate := new(big.Float).Quo(&interestRate, big.NewFloat(12*100))
	one := big.NewFloat(1)

	// 1 + r
	onePlusRate := new(big.Float).Add(one, monthlyInterestRate)

	// (1 + r)^n
	ratePower := l.calculatePower(onePlusRate, SimulationRequest.Installments)
	numerator := new(big.Float).Mul(&SimulationRequest.LoanAmount, monthlyInterestRate)
	numerator.Mul(numerator, ratePower)

	// (1 + r)^n - 1
	denominator := new(big.Float).Sub(ratePower, one)

	InstallmentValue := new(big.Float).Quo(numerator, denominator)

	// Creating instalment by month
	var loanSimulation entities.LoanSimulation
	var AmountTobePaid big.Float

	l.createInstallments(SimulationRequest, InstallmentValue, loanSimulation, AmountTobePaid)

	return entities.LoanSimulation{
		AmountTobePaid:    AmountTobePaid,
		AmountFeeTobePaid: AmountTobePaid,
		SimulationDate:    time.Now(),
		Currency:          SimulationRequest.Currency,
	}, nil
}

func (l *LoanSimulation_usecase) calculatePower(base *big.Float, exponent int) *big.Float {
	result := big.NewFloat(1)
	for i := 0; i < exponent; i++ {
		result.Mul(result, base)
	}
	return result
}

func (l *LoanSimulation_usecase) createInstallments(SimulationRequest dto.SimulationRequest_dto, InstallmentValue *big.Float, loanSimulation entities.LoanSimulation, AmountTobePaid big.Float) {
	for i := 0; i < SimulationRequest.Installments; i++ {
		installment := entities.Installment{
			InstallmentNumber:    i + 1,
			InstallmentAmount:    *InstallmentValue,
			InstallmentFeeAmount: *InstallmentValue,
			Currency:             SimulationRequest.Currency,
		}
		loanSimulation.Installments = append(loanSimulation.Installments, installment)
		AmountTobePaid.Add(&AmountTobePaid, InstallmentValue)
	}
}