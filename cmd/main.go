// @title Loan Engine API
// @version 1.0
// @description This project It's a credit simulator which allows users to consult loan conditions, based in some payments conditions.
// @host localhost:8080
// @BasePath /api
package main

import (
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"os"

	"github.com/Jonattas-21/loan-engine/internal/api/handlers"
	"github.com/Jonattas-21/loan-engine/internal/api/middlewares"
	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/cache"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/database"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/email"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/logger"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/repositories"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
)

func main() {
	//Setting up the logger
	log := logger.LogSetup()
	err := godotenv.Load("cmd/.env")

	if err != nil {
		log.Error("Error loading .env file")
		panic(err) //todo: search for panic
	}

	//Creating the router
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	//Setup cors
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	//Conecting to the database
	database := database.DatabaseNosql{Logger: log}
	mdb := database.NewDatabase()
	dbName := os.Getenv("MONGO_DB")

	//Conect to Redis cache
	rdb := cache.NewCache()
	cacheRepo := &repositories.RedisRepository{Redis: rdb, Logger: log}

	//Creating the condition usecase
	repoLoanCondition := &repositories.DefaultRepository[entities.LoanCondition]{Client: mdb, DatabaseName: dbName, CollectionName: "loan_conditions"}
	loanCondition_usecase := usecases.LoanCondition_usecase{
		LoanConditionRepository: repoLoanCondition,
		CacheRepository:         cacheRepo,
		Logger:                  log, //todo: make this a logger interface
	}

	//Init the loan conditions tiers
	err = loanCondition_usecase.InitLoanEngineConditionsData()
	if err != nil {
		log.Fatal("Error initializing loan conditions: ", err.Error())
		panic(err)
	}

	//Creating the simulation usecase
	repoLoanSimulation := &repositories.DefaultRepository[entities.LoanSimulation]{Client: mdb, DatabaseName: dbName, CollectionName: "loan_simulations"}
	emailSender := email.EmailSender{}
	loanSimulation_usecase := usecases.LoanSimulation_usecase{
		LoanCondition:            &loanCondition_usecase,
		LoanSimulationRepository: repoLoanSimulation,
		CacheRepository:          cacheRepo,
		EmailSender:              &emailSender,
		Logger:                   log,
	}

	//Creating the handlers
	repoDefault := &repositories.DefaultRepository[string]{Client: mdb, DatabaseName: dbName, CollectionName: "default"}
	dafault_handler := handlers.DefaultHandler{
		MongoRepo:       repoDefault,
		CacheRepository: cacheRepo,
	}
	loanCondition_handler := handlers.LoanConditionHandler{
		LoanCondition_usecase: &loanCondition_usecase,
		Logger:                log,
	}
	loanSimulation_handler := handlers.LoanSimulationHandler{
		LoanSimulation_usecase: loanSimulation_usecase,
		Logger:                 log,
	}

	//Defining the routes
	useAuth := os.Getenv("USE_SECURITY")

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Get("/", dafault_handler.HealthCheck)

	router.Route("/loanconditions/v1", func(r chi.Router) {
		if useAuth == "true" {
			r.Use(middlewares.Auth)
		}
		r.Post("/", loanCondition_handler.SetLoanCondition)
		r.Get("/", loanCondition_handler.GetLoanConditions)
	})

	router.Route("/loansimulations/v1", func(r chi.Router) {
		if useAuth == "true" {
			r.Use(middlewares.Auth)
		}
		r.Get("/", loanSimulation_handler.GetLoanSimulation)
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
