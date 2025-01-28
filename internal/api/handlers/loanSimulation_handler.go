package handlers

import (
	"net/http"

)

type LoanSimulationHandler struct {
	LoanSimulation_usecase usecases.LoanSimulation_usecase
}

func (h *LoanSimulationHandler) GetLoanSimulation(w http.ResponseWriter, r *http.Request) {
	//todo
}
