package service

import (
	"context"
	"marketplace/internal/entity"

	"github.com/google/uuid"
)

type PostServiceInterface interface {
	CreatePost(ctx context.Context, authorID uuid.UUID, header, content, image string, price float64) (*entity.Post, error)
	EditPost(ctx context.Context, postID uuid.UUID, header, content, image string, price float64) (*entity.Post, error)
	DeletePost(ctx context.Context, postID uuid.UUID) error
	GetPost(ctx context.Context, postID uuid.UUID) (*entity.Post, error)
	ListPosts(ctx context.Context, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error)
	ListPostsByAuthor(ctx context.Context, authorID uuid.UUID, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error)
}
