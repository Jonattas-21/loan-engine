package handlers

import (
	"encoding/json"
	"github.com/Jonattas-21/loan-engine/internal/api/dto"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
	"github.com/sirupsen/logrus"
	"net/http"
)

type LoanSimulationHandler struct {
	LoanSimulation_usecase usecases.LoanSimulation_usecase
	Logger                 *logrus.Logger
}

func (h *LoanSimulationHandler) GetLoanSimulation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loanSimulationDto []dto.SimulationRequest_dto

	if err := json.NewDecoder(r.Body).Decode(&loanSimulationDto); err != nil {
		h.Logger.Error("Error decoding loan simulation: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseSimulation, err := h.LoanSimulation_usecase.GetLoanSimulation(loanSimulationDto)
	if err != nil {
		h.Logger.Error("An internal error calculating loan: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(responseSimulation)
	if err != nil {
		h.Logger.Error("Error encoding loan simulation: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
