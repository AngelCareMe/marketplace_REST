package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) ValidateJWT(token string) (uuid.UUID, error) {
	args := m.Called(token)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockAuthService) GenerateJWT(userID uuid.UUID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GeneratePasswordHash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) VerifyPassword(hashedPassword, inputPassword string) error {
	args := m.Called(hashedPassword, inputPassword)
	return args.Error(0)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockAuthSvc := new(MockAuthService)
	logger := logrus.New()
	authHandler := NewAuthHandler(mockAuthSvc, logger)

	// Setup mock for invalid token
	mockAuthSvc.On("ValidateJWT", "invalid_token").Return(uuid.Nil, fmt.Errorf("invalid token"))

	r.Use(authHandler.AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, "success")
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockAuthSvc.AssertExpectations(t)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockAuthSvc := new(MockAuthService)
	validUserID := uuid.New()
	mockAuthSvc.On("ValidateJWT", "valid_token").Return(validUserID, nil)

	logger := logrus.New()
	authHandler := NewAuthHandler(mockAuthSvc, logger)

	r.Use(authHandler.AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		userID, ok := c.Request.Context().Value("user_id").(uuid.UUID)
		assert.True(t, ok)
		assert.Equal(t, validUserID, userID)
		c.JSON(http.StatusOK, "success")
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockAuthSvc.AssertExpectations(t)
}
