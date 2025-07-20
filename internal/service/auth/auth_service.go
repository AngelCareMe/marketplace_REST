package service

import (
	"github.com/google/uuid"
)

type AuthServiceInterface interface {
	GenerateJWT(userID uuid.UUID) (string, error)
	ValidateJWT(tokenString string) (uuid.UUID, error)
	VerifyPassword(hashedPassword, inputPassword string) error
	GeneratePasswordHash(password string) (string, error)
}
