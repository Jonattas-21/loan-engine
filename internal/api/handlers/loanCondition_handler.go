package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Jonattas-21/loan-engine/internal/api/dto"
	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
	"github.com/sirupsen/logrus"
)

type LoanConditionHandler struct {
	LoanCondition_usecase usecases.LoanCondition
	Logger                *logrus.Logger
}

// @Summary update a loan condition by name
// @Description update a loan condition by name tier1, tier2, tier3, tier4
// @Tags conditions
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /v1/loanconditions [post]
func (h *LoanConditionHandler) SetLoanCondition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var loanConditionDto dto.LoanConditionRequest_dto

	if err := json.NewDecoder(r.Body).Decode(&loanConditionDto); err != nil {
		h.Logger.Errorln("Error decoding loan condition: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Converting dto to entity, there is no need of a automapper here, yet.
	loanCondition := entities.LoanCondition{
		Name:         loanConditionDto.Name,
		InterestRate: loanConditionDto.InterestRate,
		MaxAge:       loanConditionDto.MaxAge,
		MinAge:       loanConditionDto.MinAge,
	}

	err := h.LoanCondition_usecase.SetLoanCondition(loanCondition)
	if err != nil {
		h.Logger.Errorln("An internal error setting loan condition: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode("Loan condition set successfully")
	if err != nil {
		h.Logger.Errorln("Error encoding loan condition: ", err.Error())
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Show the list of loan conditions, fees by age group
// @Description Get all conditions
// @Tags conditions
// @Accept  json
// @Produce  json
// @Success 200 {array} entities.LoanCondition
// @Router /v1/loanconditions [get]
func (h *LoanConditionHandler) GetLoanConditions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	conditions, err := h.LoanCondition_usecase.GetLoanConditions()
	if err != nil {
		h.Logger.Errorln("Error getting loan conditions: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(conditions)
	if err != nil {
		h.Logger.Errorln("Error encoding loan conditions: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
