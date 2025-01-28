package handlers

import (
	"encoding/json"
	"github.com/Jonattas-21/loan-engine/internal/api/dto"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
	"log"
	"net/http"
)

type LoanSimulationHandler struct {
	LoanSimulation_usecase usecases.LoanSimulation_usecase
}

func (h *LoanSimulationHandler) GetLoanSimulation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loanSimulationDto dto.SimulationRequest_dto

	if err := json.NewDecoder(r.Body).Decode(&loanSimulationDto); err != nil {
		log.Println("Error decoding loan simulation: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseSimulation, err := h.LoanSimulation_usecase.CalculateLoan(loanSimulationDto)
	if err != nil {
		log.Println("An internal error calculating loan: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(responseSimulation)
	if err != nil {
		log.Println("Error encoding loan simulation: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
