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
	authHandler    *httpHandler.AuthHandler
	userHandler    *httpHandler.UserHandler
	videoHandler   *httpHandler.VideoHandler
	productHandler *httpHandler.ProductHandler
	jwtSecret      string
}

func NewRouter(
	authHandler *httpHandler.AuthHandler,
	userHandler *httpHandler.UserHandler,
	videoHandler *httpHandler.VideoHandler,
	productHandler *httpHandler.ProductHandler,
	jwtSecret string,
) *Router {
	return &Router{
		authHandler:    authHandler,
		userHandler:    userHandler,
		videoHandler:   videoHandler,
		productHandler: productHandler,
		jwtSecret:      jwtSecret,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")

	// ── public ──────────────────────────────────────────────────────────────

	api.POST("/auth/register", r.authHandler.Register)
	api.POST("/auth/login", r.authHandler.Login)

	// users
	api.GET("/users", r.userHandler.GetUsers)

	// videos — read
	api.GET("/feed", r.videoHandler.GetFeed)
	api.GET("/videos/:id", r.videoHandler.GetVideo)
	api.GET("/videos/:id/stats", r.videoHandler.GetVideoStats)
	api.GET("/videos/:id/product", r.productHandler.GetProductByVideo)

	// products — read
	api.GET("/products/:id", r.productHandler.GetProduct)

	// ── protected (JWT required) ─────────────────────────────────────────────

	auth := api.Group("")
	auth.Use(middleware.AuthMiddleware(r.jwtSecret))
	{
		// videos — write
		auth.POST("/videos", r.videoHandler.CreateVideo)
		auth.POST("/videos/:id/like", r.videoHandler.LikeVideo)
		auth.POST("/videos/:id/unlike", r.videoHandler.UnlikeVideo)
		auth.POST("/videos/:id/view", r.videoHandler.TrackView)

		// products — write
		auth.POST("/products", r.productHandler.CreateProduct)
	}

	return router
}
