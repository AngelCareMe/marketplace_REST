package handler

import "github.com/gin-gonic/gin"

type UserHandlerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}
