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
	Service *services.JobService
}

// GET /jobs
func (h *JobHandler) GetJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jobs, err := h.Service.GetJobs()
	if err != nil {
		http.Error(w, "could not fetch jobs", http.StatusInternalServerError)
		return
	}

	var response []JobResponse
	for _, job := range jobs {
		response = append(response, toJobResponse(job))
	}

	json.NewEncoder(w).Encode(response)
}

// GET /jobs/{id}
func (h *JobHandler) GetJobById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	job, err := h.Service.GetByID(id)
	if err != nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(toJobResponse(*job))
}

// GET /my-jobs (protected)
func (h *JobHandler) GetMyJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value("user_id").(int)

	jobs, err := h.Service.GetJobsByUser(userID)
	if err != nil {
		http.Error(w, "could not fetch jobs", http.StatusInternalServerError)
		return
	}

	var response []JobResponse
	for _, job := range jobs {
		response = append(response, toJobResponse(job))
	}

	json.NewEncoder(w).Encode(response)
}

// POST /jobs (protected)
func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var job models.Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int)
	job.UserID = sql.NullInt64{Int64: int64(userID), Valid: true}

	err = h.Service.CreateJob(&job)
	if err != nil {
		http.Error(w, "could not create job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toJobResponse(job))
}

// PUT /jobs/{id} (protected + ownership)
func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int)

	var job models.Job
	err = json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateJob(id, userID, &job)
	if err != nil {
		http.Error(w, "not found or unauthorized", http.StatusForbidden)
		return
	}

	// Important: ensure ID + UserID are set for response
	job.ID = id
	job.UserID = sql.NullInt64{Int64: int64(userID), Valid: true}

	json.NewEncoder(w).Encode(toJobResponse(job))
}

// DELETE /jobs/{id} (protected + ownership)
func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int)

	err = h.Service.DeleteJob(id, userID)
	if err != nil {
		http.Error(w, "not found or unauthorized", http.StatusForbidden)
		return
	}

	//REST behavior: no body
	w.WriteHeader(http.StatusNoContent)
}