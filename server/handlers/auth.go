package handlers

import (
	"ToDoList/models"
	"encoding/json"
	"fmt"
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
		w.WriteHeader(400)
		w.Write([]byte("ошибка"))
		return
	}

	err = h.serv.RegisterUser(r.Context(), reg.Login, reg.Password)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("ошибка"))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("норм, создан"))
	// выводить юзера
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {

	var user models.User

	json.NewDecoder(r.Body).Decode(&user)

	token, err := h.serv.LoginUser(r.Context(), user.Login, user.Password)

	if err != nil {
		w.WriteHeader(400)

		w.Write([]byte("ошибка"))
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
