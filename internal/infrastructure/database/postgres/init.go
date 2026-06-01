package postgres

import (
	"like-api/internal/domain/repositories"
	"log"
	"os"
)

// InitRepositories reads env vars, connects to PostgreSQL, runs migrations,
// and returns all repository implementations.
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

	return newUserRepository(db),
		newLikeRepository(db),
		newVideoRepository(db),
		newProductRepository(db)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
