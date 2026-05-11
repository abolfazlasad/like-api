package repositories

import "like-api/internal/domain/entities"

type LikeRepository interface {
	Create(like entities.Like) error
	Delete(userID, postID string) error
	FindByPostID(postID string) ([]entities.Like, error)
	FindByUserID(userID string) ([]entities.Like, error)
	Exists(userID, postID string) bool
}
