package likeusecase

import (
	"like-api/internal/domain/repositories"
)

// UnlikePostInput represents input for unliking a post
type UnlikePostInput struct {
	UserID string
	PostID string
}

// UnlikePostOutput represents output for unliking a post
type UnlikePostOutput struct {
	Error      error
	LikeExists bool
}

// UnlikePostUseCase defines the interface for unliking a post
type UnlikePostUseCase interface {
	Execute(input UnlikePostInput) UnlikePostOutput
}

// unlikePostUseCaseImpl implements UnlikePostUseCase
type unlikePostUseCaseImpl struct {
	postRepo repositories.PostRepository
	likeRepo repositories.LikeRepository
}

// NewUnlikePostUseCase creates a new unlike post use case
func NewUnlikePostUseCase(
	postRepo repositories.PostRepository,
	likeRepo repositories.LikeRepository,
) UnlikePostUseCase {
	return &unlikePostUseCaseImpl{
		postRepo: postRepo,
		likeRepo: likeRepo,
	}
}

// Execute unlikes a post
func (uc *unlikePostUseCaseImpl) Execute(input UnlikePostInput) UnlikePostOutput {
	// Check if the like exists
	if !uc.likeRepo.Exists(input.UserID, input.PostID) {
		return UnlikePostOutput{
			LikeExists: false,
		}
	}

	// Delete the like
	if err := uc.likeRepo.Delete(input.UserID, input.PostID); err != nil {
		return UnlikePostOutput{
			LikeExists: true,
			Error:      err,
		}
	}

	// atomic decrement — also prevents the value from going below zero
	if err := uc.postRepo.DecrementLikes(input.PostID); err != nil {
		return UnlikePostOutput{LikeExists: true, Error: err}
	}

	return UnlikePostOutput{
		LikeExists: true,
		Error:      nil,
	}
}
