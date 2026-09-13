package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {

	ctx := context.Background()

	pool, err := pgxpool.New(
		ctx,
		"postgres://postgres:1234@localhost:5432/ToDoBase",
	)

	if err != nil {
		return nil, fmt.Errorf("не удалось создать пул: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("не удалось подключиться к бд: %w", err)
	}

	fmt.Println("Еее подключились")
	return pool, nil
}
