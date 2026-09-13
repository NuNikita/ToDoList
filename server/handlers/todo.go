package handlers

import (
	"ToDoList/database"
	"ToDoList/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

type GetTasksRequest struct {
	Creator string `json:"creator"`
}

func CreateHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	var request GetTasksRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("cant get tasks"))
		return
	}

	tasks, err := database.GetTasks(r.Context(), request.Creator, h.pool)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("cant get tasks from data base"))
		return
	}

	err = json.NewEncoder(w).Encode(tasks)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("cant encode tasks to json"))
		return
	}

}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {

	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("cant add task"))
		return
	}

	err = database.AddTask(r.Context(), task.Creator, task.Title, task.Description, h.pool)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("cant add task to data base"))
		return
	}

	w.WriteHeader(200)
	s := fmt.Sprintf("Task \"%s\" was added", task.Title)
	w.Write([]byte(s))
}
