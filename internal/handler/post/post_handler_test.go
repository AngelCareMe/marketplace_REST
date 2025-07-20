package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"marketplace/internal/entity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostService struct {
	mock.Mock
}

func (m *MockPostService) CreatePost(ctx context.Context, authorID uuid.UUID, header, content, image string, price float64) (*entity.Post, error) {
	args := m.Called(ctx, authorID, header, content, image, price)
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostService) GetPost(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostService) EditPost(ctx context.Context, id uuid.UUID, header, content, image string, price float64) (*entity.Post, error) {
	args := m.Called(ctx, id, header, content, image, price)
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *MockPostService) DeletePost(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostService) ListPosts(ctx context.Context, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error) {
	args := m.Called(ctx, page, pageSize, sortBy, filter)
	return args.Get(0).([]*entity.Post), args.Int(1), args.Error(2)
}

func (m *MockPostService) ListPostsByAuthor(ctx context.Context, authorID uuid.UUID, page, pageSize int, sortBy string, filter map[string]string) ([]*entity.Post, int, error) {
	args := m.Called(ctx, authorID, page, pageSize, sortBy, filter)
	return args.Get(0).([]*entity.Post), args.Int(1), args.Error(2)
}

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockPostSvc := new(MockPostService)
	logger := logrus.New()
	handler := NewPostHandler(mockPostSvc, nil, logger)

	r.POST("/posts", handler.CreatePost)

	reqBody := map[string]interface{}{
		"header":  "Test Post",
		"content": "This is a test post.",
		"image":   "http://example.com/image.jpg",
		"price":   99.99,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Setup mock user ID in context
	userID := uuid.New()
	ctx := context.WithValue(req.Context(), "user_id", userID)
	req = req.WithContext(ctx)

	// Setup expected response from service
	expectedPost := &entity.Post{
		ID:             uuid.New(),
		Header:         "Test Post",
		Content:        "This is a test post.",
		Image:          "http://example.com/image.jpg",
		Price:          99.99,
		AuthorID:       userID,
		CreatedAt:      time.Now(),
		AuthorUsername: "",
		IsOwnPost:      true,
	}
	mockPostSvc.On("CreatePost", ctx, userID, "Test Post", "This is a test post.", "http://example.com/image.jpg", 99.99).
		Return(expectedPost, nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
