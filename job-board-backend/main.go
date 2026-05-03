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

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	//router.Handle("/jobs", http.HandlerFunc(jobHandler.CreateJob)).Methods("POST")
	router.Handle("/jobs", http.HandlerFunc(jobHandler.GetJobs)).Methods("GET") // public endpoints return all data
	router.Handle("/jobs/{id}", http.HandlerFunc(jobHandler.GetJobById)).Methods("GET")
	router.Handle("/jobs/{id}", middleware.AuthMiddleware(http.HandlerFunc(jobHandler.UpdateJob))).Methods("PUT")
	//router.Handle("/jobs/{id}", http.HandlerFunc(jobHandler.UpdateJob)).Methods("PUT")
	router.Handle("/jobs/{id}", middleware.AuthMiddleware(http.HandlerFunc(jobHandler.DeleteJob))).Methods("DELETE")
	//router.Handle("/jobs/{id}", http.HandlerFunc(jobHandler.DeleteJob)).Methods("DELETE")
	router.Handle("/me", middleware.AuthMiddleware(http.HandlerFunc(userHandler.GetMe))).Methods("GET")
	router.Handle("/my-jobs", middleware.AuthMiddleware(http.HandlerFunc(jobHandler.GetMyJobs))).Methods("GET") // protected endpoints use the authenticated user_id to filter resources belonging to the user.

	// ---------------- START SERVER ----------------

	log.Println("Server running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", enableCORS(router)))

}
