package postgres

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

// InitRepositories reads env vars, connects to PostgreSQL, runs migrations,
// seeds the admin user if configured, and returns all repository implementations.
func InitRepositories() (
	repositories.UserRepository,
	repositories.LikeRepository,
	repositories.VideoRepository,
	repositories.ProductRepository,
) {
	dbHost := os.Getenv("DB_HOST")

	if dbHost == "" {
		panic("[db] DB_HOST not set")
	}

	log.Println("[db] connecting to postgres...")

	cfg := Config{
		Host:     dbHost,
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "appuser"),
		Password: getEnv("DB_PASSWORD", "secret"),
		DBName:   getEnv("DB_NAME", "video_commerce"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := NewDB(cfg)
	if err != nil {
		log.Fatalf("[db] failed to connect: %v", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatalf("[db] migration failed: %v", err)
	}

	userRepo := newUserRepository(db)

	seedAdmin(userRepo)

	return userRepo,
		newLikeRepository(db),
		newVideoRepository(db),
		newProductRepository(db)
}

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
func seedAdmin(userRepo repositories.UserRepository) {
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

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
