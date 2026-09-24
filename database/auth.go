package database

import (
	"ToDoList/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

func CreateUsersTable(pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login TEXT NOT NULL UNIQUE,
			passw_hash TEXT NOT NULL
		);`

	_, err := pool.Exec(context.Background(), query)

	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	return nil

}

func (b *Base) RegisterUser(ctx context.Context, login, passwHash string) error {
	query := `
	INSERT INTO users (login, passw_hash)
	VALUES ($1, $2)
`
	_, err := b.pool.Exec(ctx, query, login, passwHash)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%w: %q", ErrUserAlreadyExists, login)
		}
		return fmt.Errorf("register user %q: %w", login, err)
	}

	return nil
}

func (b *Base) GetUser(ctx context.Context, login string) (models.UserDB, error) {
	query := `
	SELECT * FROM users
	WHERE login = $1
`
	var user models.UserDB

	row := b.pool.QueryRow(ctx, query, login)

	err := row.Scan(&user.Id, &user.Login, &user.PasswHash)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.UserDB{}, ErrUserNotFound
		}
		return models.UserDB{}, fmt.Errorf("get user %q: %w", login, err)
	}

	return user, nil

}
