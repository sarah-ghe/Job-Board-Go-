package handlers

import (
	"encoding/json"
	"net/http"

	"job-board/models"
)

var jobs = []models.Job{
	{
		ID:          1,
		Title:       "Backend Developer",
		Description: "Build APIs using Go",
	},
}

func JobsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		getJobs(w, r)

	case http.MethodPost:
		createJob(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func createJob(w http.ResponseWriter, r *http.Request) {
	var newJob models.Job

	err := json.NewDecoder(r.Body).Decode(&newJob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newJob.ID = len(jobs) + 1
	jobs = append(jobs, newJob)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newJob)
}