package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"job-board/models"
	"job-board/services"

	"github.com/gorilla/mux"
)

type JobHandler struct {
	Service *services.JobService //calls the service layer to perform business logic
}

func (h *JobHandler) JobsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		h.GetJobs(w, r)

	case http.MethodPost:
		h.CreateJob(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) {

	jobs, err := h.Service.GetJobs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(jobs)
}

func (h *JobHandler) GetJobById(w http.ResponseWriter, r *http.Request) {
	// Always set response type
	w.Header().Set("Content-Type", "application/json")

	// Get ID from URL
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	// Fetch job
	job, err := h.Service.GetByID(id)
	if err != nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	// Encode response safely
	if err := json.NewEncoder(w).Encode(job); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *JobHandler) GetMyJobs(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("user_id").(int)

	jobs, err := h.Service.GetJobsByUser(userID)
	if err != nil {
		http.Error(w, "could not fetch jobs", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(jobs)
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {

	var job models.Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int)
	job.UserID = sql.NullInt64{Int64: int64(userID), Valid: true}

	err = h.Service.CreateJob(&job)
	if err != nil {
		http.Error(w, "Could not create job", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(job)
}

func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {

	//get job ID from URL
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	//get authenticaed user ID from context
	userID := r.Context().Value("user_id").(int)

	var job models.Job
	err = json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateJob(id, userID, &job)
	if err == sql.ErrNoRows {
		http.Error(w, "job not found or unauthorized", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Could not update job", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Job updated successfully",
	})
}

func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	//get authenticaed user ID from context
	userID := r.Context().Value("user_id").(int)

	//delete only if the user owns the job
	err = h.Service.DeleteJob(id, userID)
	if err == sql.ErrNoRows {
		http.Error(w, "job not found or unauthorized", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Job deleted successfully",
	})
}
