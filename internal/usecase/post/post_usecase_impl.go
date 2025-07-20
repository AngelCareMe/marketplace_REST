package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"marketplace/internal/entity"
	usecaseAuth "marketplace/internal/usecase/auth"
	usecase "marketplace/internal/usecase/user"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PostUsecase struct {
	postRepo PostRepository
	userRepo usecase.UserRepository
	authRepo usecaseAuth.AuthService
	logger   *logrus.Logger
}

func NewPostUsecase(postRepo PostRepository, userRepo usecase.UserRepository, authRepo usecaseAuth.AuthService, logger *logrus.Logger) *PostUsecase {
	return &PostUsecase{
		postRepo: postRepo,
		userRepo: userRepo,
		authRepo: authRepo,
		logger:   logger,
	}
}

func (uc *PostUsecase) Publish(ctx context.Context, authorID uuid.UUID, header, content, image string, price float64) (*entity.Post, error) {
	_, err := uc.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	post := &entity.Post{
		ID:        uuid.New(),
		Header:    header,
		Content:   content,
		Image:     image,
		Price:     price,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}

	if err := post.Validate(); err != nil {
		return nil, fmt.Errorf("validate post: %w", err)
	}

	exist, err := uc.postRepo.GetByHeaderAndContent(ctx, header, content)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("check duplicate: %w", err)
	}
	if exist != nil {
		return nil, fmt.Errorf("post with the same header and content already exists")
	}

	if err := uc.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"header":    header,
		"price":     price,
		"post_id":   post.ID,
		"author_id": authorID,
	}).Info("Post created")

	return post, nil
}

func (uc *PostUsecase) Edit(ctx context.Context, postID uuid.UUID, header, content, image string, price float64) (*entity.Post, error) {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	post, err := uc.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	if post.AuthorID != userID {
		return nil, errors.New("forbidden: not the author of the post")
	}

	if header != "" {
		post.Header = header
	}
	if content != "" {
		post.Content = content
	}
	if image != "" {
		post.Image = image
	}
	if price > 0 {
		post.Price = price
	}

	if header != "" && content != "" {
		existingPost, err := uc.postRepo.GetByHeaderAndContent(ctx, header, content)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("check duplicate: %w", err)
		}
		if existingPost != nil && existingPost.ID != postID {
			return nil, fmt.Errorf("post with the same header and content already exists")
		}
	}

	if err := post.Validate(); err != nil {
		return nil, fmt.Errorf("validate post: %w", err)
	}

	if err := uc.postRepo.Update(ctx, post); err != nil {
		return nil, fmt.Errorf("update post: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"post_id":   postID,
		"author_id": userID,
	}).Info("Post updated")

	return post, nil
}

func (uc *PostUsecase) Delete(ctx context.Context, postID uuid.UUID) error {
	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if !ok {
		return errors.New("unauthorized")
	}

	post, err := uc.postRepo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("get post by id: %w", err)
	}

	if post.AuthorID != userID {
		return errors.New("forbidden: not the author of the post")
	}

	if err := uc.postRepo.Delete(ctx, postID); err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"post_id": postID,
	}).Info("Post deleted")

	return nil
}

func (uc *PostUsecase) GetPost(ctx context.Context, postID uuid.UUID) (*entity.Post, error) {
	post, err := uc.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	user, err := uc.userRepo.GetByID(ctx, post.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	post.AuthorUsername = user.Username

	userID, ok := ctx.Value("user_id").(uuid.UUID)
	if ok {
		post.IsOwnPost = post.AuthorID == userID
	}

	uc.logger.WithFields(logrus.Fields{
		"post_id": postID,
	}).Info("Post fetched")

	return post, nil
}

func (uc *PostUsecase) ListPostsByAuthor(ctx context.Context, authorID uuid.UUID, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error) {
	_, err := uc.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, 0, fmt.Errorf("user not found: %w", err)
	}

	if sortBy == "" {
		sortBy = "created_at DESC"
	}

	posts, total, err := uc.postRepo.ListByAuthorID(ctx, authorID, page, pageSize, sortBy, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("get posts: %w", err)
	}

	userID, ok := ctx.Value("user_id").(uuid.UUID)
	for _, post := range posts {
		user, err := uc.userRepo.GetByID(ctx, post.AuthorID)
		if err != nil {
			return nil, 0, fmt.Errorf("get user: %w", err)
		}
		post.AuthorUsername = user.Username
		if ok {
			post.IsOwnPost = post.AuthorID == userID
		}
	}

	uc.logger.WithFields(logrus.Fields{
		"author_id":   authorID,
		"page":        page,
		"page_size":   pageSize,
		"sort_by":     sortBy,
		"filter":      filter,
		"total_posts": total,
	}).Info("Posts by author listed")

	return posts, total, nil
}

func (uc *PostUsecase) ListPosts(ctx context.Context, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error) {
	if sortBy == "" {
		sortBy = "created_at DESC"
	}

	posts, total, err := uc.postRepo.ListPosts(ctx, page, pageSize, sortBy, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("get posts: %w", err)
	}

	userID, ok := ctx.Value("user_id").(uuid.UUID)
	for _, post := range posts {
		user, err := uc.userRepo.GetByID(ctx, post.AuthorID)
		if err != nil {
			return nil, 0, fmt.Errorf("get user: %w", err)
		}
		post.AuthorUsername = user.Username
		if ok {
			post.IsOwnPost = post.AuthorID == userID
		}
	}

	uc.logger.WithFields(logrus.Fields{
		"page":        page,
		"page_size":   pageSize,
		"sort_by":     sortBy,
		"filter":      filter,
		"total_posts": total,
	}).Info("Posts listed")

	return posts, total, nil
}
