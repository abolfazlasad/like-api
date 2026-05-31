package main

import (
	"os"

	authusecase "like-api/internal/application/usecases/auth"
	likeusecase "like-api/internal/application/usecases/like"
	postusecase "like-api/internal/application/usecases/post"
	userusecase "like-api/internal/application/usecases/user"
	"like-api/internal/infrastructure/database/memory"
	"like-api/internal/infrastructure/router"
	httpHandler "like-api/internal/interfaces/http"

	_ "like-api/docs"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-change-in-production"
	}

	// repositories
	userRepo := memory.NewUserRepository()
	postRepo := memory.NewPostRepository()
	likeRepo := memory.NewLikeRepository()

	// Initialize with hardcoded data
	memory.InitData(userRepo, postRepo, likeRepo)

	// auth use cases
	registerUC := authusecase.NewRegisterUseCase(userRepo, jwtSecret)
	loginUC := authusecase.NewLoginUseCase(userRepo, jwtSecret)

	// existing use cases
	getUsersUseCase := userusecase.NewGetUsersUseCase(userRepo)
	getPostsUseCase := postusecase.NewGetPostsUseCase(postRepo)
	getUserPostsUseCase := postusecase.NewGetUserPostsUseCase(userRepo, postRepo)
	getPostLikesUseCase := postusecase.NewGetPostLikesUseCase(postRepo, likeRepo)
	likePostUseCase := likeusecase.NewLikePostUseCase(userRepo, postRepo, likeRepo)
	unlikePostUseCase := likeusecase.NewUnlikePostUseCase(postRepo, likeRepo)
	getUserLikedPostsUseCase := likeusecase.NewGetUserLikedPostsUseCase(userRepo, postRepo, likeRepo)

	// handlers
	authHandler := httpHandler.NewAuthHandler(registerUC, loginUC)
	userHandler := httpHandler.NewUserHandler(getUsersUseCase)
	postHandler := httpHandler.NewPostHandler(getPostsUseCase, getUserPostsUseCase, getPostLikesUseCase)
	likeHandler := httpHandler.NewLikeHandler(likePostUseCase, unlikePostUseCase, getUserLikedPostsUseCase)

	// router
	r := router.NewRouter(authHandler, userHandler, postHandler, likeHandler, jwtSecret)
	r.Setup().Run(":8080")
}
