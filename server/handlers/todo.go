package handlers

import (
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

type Handler struct {
	serv *service.Service
}

type GetFromJsonDescription struct {
	Description string `json:"description"`
}

func CreateHandler(serv *service.Service) *Handler {
	return &Handler{serv: serv}
}

func MiddlewareParseToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		const prefix = "Bearer "

		if !strings.HasPrefix(authHeader, prefix) {
			w.WriteHeader(400)
			w.Write([]byte("Bad token"))
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

	tokenStr := r.Context().Value("token").(string)

	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("cant add task"))
		return
	}

	err = h.serv.AddTask(r.Context(), tokenStr, task.Title, task.Description)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("cant add task to data base: " + err.Error()))
		return
	}

	w.WriteHeader(200)
	s := fmt.Sprintf("Task \"%s\" was added", task.Title)
	w.Write([]byte(s))
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {

	tokenStr := r.Context().Value("token").(string)

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id must be integer")
		return
	}

	err = h.serv.DeleteTask(r.Context(), id, tokenStr)

	if err != nil {
		if errors.Is(err, service.ErrNotYourTask) {
			w.WriteHeader(404)
			fmt.Println(err)
			fmt.Fprint(w, "It's not your task")
			return
		} else if errors.Is(err, database.ErrTaskNotFound) {
			w.WriteHeader(404)
			fmt.Println(err)
			fmt.Fprint(w, "There is no such id")
			return
		} else {

			w.WriteHeader(500)
			fmt.Fprint(w, "Couldn't delete task")
			return
		}
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "Task with id %d was deleted", id)

}

func (h *Handler) EditDescriptionTask(w http.ResponseWriter, r *http.Request) {

	tokenStr := r.Context().Value("token").(string)

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id must be integer")
		return
	}

	var description GetFromJsonDescription

	err = json.NewDecoder(r.Body).Decode(&description)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Error")
		return
	}

	err = h.serv.EditDescriptionTask(r.Context(), id, description.Description, tokenStr)

	if err != nil {
		if errors.Is(err, service.ErrNotYourTask) {
			w.WriteHeader(404)
			fmt.Println(err)
			fmt.Fprint(w, "It's not your task")
			return
		} else if errors.Is(err, database.ErrTaskNotFound) {
			w.WriteHeader(404)
			fmt.Println(err)
			fmt.Fprint(w, "There is no such id")
			return
		} else {

			w.WriteHeader(500)
			fmt.Fprint(w, "Couldn't delete task")
			return
		}
	}

	fmt.Fprintf(w, "Description of the task with id %d was edited: %s", id, description.Description)
}

func (h *Handler) CompleteTask(w http.ResponseWriter, r *http.Request) {

	tokenStr := r.Context().Value("token").(string)

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id must be integer")
		return
	}

	err = h.serv.CompleteTask(r.Context(), id, tokenStr)

	if err != nil {
		if errors.Is(err, service.ErrNotYourTask) {
			w.WriteHeader(404)
			fmt.Println(err)
			fmt.Fprint(w, "It's not your task")
			return
		} else if errors.Is(err, database.ErrTaskNotFound) {
			w.WriteHeader(404)
			fmt.Println(err)
			fmt.Fprint(w, "There is no such id")
			return
		} else {

			w.WriteHeader(500)
			fmt.Fprint(w, "Couldn't complete task")
			return
		}
	}

	task, err := h.serv.GetTask(r.Context(), id, tokenStr)

	if err != nil {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Task with id %d, completion changed", id)
		return
	}

	fmt.Fprintf(w, "Task with id %d, completed: %v", id, task.Completed)

}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Context().Value("token").(string)

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id must be integer")
		return
	}

	task, err := h.serv.GetTask(r.Context(), id, tokenStr)

	if err != nil {

		if errors.Is(err, service.ErrNotYourTask) {
			w.WriteHeader(404)
			fmt.Println(err)
			fmt.Fprint(w, "It's not your task")
			return
		} else if errors.Is(err, database.ErrTaskNotFound) {
			w.WriteHeader(404)
			fmt.Println(err)
			fmt.Fprint(w, "There is no such id")
			return
		} else {
			w.WriteHeader(500)
			fmt.Fprint(w, "Couldn't get task")
			return
		}
	}

	err = json.NewEncoder(w).Encode(task)

	if err != nil {
		fmt.Println("error:", err)
	}

}
