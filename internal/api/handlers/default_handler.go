package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Jonattas-21/loan-engine/internal/domain/interfaces"
)

type DefaultHandler struct {
	MongoRepo       interfaces.Repository[string]
	CacheRepository interfaces.CacheRepository
}

func (d *DefaultHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	//ping db and other services
	var message string = fmt.Sprintf("Api said: I am alive now %v!\n", time.Now().Format("2006-01-02 15:04:05"))

	err := d.MongoRepo.Ping()
	message += "MongoDB said:\n"
	if err != nil {
		message += fmt.Sprintf("Error during ping in DB: %v\n", err.Error())
	} else {
		message += fmt.Sprintf("MongoDB is alive now %v!\n", time.Now().Format("2006-01-02 15:04:05"))
	}

	err = d.CacheRepository.Ping()
	message += "RedisCache said:\n"
	if err != nil {
		message += fmt.Sprintf("Error during ping in Cache: %v\n", err.Error())
	} else {
		message += fmt.Sprintf("Cache is alive now %v!\n", time.Now().Format("2006-01-02 15:04:05"))
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
