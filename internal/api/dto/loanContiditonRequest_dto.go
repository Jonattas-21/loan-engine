package dto

import "time"

type LoanConditionRequest_dto struct {
	Name         string 
	InterestRate float64  
	MinAge       int  
	MaxAge       int  
	ModifiedDate time.Time 
}
