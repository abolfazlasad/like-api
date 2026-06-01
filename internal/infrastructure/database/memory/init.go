package memory

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

	return userRepo, likeRepo, videoRepo, productRepo
}

func InitData(
	userRepo repositories.UserRepository,
	likeRepo repositories.LikeRepository,
	videoRepo repositories.VideoRepository,
	productRepo repositories.ProductRepository,
) {
	hash, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	users := []entities.User{
		{
			ID:        "user-admin",
			Username:  "admin",
			Name:      "Admin",
			Email:     "admin@example.com",
			Password:  string(hash),
			Role:      entities.RoleAdmin,
			CreatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ID:        "user1",
			Username:  "john_doe",
			Name:      "John Doe",
			Email:     "john@example.com",
			Password:  string(hash),
			Role:      entities.RoleUser,
			CreatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ID:        "user2",
			Username:  "jane_smith",
			Name:      "Jane Smith",
			Email:     "jane@example.com",
			Password:  string(hash),
			Role:      entities.RoleUser,
			CreatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ID:        "user3",
			Username:  "bob_wilson",
			Name:      "Bob Wilson",
			Email:     "bob@example.com",
			Password:  string(hash),
			Role:      entities.RoleUser,
			CreatedAt: time.Now().Format(time.RFC3339),
		},
	}
	for _, u := range users {
		userRepo.Create(u)
	}

	now := time.Now()
	videos := []entities.Video{
		{
			ID: "video1", UserID: "user1",
			Title:       "Go Clean Architecture",
			Description: "Building scalable Go services",
			VideoURL:    "https://example.com/videos/go-clean-arch.mp4",
			LikesCount:  0, ViewsCount: 0,
			CreatedAt: now.Add(-2 * time.Hour).Format(time.RFC3339),
		},
		{
			ID: "video2", UserID: "user2",
			Title:       "Cursor Pagination Deep Dive",
			Description: "Why offset pagination does not scale",
			VideoURL:    "https://example.com/videos/cursor-pagination.mp4",
			LikesCount:  0, ViewsCount: 0,
			CreatedAt: now.Add(-1 * time.Hour).Format(time.RFC3339),
		},
		{
			ID: "video3", UserID: "user3",
			Title:       "Redis for View Tracking",
			Description: "Deduplicating views at scale",
			VideoURL:    "https://example.com/videos/redis-views.mp4",
			LikesCount:  0, ViewsCount: 0,
			CreatedAt: now.Format(time.RFC3339),
		},
	}
	for _, v := range videos {
		videoRepo.Create(v)
	}

	products := []entities.Product{
		{ID: "product1", VideoID: "video1", Name: "Go in Action Book", Price: 39.99, ImageURL: "https://example.com/images/go-book.jpg"},
		{ID: "product2", VideoID: "video2", Name: "Backend Course", Price: 99.00, ImageURL: "https://example.com/images/course.jpg"},
	}
	for _, p := range products {
		productRepo.Create(p)
	}

	preLikes := []entities.Like{
		{ID: uuid.New().String(), UserID: "user2", PostID: "video1", CreatedAt: time.Now().Format(time.RFC3339)},
		{ID: uuid.New().String(), UserID: "user3", PostID: "video1", CreatedAt: time.Now().Format(time.RFC3339)},
	}
	for _, like := range preLikes {
		likeRepo.Create(like)
		videoRepo.IncrementLikes(like.PostID)
	}
}
