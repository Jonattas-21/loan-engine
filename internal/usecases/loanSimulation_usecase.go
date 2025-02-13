package usecases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"time"

	"math"

	"github.com/Jonattas-21/loan-engine/internal/api/dto"
	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/domain/interfaces"
)

type LoanSimulation interface {
	GetLoanSimulation(SimulationRequests []dto.SimulationRequest_dto) ([]entities.LoanCondition, []error)
	CalculateLoan(SimulationRequest dto.SimulationRequest_dto) (entities.LoanCondition, error)
	TruncateToTwoDecimals(value float64) float64
	CalculatePower(base *big.Float, exponent int) *big.Float
	SendLoanSimulationEmailMessage(loanSimulation entities.LoanSimulation) error
	CreateInstallments(simulationRequest dto.SimulationRequest_dto, installmentValue *big.Float) []entities.Installment
}

type LoanSimulation_usecase struct {
	LoanSimulationRepository interfaces.Repository[entities.LoanSimulation]
	CacheRepository          interfaces.CacheRepository
	EmailSender              interfaces.EmailSender
	LoanCondition            LoanCondition
	Logger                   interfaces.Log
	QueuePublisher           interfaces.Queue
}

func (l *LoanSimulation_usecase) GetLoanSimulation(SimulationRequests []dto.SimulationRequest_dto) ([]entities.LoanSimulation, []string) {
	var simulationResponses []entities.LoanSimulation
	var errorsResponse []string
	simulatorChan := make(chan entities.LoanSimulation)
	errorChan := make(chan error)
	doneChan := make(chan bool)

	//loop through all simulation requests
	for _, simulationRequest := range SimulationRequests {
		go func(simulationRequest dto.SimulationRequest_dto) {
			var loanSimulation entities.LoanSimulation

			//validate request
			errors := l.ValidateSimulationRequest(simulationRequest)
			if errors != nil {
				errorChan <- fmt.Errorf("[email:%v] error validating simulation request: %v",simulationRequest.Email, errors )
				return
			}

			keyRedis := fmt.Sprintf("simulation_%v_%v_%v", simulationRequest.Email, simulationRequest.LoanAmount, simulationRequest.Installments)

			//check if the request is in cache
			value, err := l.CacheRepository.Get(keyRedis)
			if err == nil {
				err = json.Unmarshal([]byte(value), &loanSimulation)
				if err != nil {
					l.Logger.Errorln(fmt.Sprintf("[email:%v] Error unmarshalling loan simulation from cache", simulationRequest.Email), err.Error())
				} else {
					//send email
					err = l.SendLoanSimulationEmailMessage(loanSimulation)
					if err != nil {
						//if email fails, let's just log the error and continue
						l.Logger.Errorln(fmt.Sprintf("[email:%v] Error sending email for loan simulation", simulationRequest.Email), err.Error())
					}
					simulatorChan <- loanSimulation
					return
				}
			}

			//calculate loan if not in cache
			simulationResponse, err := l.CalculateLoan(simulationRequest)
			if err != nil {
				errorChan <- fmt.Errorf("[email:%v] error calculating loan, %v", simulationRequest.Email ,err.Error())
				return
			}

			// Save in cache, if not, let's just log the error and continue
			jsonConditions, err := json.Marshal(simulationResponse)
			if err != nil {
				l.Logger.Errorln(fmt.Sprintf("[email:%v] error marshalling loan conditions: %v", simulationRequest.Email, err.Error()))
			} else {
				err = l.CacheRepository.Set(keyRedis, jsonConditions, time.Minute*5)
				if err != nil {
					l.Logger.Errorln(fmt.Sprintf("[email:%v] error setting loan simulation in cache", simulationRequest.Email), err.Error())
				}
			}

			err = l.LoanSimulationRepository.SaveItemCollection(simulationResponse)
			if err != nil {
				l.Logger.Errorln(fmt.Sprintf("[email:%v] Error saving loan simulation", simulationRequest.Email), err.Error())
				errorChan <- fmt.Errorf("error saving loan simulation, %v", err.Error())
			}

			//send email
			err = l.SendLoanSimulationEmailMessage(loanSimulation)
			if err != nil {
				//if email fails, let's just log the error and continue
				l.Logger.Errorln(fmt.Sprintf("[email:%v] Error sending email for loan simulation", simulationRequest.Email), err.Error())
			}
			simulatorChan <- simulationResponse
		}(simulationRequest)
	}

	// Collect async results
	go func() {
		for i := 0; i < len(SimulationRequests); i++ {
			select {
			case res := <-simulatorChan:
				//sync response
				simulationResponses = append(simulationResponses, res)

				//publish in queue, async response
				// jsonConditions, err := json.Marshal(res)
				// if err != nil {
				// 	l.Logger.Errorln("Error marshalling loan conditions: ", err.Error())
				// }
				// l.QueuePublisher.PublishMessage(os.Getenv("RABBITMQ_PUBLISH_QUEUE") , string(jsonConditions))

			case err := <-errorChan:
				errorFormated := fmt.Errorf("error processing simulation: %v for request: %v", err.Error(), SimulationRequests[i])
				l.Logger.Errorln(errorFormated)
				errorsResponse = append(errorsResponse, errorFormated.Error())
			}
		}
		doneChan <- true
	}()

	// Wait for completion all the simulations
	<-doneChan

	return simulationResponses, errorsResponse
}

func (l *LoanSimulation_usecase) CalculateLoan(SimulationRequest dto.SimulationRequest_dto) (entities.LoanSimulation, error) {
	//get fee conditions
	conditions, err := l.LoanCondition.GetLoanConditions()
	if err != nil {
		l.Logger.Errorln(fmt.Printf("Error getting loan conditions, %v", err.Error()))
		return entities.LoanSimulation{}, fmt.Errorf("error getting loan conditions, %v", err.Error())
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
		if age >= condition.MinAge && age <= condition.MaxAge {
			interestRate = *big.NewFloat(condition.InterestRate)
		}
	}

	//check if interest rate was found
	interestRateFloat, _ := interestRate.Float64()
	if interestRateFloat == 0 {
		return entities.LoanSimulation{}, fmt.Errorf("interest rate not found for age %v", age)
	}

	l.Logger.Infoln(fmt.Sprintf("input fro calc: rate %v, age %v, instalmentsN %v, pv %v", interestRateFloat, age, SimulationRequest.Installments, SimulationRequest.LoanAmount))

	//calculate instalment from a given loan value
	monthlyInterestRate := new(big.Float).Quo(&interestRate, big.NewFloat(12*100))
	one := big.NewFloat(1)
	l.Logger.Infoln(fmt.Sprintf("monthlyInterestRate %v", monthlyInterestRate))

	// 1 + r
	onePlusRate := new(big.Float).Add(one, monthlyInterestRate)
	l.Logger.Infoln(fmt.Sprintf("onePlusRate %v", onePlusRate))

	// (1 + r)^n
	ratePower := l.CalculatePower(onePlusRate, SimulationRequest.Installments)
	l.Logger.Infoln(fmt.Sprintf("ratePower %v", ratePower))
	loanAmountBig := big.NewFloat(SimulationRequest.LoanAmount)
	l.Logger.Infoln(fmt.Sprintf("loanAmountBig %v", loanAmountBig))
	numerator := new(big.Float).Mul(loanAmountBig, monthlyInterestRate)
	l.Logger.Infoln(fmt.Sprintf("numerator %v", numerator))
	numerator.Mul(numerator, ratePower)

	// (1 + r)^n - 1
	denominator := new(big.Float).Sub(ratePower, one)
	l.Logger.Infoln(fmt.Sprintf("denominator %v", denominator))
	InstallmentValue := new(big.Float).Quo(numerator, denominator)
	l.Logger.Infoln(fmt.Sprintf("InstallmentValue %v", InstallmentValue))

	// Creating instalment by month
	totalAmountTobePaid := new(big.Float).Mul(InstallmentValue, big.NewFloat(float64(SimulationRequest.Installments)))
	totalAmountTobePaid_float, _ := totalAmountTobePaid.Float64()

	//calculate fee
	amountFeeTobePaid := totalAmountTobePaid_float - SimulationRequest.LoanAmount

	return entities.LoanSimulation{
		LoanAmount:          l.TruncateToTwoDecimals(SimulationRequest.LoanAmount),
		AmountTobePaid:      l.TruncateToTwoDecimals(totalAmountTobePaid_float),
		AmountFeeTobePaid:   l.TruncateToTwoDecimals(amountFeeTobePaid),
		FeeAmountPercentage: interestRateFloat,
		TotalInstallments:   SimulationRequest.Installments,
		SimulationDate:      time.Now(),
		Currency:            SimulationRequest.Currency,
		Email:               SimulationRequest.Email,
		Installments:        l.CreateInstallments(SimulationRequest, InstallmentValue),
	}, nil
}

func (l *LoanSimulation_usecase) CalculatePower(base *big.Float, exponent int) *big.Float {
	result := big.NewFloat(1)
	for i := 0; i < exponent; i++ {
		result.Mul(result, base)
	}
	return result
}

func (l *LoanSimulation_usecase) TruncateToTwoDecimals(value float64) float64 {
	factor := math.Pow(10, 2)
	truncatedValue := math.Trunc(value*factor) / factor
	return truncatedValue
}

func (l *LoanSimulation_usecase) CreateInstallments(simulationRequest dto.SimulationRequest_dto, installmentValue *big.Float) []entities.Installment {
	var returnInstallments []entities.Installment
	for i := 0; i < simulationRequest.Installments; i++ {
		installmentValueFloat, _ := installmentValue.Float64()
		installment := entities.Installment{
			InstallmentNumber:    i + 1,
			InstallmentAmount:    l.TruncateToTwoDecimals(installmentValueFloat),
			//InstallmentFeeAmount: l.TruncateToTwoDecimals(installmentValueFloat), //todo
			Currency:             simulationRequest.Currency,
		}
		returnInstallments = append(returnInstallments, installment)
	}
	return returnInstallments
}

func (l *LoanSimulation_usecase) SendLoanSimulationEmailMessage(loanSimulation entities.LoanSimulation) error {

	// Read the template file
	tmpl, err := template.ParseFiles("internal/infrastructure/email/templates/sendLoanSimulation.html") //could be readed on init, one time.
	if err != nil {
		l.Logger.Errorln(fmt.Printf("Error reading email template, %v", err.Error()))
		return fmt.Errorf("error reading email template, %v", err.Error())
	}

	// Generate the HTML content
	var htmlContent bytes.Buffer
	err = tmpl.Execute(&htmlContent, loanSimulation)
	if err != nil {
		l.Logger.Errorln(fmt.Printf("Error executing email template, %v, simulation for email %v", err.Error(), loanSimulation.Email))
		return fmt.Errorf("error executing email template, %v, simulation for email %v", err.Error(), loanSimulation.Email)
	}

	err = l.EmailSender.SendMail(fmt.Sprintf("Loan simulation %v", time.Now().Format("2006-01-02 15:04:05")), htmlContent.String(), loanSimulation.Email)
	if err != nil {
		l.Logger.Errorln(fmt.Printf("Error sending email, %v, simulation for email %v", err.Error(), loanSimulation.Email))
		return fmt.Errorf("error sending email, %v, simulation for email %v", err.Error(), loanSimulation.Email)
	}

	return nil
}

func (l *LoanSimulation_usecase) ValidateSimulationRequest(SimulationRequest dto.SimulationRequest_dto) []string {

	var errors []string

	if SimulationRequest.BithDate.IsZero() {
		errors = append(errors, "Birthdate is required")
	}

	if SimulationRequest.LoanAmount <= 0 {
		errors = append(errors, "Loan amount is required above 0")
	}

	if SimulationRequest.Installments <= 0 {
		errors = append(errors, "Installments is required above 0")
	}

	if SimulationRequest.Currency != "R$" && SimulationRequest.Currency != "U$" {
		errors = append(errors, "Currency is required R$ or U$")
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
