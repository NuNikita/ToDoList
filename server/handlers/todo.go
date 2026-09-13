package handlers

import (
	"ToDoList/database"
	"ToDoList/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

type GetCreator struct {
	Creator string `json:"creator"`
}

type GetDescription struct {
	Description string `json:"description"`
}

func CreateHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	var request GetCreator

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

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id must be integer")
		return
	}

	err = database.DeleteTask(r.Context(), id, h.pool)

	if err != nil {
		if errors.Is(err, database.ErrTaskNotFound) {
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
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id must be integer")
		return
	}

	var description GetDescription

	err = json.NewDecoder(r.Body).Decode(&description)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Error")
		return
	}

	err = database.EditDescriptionTask(r.Context(), id, description.Description, h.pool)

	if err != nil {
		if errors.Is(err, database.ErrTaskNotFound) {
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
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id must be integer")
		return
	}

	err = database.CompleteTask(r.Context(), id, h.pool)

	if err != nil {
		if errors.Is(err, database.ErrTaskNotFound) {
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

	task, err := database.GetTask(r.Context(), id, h.pool)

	if err != nil {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Task with id %d, completion changed", id)
		return
	}

	fmt.Fprintf(w, "Task with id %d, completed: %v", id, task.Completed)

}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request)  {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "id must be integer")
		return
	}

	task, err := database.GetTask(r.Context(), id, h.pool)

	if err != nil {
		if errors.Is(err, database.ErrTaskNotFound) {
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

	if err != nil{
		fmt.Println("error:",err)
	}

}