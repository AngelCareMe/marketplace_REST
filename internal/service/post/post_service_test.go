package service

import (
	"context"
	"testing"
	"time"

	"marketplace/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostUseCase struct {
	mock.Mock
}

func (m *MockPostUseCase) Publish(ctx context.Context, authorID uuid.UUID, header, content, image string, price float64) (*entity.Post, error) {
	args := m.Called(ctx, authorID, header, content, image, price)
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostUseCase) Edit(ctx context.Context, postID uuid.UUID, header, content, image string, price float64) (*entity.Post, error) {
	args := m.Called(ctx, postID, header, content, image, price)
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostUseCase) Delete(ctx context.Context, postID uuid.UUID) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockPostUseCase) GetPost(ctx context.Context, postID uuid.UUID) (*entity.Post, error) {
	args := m.Called(ctx, postID)
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostUseCase) ListPosts(ctx context.Context, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error) {
	args := m.Called(ctx, page, pageSize, sortBy, filter)
	return args.Get(0).([]*entity.Post), args.Int(1), args.Error(2)
}

func (m *MockPostUseCase) ListPostsByAuthor(ctx context.Context, authorID uuid.UUID, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error) {
	args := m.Called(ctx, authorID, page, pageSize, sortBy, filter)
	return args.Get(0).([]*entity.Post), args.Int(1), args.Error(2)
}
func TestEditPost(t *testing.T) {
	mockUsecase := new(MockPostUseCase)
	logger := logrus.New()
	postService := NewPostService(mockUsecase, logger)

	postID := uuid.New()
	header := "Updated Header"
	content := "Updated Content"
	image := "http://example.com/new-image.jpg"
	price := 89.99

	expectedPost := &entity.Post{
		ID:        postID,
		Header:    header,
		Content:   content,
		Image:     image,
		Price:     price,
		AuthorID:  uuid.New(),
		CreatedAt: time.Now(),
	}

	mockUsecase.On("Edit", mock.Anything, postID, header, content, image, price).
		Return(expectedPost, nil)

	result, err := postService.EditPost(context.Background(), postID, header, content, image, price)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, result)
	mockUsecase.AssertExpectations(t)
}
func TestDeletePost(t *testing.T) {
	mockUsecase := new(MockPostUseCase)
	logger := logrus.New()
	postService := NewPostService(mockUsecase, logger)

	postID := uuid.New()

	mockUsecase.On("Delete", mock.Anything, postID).Return(nil)

	err := postService.DeletePost(context.Background(), postID)
	assert.NoError(t, err)
	mockUsecase.AssertExpectations(t)
}

func TestCreatePost(t *testing.T) {
	mockUsecase := new(MockPostUseCase)
	logger := logrus.New()
	postService := NewPostService(mockUsecase, logger)

	authorID := uuid.New()
	header := "Test Post"
	content := "This is a test post."
	image := "http://example.com/image.jpg"
	price := 99.99

	expectedPost := &entity.Post{
		ID:        uuid.New(),
		Header:    header,
		Content:   content,
		Image:     image,
		Price:     price,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
	}

	mockUsecase.On("Publish", mock.Anything, authorID, header, content, image, price).
		Return(expectedPost, nil)

	result, err := postService.CreatePost(context.Background(), authorID, header, content, image, price)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, result)
	mockUsecase.AssertExpectations(t)
}
