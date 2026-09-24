package service

import (
	"ToDoList/auth"
	"ToDoList/database"
	"ToDoList/models"
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	tasks      map[int]models.Task
	users      map[string]models.UserDB
	nextTaskID int
	nextUserID int
}

var (
	errFakeUserNotFound = errors.New("user not found")
	errFakeLoginExists  = errors.New("login already exists")
)

func (f *fakeRepository) GetTasks(_ context.Context, userID int) ([]models.Task, error) {
	var tasks []models.Task
	for _, task := range f.tasks {
		if task.UserId == userID {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (f *fakeRepository) GetTask(_ context.Context, id int) (models.Task, error) {
	task, ok := f.tasks[id]
	if !ok {
		return models.Task{}, database.ErrTaskNotFound
	}
	return task, nil
}

func (f *fakeRepository) AddTask(_ context.Context, userID int, title, description string) error {
	if f.nextTaskID == 0 {
		f.nextTaskID = 1
	}
	id := f.nextTaskID
	f.nextTaskID++
	f.tasks[id] = models.Task{Id: id, Title: title, Description: description, UserId: userID}
	return nil
}

func (f *fakeRepository) DeleteTask(_ context.Context, id int) error {
	if _, ok := f.tasks[id]; !ok {
		return database.ErrTaskNotFound
	}
	delete(f.tasks, id)
	return nil
}

func (f *fakeRepository) EditDescriptionTask(_ context.Context, id int, description string) error {
	task, ok := f.tasks[id]
	if !ok {
		return database.ErrTaskNotFound
	}
	task.Description = description
	f.tasks[id] = task
	return nil
}

func (f *fakeRepository) CompleteTask(_ context.Context, id int) error {
	task, ok := f.tasks[id]
	if !ok {
		return database.ErrTaskNotFound
	}
	task.Completed = !task.Completed
	f.tasks[id] = task
	return nil
}

func (f *fakeRepository) RegisterUser(_ context.Context, login, passwordHash string) error {
	if _, ok := f.users[login]; ok {
		return errFakeLoginExists
	}
	if f.nextUserID == 0 {
		f.nextUserID = 1
	}
	f.users[login] = models.UserDB{Id: f.nextUserID, Login: login, PasswHash: passwordHash}
	f.nextUserID++
	return nil
}

func (f *fakeRepository) GetUser(_ context.Context, login string) (models.UserDB, error) {
	user, ok := f.users[login]
	if !ok {
		return models.UserDB{}, errFakeUserNotFound
	}
	return user, nil
}

func TestAddTask_EmptyTitle(t *testing.T) {
	repo := &fakeRepository{tasks: map[int]models.Task{}, users: map[string]models.UserDB{}}
	serv := CreateService(repo)

	err := serv.AddTask(context.Background(), "", "   ", "описание")

	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("ожидали ошибку %v, получили %v", ErrEmptyTitle, err)
	}
}

func TestAddTask_Success(t *testing.T) {
	repo := &fakeRepository{tasks: map[int]models.Task{}, users: map[string]models.UserDB{}}
	token, err := auth.GenerateToken(1)

	if err != nil {
		t.Fatal(err)
	}

	serv := CreateService(repo)

	err = serv.AddTask(context.Background(), token, "Title1", "Description1")

	if err != nil {
		t.Fatal(err)
	}

	task, ok := repo.tasks[1]

	if !ok {
		t.Fatal("Задача не добавлена")
	}

	if task.Title != "Title1" {
		t.Fatalf("Expected Title1 but got %v", task.Title)
	}

	if task.Description != "Description1" {
		t.Fatalf("Expected Description1 but got %v", task.Description)
	}

	if task.UserId != 1 {
		t.Fatalf("Expected user ID 1 but got %v", task.UserId)
	}

}

func TestGetTasks_ReturnsOnlyCurrentUserTasks(t *testing.T) {
	repo := &fakeRepository{tasks: map[int]models.Task{}, users: map[string]models.UserDB{}}
	repo.tasks[1] = models.Task{Title: "Title1", Description: "Description1", UserId: 1}
	repo.tasks[2] = models.Task{Title: "Title2", Description: "Description2", UserId: 2}
	serv := CreateService(repo)
	token, err := auth.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	tasks, err := serv.GetTasks(context.Background(), token)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task but got %v", len(tasks))
	}

	if tasks[0].Title != "Title1" || tasks[0].UserId != 1 {
		t.Fatalf("вернулась не та задача: %+v", tasks[0])
	}

}

func TestGetTask_NotOwner(t *testing.T) {
	repo := &fakeRepository{
		tasks: map[int]models.Task{
			1: {Id: 1, Title: "Чужая задача", UserId: 2},
		},
	}

	serv := CreateService(repo)

	token, err := auth.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = serv.GetTask(context.Background(), 1, token)

	if !errors.Is(err, ErrNotYourTask) {
		t.Fatalf("ожидали %v, получили %v", ErrNotYourTask, err)
	}
}

func TestDeleteTask_NotOwner(t *testing.T) {
	repo := &fakeRepository{
		tasks: map[int]models.Task{
			1: {Id: 1, Title: "Чужая задача", UserId: 2},
		}}
	serv := CreateService(repo)
	token, err := auth.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	err = serv.DeleteTask(context.Background(), 1, token)

	if !errors.Is(err, ErrNotYourTask) {
		t.Fatalf("ожидали %v, получили %v", ErrNotYourTask, err)
	}

	if _, ok := repo.tasks[1]; !ok {
		t.Fatal("чужая задача была удалена")
	}

}

func TestEditDescriptionTask_Success(t *testing.T) {
	repo := &fakeRepository{
		tasks: map[int]models.Task{
			1: {Id: 1, Title: "Моя задача", Description: "Старое", UserId: 1},
		},
		users: map[string]models.UserDB{},
	}

	serv := CreateService(repo)

	token, err := auth.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	err = serv.EditDescriptionTask(context.Background(), 1, "Новое", token)
	if err != nil {
		t.Fatal(err)
	}

	if repo.tasks[1].Description != "Новое" {
		t.Fatalf("ожидали новое описание, получили %q", repo.tasks[1].Description)
	}
}

func TestCompleteTask_Success(t *testing.T) {
	repo := &fakeRepository{
		tasks: map[int]models.Task{
			1: {Id: 1, Title: "Моя задача", Completed: false, UserId: 1},
		},
		users: map[string]models.UserDB{},
	}

	serv := CreateService(repo)

	token, err := auth.GenerateToken(1)
	if err != nil {
		t.Fatal(err)
	}

	err = serv.CompleteTask(context.Background(), 1, token)
	if err != nil {
		t.Fatal(err)
	}

	if !repo.tasks[1].Completed {
		t.Fatal("статус задачи не изменился на completed")
	}
}

func TestRegisterUser_EmptyFields(t *testing.T) {

	type test struct {
		name     string
		login    string
		password string
	}

	tests := []test{
		{"пустой логин", "", "password"},
		{"логин из пробелов", "   ", "password"},
		{"пустой пароль", "user", ""},
		{"пароль из пробелов", "user", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepository{
				tasks: map[int]models.Task{},
				users: map[string]models.UserDB{},
			}
			serv := CreateService(repo)

			err := serv.RegisterUser(context.Background(), tt.login, tt.password)

			if !errors.Is(err, ErrEmptyField) {
				t.Fatalf("ожидали %v, получили %v", ErrEmptyField, err)
			}
		})
	}
}

func TestRegisterUser_Success(t *testing.T) {
	repo := &fakeRepository{
		tasks: map[int]models.Task{},
		users: map[string]models.UserDB{},
	}
	serv := CreateService(repo)

	err := serv.RegisterUser(context.Background(), "user1", "my-password")
	if err != nil {
		t.Fatal(err)
	}

	user, ok := repo.users["user1"]
	if !ok {
		t.Fatal("пользователь не сохранён")
	}

	if user.PasswHash == "my-password" {
		t.Fatal("пароль сохранён в открытом виде")
	}
}

func TestLoginUser_WrongPassword(t *testing.T) {
	repo := &fakeRepository{
		tasks: map[int]models.Task{},
		users: map[string]models.UserDB{},
	}
	serv := CreateService(repo)

	err := serv.RegisterUser(context.Background(), "user1", "correct-password")
	if err != nil {
		t.Fatal(err)
	}

	_, err = serv.LoginUser(context.Background(), "user1", "wrong-password")

	if !errors.Is(err, ErrWrongPassword) {
		t.Fatalf("ожидали %v, получили %v", ErrWrongPassword, err)
	}
}

func TestLoginUser_Success(t *testing.T) {
	repo := &fakeRepository{
		tasks: map[int]models.Task{},
		users: map[string]models.UserDB{},
	}
	serv := CreateService(repo)

	err := serv.RegisterUser(context.Background(), "user1", "correct-password")
	if err != nil {
		t.Fatal(err)
	}

	token, err := serv.LoginUser(context.Background(), "user1", "correct-password")
	if err != nil {
		t.Fatal(err)
	}

	userID, err := auth.ParseToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if userID != repo.users["user1"].Id {
		t.Fatalf("ожидали ID %d, получили %d", repo.users["user1"].Id, userID)
	}
}
