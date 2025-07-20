package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) GenerateJWT(userID uuid.UUID) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthUseCase) ValidateJWT(tokenString string) (uuid.UUID, error) {
	args := m.Called(tokenString)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockAuthUseCase) VerifyPassword(hashedPassword, inputPassword string) error {
	args := m.Called(hashedPassword, inputPassword)
	return args.Error(0)
}

func (m *MockAuthUseCase) GeneratePasswordHash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func TestGenerateJWT(t *testing.T) {
	mockUsecase := new(MockAuthUseCase)
	logger := logrus.New()
	authService := NewAuthService(mockUsecase, logger)

	userID := uuid.New()
	expectedToken := "fake-jwt-token"

	mockUsecase.On("GenerateJWT", userID).Return(expectedToken, nil)

	token, err := authService.GenerateJWT(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockUsecase.AssertExpectations(t)
}

func TestValidateJWT(t *testing.T) {
	mockUsecase := new(MockAuthUseCase)
	logger := logrus.New()
	authService := NewAuthService(mockUsecase, logger)

	token := "valid-jwt-token"
	expectedUserID := uuid.New()

	mockUsecase.On("ValidateJWT", token).Return(expectedUserID, nil)

	userID, err := authService.ValidateJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)
	mockUsecase.AssertExpectations(t)
}

func TestVerifyPassword(t *testing.T) {
	mockUsecase := new(MockAuthUseCase)
	logger := logrus.New()
	authService := NewAuthService(mockUsecase, logger)

	hashedPassword := "hashed-pass"
	inputPassword := "input-pass"

	mockUsecase.On("VerifyPassword", hashedPassword, inputPassword).Return(nil)

	err := authService.VerifyPassword(hashedPassword, inputPassword)
	assert.NoError(t, err)
	mockUsecase.AssertExpectations(t)
}

func TestGeneratePasswordHash(t *testing.T) {
	mockUsecase := new(MockAuthUseCase)
	logger := logrus.New()
	authService := NewAuthService(mockUsecase, logger)

	password := "SecurePass123!"
	expectedHash := "hashed-password"

	mockUsecase.On("GeneratePasswordHash", password).Return(expectedHash, nil)

	hash, err := authService.GeneratePasswordHash(password)
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, hash)
	mockUsecase.AssertExpectations(t)
}
