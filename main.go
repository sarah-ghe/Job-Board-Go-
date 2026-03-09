package main

import (
	"job-board/config"
	"job-board/handlers"
	"job-board/repositories"
	"job-board/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	
	// Initialize database connection
	config.ConnectDB()

	// Create repository
	jobRepo := &repositories.PostgresJobRepository{
		DB: config.DB,
	}

	// Create service
	jobService := &services.JobService{
		Repo: jobRepo,
	}

	// Create handler
	jobHandler := &handlers.JobHandler{
		Service: jobService,
	}

	// Initialize router
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/jobs", jobHandler.CreateJob).Methods("POST")
	router.HandleFunc("/jobs", jobHandler.GetJobs).Methods("GET")
	router.HandleFunc("/jobs/{id}", jobHandler.UpdateJob).Methods("PUT")
	router.HandleFunc("/jobs/{id}", jobHandler.DeleteJob).Methods("DELETE")

	// Start server
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}