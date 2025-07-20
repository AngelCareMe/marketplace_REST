package service

import (
	"context"
	"fmt"
	"marketplace/internal/entity"
	usecase "marketplace/internal/usecase/user"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	userUsecase usecase.UserUseCaseRepo
	logger      *logrus.Logger
}

func NewUserService(userUsecase usecase.UserUseCaseRepo, logger *logrus.Logger) *UserService {
	return &UserService{
		userUsecase: userUsecase,
		logger:      logger,
	}
}

func (s *UserService) Register(ctx context.Context, username, password string) (*entity.UserDTO, string, error) {
	if username == "" || password == "" {
		return nil, "", fmt.Errorf("username and password are required")
	}

	userDTO, token, err := s.userUsecase.Register(ctx, username, password)
	if err != nil {
		s.logger.WithError(err).Error("Failed to register user")
		return nil, "", err
	}

	s.logger.WithFields(logrus.Fields{
		"username": username,
		"user_id":  userDTO.ID,
	}).Info("User registered successfully")

	return userDTO, token, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*entity.UserDTO, string, error) {
	if username == "" || password == "" {
		return nil, "", fmt.Errorf("username and password are required")
	}

	user, token, err := s.userUsecase.Login(ctx, username, password)
	if err != nil {
		s.logger.WithError(err).Error("Failed to login user")
		return nil, "", err
	}

	s.logger.WithFields(logrus.Fields{
		"username": username,
		"user_id":  user.ID,
	}).Info("User logged in successfully")

	return user, token, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*entity.UserDTO, error) {
	user, err := s.userUsecase.GetByID(ctx, id)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get user")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User fetched successfully")

	return user.ToDTO(), nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, username, password string) error {
	if username == "" && password == "" {
		return fmt.Errorf("no fields to update")
	}

	if err := s.userUsecase.Update(ctx, id, username, password); err != nil {
		s.logger.WithError(err).Error("Failed to update user")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User updated successfully")

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.userUsecase.Delete(ctx, id); err != nil {
		s.logger.WithError(err).Error("Failed to delete user")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User deleted successfully")

	return nil
}
