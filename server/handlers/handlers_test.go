package handlers

import (
	"ToDoList/auth"
	"ToDoList/models"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeService struct {
	tasks []models.Task
	task  models.Task
	token string

	getTasksErr        error
	getTaskErr         error
	addTaskErr         error
	deleteTaskErr      error
	editDescriptionErr error
	completeTaskErr    error
	loginErr           error
	registerErr        error

	receivedToken       string
	receivedTaskID      int
	receivedTitle       string
	receivedDescription string
	receivedLogin       string
	receivedPassword    string
	registerCalls       int
}

func (f *fakeService) GetTasks(_ context.Context, token string) ([]models.Task, error) {
	f.receivedToken = token
	return f.tasks, f.getTasksErr
}

func (f *fakeService) GetTask(_ context.Context, id int, token string) (models.Task, error) {
	f.receivedTaskID = id
	f.receivedToken = token
	return f.task, f.getTaskErr
}

func (f *fakeService) AddTask(_ context.Context, token, title, description string) error {
	f.receivedToken = token
	f.receivedTitle = title
	f.receivedDescription = description
	return f.addTaskErr
}

func (f *fakeService) DeleteTask(_ context.Context, id int, token string) error {
	f.receivedTaskID = id
	f.receivedToken = token
	return f.deleteTaskErr
}

func (f *fakeService) EditDescriptionTask(_ context.Context, id int, description, token string) error {
	f.receivedTaskID = id
	f.receivedDescription = description
	f.receivedToken = token
	return f.editDescriptionErr
}

func (f *fakeService) CompleteTask(_ context.Context, id int, token string) error {
	f.receivedTaskID = id
	f.receivedToken = token
	return f.completeTaskErr
}

func (f *fakeService) LoginUser(_ context.Context, login, password string) (string, error) {
	f.receivedLogin = login
	f.receivedPassword = password
	return f.token, f.loginErr
}

func (f *fakeService) RegisterUser(_ context.Context, login, password string) error {
	f.registerCalls++
	f.receivedLogin = login
	f.receivedPassword = password
	return f.registerErr
}

func TestMiddlewareParseToken_NoBearerToken(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	res := httptest.NewRecorder()

	MiddlewareParseToken(next).ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("ожидали статус %d, получили %d", http.StatusUnauthorized, res.Code)
	}

	if called {
		t.Fatal("запрос дошёл до следующего хэндлера без токена")
	}

}

func TestMiddlewareParseToken_WithBearerToken(t *testing.T) {
	called := false
	token, err := auth.GenerateToken(1)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		if r.Context().Value("token") != token {
			t.Fatal("токен не добавлен в context")
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()
	MiddlewareParseToken(next).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("ожидали статус %d, получили %d", http.StatusOK, res.Code)
	}
	if !called {
		t.Fatal("запрос не дошёл до следующего хэндлера")
	}
}

func TestRegisterUser_InvalidJSON(t *testing.T) {
	serv := &fakeService{}
	handler := CreateHandler(serv)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("{invalid"))
	res := httptest.NewRecorder()

	handler.RegisterUser(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("ожидали статус %d, получили %d", http.StatusBadRequest, res.Code)
	}
	if serv.registerCalls != 0 {
		t.Fatal("сервис был вызван при невалидном JSON")
	}
}

func TestRegisterUser_Success(t *testing.T) {
	serv := &fakeService{}
	handler := CreateHandler(serv)

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"login":"user1","password":"password"}`))
	res := httptest.NewRecorder()

	handler.RegisterUser(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("ожидали статус %d, получили %d", http.StatusCreated, res.Code)
	}
	if serv.receivedLogin != "user1" || serv.receivedPassword != "password" {
		t.Fatalf("сервис получил неверные данные: login=%q password=%q", serv.receivedLogin, serv.receivedPassword)
	}
}

func TestLoginUser_Success(t *testing.T) {
	serv := &fakeService{token: "test-token"}
	handler := CreateHandler(serv)

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"login":"user1","password":"password"}`))
	res := httptest.NewRecorder()

	handler.LoginUser(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("ожидали статус %d, получили %d", http.StatusOK, res.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response["token"] != "test-token" {
		t.Fatalf("ожидали токен %q, получили %q", "test-token", response["token"])
	}
}
