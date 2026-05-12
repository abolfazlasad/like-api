package postusecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
)

// GetPostsInput represents input for get all posts use case
type GetPostsInput struct{}

// GetPostsOutput represents output for get all posts use case
type GetPostsOutput struct {
	Posts []entities.Post
	Error error
}

// GetPostsUseCase defines the interface for getting all posts
type GetPostsUseCase interface {
	Execute(input GetPostsInput) GetPostsOutput
}

// getPostsUseCaseImpl implements GetPostsUseCase
type getPostsUseCaseImpl struct {
	postRepo repositories.PostRepository
}

// NewGetPostsUseCase creates a new get all posts use case
func NewGetPostsUseCase(postRepo repositories.PostRepository) GetPostsUseCase {
	return &getPostsUseCaseImpl{
		postRepo: postRepo,
	}
}

// Execute gets all posts
func (uc *getPostsUseCaseImpl) Execute(input GetPostsInput) GetPostsOutput {
	posts, err := uc.postRepo.FindAll()
	if err != nil {
		return GetPostsOutput{
			Posts: nil,
			Error: err,
		}
	}

	return GetPostsOutput{
		Posts: posts,
		Error: nil,
	}
}
