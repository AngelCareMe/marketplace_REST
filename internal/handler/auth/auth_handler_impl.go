package handler

import (
	"context"
	service "marketplace/internal/service/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authSvc service.AuthServiceInterface
	logger  *logrus.Logger
}

func NewAuthHandler(authSvc service.AuthServiceInterface, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
		logger:  logger,
	}
}

func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			h.logger.Warn("Authorization header is missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			h.logger.Warn("Invalid Authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		userID, err := h.authSvc.ValidateJWT(parts[1])
		if err != nil {
			h.logger.WithError(err).Error("Failed to validate JWT")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "user_id", userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (h *AuthHandler) OwnerMiddleware(paramID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Request.Context().Value("user_id").(uuid.UUID)
		if !ok {
			h.logger.Error("Failed to get user_id from context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		idStr := c.Param(paramID)
		resourceID, err := uuid.Parse(idStr)
		if err != nil {
			h.logger.WithError(err).Error("Invalid resource ID")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
			return
		}

		if userID != resourceID {
			h.logger.WithFields(logrus.Fields{
				"user_id":     userID,
				"resource_id": resourceID,
			}).Warn("User is not the owner of the resource")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		c.Next()
	}
}
