package dto

import "github.com/Jonattas-21/loan-engine/internal/domain/entities"

type LoanSimulationResponse_dto struct {
	LoanSimulations  []entities.LoanSimulation
	ErrorSimulations []string
}
