package usecase

import (
	"context"
	"marketplace/internal/entity"

	"github.com/google/uuid"
)

type UserUseCaseRepo interface {
	Register(ctx context.Context, username, password string) (*entity.UserDTO, string, error)
	Login(ctx context.Context, username, password string) (*entity.UserDTO, string, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	Update(ctx context.Context, id uuid.UUID, username, password string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
