package handlers

import (
	"encoding/json"
	"log"

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
// @Description update a loan condition by name
// @Tags conditions
// @Accept  json
// @Produce  json
// @Router /loanconditions/v1 [post]
func (h *LoanConditionHandler) SetLoanCondition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var loanConditionDto dto.LoanConditionRequest_dto

	log.Println("Setting loan condition, body is: ", r.Body)
	h.Logger.Infoln("Setting loan condition, body is: ", r.Body)

	if err := json.NewDecoder(r.Body).Decode(&loanConditionDto); err != nil {
		h.Logger.Errorln("Error decoding loan condition: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.Logger.Infoln("Loan condition dto is: ", loanConditionDto)
	log.Println("Loan condition dto is: ", loanConditionDto)

	// Converting dto to entity, there is no need of a automapper here, yet.
	loanCondition := entities.LoanCondition{
		Name:         loanConditionDto.Name,
		InterestRate: loanConditionDto.InterestRate,
		MaxAge:       loanConditionDto.MaxAge,
		MinAge:       loanConditionDto.MinAge,
	}

	err := h.LoanCondition_usecase.SetLoanCondition(loanCondition)
	if err != nil {
		h.Logger.Errorln("An internal error setting loan condition: ", err.Error())
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
// @Router /loanconditions/v1 [get]
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
