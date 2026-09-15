package database

import (
	"ToDoList/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUsersTable(pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login TEXT NOT NULL UNIQUE,
			passw_hash TEXT NOT NULL
		);`

	_, err := pool.Exec(context.Background(), query)

	return err

}

func (b *Base) RegisterUser(ctx context.Context, login, passwHash string) error {
	query := `
	INSERT INTO users (login, passw_hash)
	VALUES ($1, $2)
`
	_, err := b.pool.Exec(ctx, query, login, passwHash)

	return err
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
		return models.UserDB{}, err
	}

	return user, nil

}
