package adapter

import (
	"context"
	"marketplace/internal/entity"

	"github.com/google/uuid"
)

type PostAdapterInterface interface {
	Create(ctx context.Context, post *entity.Post) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Post, error)
	ListByAuthorID(ctx context.Context, authorID uuid.UUID, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error)
	ListPosts(ctx context.Context, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error)
	GetByHeaderAndContent(ctx context.Context, header, content string) (*entity.Post, error)
	Update(ctx context.Context, post *entity.Post) error
	Delete(ctx context.Context, id uuid.UUID) error
}
