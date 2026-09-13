package models

import "time"

type Task struct {
	Id          int       `json:"id"`
	Creator     string    `json:"creator"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}
