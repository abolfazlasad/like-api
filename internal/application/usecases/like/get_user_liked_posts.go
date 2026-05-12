package likeusecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
)

// GetUserLikedPostsInput represents input for getting user's liked posts
type GetUserLikedPostsInput struct {
	UserID string
}

// GetUserLikedPostsOutput represents output for getting user's liked posts
type GetUserLikedPostsOutput struct {
	Posts      []entities.Post
	Error      error
	UserExists bool
}

// GetUserLikedPostsUseCase defines the interface for getting user's liked posts
type GetUserLikedPostsUseCase interface {
	Execute(input GetUserLikedPostsInput) GetUserLikedPostsOutput
}

// getUserLikedPostsUseCaseImpl implements GetUserLikedPostsUseCase
type getUserLikedPostsUseCaseImpl struct {
	userRepo repositories.UserRepository
	postRepo repositories.PostRepository
	likeRepo repositories.LikeRepository
}

// NewGetUserLikedPostsUseCase creates a new get user liked posts use case
func NewGetUserLikedPostsUseCase(
	userRepo repositories.UserRepository,
	postRepo repositories.PostRepository,
	likeRepo repositories.LikeRepository,
) GetUserLikedPostsUseCase {
	return &getUserLikedPostsUseCaseImpl{
		userRepo: userRepo,
		postRepo: postRepo,
		likeRepo: likeRepo,
	}
}

// Execute gets all posts liked by a specific user
func (uc *getUserLikedPostsUseCaseImpl) Execute(input GetUserLikedPostsInput) GetUserLikedPostsOutput {
	// Check if user exists
	if !uc.userRepo.Exists(input.UserID) {
		return GetUserLikedPostsOutput{
			Posts:      nil,
			Error:      nil,
			UserExists: false,
		}
	}

	// Get all likes by user
	likes, err := uc.likeRepo.FindByUserID(input.UserID)
	if err != nil {
		return GetUserLikedPostsOutput{
			Posts:      nil,
			Error:      err,
			UserExists: true,
		}
	}

	// Get the actual posts
	likedPosts := make([]entities.Post, 0)
	for _, like := range likes {
		if post, err := uc.postRepo.FindByID(like.PostID); err == nil {
			likedPosts = append(likedPosts, post)
		}
	}

	return GetUserLikedPostsOutput{
		Posts:      likedPosts,
		Error:      nil,
		UserExists: true,
	}
}
