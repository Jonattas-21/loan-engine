package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Jonattas-21/loan-engine/internal/api/dto"
	_ "github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
	"github.com/Jonattas-21/loan-engine/internal/domain/interfaces"
	"strings"
)

type LoanConditionHandler struct {
	LoanCondition_usecase usecases.LoanCondition
	Logger                interfaces.Log
}

// @Summary update a loan condition by name
// @Description update a loan condition by name tier1, tier2, tier3, tier4 and for now it's only possible to update the interest rate
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

	err, validations := h.LoanCondition_usecase.SetLoanCondition(loanConditionDto)
	if err != nil {
		h.Logger.Errorln("An internal error setting loan condition: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if validations != nil {
		err = json.NewEncoder(w).Encode(strings.Join(validations, ", "))
		if err != nil {
			h.Logger.Errorln("Error encoding loan condition: ", err.Error())
		}
		http.Error(w, "Error setting loan condition", http.StatusBadRequest)
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
