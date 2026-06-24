package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Blue-Onion/ArtmeisterBackend/handler/user"
	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
	"github.com/Blue-Onion/ArtmeisterBackend/model"
	"github.com/Blue-Onion/ArtmeisterBackend/utlis"
	"github.com/google/uuid"
)

type MockRepo struct {
	database.UserRepository
	Users []database.User
}

func (m *MockRepo) CreateUser(ctx context.Context, arg database.CreateUserParams) (uuid.UUID, error) {
	now := time.Now()
	u := database.User{
		ID:        uuid.New(),
		Name:      arg.Name,
		Email:     arg.Email,
		Password:  arg.Password,
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.Users = append(m.Users, u)
	return u.ID, nil
}

func (m *MockRepo) GetUserByEmail(ctx context.Context, email string) (database.GetUserByEmailRow, error) {
	for _, u := range m.Users {
		if u.Email == email {
			return database.GetUserByEmailRow{
				ID:       u.ID,
				Name:     u.Name,
				Email:    u.Email,
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

	var response struct {
		Success bool
		Data    map[string]string
	}
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Data["ID"] == "" {
		t.Errorf("expected non-empty ID in response")
	}
}

func TestHandleLogin(t *testing.T) {
	mockRepo := &MockRepo{}
	h := &user.Handler{Repo: mockRepo}

	// First create a user to login with
	password := "password123"
	hash, _ := utlis.HashPassword(password)
	userData := database.CreateUserParams{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: hash,
	}
	mockRepo.CreateUser(context.Background(), userData)

	loginData := model.AuthenticateUser{
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
		if c.Name == "authToken" && c.MaxAge == -1 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("authToken cookie removal not found in response")
	}
}
