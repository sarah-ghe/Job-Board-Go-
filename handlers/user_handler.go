package handlers

import (
	"encoding/json"
	"job-board/models"
	"job-board/services"
	"net/http"
)

type UserHandler struct {
	Service *services.UserService
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = h.Service.Register(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}




func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.Service.Login(request.Email, request.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	user, _ := h.Service.GetByEmail(request.Email)

	response := models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	}

	json.NewEncoder(w).Encode(response)
}





func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("user_id").(int)

	user, err := h.Service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// used the dto to not return the password field in the response
	response := models.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	json.NewEncoder(w).Encode(response)
}