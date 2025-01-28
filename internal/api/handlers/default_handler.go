package handlers

import (
	"fmt"
	"net/http"
	"time"
)

type DefaultHandler struct {
}

func (d *DefaultHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	//ping db and other services
	w.Write([]byte(fmt.Sprintf("I am alive now %v!", time.Now().Format("2006-01-02 15:04:05"))))
}
