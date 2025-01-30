package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_"github.com/Jonattas-21/loan-engine/internal/api/dto"
	"github.com/Jonattas-21/loan-engine/internal/domain/interfaces"
	"github.com/Jonattas-21/loan-engine/package/auth"
)

type DefaultHandler struct {
	MongoRepo       interfaces.Repository[string]
	CacheRepository interfaces.CacheRepository
}

// @Summary Check if the application is running
// @Description Check if the application is running and connected to the database and cache
// @Tags default
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router / [get]
func (d *DefaultHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	//ping db and other services
	var message string = fmt.Sprintf("Api said: I am alive now %v!\n", time.Now().Format("2006-01-02 15:04:05"))

	err := d.MongoRepo.Ping()
	message += "MongoDB said:\n"
	if err != nil {
		message += fmt.Sprintf("- I got an error during ping in DB: %v\n", err.Error())
	} else {
		message += fmt.Sprintf("- I'm alive now %v!\n", time.Now().Format("2006-01-02 15:04:05"))
	}

	err = d.CacheRepository.Ping()
	message += "RedisCache said:\n"
	if err != nil {
		message += fmt.Sprintf("- I got an error during ping in Cache: %v\n", err.Error())
	} else {
		message += fmt.Sprintf("- I'm alive now %v!\n", time.Now().Format("2006-01-02 15:04:05"))
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary login in the application
// @Description login in the application and get a token
// @Tags default
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.TokenResponse_dto
// @Router /v1/auth/token [post]
func (d *DefaultHandler) GetToken(w http.ResponseWriter, r *http.Request) {

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Get username and password from form data
	username := r.FormValue("username")
	password := r.FormValue("password")

	log.Printf("username: %s, password: %s", username, password)

	// Get token from Keycloak
	tokenResponse, err := auth.GetTokenFromKeycloak(username, password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting token: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the token response as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tokenResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}
