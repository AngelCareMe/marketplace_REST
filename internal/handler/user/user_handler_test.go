package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"marketplace/internal/entity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, username, password string) (*entity.UserDTO, string, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(*entity.UserDTO), args.String(1), args.Error(2)
}

func (m *MockUserService) Login(ctx context.Context, username, password string) (*entity.UserDTO, string, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(*entity.UserDTO), args.String(1), args.Error(2)
}

func (m *MockUserService) GetUser(ctx context.Context, id uuid.UUID) (*entity.UserDTO, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.UserDTO), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uuid.UUID, username, password string) error {
	args := m.Called(ctx, id, username, password)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestRegisterUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockUserSvc := new(MockUserService)
	logger := logrus.New()
	handler := NewUserHandler(mockUserSvc, logger)

	r.POST("/users/register", handler.Register)

	reqBody := map[string]string{
		"username": "testuser",
		"password": "SecurePass123!",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Настройка ожидаемого поведения мока
	user := &entity.UserDTO{
		ID:       uuid.New(),
		Username: "testuser",
	}
	token := "fake-jwt-token"
	mockUserSvc.On("Register", mock.Anything, "testuser", "SecurePass123!").
		Return(user, token, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUserSvc.AssertExpectations(t)
}
