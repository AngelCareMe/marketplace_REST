package handler

import (
	handlerAuth "marketplace/internal/handler/auth"
	handlerPost "marketplace/internal/handler/post"
	handlerUser "marketplace/internal/handler/user"

	"github.com/gin-gonic/gin"
)

type Router struct {
	userHandler handlerUser.UserHandlerInterface
	postHandler handlerPost.PostHandlerInterface
	authHandler handlerAuth.AuthHandlerInterface
}

func NewRouter(userHandler handlerUser.UserHandlerInterface, postHandler handlerPost.PostHandlerInterface, authHandler handlerAuth.AuthHandlerInterface) *Router {
	return &Router{
		userHandler: userHandler,
		postHandler: postHandler,
		authHandler: authHandler,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	ginRouter := gin.New()
	ginRouter.Use(gin.Logger(), gin.Recovery())

	ginRouter.POST("/users/register", r.userHandler.Register)
	ginRouter.POST("/users/login", r.userHandler.Login)
	ginRouter.GET("/posts/:id", r.postHandler.GetPost)
	ginRouter.GET("/posts", r.postHandler.ListPosts)

	private := ginRouter.Group("/", r.authHandler.AuthMiddleware())
	{
		private.GET("/users/:id", r.userHandler.GetUser)
		private.PUT("/users/:id", r.authHandler.OwnerMiddleware("id"), r.userHandler.UpdateUser)
		private.DELETE("/users/:id", r.authHandler.OwnerMiddleware("id"), r.userHandler.DeleteUser)
		private.POST("/posts", r.postHandler.CreatePost)
		private.PUT("/posts/:id", r.postHandler.EditPost)
		private.DELETE("/posts/:id", r.postHandler.DeletePost)
		private.GET("/users/:id/posts", r.postHandler.ListPostsByAuthor)
	}

	return ginRouter
}
