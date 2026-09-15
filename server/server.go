package server

import (
	"ToDoList/server/handlers"
	"ToDoList/service"
	"fmt"
	"net/http"
	"strconv"
)

func StartServer(port int, serv *service.Service) {

	handler := handlers.CreateHandler(serv)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", handlers.MiddlewareParseToken(handler.AddTask))
	mux.HandleFunc("GET /tasks", handlers.MiddlewareParseToken(handler.GetTasks))
	mux.HandleFunc("GET /tasks/{id}", handlers.MiddlewareParseToken(handler.GetTask))
	mux.HandleFunc("DELETE /tasks/{id}", handlers.MiddlewareParseToken(handler.DeleteTask))
	mux.HandleFunc("PATCH /tasks/{id}/description", handlers.MiddlewareParseToken(handler.EditDescriptionTask))
	mux.HandleFunc("PATCH /tasks/{id}/complete", handlers.MiddlewareParseToken(handler.CompleteTask))
	mux.HandleFunc("POST /login", handler.LoginUser)
	mux.HandleFunc("POST /register", handler.RegisterUser)

	portStr := ":" + strconv.Itoa(port)
	fmt.Println("The server is starting on port", port)

	if err := http.ListenAndServe(portStr, mux); err != nil {
		fmt.Println("error:", err)
		return
	}

}
