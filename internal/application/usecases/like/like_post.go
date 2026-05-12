package likeusecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"

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
	post, err := uc.postRepo.FindByID(input.PostID)
	if err != nil {
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
	likeID := uuid.New().String()
	like := entities.Like{
		ID:        likeID,
		UserID:    input.UserID,
		PostID:    input.PostID,
		CreatedAt: "2024-01-20T10:00:00Z", // In real app, use time.Now()
	}

	if err := uc.likeRepo.Create(like); err != nil {
		return LikePostOutput{
			UserExists:   true,
			PostExists:   true,
			AlreadyLiked: false,
			Error:        err,
		}
	}

	// Update post like count
	post.Likes++
	if err := uc.postRepo.Update(post); err != nil {
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
