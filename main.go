package main

import (
	"ToDoList/database"
	"ToDoList/server"
	"ToDoList/service"
	"fmt"
)

// Обернуть ошибки все

func main() {
	port := 8080

	pool, err := database.Connect()

	if err != nil {
		fmt.Println(err)
		return
	}

	defer pool.Close()

	err = database.CreateUsersTable(pool)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = database.CreateTodosTable(pool)

	if err != nil {
		fmt.Println(err)
		return
	}

	base := database.CreateBasePool(pool)

	serv := service.CreateService(base)

	server.StartServer(port, serv)

}
