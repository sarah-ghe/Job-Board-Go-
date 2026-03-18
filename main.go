package main

import (
	"job-board/config"
	"job-board/handlers"
	"job-board/middleware"
	"job-board/repositories"
	"job-board/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Initialize database connection
	config.ConnectDB()

	// ---------------- JOB SETUP ----------------

	jobRepo := &repositories.PostgresJobRepository{
		DB: config.DB,
	}

	jobService := &services.JobService{
		Repo: jobRepo,
	}

	jobHandler := &handlers.JobHandler{
		Service: jobService,
	}

	// ---------------- USER SETUP ----------------

	userRepo := &repositories.PostgresUserRepository{
		DB: config.DB,
	}

	userService := &services.UserService{
		Repo: userRepo,
	}

	userHandler := &handlers.UserHandler{
		Service: userService,
	}

	// ---------------- ROUTER ----------------

	router := mux.NewRouter()

	// -------- PUBLIC ROUTES --------

	router.HandleFunc("/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	// -------- PROTECTED ROUTES --------

	router.Handle("/jobs", middleware.AuthMiddleware(http.HandlerFunc(jobHandler.CreateJob))).Methods("POST")
	router.Handle("/jobs", middleware.AuthMiddleware(http.HandlerFunc(jobHandler.GetJobs))).Methods("GET")
	router.Handle("/jobs/{id}", middleware.AuthMiddleware(http.HandlerFunc(jobHandler.UpdateJob))).Methods("PUT")
	router.Handle("/jobs/{id}", middleware.AuthMiddleware(http.HandlerFunc(jobHandler.DeleteJob))).Methods("DELETE")
	router.Handle("/me", middleware.AuthMiddleware(http.HandlerFunc(userHandler.GetMe))).Methods("GET")

	// ---------------- START SERVER ----------------

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}