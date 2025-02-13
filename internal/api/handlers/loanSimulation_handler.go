package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Jonattas-21/loan-engine/internal/api/dto"
	_ "github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/domain/interfaces"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
)

type LoanSimulationHandler struct {
	LoanSimulation_usecase usecases.LoanSimulation_usecase
	Logger                 interfaces.Log
}

// @Summary  Get a plenty of loan simulations
// @Description Get a plenty of loan simulations
// @Tags simulation
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.LoanSimulationResponse_dto
// @Router /v1/loansimulations [post]
func (h *LoanSimulationHandler) GetLoanSimulation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Infoln("Received request to calculate loan simulation: ", r.Body)
	var loanSimulationDto []dto.SimulationRequest_dto

	if err := json.NewDecoder(r.Body).Decode(&loanSimulationDto); err != nil {
		h.Logger.Errorln("Error decoding loan simulation: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.Logger.Infoln("Calculating loan simulation: ", loanSimulationDto)
	responseSimulation, errs := h.LoanSimulation_usecase.GetLoanSimulation(loanSimulationDto)

	reponse := dto.LoanSimulationResponse_dto{
		LoanSimulations: responseSimulation,
		ErrorSimulations: errs,
	}

	err := json.NewEncoder(w).Encode(reponse)
	if err != nil {
		h.Logger.Errorln("Error encoding loan simulation: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
