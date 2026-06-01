// Package main is the entry point for the Video Commerce & Social Feed API.
//
//	@title			Video Commerce & Social Feed API
//	@version		1.0
//	@description	Backend for a video-based social commerce platform.
//	@description
//	@description	## Authentication
//	@description	This API uses **JWT Bearer** authentication. Obtain a token via
//	@description	`POST /api/v1/auth/register` or `POST /api/v1/auth/login`, then pass it
//	@description	in the `Authorization` header:
//	@description	```
//	@description	Authorization: Bearer <token>
//	@description	```
//	@description
//	@description	## Roles
//	@description	| Role  | Description |
//	@description	|-------|-------------|
//	@description	| user  | Default role assigned on registration |
//	@description	| admin | Required for admin-only endpoints (e.g. `GET /api/v1/admin/users`) |
//	@description
//	@description	## Optional Authentication
//	@description	Some endpoints (`GET /feed`, `GET /videos/:id`, `POST /videos/:id/view`,
//	@description	`GET /videos/:id/stats`) accept but do not require a JWT. Anonymous
//	@description	requests are served normally; authenticated requests may receive
//	@description	personalised responses in future versions.
//
//	@contact.name	API Support
//
//	@host		localhost:8080
//	@BasePath	/
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter your JWT as: **Bearer &lt;token&gt;**
package main

import (
	"log"
	"os"

	authusecase "like-api/internal/application/usecases/auth"
	productusecase "like-api/internal/application/usecases/product"
	userusecase "like-api/internal/application/usecases/user"
	videousecase "like-api/internal/application/usecases/video"
	"like-api/internal/domain/repositories"
	"like-api/internal/infrastructure/database/memory"
	"like-api/internal/infrastructure/database/postgres"
	"like-api/internal/infrastructure/router"
	httpHandler "like-api/internal/interfaces/http"

	_ "like-api/docs"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production"
	}

	var (
		userRepo    repositories.UserRepository
		likeRepo    repositories.LikeRepository
		videoRepo   repositories.VideoRepository
		productRepo repositories.ProductRepository
	)

	if os.Getenv("DB_TYPE") == "memory" {
		userRepo, likeRepo, videoRepo, productRepo = memory.InitRepositories()
	} else {
		userRepo, likeRepo, videoRepo, productRepo = postgres.InitRepositories()
	}

	// ── auth ──────────────────────────────────────────────────────────────────
	registerUC := authusecase.NewRegisterUseCase(userRepo, jwtSecret)
	loginUC := authusecase.NewLoginUseCase(userRepo, jwtSecret)

	// ── user ──────────────────────────────────────────────────────────────────
	getUsersUseCase := userusecase.NewGetUsersUseCase(userRepo)

	// ── video ─────────────────────────────────────────────────────────────────
	createVideoUC := videousecase.NewCreateVideoUseCase(userRepo, videoRepo)
	getVideoUC := videousecase.NewGetVideoUseCase(videoRepo)
	getFeedUC := videousecase.NewGetFeedUseCase(videoRepo)
	likeVideoUC := videousecase.NewLikeVideoUseCase(userRepo, videoRepo, likeRepo)
	unlikeVideoUC := videousecase.NewUnlikeVideoUseCase(videoRepo, likeRepo)
	trackViewUC := videousecase.NewTrackViewUseCase(videoRepo)
	getVideoStatsUC := videousecase.NewGetVideoStatsUseCase(videoRepo)

	// ── product ───────────────────────────────────────────────────────────────
	createProductUC := productusecase.NewCreateProductUseCase(videoRepo, productRepo)
	getProductUC := productusecase.NewGetProductUseCase(productRepo)
	getProductByVideoUC := productusecase.NewGetProductByVideoUseCase(videoRepo, productRepo)

	// ── handlers ──────────────────────────────────────────────────────────────
	authHandler := httpHandler.NewAuthHandler(registerUC, loginUC)
	userHandler := httpHandler.NewUserHandler(getUsersUseCase)
	videoHandler := httpHandler.NewVideoHandler(
		createVideoUC, getVideoUC, getFeedUC,
		likeVideoUC, unlikeVideoUC, trackViewUC, getVideoStatsUC,
	)
	productHandler := httpHandler.NewProductHandler(createProductUC, getProductUC, getProductByVideoUC)

	// ── router ────────────────────────────────────────────────────────────────
	r := router.NewRouter(
		authHandler, userHandler,
		videoHandler, productHandler,
		jwtSecret,
	)
	r.Setup().Run(":8080")
}
