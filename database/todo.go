package database

import (
	"ToDoList/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrTaskNotFound error = errors.New("task not found")

type Base struct {
	pool *pgxpool.Pool
}

func CreateBasePool(pool *pgxpool.Pool) *Base {
	return &Base{pool: pool}
}

func CreateTodosTable(pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			user_id INTEGER NOT NULL,
	                                 
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`

	_, err := pool.Exec(context.Background(), query)

	return err

}

func (b *Base) GetTasks(ctx context.Context, userId int) ([]models.Task, error) {
	query := `
				SELECT id, title, description, completed, created_at, user_id
				FROM todos
				WHERE user_id = $1
				ORDER BY created_at DESC
`

	rows, err := b.pool.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		err = rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.UserId,
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

func (b *Base) GetTask(ctx context.Context, id int) (models.Task, error) {
	query := `
				SELECT id, title, description, completed, created_at, user_id
				FROM todos
				WHERE id = $1
`

	var task models.Task
	row := b.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UserId,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, err
	}

	return task, nil
}

func (b *Base) AddTask(ctx context.Context, userId int, title, descr string) error {
	query := `
	INSERT INTO todos (title, description, user_id)
	VALUES ($1, $2, $3)
`

	_, err := b.pool.Exec(ctx, query, title, descr, userId)

	return err
}

func (b *Base) DeleteTask(ctx context.Context, id int) error {
	query := `
		DELETE FROM todos
		WHERE id=$1
`
	result, err := b.pool.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (b *Base) EditDescriptionTask(ctx context.Context, id int, descr string) error {
	query := `
		UPDATE todos
		SET description = $1
		WHERE ID = $2
`
	result, err := b.pool.Exec(ctx, query, descr, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}

func (b *Base) CompleteTask(ctx context.Context, id int) error {
	query := `
		UPDATE todos
		SET completed = NOT completed
		WHERE ID = $1
`
	result, err := b.pool.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil

}

func (b *Base) GetIdTasks(ctx context.Context, userId int) ([]int, error) {
	query := `
		SELECT id
		FROM todos
		WHERE user_id = $1
`
	rows, err := b.pool.Query(ctx, query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var ids []int

	for rows.Next() {

		var id int
		err = rows.Scan(&id)

		if err != nil {
			return nil, err
		}

		ids = append(ids, id)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil

}
