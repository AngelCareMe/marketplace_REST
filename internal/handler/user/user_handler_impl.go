package handler

import (
	service "marketplace/internal/service/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userSvc service.UserServiceInterface
	logger  *logrus.Logger
}

func NewUserHandler(userSvc service.UserServiceInterface, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
		logger:  logger,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=8,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid register request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, token, err := h.userSvc.Register(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.WithError(err).Error("Failed to register user")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("User registered via handler")
	c.JSON(http.StatusCreated, gin.H{"user": user, "token": token})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=8,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid login request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, token, err := h.userSvc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.WithError(err).Error("Failed to login user")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("User logged in via handler")
	c.JSON(http.StatusOK, gin.H{"user": user, "token": token})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userSvc.GetUser(c.Request.Context(), id)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get user")
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User fetched via handler")
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"omitempty,min=3,max=50"`
		Password string `json:"password" binding:"omitempty,min=8,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Invalid update user request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userSvc.UpdateUser(c.Request.Context(), id, req.Username, req.Password); err != nil {
		h.logger.WithError(err).Error("Failed to update user")
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User updated via handler")
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.WithError(err).Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userSvc.DeleteUser(c.Request.Context(), id); err != nil {
		h.logger.WithError(err).Error("Failed to delete user")
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User deleted via handler")
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
