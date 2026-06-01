package memory

import (
	"like-api/internal/domain/repositories"
	"like-api/internal/infrastructure/database"
)

func InitRepositories() (
	repositories.UserRepository,
	repositories.LikeRepository,
	repositories.VideoRepository,
	repositories.ProductRepository,
) {
	userRepo := NewUserRepository()
	likeRepo := NewLikeRepository()
	videoRepo := NewVideoRepository()
	productRepo := NewProductRepository()

	database.SeedAdmin(userRepo)

	return userRepo, likeRepo, videoRepo, productRepo
}
