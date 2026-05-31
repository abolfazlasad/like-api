package repositories

import "like-api/internal/domain/entities"

// VideoRepository defines the data-access contract for videos.
type VideoRepository interface {
	Create(video entities.Video) error
	FindByID(id string) (entities.Video, error)
	FindAll(cursor string, limit int) ([]entities.Video, string, error)
	FindByUserID(userID string) ([]entities.Video, error)
	Exists(id string) bool

	// IncrementLikes and DecrementLikes are atomic operations.
	IncrementLikes(videoID string) error
	DecrementLikes(videoID string) error

	// IncrementViews is an atomic operation used by the view-tracking service.
	IncrementViews(videoID string) error
}
