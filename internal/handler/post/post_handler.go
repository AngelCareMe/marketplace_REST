package handler

import "github.com/gin-gonic/gin"

type PostHandlerInterface interface {
	CreatePost(c *gin.Context)
	GetPost(c *gin.Context)
	EditPost(c *gin.Context)
	DeletePost(c *gin.Context)
	ListPosts(c *gin.Context)
	ListPostsByAuthor(c *gin.Context)
}
