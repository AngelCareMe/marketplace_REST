package service

import (
	"context"
	"marketplace/internal/entity"

	"github.com/google/uuid"
)

type UserServiceInterface interface {
	Register(ctx context.Context, username, password string) (*entity.UserDTO, string, error)
	Login(ctx context.Context, username, password string) (*entity.UserDTO, string, error)
	GetUser(ctx context.Context, id uuid.UUID) (*entity.UserDTO, error)
	UpdateUser(ctx context.Context, id uuid.UUID, username, password string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
