package database

import (
	"ToDoList/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTodosTable(pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			creator TEXT NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`

	_, err := pool.Exec(context.Background(), query)

	return err

}

func GetTasks(ctx context.Context, creator string, pool *pgxpool.Pool) ([]models.Task, error) {
	query := `
				SELECT id, creator, title, description, completed, created_at
				FROM todos
				WHERE creator = $1
				ORDER BY created_at DESC
`

	rows, err := pool.Query(ctx, query, creator)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		err = rows.Scan(
			&task.Id,
			&task.Creator,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil

}

func AddTask(ctx context.Context, creator, title, descr string, pool *pgxpool.Pool) error {
	query := `INSERT INTO todos (title, description, creator) VALUES ($1, $2, $3)`

	_, err := pool.Exec(ctx, query, title, descr, creator)

	return err
}
