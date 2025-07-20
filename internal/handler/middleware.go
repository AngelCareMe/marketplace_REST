package handler

import "github.com/gin-gonic/gin"

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) LoggerMiddleware() gin.HandlerFunc {
	return gin.Logger()
}

func (m *Middleware) RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}
