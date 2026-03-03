package main

import (
	"job-board/config"
	"job-board/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()

	r := mux.NewRouter()

	// Route for listing and creating jobs. JobsHandler itself switches on method.
	r.HandleFunc("/jobs", handlers.JobsHandler).Methods(http.MethodGet, http.MethodPost)

	// Routes for updating and deleting a single job by id.
	r.HandleFunc("/jobs/{id}", handlers.UpdateJob).Methods(http.MethodPut)
	r.HandleFunc("/jobs/{id}", handlers.DeleteJob).Methods(http.MethodDelete)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
