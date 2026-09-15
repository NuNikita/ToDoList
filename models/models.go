package models

import "time"

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UserId      int       `json:"user_id"`
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserDB struct {
	Id        int
	Login     string
	PasswHash string
}
