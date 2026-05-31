package repositories

import "like-api/internal/domain/entities"

type PostRepository interface {
	Create(post entities.Post) error
	FindAll() ([]entities.Post, error)
	FindByID(id string) (entities.Post, error)
	FindByUserID(userID string) ([]entities.Post, error)
	Update(post entities.Post) error
	Exists(id string) bool

	// IncrementLikes and DecrementLikes are atomic operations and are safe
	// from race conditions. They should be used instead of read-modify-write
	// patterns like fetching a post and doing likes++.
	IncrementLikes(postID string) error
	DecrementLikes(postID string) error
}
