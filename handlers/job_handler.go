package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"job-board/services"
	"job-board/models"
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



func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {

	var job models.Job

	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.CreateJob(&job)
	userID := r.Context().Value("user_id").(int)
	job.UserID = userID
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(job)
}



func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	var job models.Job

	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateJob(id, &job)

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

	id := mux.Vars(r)["id"]

	err := h.Service.DeleteJob(id)

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