package main

import (
	"ToDoList/auth"
	"ToDoList/database"
	"ToDoList/server"
	"ToDoList/service"
	"fmt"
)

func main() {

	err := auth.CheckSecretKey()
	if err != nil {
		fmt.Println(err)
		return
	}

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

	if err := server.StartServer(port, serv); err != nil {
		fmt.Println(err)
	}

}
