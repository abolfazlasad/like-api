package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	httpHandler "like-api/internal/interfaces/http"
)

type Router struct {
	userHandler *httpHandler.UserHandler
	postHandler *httpHandler.PostHandler
	likeHandler *httpHandler.LikeHandler
}

func NewRouter(
	userHandler *httpHandler.UserHandler,
	postHandler *httpHandler.PostHandler,
	likeHandler *httpHandler.LikeHandler,
) *Router {
	return &Router{
		userHandler: userHandler,
		postHandler: postHandler,
		likeHandler: likeHandler,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	{
		// User routes
		api.GET("/users", r.userHandler.GetUsers)
		api.GET("/users/:user_id/posts", r.postHandler.GetUserPosts)
		api.GET("/users/:user_id/liked-posts", r.likeHandler.GetUserLikedPosts)

		// Post routes
		api.GET("/posts", r.postHandler.GetPosts)
		api.GET("/posts/:post_id/likes", r.postHandler.GetPostLikes)

		// Like routes
		api.POST("/likes", r.likeHandler.LikePost)
		api.DELETE("/likes", r.likeHandler.UnlikePost)
	}

	return router
}
