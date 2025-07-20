package usecase

import "github.com/google/uuid"

type AuthService interface {
	GeneratePasswordHash(password string) (string, error)
	VerifyPassword(hashedPassword, inputPassword string) error
	GenerateJWT(userID uuid.UUID) (string, error)
	ValidateJWT(tokenString string) (uuid.UUID, error)
}
