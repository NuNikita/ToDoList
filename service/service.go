package service

import (
	"ToDoList/auth"
	"ToDoList/database"
	"ToDoList/models"
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type repository interface {
	GetTasks(ctx context.Context, userID int) ([]models.Task, error)
	GetTask(ctx context.Context, id int) (models.Task, error)
	AddTask(ctx context.Context, userID int, title, description string) error
	DeleteTask(ctx context.Context, id int) error
	EditDescriptionTask(ctx context.Context, id int, description string) error
	CompleteTask(ctx context.Context, id int) error

	RegisterUser(ctx context.Context, login, passwordHash string) error
	GetUser(ctx context.Context, login string) (models.UserDB, error)
}

type Service struct {
	base repository
}

func CreateService(base repository) *Service {
	return &Service{base: base}
}

func hashPassword(passw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passw), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

func checkPassword(hash, passw string) bool {
	byteHash := []byte(hash)
	bytePassw := []byte(passw)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePassw)

	if err != nil {
		return false
	}
	return true

}

var ErrEmptyTitle = errors.New("empty title")
var ErrWrongPassword = errors.New("wrong password")
var ErrEmptyField = errors.New("empty field")
var ErrNotYourTask = errors.New("not your task")

func (s *Service) isYourTask(ctx context.Context, id int, token string) (bool, error) {
	userId, err := auth.ParseToken(token)

	if err != nil {
		return false, fmt.Errorf("parse authentication token: %w", err)
	}

	task, err := s.base.GetTask(ctx, id)

	if err != nil {
		return false, fmt.Errorf("get task %d for ownership check: %w", id, err)
	}

	return task.UserId == userId, nil

}

func (s *Service) GetTasks(ctx context.Context, token string) ([]models.Task, error) {
	userId, err := auth.ParseToken(token)

	if err != nil {
		return nil, fmt.Errorf("parse authentication token: %w", err)
	}

	tasks, err := s.base.GetTasks(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get tasks for user %d: %w", userId, err)
	}

	return tasks, nil
}

func (s *Service) GetTask(ctx context.Context, id int, token string) (models.Task, error) {

	userId, err := auth.ParseToken(token)

	if err != nil {
		return models.Task{}, fmt.Errorf("parse authentication token: %w", err)
	}

	task, err := s.base.GetTask(ctx, id)

	if err != nil {
		return models.Task{}, fmt.Errorf("get task %d: %w", id, err)
	}

	if task.UserId == userId {
		return task, nil
	}

	return models.Task{}, ErrNotYourTask

}

func (s *Service) AddTask(ctx context.Context, token, title, descr string) error {

	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}

	userId, err := auth.ParseToken(token)

	if err != nil {
		return fmt.Errorf("parse authentication token: %w", err)
	}

	if err := s.base.AddTask(ctx, userId, title, descr); err != nil {
		return fmt.Errorf("add task for user %d: %w", userId, err)
	}

	return nil
}

func (s *Service) withOwnershipCheck(ctx context.Context, id int, token string, action func() error) error {
	res, err := s.isYourTask(ctx, id, token)
	if err != nil {
		return fmt.Errorf("check task %d ownership: %w", id, err)
	}

	if !res {
		return ErrNotYourTask
	}
	if err = action(); err != nil {
		return fmt.Errorf("apply action to task %d: %w", id, err)
	}

	return nil

}

func (s *Service) DeleteTask(ctx context.Context, id int, token string) error {

	return s.withOwnershipCheck(ctx, id, token, func() error {
		return s.base.DeleteTask(ctx, id)
	})

}

func (s *Service) EditDescriptionTask(ctx context.Context, id int, descr, token string) error {
	return s.withOwnershipCheck(ctx, id, token, func() error {
		return s.base.EditDescriptionTask(ctx, id, descr)
	})

}

func (s *Service) CompleteTask(ctx context.Context, id int, token string) error {
	return s.withOwnershipCheck(ctx, id, token, func() error {
		return s.base.CompleteTask(ctx, id)
	})

}

func (s *Service) RegisterUser(ctx context.Context, login, passw string) error {
	if strings.TrimSpace(login) == "" {
		return ErrEmptyField
	}
	if strings.TrimSpace(passw) == "" {
		return ErrEmptyField
	}

	passwHash, err := hashPassword(passw)

	if err != nil {
		return fmt.Errorf("hash password for user %q: %w", login, err)
	}

	if err := s.base.RegisterUser(ctx, login, passwHash); err != nil {
		return fmt.Errorf("register user %q: %w", login, err)
	}

	return nil
}

func (s *Service) LoginUser(ctx context.Context, login, passw string) (string, error) {
	user, err := s.base.GetUser(ctx, login)

	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			return "", ErrWrongPassword
		}
		return "", fmt.Errorf("get user %q for login: %w", login, err)
	}

	res := checkPassword(user.PasswHash, passw)

	if !res {
		return "", ErrWrongPassword
	}
	token, err := auth.GenerateToken(user.Id)
	if err != nil {
		return "", fmt.Errorf("generate token for user %q: %w", login, err)
	}

	return token, nil

}
