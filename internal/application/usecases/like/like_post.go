package likeusecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
)

// LikePostInput represents input for liking a post
type LikePostInput struct {
	UserID string
	PostID string
}

// LikePostOutput represents output for liking a post
type LikePostOutput struct {
	Like         entities.Like
	Error        error
	UserExists   bool
	PostExists   bool
	AlreadyLiked bool
}

// LikePostUseCase defines the interface for liking a post
type LikePostUseCase interface {
	Execute(input LikePostInput) LikePostOutput
}

// LikePostUseCaseImpl implements LikePostUseCase
type likePostUseCaseImpl struct {
	userRepo repositories.UserRepository
	postRepo repositories.PostRepository
	likeRepo repositories.LikeRepository
}

// NewLikePostUseCase creates a new like post use case
func NewLikePostUseCase(
	userRepo repositories.UserRepository,
	postRepo repositories.PostRepository,
	likeRepo repositories.LikeRepository,
) LikePostUseCase {
	return &likePostUseCaseImpl{
		userRepo: userRepo,
		postRepo: postRepo,
		likeRepo: likeRepo,
	}
}

// Execute likes a post
func (uc *likePostUseCaseImpl) Execute(input LikePostInput) LikePostOutput {
	// Check if user exists
	if !uc.userRepo.Exists(input.UserID) {
		return LikePostOutput{
			UserExists: false,
		}
	}

	// Check if post exists
	if !uc.postRepo.Exists(input.PostID) {
		return LikePostOutput{
			UserExists: true,
			PostExists: false,
		}
	}

	// Check if already liked
	if uc.likeRepo.Exists(input.UserID, input.PostID) {
		return LikePostOutput{
			UserExists:   true,
			PostExists:   true,
			AlreadyLiked: true,
		}
	}

	// Create like
	like := entities.Like{
		ID:        uuid.New().String(),
		UserID:    input.UserID,
		PostID:    input.PostID,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := uc.likeRepo.Create(like); err != nil {
		return LikePostOutput{
			UserExists:   true,
			PostExists:   true,
			AlreadyLiked: false,
			Error:        err,
		}
	}

	// atomic increment — no race condition
	if err := uc.postRepo.IncrementLikes(input.PostID); err != nil {
		return LikePostOutput{
			UserExists:   true,
			PostExists:   true,
			AlreadyLiked: false,
			Error:        err,
		}
	}

	return LikePostOutput{
		Like:         like,
		UserExists:   true,
		PostExists:   true,
		AlreadyLiked: false,
		Error:        nil,
	}
}
