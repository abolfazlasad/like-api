package main

import (
	"os"

	authusecase "like-api/internal/application/usecases/auth"
	productusecase "like-api/internal/application/usecases/product"
	userusecase "like-api/internal/application/usecases/user"
	videousecase "like-api/internal/application/usecases/video"
	"like-api/internal/infrastructure/database/postgres"
	"like-api/internal/infrastructure/router"
	httpHandler "like-api/internal/interfaces/http"

	_ "like-api/docs"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production"
	}

	// userRepo, likeRepo, videoRepo, productRepo := memory.InitRepositories(true)
	userRepo, likeRepo, videoRepo, productRepo := postgres.InitRepositories()

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
