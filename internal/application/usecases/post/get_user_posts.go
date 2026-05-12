package postusecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
)

// GetUserPostsInput represents input for getting posts by user
type GetUserPostsInput struct {
	UserID string
}

// GetUserPostsOutput represents output for getting posts by user
type GetUserPostsOutput struct {
	Posts      []entities.Post
	Error      error
	UserExists bool
}

// GetUserPostsUseCase defines the interface for getting posts by user
type GetUserPostsUseCase interface {
	Execute(input GetUserPostsInput) GetUserPostsOutput
}

// getUserPostsUseCaseImpl implements GetUserPostsUseCase
type getUserPostsUseCaseImpl struct {
	userRepo repositories.UserRepository
	postRepo repositories.PostRepository
}

// NewGetUserPostsUseCase creates a new get user posts use case
func NewGetUserPostsUseCase(userRepo repositories.UserRepository, postRepo repositories.PostRepository) GetUserPostsUseCase {
	return &getUserPostsUseCaseImpl{
		userRepo: userRepo,
		postRepo: postRepo,
	}
}

// Execute gets posts by user ID
func (uc *getUserPostsUseCaseImpl) Execute(input GetUserPostsInput) GetUserPostsOutput {
	// Check if user exists
	if !uc.userRepo.Exists(input.UserID) {
		return GetUserPostsOutput{
			Posts:      nil,
			Error:      nil,
			UserExists: false,
		}
	}

	// Get posts by user
	posts, err := uc.postRepo.FindByUserID(input.UserID)
	if err != nil {
		return GetUserPostsOutput{
			Posts:      nil,
			Error:      err,
			UserExists: true,
		}
	}

	return GetUserPostsOutput{
		Posts:      posts,
		Error:      nil,
		UserExists: true,
	}
}
