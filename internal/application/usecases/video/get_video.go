package videousecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
)

type GetVideoInput struct {
	VideoID string
}

type GetVideoOutput struct {
	Video       entities.Video
	Error       error
	VideoExists bool
}

type GetVideoUseCase interface {
	Execute(input GetVideoInput) GetVideoOutput
}

type getVideoUseCaseImpl struct {
	videoRepo repositories.VideoRepository
}

func NewGetVideoUseCase(videoRepo repositories.VideoRepository) GetVideoUseCase {
	return &getVideoUseCaseImpl{videoRepo: videoRepo}
}

func (uc *getVideoUseCaseImpl) Execute(input GetVideoInput) GetVideoOutput {
	video, err := uc.videoRepo.FindByID(input.VideoID)
	if err != nil {
		return GetVideoOutput{VideoExists: false}
	}
	return GetVideoOutput{Video: video, VideoExists: true}
}
