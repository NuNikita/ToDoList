package database

import (
	"ToDoList/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTaskNotFound error = errors.New("task not found")

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

func GetTask(ctx context.Context, id int, pool *pgxpool.Pool) (models.Task, error) {
	query := `
				SELECT id, creator, title, description, completed, created_at
				FROM todos
				WHERE id = $1
`

	var task models.Task
	row := pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&task.Id,
		&task.Creator,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, err
	}

	return task, nil
}

func AddTask(ctx context.Context, creator, title, descr string, pool *pgxpool.Pool) error {
	query := `
	INSERT INTO todos (title, description, creator)
	VALUES ($1, $2, $3)
`

	_, err := pool.Exec(ctx, query, title, descr, creator)

	return err
}

func DeleteTask(ctx context.Context, id int, pool *pgxpool.Pool) error {
	query := `
		DELETE FROM todos
		WHERE id=$1
`
	result, err := pool.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func EditDescriptionTask(ctx context.Context, id int, descr string, pool *pgxpool.Pool) error {
	query := `
		UPDATE todos
		SET description = $1
		WHERE ID = $2
`
	result, err := pool.Exec(ctx, query, descr, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func CompleteTask(ctx context.Context, id int, pool *pgxpool.Pool) error {
	query := `
		UPDATE todos
		SET completed = NOT completed
		WHERE ID = $1
`
	result, err := pool.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil

}
