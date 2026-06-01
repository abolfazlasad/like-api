package database

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// seedAdmin creates the admin user from environment variables if it does not
// already exist. All four env vars must be present; if any is missing the seed
// step is silently skipped so existing deployments are unaffected.
//
// Required env vars:
//
//	ADMIN_EMAIL
//	ADMIN_PASSWORD   (plain-text; will be bcrypt-hashed before storage)
//	ADMIN_NAME
//	ADMIN_USERNAME
func SeedAdmin(userRepo repositories.UserRepository) {
	email := strings.TrimSpace(os.Getenv("ADMIN_EMAIL"))
	password := os.Getenv("ADMIN_PASSWORD")
	name := strings.TrimSpace(os.Getenv("ADMIN_NAME"))
	username := strings.TrimSpace(os.Getenv("ADMIN_USERNAME"))

	// All four must be set; partial config is ignored.
	if email == "" || password == "" || name == "" || username == "" {
		log.Println("[db] admin seed skipped: ADMIN_EMAIL / ADMIN_PASSWORD / ADMIN_NAME / ADMIN_USERNAME not fully set")
		return
	}

	if userRepo.ExistsByEmail(email) {
		log.Printf("[db] admin user %q already exists — skipping seed", email)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("[db] failed to hash admin password: %v", err)
	}

	admin := entities.User{
		ID:        uuid.New().String(),
		Username:  username,
		Name:      name,
		Email:     email,
		Password:  string(hashed),
		Role:      entities.RoleAdmin,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := userRepo.Create(admin); err != nil {
		log.Fatalf("[db] failed to create admin user: %v", err)
	}

	log.Printf("[db] admin user %q created successfully", email)
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
	hashAdmin, errAdmin := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if errAdmin != nil {
		panic(errAdmin)
	}

	users := []entities.User{
		{
			ID:        "user-admin",
			Username:  "admin",
			Name:      "Admin",
			Email:     "admin@example.com",
			Password:  string(hashAdmin),
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
