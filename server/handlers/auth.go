package handlers

import (
	"ToDoList/database"
	"ToDoList/models"
	"ToDoList/service"
	"encoding/json"
	"errors"
	"net/http"
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var reg RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&reg)

	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	err = h.serv.RegisterUser(r.Context(), reg.Login, reg.Password)

	if err != nil {
		if errors.Is(err, service.ErrEmptyField) {
			http.Error(w, "login and password are required", http.StatusBadRequest)
			return
		}
		if errors.Is(err, database.ErrUserAlreadyExists) {
			http.Error(w, "login is already in use", http.StatusConflict)
			return
		}
		http.Error(w, "could not register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	token, err := h.serv.LoginUser(r.Context(), user.Login, user.Password)

	if err != nil {
		if errors.Is(err, service.ErrWrongPassword) {
			http.Error(w, "invalid login or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "could not log in", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		http.Error(w, "could not encode response", http.StatusInternalServerError)
	}

}
