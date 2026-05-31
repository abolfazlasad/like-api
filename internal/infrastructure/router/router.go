package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"like-api/internal/infrastructure/middleware"
	httpHandler "like-api/internal/interfaces/http"
)

type Router struct {
	authHandler *httpHandler.AuthHandler
	userHandler *httpHandler.UserHandler
	postHandler *httpHandler.PostHandler
	likeHandler *httpHandler.LikeHandler
	jwtSecret   string
}

func NewRouter(
	authHandler *httpHandler.AuthHandler,
	userHandler *httpHandler.UserHandler,
	postHandler *httpHandler.PostHandler,
	likeHandler *httpHandler.LikeHandler,
	jwtSecret string,
) *Router {
	return &Router{
		authHandler: authHandler,
		userHandler: userHandler,
		postHandler: postHandler,
		likeHandler: likeHandler,
		jwtSecret:   jwtSecret,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")

	// public
	api.POST("/auth/register", r.authHandler.Register)
	api.POST("/auth/login", r.authHandler.Login)

	api.GET("/users", r.userHandler.GetUsers)
	api.GET("/users/:user_id/posts", r.postHandler.GetUserPosts)
	api.GET("/users/:user_id/liked-posts", r.likeHandler.GetUserLikedPosts)
	api.GET("/posts", r.postHandler.GetPosts)
	api.GET("/posts/:post_id/likes", r.postHandler.GetPostLikes)

	// protected
	likes := api.Group("/likes")
	likes.Use(middleware.AuthMiddleware(r.jwtSecret))
	{
		likes.POST("", r.likeHandler.LikePost)
		likes.DELETE("", r.likeHandler.UnlikePost)
	}

	return router
}
