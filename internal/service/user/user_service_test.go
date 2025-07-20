package service

import (
	"context"
	"testing"

	"marketplace/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Register(ctx context.Context, username, password string) (*entity.UserDTO, string, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(*entity.UserDTO), args.String(1), args.Error(2)
}

func (m *MockUserUseCase) Login(ctx context.Context, username, password string) (*entity.UserDTO, string, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(*entity.UserDTO), args.String(1), args.Error(2)
}

func (m *MockUserUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) Update(ctx context.Context, id uuid.UUID, username, password string) error {
	args := m.Called(ctx, id, username, password)
	return args.Error(0)
}

func (m *MockUserUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestRegister(t *testing.T) {
	mockUsecase := new(MockUserUseCase)
	logger := logrus.New()
	userService := NewUserService(mockUsecase, logger)

	username := "testuser"
	password := "SecurePass123!"
	expectedUser := &entity.UserDTO{
		ID:       uuid.New(),
		Username: username,
	}
	expectedToken := "fake-jwt-token"

	mockUsecase.On("Register", mock.Anything, username, password).
		Return(expectedUser, expectedToken, nil)

	user, token, err := userService.Register(context.Background(), username, password)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.Equal(t, expectedToken, token)
	mockUsecase.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	mockUsecase := new(MockUserUseCase)
	logger := logrus.New()
	userService := NewUserService(mockUsecase, logger)

	username := "testuser"
	password := "SecurePass123!"
	expectedUser := &entity.UserDTO{
		ID:       uuid.New(),
		Username: username,
	}
	expectedToken := "fake-jwt-token"

	mockUsecase.On("Login", mock.Anything, username, password).
		Return(expectedUser, expectedToken, nil)

	user, token, err := userService.Login(context.Background(), username, password)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.Equal(t, expectedToken, token)
	mockUsecase.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	mockUsecase := new(MockUserUseCase)
	logger := logrus.New()
	userService := NewUserService(mockUsecase, logger)

	userID := uuid.New()
	expectedUser := &entity.User{
		ID:       userID,
		Username: "testuser",
	}

	mockUsecase.On("GetByID", mock.Anything, userID).Return(expectedUser, nil)

	userDTO, err := userService.GetUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, userDTO)
	assert.Equal(t, expectedUser.ID, userDTO.ID)
	assert.Equal(t, expectedUser.Username, userDTO.Username)
	mockUsecase.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockUsecase := new(MockUserUseCase)
	logger := logrus.New()
	userService := NewUserService(mockUsecase, logger)

	userID := uuid.New()
	username := "new_username"
	password := "NewPass123!"

	mockUsecase.On("Update", mock.Anything, userID, username, password).Return(nil)

	err := userService.UpdateUser(context.Background(), userID, username, password)
	assert.NoError(t, err)
	mockUsecase.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockUsecase := new(MockUserUseCase)
	logger := logrus.New()
	userService := NewUserService(mockUsecase, logger)

	userID := uuid.New()

	mockUsecase.On("Delete", mock.Anything, userID).Return(nil)

	err := userService.DeleteUser(context.Background(), userID)
	assert.NoError(t, err)
	mockUsecase.AssertExpectations(t)
}
