package handlers

import(
	"net/http"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
)

type LoanConditionHandler struct {
	LoanCondition_usecase usecases.LoanCondition
}

func (h *LoanConditionHandler) SetLoanConditions(w http.ResponseWriter, r *http.Request) {
	//todo
}

func (h *LoanConditionHandler) GetLoanConditions(w http.ResponseWriter, r *http.Request) {
	//todo
}