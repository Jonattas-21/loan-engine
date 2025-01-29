package handlers

import (
	"encoding/json"
	"log"
	"math/big"
	"net/http"

	"github.com/Jonattas-21/loan-engine/internal/api/dto"
	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
)

type LoanConditionHandler struct {
	LoanCondition_usecase usecases.LoanCondition
}

func (h *LoanConditionHandler) SetLoanCondition(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var loanConditionDto dto.LoanConditionRequest_dto

	if err := json.NewDecoder(r.Body).Decode(&loanConditionDto); err != nil {
		log.Println("Error decoding loan condition: ", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Converting dto to entity, there is no need of a automapper here, yet.
	loanCondition := entities.LoanCondition{
		Name:         loanConditionDto.Name,
		InterestRate: *big.NewFloat(loanConditionDto.InterestRate),
		MaxAge:       loanConditionDto.MaxAge,
		MinAge:       loanConditionDto.MinAge,
	}

	err := h.LoanCondition_usecase.SetLoanCondition(loanCondition)
	if err != nil {
		log.Println("An internal error setting loan condition: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *LoanConditionHandler) GetLoanConditions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	conditions, err := h.LoanCondition_usecase.GetLoanConditions()
	if err != nil {
		log.Println("Error getting loan conditions: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(conditions)
	if err != nil {
		log.Println("Error encoding loan conditions: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}