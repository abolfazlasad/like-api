package videousecase

import (
	"like-api/internal/domain/repositories"
)

type UnlikeVideoInput struct {
	UserID  string
	VideoID string
}

type UnlikeVideoOutput struct {
	Error      error
	LikeExists bool
}

type UnlikeVideoUseCase interface {
	Execute(input UnlikeVideoInput) UnlikeVideoOutput
}

type unlikeVideoUseCaseImpl struct {
	videoRepo repositories.VideoRepository
	likeRepo  repositories.LikeRepository
}

func NewUnlikeVideoUseCase(
	videoRepo repositories.VideoRepository,
	likeRepo repositories.LikeRepository,
) UnlikeVideoUseCase {
	return &unlikeVideoUseCaseImpl{
		videoRepo: videoRepo,
		likeRepo:  likeRepo,
	}
}

func (uc *unlikeVideoUseCaseImpl) Execute(input UnlikeVideoInput) UnlikeVideoOutput {
	if !uc.likeRepo.Exists(input.UserID, input.VideoID) {
		return UnlikeVideoOutput{LikeExists: false}
	}

	if err := uc.likeRepo.Delete(input.UserID, input.VideoID); err != nil {
		return UnlikeVideoOutput{LikeExists: true, Error: err}
	}

	if err := uc.videoRepo.DecrementLikes(input.VideoID); err != nil {
		return UnlikeVideoOutput{LikeExists: true, Error: err}
	}

	return UnlikeVideoOutput{LikeExists: true}
}
