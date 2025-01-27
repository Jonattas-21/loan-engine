package dto

import "time"

type loanConditionRequest_dto struct {
	Name         string 
	InterestRate int64  
	MinAge       int64  
	MaxAge       int64  
	ModifiedDate time.Time 
}
