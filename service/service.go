package service

import (
	"ToDoList/auth"
	"ToDoList/database"
	"ToDoList/models"
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	base *database.Base
}

func CreateService(base *database.Base) *Service {
	return &Service{base: base}
}

func hashPassword(passw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
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
		return false, err
	}

	task, err := s.base.GetTask(ctx, id)

	if err != nil {
		return false, err
	}

	return task.UserId == userId, nil

}

func (s *Service) GetTasks(ctx context.Context, token string) ([]models.Task, error) {
	userId, err := auth.ParseToken(token)

	if err != nil {
		return nil, err
	}

	tasks, err := s.base.GetTasks(ctx, userId)

	return tasks, err
}

func (s *Service) GetTask(ctx context.Context, id int, token string) (models.Task, error) {

	userId, err := auth.ParseToken(token)

	if err != nil {
		return models.Task{}, err
	}

	task, err := s.base.GetTask(ctx, id)

	if err != nil {
		return models.Task{}, err
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
		return err
	}

	return s.base.AddTask(ctx, userId, title, descr)
}

func (s *Service) withOwnershipCheck(ctx context.Context, id int, token string, action func() error) error {
	res, err := s.isYourTask(ctx, id, token)
	if err != nil {
		return err
	}

	if !res {
		return ErrNotYourTask
	}
	return action()

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
		return err
	}

	return s.base.RegisterUser(ctx, login, passwHash)
}

func (s *Service) LoginUser(ctx context.Context, login, passw string) (string, error) {
	user, err := s.base.GetUser(ctx, login)

	if err != nil {
		return "", err
	}

	res := checkPassword(user.PasswHash, passw)

	if !res {
		return "", ErrWrongPassword
	}
	return auth.GenerateToken(user.Id)

}
