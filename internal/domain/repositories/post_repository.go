package repositories

import "like-api/internal/domain/entities"

type PostRepository interface {
	Create(post entities.Post) error
	FindAll() ([]entities.Post, error)
	FindByID(id string) (entities.Post, error)
	FindByUserID(userID string) ([]entities.Post, error)
	Update(post entities.Post) error
	Exists(id string) bool
}
