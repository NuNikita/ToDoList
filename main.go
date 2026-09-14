package main

import (
	"ToDoList/database"
	"ToDoList/server"
	"ToDoList/service"
	"fmt"
)

func main() {
	port := 8080

	pool, err := database.Connect()

	if err != nil {
		fmt.Println(err)
		return
	}

	defer pool.Close()

	err = database.CreateTodosTable(pool)

	if err != nil {
		fmt.Println(err)
		return
	}

	serv := service.CreateService(pool)

	server.StartServer(port, serv)

}
