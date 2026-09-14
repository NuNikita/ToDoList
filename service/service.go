package service

import (
	"ToDoList/database"
	"ToDoList/models"
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool *pgxpool.Pool
}

func CreateService(pool *pgxpool.Pool) Service {
	return Service{pool: pool}
}

var ErrEmptyCreator = errors.New("empty name")
var ErrEmptyTitle = errors.New("empty title")

func (s *Service) GetTasks(ctx context.Context, creator string) ([]models.Task, error) {
	if strings.TrimSpace(creator) == "" {
		return nil, ErrEmptyCreator
	}

	tasks, err := database.GetTasks(ctx, creator, s.pool)
	return tasks, err
}

func (s *Service) GetTask(ctx context.Context, id int) (models.Task, error) {
	return database.GetTask(ctx, id, s.pool)
}

func (s *Service) AddTask(ctx context.Context, creator, title, descr string) error {
	if strings.TrimSpace(creator) == "" {
		return ErrEmptyCreator
	}
	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}

	return database.AddTask(ctx, creator, title, descr, s.pool)
}

func (s *Service) DeleteTask(ctx context.Context, id int) error {
	return database.DeleteTask(ctx, id, s.pool)
}

func (s *Service) EditDescriptionTask(ctx context.Context, id int, descr string) error {
	return database.EditDescriptionTask(ctx, id, descr, s.pool)
}

func (s *Service) CompleteTask(ctx context.Context, id int) error {
	return database.CompleteTask(ctx, id, s.pool)
}
