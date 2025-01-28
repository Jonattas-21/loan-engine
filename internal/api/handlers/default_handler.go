package handlers

import (
	"net/http"
)

type DefaultHandler struct {
}


func (d *DefaultHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	//ping db and other services
	w.Write([]byte("I am alive!"))
}