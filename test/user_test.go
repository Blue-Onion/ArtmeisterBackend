package test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Blue-Onion/RestApi-Go/handler/user"
	"github.com/Blue-Onion/RestApi-Go/internal/database"
	"github.com/Blue-Onion/RestApi-Go/model"
	"github.com/Blue-Onion/RestApi-Go/utlis"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockRepo struct {
	database.UserRepository
	Users []database.User
}

func (m *MockRepo) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.CreateUserRow, error) {
	user := database.User{
		ID:        arg.ID,
		Name:      arg.Name,
		Email:     arg.Email,
		Password:  arg.Password,
		Createdat: arg.Createdat,
		Updatedat: arg.Updatedat,
	}
	m.Users = append(m.Users, user)
	return database.CreateUserRow{
		ID:        user.ID,
		Name:      user.Name,
		Createdat: user.Createdat,
		Updatedat: user.Updatedat,
	}, nil
}

func (m *MockRepo) GetUserByEmail(ctx context.Context, email string) (database.GetUserByEmailRow, error) {
	for _, u := range m.Users {
		if u.Email == email {
			return database.GetUserByEmailRow{
				ID:       u.ID,
				Password: u.Password,
			}, nil
		}
	}
	return database.GetUserByEmailRow{}, nil
}

func TestHandleCreateUser(t *testing.T) {
	mockRepo := &MockRepo{}
	h := &user.Handler{Repo: mockRepo}

	userData := model.CreateUser{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(userData)

	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.HandleCreateUser(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdUser database.CreateUserRow
	json.NewDecoder(rr.Body).Decode(&createdUser)

	if createdUser.Name != userData.Name {
		t.Errorf("handler returned unexpected body: got %v want %v", createdUser.Name, userData.Name)
	}
}

func TestHandleLogin(t *testing.T) {
	mockRepo := &MockRepo{}
	h := &user.Handler{Repo: mockRepo}

	// First create a user to login with
	password := "password123"
	hash, _ := utlis.HashPassword(password)
	userData := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  hash,
		Createdat: time.Now(),
		Updatedat: time.Now(),
	}
	mockRepo.CreateUser(context.Background(), userData)

	loginData := model.AutheticateUser{
		Email:    "test@example.com",
		Password: password,
	}
	body, _ := json.Marshal(loginData)

	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.HandleLogin(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check for cookie
	cookies := rr.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "authToken" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("authToken cookie not found in response")
	}
}

func TestHandleLogOut(t *testing.T) {
	mockRepo := &MockRepo{}
	h := &user.Handler{Repo: mockRepo}

	req, _ := http.NewRequest("POST", "/api/logOut", nil)
	rr := httptest.NewRecorder()

	h.HandleLogOut(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check for cookie removal
	cookies := rr.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "auth_token" && c.MaxAge == -1 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("auth_token cookie removal not found in response")
	}
}
