package server

import (
	"ToDoList/server/handlers"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartServer(port int, pool *pgxpool.Pool) {

	handler := handlers.CreateHandler(pool)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", handler.AddTask)
	mux.HandleFunc("GET /tasks", handler.GetTasks)

	portStr := ":" + strconv.Itoa(port)

	fmt.Println("The server is starting on port", port)

	if err := http.ListenAndServe(portStr, mux); err != nil {
		fmt.Println("error:", err)
		return
	}

}
