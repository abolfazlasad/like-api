package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"like-api/internal/infrastructure/middleware"
	httpHandler "like-api/internal/interfaces/http"
)

// Router wires HTTP handlers to Gin routes.
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

// Setup registers all routes and returns the configured Gin engine.
//
// Route groups:
//
//	public          — no authentication required
//	optionalAuth    — JWT parsed when present; anonymous requests still served
//	protected       — valid JWT required (any role)
//	admin           — valid JWT required + role must be "admin"
func (r *Router) Setup() *gin.Engine {
	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins: true, // Replace for production
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization",
		},
	}))

	// Swagger UI
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := engine.Group("/api/v1")

	// ── public ───────────────────────────────────────────────────────────────
	public := api.Group("")
	{
		// Auth
		public.POST("/auth/register", r.authHandler.Register)
		public.POST("/auth/login", r.authHandler.Login)

		// video
		public.GET("/videos/:id", r.videoHandler.GetVideo)
		public.GET("/videos/:id/stats", r.videoHandler.GetVideoStats)

		// Products
		public.GET("/products/:id", r.productHandler.GetProduct)
		public.GET("/videos/:id/product", r.productHandler.GetProductByVideo)
	}

	// ── optional auth (JWT decoded when present, never rejected) ─────────────
	optionalAuth := api.Group("")
	optionalAuth.Use(middleware.OptionalAuthMiddleware(r.jwtSecret))
	{
		// Feed & video reads — available to anonymous users too;
		// when authenticated the userID is available for personalisation.
		optionalAuth.GET("/feed", r.videoHandler.GetFeed)

		// View tracking — authenticated view counts are deduplicated per user;
		// anonymous views are deduplicated by IP (handler falls back to IP).
		optionalAuth.POST("/videos/:id/view", r.videoHandler.TrackView)
	}

	// ── protected — valid JWT required (role: user or admin) ─────────────────
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(r.jwtSecret))
	{
		// Videos — write
		protected.POST("/videos", r.videoHandler.CreateVideo)
		protected.POST("/videos/:id/like", r.videoHandler.LikeVideo)
		protected.POST("/videos/:id/unlike", r.videoHandler.UnlikeVideo)

		// Products — write
		protected.POST("/products", r.productHandler.CreateProduct)
	}

	// ── admin — valid JWT required with role: admin ───────────────────────────
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(r.jwtSecret))
	admin.Use(middleware.AdminMiddleware())
	{
		admin.GET("/users", r.userHandler.GetUsers)
	}

	return engine
}
