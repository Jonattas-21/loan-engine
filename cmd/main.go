// @title Loan Engine API
// @version 1.0
// @description This project It's a credit simulator which allows users to consult loan conditions, based in some payments conditions.
// @host localhost:8088
// @BasePath /api
package main

import (
	"github.com/joho/godotenv"

	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/Jonattas-21/loan-engine/internal/api/handlers"
	"github.com/Jonattas-21/loan-engine/internal/api/middlewares"
	"github.com/Jonattas-21/loan-engine/internal/domain/entities"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/cache"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/database"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/email"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/logger"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/queue"
	"github.com/Jonattas-21/loan-engine/internal/infrastructure/repositories"
	"github.com/Jonattas-21/loan-engine/internal/usecases"

	_ "github.com/Jonattas-21/loan-engine/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"context"

	"os/signal"
	"syscall"
)

func main() {
	//Setting up the logger
	log := logger.LogSetup()
	err := godotenv.Load("cmd/.env")

	if err != nil {
		log.Errorln("Error loading .env file")
		panic(err)
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
		MaxAge:           300,
	}))

	//Conecting to the database
	database := database.DatabaseNosql{Logger: log}
	log.Infoln("Connecting to the database " + os.Getenv("MONGO_HOST"))
	mdb := database.NewDatabase()
	log.Infoln("Connected to the database ")
	dbName := os.Getenv("MONGO_DB")

	//Conect to Redis cache
	rdb := cache.NewCache()
	cacheRepo := &repositories.RedisRepository{Redis: rdb, Logger: log}

	//Creating the condition usecase
	repoLoanCondition := &repositories.DefaultRepository[entities.LoanCondition]{Client: mdb, DatabaseName: dbName, CollectionName: "loan_conditions"}
	loanCondition_usecase := usecases.LoanCondition_usecase{
		LoanConditionRepository: repoLoanCondition,
		CacheRepository:         cacheRepo,
		Logger:                  log, //todo future: make this a logger interface
	}

	//Init the loan conditions tiers
	err = loanCondition_usecase.InitLoanEngineConditionsData()
	if err != nil {
		log.Fatalln("Error initializing loan conditions: ", err.Error())
		panic(err)
	}

	//Creating the queue
	queue := queue.RabbitMQ{Logger: log}
	// err = queue.CreateQueue(os.Getenv("RABBITMQ_PUBLISH_QUEUE"))
	// if err != nil {
	// 	log.Fatalln("Error creating queue: ", err.Error())
	// 	panic(err)
	// }

	//Creating the simulation usecase
	repoLoanSimulation := &repositories.DefaultRepository[entities.LoanSimulation]{Client: mdb, DatabaseName: dbName, CollectionName: "loan_simulations"}
	emailSender := email.EmailSender{}
	loanSimulation_usecase := usecases.LoanSimulation_usecase{
		LoanCondition:            &loanCondition_usecase,
		LoanSimulationRepository: repoLoanSimulation,
		CacheRepository:          cacheRepo,
		EmailSender:              &emailSender,
		Logger:                   log,
		QueuePublisher:           &queue,
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
	router.Get("/api/", dafault_handler.HealthCheck)

	router.Route("/api/v1/auth/", func(r chi.Router) {
		r.Post("/token", dafault_handler.GetToken)
	})

	router.Route("/api/v1/loanconditions/", func(r chi.Router) {
		if useAuth == "true" {
			r.Use(middlewares.Auth)
		}
		r.Post("/", loanCondition_handler.SetLoanCondition)
		r.Get("/", loanCondition_handler.GetLoanConditions)
	})

	router.Route("/api/v1/loansimulations/", func(r chi.Router) {
		if useAuth == "true" {
			r.Use(middlewares.Auth)
		}
		r.Get("/", loanSimulation_handler.GetLoanSimulation)
	})

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	// Create context with cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for shutdown signals
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Infoln("Shutting down server...")

		// Close connections
		_ = mdb.Disconnect(ctx)
		_ = rdb.Close()

		log.Infoln("Database and cache connections are closed.")
		cancel()
	}()

	// Run the serverb
	server := &http.Server{
		Addr:    ":8088",
		Handler: router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Infoln("Server running...")

	// Wait for the context to be canceled (shutdown signal received)
	<-ctx.Done()
	log.Infoln("Server stopped.")
}
