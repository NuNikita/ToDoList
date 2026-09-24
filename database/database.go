package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {

	ctx := context.Background()

	pool, err := pgxpool.New(
		ctx,
		os.Getenv("DATABASE_URL"),
	)

	if err != nil {
		return nil, fmt.Errorf("create database connection pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	fmt.Println("Еее подключились")
	return pool, nil
}
