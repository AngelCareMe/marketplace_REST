package handler

import (
	servicePost "marketplace/internal/service/post"
	serviceUser "marketplace/internal/service/user"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PostHandler struct {
	postSvc servicePost.PostServiceInterface
	userSvc serviceUser.UserServiceInterface
	logger  *logrus.Logger
}

func NewPostHandler(postSvc servicePost.PostServiceInterface, userSvc serviceUser.UserServiceInterface, logger *logrus.Logger) *PostHandler {
	return &PostHandler{
		postSvc: postSvc,
		userSvc: userSvc,
		logger:  logger,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req struct {
		Header  string  `json:"header" binding:"required,min=1,max=100"`
		Content string  `json:"content" binding:"required,min=1,max=1000"`
		Image   string  `json:"image" binding:"omitempty,url"`
		Price   float64 `json:"price" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid create post request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, ok := c.Request.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Error("Failed to get user_id from context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	post, err := h.postSvc.CreatePost(c.Request.Context(), userID, req.Header, req.Content, req.Image, req.Price)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create post")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"post_id":   post.ID,
		"author_id": userID,
	}).Info("Post created via handler")
	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).Error("Invalid post ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.postSvc.GetPost(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get post")
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	h.logger.WithFields(logrus.Fields{
		"post_id": id,
	}).Info("Post fetched via handler")
	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) EditPost(c *gin.Context) {
	var req struct {
		Header  string  `json:"header" binding:"omitempty,min=1,max=100"`
		Content string  `json:"content" binding:"omitempty,min=1,max=1000"`
		Image   string  `json:"image" binding:"omitempty,url"`
		Price   float64 `json:"price" binding:"omitempty,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid edit post request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).Error("Invalid post ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userID, ok := c.Request.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Error("Failed to get user_id from context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	post, err := h.postSvc.GetPost(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get post for edit")
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}
	if post.AuthorID != userID {
		h.logger.WithFields(logrus.Fields{
			"user_id": userID,
			"post_id": id,
		}).Warn("User is not the author of the post")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	updatedPost, err := h.postSvc.EditPost(c.Request.Context(), id, req.Header, req.Content, req.Image, req.Price)
	if err != nil {
		h.logger.WithError(err).Error("Failed to edit post")
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	h.logger.WithFields(logrus.Fields{
		"post_id": id,
	}).Info("Post edited via handler")
	c.JSON(http.StatusOK, updatedPost)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).Error("Invalid post ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userID, ok := c.Request.Context().Value("user_id").(uuid.UUID)
	if !ok {
		h.logger.Error("Failed to get user_id from context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	post, err := h.postSvc.GetPost(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get post for delete")
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}
	if post.AuthorID != userID {
		h.logger.WithFields(logrus.Fields{
			"user_id": userID,
			"post_id": id,
		}).Warn("User is not the author of the post")
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if err := h.postSvc.DeletePost(c.Request.Context(), id); err != nil {
		h.logger.WithError(err).Error("Failed to delete post")
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	h.logger.WithFields(logrus.Fields{
		"post_id": id,
	}).Info("Post deleted via handler")
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func (h *PostHandler) ListPosts(c *gin.Context) {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")
	sortBy := c.Query("sortBy")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	filter := make(map[string]string)
	if minPrice != "" {
		filter["min_price"] = minPrice
	}
	if maxPrice != "" {
		filter["max_price"] = maxPrice
	}

	posts, total, err := h.postSvc.ListPosts(c.Request.Context(), page, pageSize, sortBy, filter)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list posts")
		if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	userID, _ := c.Request.Context().Value("user_id").(uuid.UUID)
	for _, post := range posts {
		post.IsOwnPost = post.AuthorID == userID
		if post.IsOwnPost {
			user, err := h.userSvc.GetUser(c.Request.Context(), post.AuthorID)
			if err == nil {
				post.AuthorUsername = user.Username
			}
		}
	}

	h.logger.WithFields(logrus.Fields{
		"page":        page,
		"page_size":   pageSize,
		"total_posts": total,
	}).Info("Posts listed via handler")
	c.JSON(http.StatusOK, gin.H{
		"posts":     posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *PostHandler) ListPostsByAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	pageStr := c.Query("page")
	pageSizeStr := c.Query("pageSize")
	sortBy := c.Query("sortBy")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	filter := make(map[string]string)
	if minPrice != "" {
		filter["min_price"] = minPrice
	}
	if maxPrice != "" {
		filter["max_price"] = maxPrice
	}

	posts, total, err := h.postSvc.ListPostsByAuthor(c.Request.Context(), id, page, pageSize, sortBy, filter)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list posts by author")
		if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	userID, _ := c.Request.Context().Value("user_id").(uuid.UUID)
	for _, post := range posts {
		post.IsOwnPost = post.AuthorID == userID
		if post.IsOwnPost {
			user, err := h.userSvc.GetUser(c.Request.Context(), post.AuthorID)
			if err == nil {
				post.AuthorUsername = user.Username
			}
		}
	}

	h.logger.WithFields(logrus.Fields{
		"author_id":   id,
		"page":        page,
		"page_size":   pageSize,
		"total_posts": total,
	}).Info("Posts by author listed via handler")
	c.JSON(http.StatusOK, gin.H{
		"posts":     posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
