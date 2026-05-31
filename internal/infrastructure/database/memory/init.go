package memory

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func InitRepositories(
	initDefaultData bool,
) (
	repositories.UserRepository,
	repositories.PostRepository,
	repositories.LikeRepository,
) {
	userRepo := newUserRepository()
	postRepo := newPostRepository()
	likeRepo := newLikeRepository()

	if initDefaultData {
		initData(userRepo, postRepo, likeRepo)
	}

	return userRepo, postRepo, likeRepo
}

func initData(
	userRepo repositories.UserRepository,
	postRepo repositories.PostRepository,
	likeRepo repositories.LikeRepository,
) {
	// Hardcoded users
	hash, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	users := []entities.User{
		{ID: "user1", Username: "john_doe", Name: "John Doe", Email: "john@example.com", Password: string(hash), CreatedAt: time.Now().Format(time.RFC3339)},
		{ID: "user2", Username: "jane_smith", Name: "Jane Smith", Email: "jane@example.com", Password: string(hash), CreatedAt: time.Now().Format(time.RFC3339)},
		{ID: "user3", Username: "bob_wilson", Name: "Bob Wilson", Email: "bob@example.com", Password: string(hash), CreatedAt: time.Now().Format(time.RFC3339)},
	}

	for _, user := range users {
		userRepo.Create(user)
	}

	// Hardcoded posts for each user
	posts := []entities.Post{
		// John's posts
		{ID: "post1", UserID: "user1", Title: "First Post!", Content: "Hello everyone! This is my first post.", Likes: 0},
		{ID: "post2", UserID: "user1", Title: "Go Programming", Content: "Learning Go is amazing!", Likes: 0},

		// Jane's posts
		{ID: "post3", UserID: "user2", Title: "Travel Diary", Content: "Just visited Paris! 🗼", Likes: 0},
		{ID: "post4", UserID: "user2", Title: "Cooking Tips", Content: "Here's my favorite pasta recipe", Likes: 0},
		{ID: "post5", UserID: "user2", Title: "Fitness Journey", Content: "Day 30 of my workout challenge", Likes: 0},

		// Bob's posts
		{ID: "post6", UserID: "user3", Title: "Tech News", Content: "New iPhone just dropped!", Likes: 0},
		{ID: "post7", UserID: "user3", Title: "Gaming", Content: "Just finished Elden Ring", Likes: 0},
	}

	for _, post := range posts {
		postRepo.Create(post)
	}

	// Add some pre-existing likes for demo
	preLikes := []entities.Like{
		{ID: uuid.New().String(), UserID: "user1", PostID: "post1", CreatedAt: "2024-01-16T09:00:00Z"},
		{ID: uuid.New().String(), UserID: "user2", PostID: "post2", CreatedAt: "2024-01-15T10:00:00Z"},
		{ID: uuid.New().String(), UserID: "user3", PostID: "post3", CreatedAt: "2024-01-15T11:00:00Z"},
	}

	for _, like := range preLikes {
		likeRepo.Create(like)
		// Update post like count
		if post, err := postRepo.FindByID(like.PostID); err == nil {
			post.Likes++
			postRepo.Update(post)
		}
	}
}
