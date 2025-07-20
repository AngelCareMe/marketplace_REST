package service

import (
	"context"
	"fmt"
	"marketplace/internal/entity"
	usecasePost "marketplace/internal/usecase/post"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PostService struct {
	postUsecase usecasePost.PostUseCaseRepo
	logger      *logrus.Logger
}

func NewPostService(postUsecase usecasePost.PostUseCaseRepo, logger *logrus.Logger) *PostService {
	return &PostService{
		postUsecase: postUsecase,
		logger:      logger,
	}
}

func (s *PostService) CreatePost(ctx context.Context, authorID uuid.UUID, header, content, image string, price float64) (*entity.Post, error) {
	if header == "" || content == "" || price <= 0 {
		return nil, fmt.Errorf("header, content, and valid price are required")
	}

	post, err := s.postUsecase.Publish(ctx, authorID, header, content, image, price)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create post")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"post_id":   post.ID,
		"author_id": authorID,
		"header":    header,
	}).Info("Post created successfully")

	return post, nil
}

func (s *PostService) EditPost(ctx context.Context, postID uuid.UUID, header, content, image string, price float64) (*entity.Post, error) {
	if header == "" && content == "" && image == "" && price <= 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	post, err := s.postUsecase.Edit(ctx, postID, header, content, image, price)
	if err != nil {
		s.logger.WithError(err).Error("Failed to edit post")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"post_id": postID,
	}).Info("Post edited successfully")

	return post, nil
}

func (s *PostService) DeletePost(ctx context.Context, postID uuid.UUID) error {
	if err := s.postUsecase.Delete(ctx, postID); err != nil {
		s.logger.WithError(err).Error("Failed to delete post")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"post_id": postID,
	}).Info("Post deleted successfully")

	return nil
}

func (s *PostService) GetPost(ctx context.Context, postID uuid.UUID) (*entity.Post, error) {
	post, err := s.postUsecase.GetPost(ctx, postID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get post")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"post_id": postID,
	}).Info("Post fetched successfully")

	return post, nil
}

func (s *PostService) ListPosts(ctx context.Context, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error) {
	if page < 1 || pageSize < 1 {
		return nil, 0, fmt.Errorf("invalid pagination parameters")
	}

	posts, total, err := s.postUsecase.ListPosts(ctx, page, pageSize, sortBy, filter)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list posts")
		return nil, 0, err
	}

	s.logger.WithFields(logrus.Fields{
		"page":        page,
		"page_size":   pageSize,
		"total_posts": total,
	}).Info("Posts listed successfully")

	return posts, total, nil
}

func (s *PostService) ListPostsByAuthor(ctx context.Context, authorID uuid.UUID, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error) {
	if page < 1 || pageSize < 1 {
		return nil, 0, fmt.Errorf("invalid pagination parameters")
	}

	posts, total, err := s.postUsecase.ListPostsByAuthor(ctx, authorID, page, pageSize, sortBy, filter)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list posts by author")
		return nil, 0, err
	}

	s.logger.WithFields(logrus.Fields{
		"author_id":   authorID,
		"page":        page,
		"page_size":   pageSize,
		"total_posts": total,
	}).Info("Posts by author listed successfully")

	return posts, total, nil
}
