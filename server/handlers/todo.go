package handlers

import (
	"ToDoList/auth"
	"ToDoList/database"
	"ToDoList/models"
	"ToDoList/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type serviceApi interface {
	GetTasks(ctx context.Context, token string) ([]models.Task, error)
	GetTask(ctx context.Context, id int, token string) (models.Task, error)
	AddTask(ctx context.Context, token, title, descr string) error
	DeleteTask(ctx context.Context, id int, token string) error
	EditDescriptionTask(ctx context.Context, id int, descr, token string) error
	CompleteTask(ctx context.Context, id int, token string) error
	LoginUser(ctx context.Context, login, passw string) (string, error)
	RegisterUser(ctx context.Context, login, passw string) error
}

type Handler struct {
	serv serviceApi
}

type GetFromJsonDescription struct {
	Description string `json:"description"`
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidToken):
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
	case errors.Is(err, service.ErrEmptyTitle), errors.Is(err, service.ErrEmptyField):
		http.Error(w, "invalid request data", http.StatusBadRequest)
	case errors.Is(err, service.ErrNotYourTask):
		http.Error(w, "access to this task is forbidden", http.StatusForbidden)
	case errors.Is(err, database.ErrTaskNotFound):
		http.Error(w, "task not found", http.StatusNotFound)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func CreateHandler(serv serviceApi) *Handler {
	return &Handler{serv: serv}
}

func MiddlewareParseToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		const prefix = "Bearer "

		if !strings.HasPrefix(authHeader, prefix) || strings.TrimSpace(strings.TrimPrefix(authHeader, prefix)) == "" {
			http.Error(w, "authorization token is required", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, prefix)

		ctx := context.WithValue(r.Context(), "token", tokenStr)

		r = r.WithContext(ctx)

		next(w, r)

	})
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {

	tokenStr := r.Context().Value("token").(string)

	tasks, err := h.serv.GetTasks(r.Context(), tokenStr)

	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		// Headers may already be sent by Encoder; the error is still handled instead of ignored.
		fmt.Println("encode tasks response:", err)
		return
	}
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {

	tokenStr := r.Context().Value("token").(string)

	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	err = h.serv.AddTask(r.Context(), tokenStr, task.Title, task.Description)

	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	s := fmt.Sprintf("Task \"%s\" was added", task.Title)
	w.Write([]byte(s))
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {

	tokenStr := r.Context().Value("token").(string)

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "task id must be an integer", http.StatusBadRequest)
		return
	}

	err = h.serv.DeleteTask(r.Context(), id, tokenStr)

	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *Handler) EditDescriptionTask(w http.ResponseWriter, r *http.Request) {

	tokenStr := r.Context().Value("token").(string)

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "task id must be an integer", http.StatusBadRequest)
		return
	}

	var description GetFromJsonDescription

	err = json.NewDecoder(r.Body).Decode(&description)

	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	err = h.serv.EditDescriptionTask(r.Context(), id, description.Description, tokenStr)

	if err != nil {
		writeServiceError(w, err)
		return
	}

	fmt.Fprintf(w, "Description of the task with id %d was edited: %s", id, description.Description)
}

func (h *Handler) CompleteTask(w http.ResponseWriter, r *http.Request) {

	tokenStr := r.Context().Value("token").(string)

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "task id must be an integer", http.StatusBadRequest)
		return
	}

	err = h.serv.CompleteTask(r.Context(), id, tokenStr)

	if err != nil {
		writeServiceError(w, err)
		return
	}

	task, err := h.serv.GetTask(r.Context(), id, tokenStr)

	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		fmt.Println("encode completed task response:", err)
	}

}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Context().Value("token").(string)

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "task id must be an integer", http.StatusBadRequest)
		return
	}

	task, err := h.serv.GetTask(r.Context(), id, tokenStr)

	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		// Headers may already be sent by Encoder; the error is still handled instead of ignored.
		fmt.Println("encode task response:", err)
	}

}
