package handlers

import (
	"encoding/json"
	"github.com/Jonattas-21/loan-engine/internal/api/dto"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
	"github.com/sirupsen/logrus"
	"net/http"
	_"github.com/Jonattas-21/loan-engine/internal/domain/entities"
)

type LoanSimulationHandler struct {
	LoanSimulation_usecase usecases.LoanSimulation_usecase
	Logger                 *logrus.Logger
}

// @Summary  Get a plenty of loan simulations
// @Description Get a plenty of loan simulations
// @Tags simulation
// @Accept  json
// @Produce  json
// @Success 200 {array} entities.LoanSimulation
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
	responseSimulation, err := h.LoanSimulation_usecase.GetLoanSimulation(loanSimulationDto)
	if err != nil {
		h.Logger.Errorln("An internal error calculating loan: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infoln("Loan simulation calculated: ", responseSimulation)
	err = json.NewEncoder(w).Encode(responseSimulation)
	if err != nil {
		h.Logger.Errorln("Error encoding loan simulation: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
