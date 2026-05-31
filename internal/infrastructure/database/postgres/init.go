package postgres

import (
	"like-api/internal/domain/repositories"
	"log"
	"os"
)

// InitRepositories بر اساس env متغیرها، postgres یا memory انتخاب می‌کنه
func InitRepositories() (
	repositories.UserRepository,
	repositories.PostRepository,
	repositories.LikeRepository,
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
		newPostRepository(db),
		newLikeRepository(db)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
