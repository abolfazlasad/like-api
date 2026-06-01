package videousecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
)

// GetFeedInput holds pagination parameters.
// Cursor is the RFC3339 created_at of the last item from the previous page.
// An empty Cursor means "start from the beginning" (newest items first).
type GetFeedInput struct {
	Cursor string
	Limit  int
}

// GetFeedOutput carries the page of videos plus the cursor for the next page.
// NextCursor is empty when there are no more pages.
type GetFeedOutput struct {
	Videos     []entities.Video
	NextCursor string
	Error      error
}

type GetFeedUseCase interface {
	Execute(input GetFeedInput) GetFeedOutput
}

type getFeedUseCaseImpl struct {
	videoRepo repositories.VideoRepository
}

func NewGetFeedUseCase(videoRepo repositories.VideoRepository) GetFeedUseCase {
	return &getFeedUseCaseImpl{videoRepo: videoRepo}
}

func (uc *getFeedUseCaseImpl) Execute(input GetFeedInput) GetFeedOutput {
	videos, nextCursor, err := uc.videoRepo.FindAll(input.Cursor, input.Limit)
	if err != nil {
		return GetFeedOutput{Error: err}
	}

	return GetFeedOutput{
		Videos:     videos,
		NextCursor: nextCursor,
	}
}
