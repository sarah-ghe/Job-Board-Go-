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
	job.UserID = userID

	err = h.Service.CreateJob(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(job)
}

func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	userID := r.Context().Value("user_id").(int)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	var job models.Job

	err = json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateJob(id, userID, &job)

	if err == sql.ErrNoRows {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(job)
}

func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	userID := r.Context().Value("user_id").(int)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid job id", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteJob(id, userID)

	if err == sql.ErrNoRows {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
