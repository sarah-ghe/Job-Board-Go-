package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"job-board/config"
	"job-board/models"
)

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

	rows, err := config.DB.Query("SELECT id, title, description FROM jobs")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var jobs []models.Job

	for rows.Next() {
		var job models.Job

		err := rows.Scan(&job.ID, &job.Title, &job.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jobs = append(jobs, job)
	}

	json.NewEncoder(w).Encode(jobs)
}

func createJob(w http.ResponseWriter, r *http.Request) {

	var newJob models.Job

	err := json.NewDecoder(r.Body).Decode(&newJob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO jobs (title, description)
	VALUES ($1, $2)
	RETURNING id
	`

	err = config.DB.QueryRow(
		query,
		newJob.Title,
		newJob.Description,
	).Scan(&newJob.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newJob)
}

func UpdateJob(w http.ResponseWriter, r *http.Request) {

	var updatedJob models.Job

	err := json.NewDecoder(r.Body).Decode(&updatedJob)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := mux.Vars(r)["id"]

	query := `
	UPDATE jobs
	SET title = $1, description = $2
	WHERE id = $3
	RETURNING id
	`

	err = config.DB.QueryRow(
		query,
		updatedJob.Title,
		updatedJob.Description,
		id,
	).Scan(&updatedJob.ID)

	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(updatedJob)
}

func DeleteJob(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	query := `
	DELETE FROM jobs
	WHERE id = $1
	`

	result, err := config.DB.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
