package videousecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type CreateVideoInput struct {
	UserID      string
	Title       string
	Description string
	VideoURL    string
}

type CreateVideoOutput struct {
	Video entities.Video
	Error error
}

type CreateVideoUseCase interface {
	Execute(input CreateVideoInput) CreateVideoOutput
}

type createVideoUseCaseImpl struct {
	userRepo  repositories.UserRepository
	videoRepo repositories.VideoRepository
}

func NewCreateVideoUseCase(
	userRepo repositories.UserRepository,
	videoRepo repositories.VideoRepository,
) CreateVideoUseCase {
	return &createVideoUseCaseImpl{
		userRepo:  userRepo,
		videoRepo: videoRepo,
	}
}

func (uc *createVideoUseCaseImpl) Execute(input CreateVideoInput) CreateVideoOutput {
	video := entities.Video{
		ID:          uuid.New().String(),
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Description,
		VideoURL:    input.VideoURL,
		LikesCount:  0,
		ViewsCount:  0,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339),
	}

	if err := uc.videoRepo.Create(video); err != nil {
		return CreateVideoOutput{Error: err}
	}

	return CreateVideoOutput{Video: video}
}
