package usecases

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"time"

	"github.com/Jonattas-21/loan-engine/internal/api/dto"
	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/domain/interfaces"
)

type LoanSimulation interface {
	GetLoanSimulation(SimulationRequests []dto.SimulationRequest_dto) ([]entities.LoanCondition, error)
	CalculateLoan(SimulationRequest dto.SimulationRequest_dto) (entities.LoanCondition, error)
}

type LoanSimulation_usecase struct {
	LoanSimulationRepository interfaces.Repository[entities.LoanSimulation]
	CacheRepository          interfaces.CacheRepository
	EmailSender              interfaces.EmailSender
	LoanCondition            LoanCondition
}

func (l *LoanSimulation_usecase) GetLoanSimulation(SimulationRequests []dto.SimulationRequest_dto) ([]entities.LoanSimulation, error) {

	var simulationResponses []entities.LoanSimulation
	var loanSimulation entities.LoanSimulation

	//loop through all simulation requests
	//todo async
	for _, simulationRequest := range SimulationRequests {

		//check if the request is in cache
		value, err := l.CacheRepository.Get(simulationRequest.Email)
		if err == nil {
			err = json.Unmarshal([]byte(value.(string)), &loanSimulation)
			if err != nil {
				log.Println(fmt.Sprintf("Error unmarshalling loan simulation from cache from email: %v ", simulationRequest.Email), err.Error())
			} else {
				//send email
				message, err := l.formatEmailMessage(loanSimulation)
				if err == nil {
					//the error it's been logged before, but the process should continue
					l.EmailSender.SendMail(fmt.Sprintf("Loan simulation %v", time.Now().Format("2006-01-02 15:04:05")), message, loanSimulation.Email)
				}
				simulationResponses = append(simulationResponses, loanSimulation)

				continue
			}
		}

		//calculate loan if not in cache
		simulationResponse, err := l.CalculateLoan(simulationRequest)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error calculating loan, %v", err.Error()))
		}
		l.CacheRepository.Set(simulationResponse.Email, simulationRequest, time.Second*5)
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
		Email:             SimulationRequest.Email,
	}, nil
}

func (l *LoanSimulation_usecase) calculatePower(base *big.Float, exponent int) *big.Float {
	result := big.NewFloat(1)
	for i := 0; i < exponent; i++ {
		result.Mul(result, base)
	}
	return result
}

func (l *LoanSimulation_usecase) createInstallments(simulationRequest dto.SimulationRequest_dto, installmentValue *big.Float, loanSimulation entities.LoanSimulation, amountTobePaid big.Float) {
	for i := 0; i < simulationRequest.Installments; i++ {
		installment := entities.Installment{
			InstallmentNumber:    i + 1,
			InstallmentAmount:    *installmentValue,
			InstallmentFeeAmount: *installmentValue,
			Currency:             simulationRequest.Currency,
		}
		loanSimulation.Installments = append(loanSimulation.Installments, installment)
		amountTobePaid.Add(&amountTobePaid, installmentValue)
	}
}

func (l *LoanSimulation_usecase) formatEmailMessage(loanSimulation entities.LoanSimulation) (string, error) {
	// Read the template file
	tmpl, err := template.ParseFiles("internal/infrastructure/email/templates/sendLoanSimulation.html") //could be readed on init, one time.
	if err != nil {
		log.Println(fmt.Sprintf("Error reading email template, %v", err.Error()))
		return "", errors.New(fmt.Sprintf("Error reading email template, %v", err.Error()))
	}

	// Generate the HTML content
	var htmlContent bytes.Buffer
	err = tmpl.Execute(&htmlContent, loanSimulation)
	if err != nil {
		log.Println(fmt.Sprintf("Error executing email template, %v, simulation for email %v", err.Error(), loanSimulation.Email))
	}

	return htmlContent.String(), nil
}
