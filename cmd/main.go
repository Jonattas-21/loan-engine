package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"os"

	"github.com/Jonattas-21/loan-engine/internal/api/handlers"
	"github.com/Jonattas-21/loan-engine/internal/api/middlewares"
	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	//"github.com/Jonattas-21/loan-engine/internal/infrastructure/cache"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/database"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/repositories"
	"github.com/Jonattas-21/loan-engine/internal/usecases"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err) //todo: search for panic
	}

	useAuth := os.Getenv("USE_SECURITY")

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
	mdb := database.NewDatabase()
	dbName := os.Getenv("MONGO_DB")

	//Conect to Redis cache
	//rdb := cache.NewCache()

	//Creating the usecases
	repoLoanCondition := &repositories.DefaultRepository[entities.LoanCondition]{Client: mdb, DatabaseName: dbName, CollectionName: "loan-conditions"}
	loanCondition_usecase := usecases.LoanCondition_usecase{
		LoanConditionRepository: repoLoanCondition,
	}

	repoLoanSimulation := &repositories.DefaultRepository[entities.LoanSimulation]{Client: mdb, DatabaseName: dbName, CollectionName: "loan-simulations"}
	loanSimulation_usecase := usecases.LoanSimulation_usecase{
		LoanCondition:            &loanCondition_usecase,
		LoanSimulationRepository: repoLoanSimulation,
	}

	//Creating the handlers
	dafault_handler := handlers.DefaultHandler{}
	loanCondition_handler := handlers.LoanConditionHandler{
		LoanCondition_usecase: &loanCondition_usecase,
	}
	loanSimulation_handler := handlers.LoanSimulationHandler{
		LoanSimulation_usecase: loanSimulation_usecase,
	}

	//Defining the routes
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Get("/", dafault_handler.HealthCheck)

	router.Route("/loanconditions", func(r chi.Router) {
		if useAuth == "true" {
			r.Use(middlewares.Auth)
		}
		r.Post("/", loanCondition_handler.SetLoanCondition)
		r.Get("/", loanCondition_handler.GetLoanConditions)
	})

	router.Route("/loansimulations", func(r chi.Router) {
		if useAuth == "true" {
			r.Use(middlewares.Auth)
		}
		r.Get("/", loanSimulation_handler.GetLoanSimulation)
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
