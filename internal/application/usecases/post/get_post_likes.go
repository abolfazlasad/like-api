package postusecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
)

// GetPostLikesInput represents input for getting likes of a post
type GetPostLikesInput struct {
	PostID string
}

// GetPostLikesOutput represents output for getting likes of a post
type GetPostLikesOutput struct {
	Likes      []entities.Like
	Error      error
	PostExists bool
}

// GetPostLikesUseCase defines the interface for getting post likes
type GetPostLikesUseCase interface {
	Execute(input GetPostLikesInput) GetPostLikesOutput
}

// getPostLikesUseCaseImpl implements GetPostLikesUseCase
type getPostLikesUseCaseImpl struct {
	postRepo repositories.PostRepository
	likeRepo repositories.LikeRepository
}

// NewGetPostLikesUseCase creates a new get post likes use case
func NewGetPostLikesUseCase(postRepo repositories.PostRepository, likeRepo repositories.LikeRepository) GetPostLikesUseCase {
	return &getPostLikesUseCaseImpl{
		postRepo: postRepo,
		likeRepo: likeRepo,
	}
}

// Execute gets all likes for a specific post
func (uc *getPostLikesUseCaseImpl) Execute(input GetPostLikesInput) GetPostLikesOutput {
	// Check if post exists
	if !uc.postRepo.Exists(input.PostID) {
		return GetPostLikesOutput{
			Likes:      nil,
			Error:      nil,
			PostExists: false,
		}
	}

	// Get likes for the post
	likes, err := uc.likeRepo.FindByPostID(input.PostID)
	if err != nil {
		return GetPostLikesOutput{
			Likes:      nil,
			Error:      err,
			PostExists: true,
		}
	}

	return GetPostLikesOutput{
		Likes:      likes,
		Error:      nil,
		PostExists: true,
	}
}
