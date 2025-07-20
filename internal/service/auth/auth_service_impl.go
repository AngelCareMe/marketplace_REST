package service

import (
	"fmt"
	usecase "marketplace/internal/usecase/auth"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	authRepo usecase.AuthService
	logger   *logrus.Logger
}

func NewAuthService(authRepo usecase.AuthService, logger *logrus.Logger) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		logger:   logger,
	}
}

func (s *AuthService) GenerateJWT(userID uuid.UUID) (string, error) {
	if userID == uuid.Nil {
		return "", fmt.Errorf("userID cannot be empty")
	}

	token, err := s.authRepo.GenerateJWT(userID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate JWT")
		return "", err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": userID,
	}).Info("JWT generated successfully")

	return token, nil
}

func (s *AuthService) ValidateJWT(tokenString string) (uuid.UUID, error) {
	if tokenString == "" {
		return uuid.Nil, fmt.Errorf("token cannot be empty")
	}

	userID, err := s.authRepo.ValidateJWT(tokenString)
	if err != nil {
		s.logger.WithError(err).Error("Failed to validate JWT")
		return uuid.Nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": userID,
	}).Info("JWT validated successfully")

	return userID, nil
}

func (s *AuthService) VerifyPassword(hashedPassword, inputPassword string) error {
	if hashedPassword == "" || inputPassword == "" {
		return fmt.Errorf("hashed password and input password cannot be empty")
	}

	err := s.authRepo.VerifyPassword(hashedPassword, inputPassword)
	if err != nil {
		s.logger.WithError(err).Error("Failed to verify password")
		return err
	}

	s.logger.Info("Password verified successfully")
	return nil
}

func (s *AuthService) GeneratePasswordHash(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	hashedPassword, err := s.authRepo.GeneratePasswordHash(password)
	if err != nil {
		s.logger.WithError(err).Error("Failed to generate password hash")
		return "", err
	}

	s.logger.Info("Password hash generated successfully")
	return hashedPassword, nil
}
