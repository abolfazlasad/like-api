package main

import (
	"os"

	authusecase "like-api/internal/application/usecases/auth"
	likeusecase "like-api/internal/application/usecases/like"
	postusecase "like-api/internal/application/usecases/post"
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

	// TODO
	// userRepo, postRepo, likeRepo := memory.InitRepos(true)
	userRepo, postRepo, likeRepo, videoRepo, productRepo := postgres.InitRepositories()

	// ── auth ──────────────────────────────────────────────────────────────────
	registerUC := authusecase.NewRegisterUseCase(userRepo, jwtSecret)
	loginUC := authusecase.NewLoginUseCase(userRepo, jwtSecret)

	// ── legacy post / like ────────────────────────────────────────────────────
	getUsersUseCase := userusecase.NewGetUsersUseCase(userRepo)
	getPostsUseCase := postusecase.NewGetPostsUseCase(postRepo)
	getUserPostsUseCase := postusecase.NewGetUserPostsUseCase(userRepo, postRepo)
	getPostLikesUseCase := postusecase.NewGetPostLikesUseCase(postRepo, likeRepo)
	likePostUseCase := likeusecase.NewLikePostUseCase(userRepo, postRepo, likeRepo)
	unlikePostUseCase := likeusecase.NewUnlikePostUseCase(postRepo, likeRepo)
	getUserLikedPostsUseCase := likeusecase.NewGetUserLikedPostsUseCase(userRepo, postRepo, likeRepo)

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
	postHandler := httpHandler.NewPostHandler(getPostsUseCase, getUserPostsUseCase, getPostLikesUseCase)
	likeHandler := httpHandler.NewLikeHandler(likePostUseCase, unlikePostUseCase, getUserLikedPostsUseCase)
	videoHandler := httpHandler.NewVideoHandler(
		createVideoUC, getVideoUC, getFeedUC,
		likeVideoUC, unlikeVideoUC, trackViewUC, getVideoStatsUC,
	)
	productHandler := httpHandler.NewProductHandler(createProductUC, getProductUC, getProductByVideoUC)

	// ── router ────────────────────────────────────────────────────────────────
	r := router.NewRouter(
		authHandler, userHandler, postHandler, likeHandler,
		videoHandler, productHandler,
		jwtSecret,
	)
	r.Setup().Run(":8080")
}
