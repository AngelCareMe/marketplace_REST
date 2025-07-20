package handler

import "github.com/gin-gonic/gin"

type AuthHandlerInterface interface {
	AuthMiddleware() gin.HandlerFunc
	OwnerMiddleware(paramID string) gin.HandlerFunc
}
