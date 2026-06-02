package videousecase

import (
	"like-api/internal/domain/entities"
	"like-api/internal/domain/repositories"
	"time"

	"github.com/google/uuid"
)

type LikeVideoInput struct {
	UserID  string
	VideoID string
}

type LikeVideoOutput struct {
	Like         entities.Like
	Error        error
	UserExists   bool
	VideoExists  bool
	AlreadyLiked bool
}

type LikeVideoUseCase interface {
	Execute(input LikeVideoInput) LikeVideoOutput
}

type likeVideoUseCaseImpl struct {
	userRepo  repositories.UserRepository
	videoRepo repositories.VideoRepository
	likeRepo  repositories.LikeRepository
}

func NewLikeVideoUseCase(
	userRepo repositories.UserRepository,
	videoRepo repositories.VideoRepository,
	likeRepo repositories.LikeRepository,
) LikeVideoUseCase {
	return &likeVideoUseCaseImpl{
		userRepo:  userRepo,
		videoRepo: videoRepo,
		likeRepo:  likeRepo,
	}
}

func (uc *likeVideoUseCaseImpl) Execute(input LikeVideoInput) LikeVideoOutput {
	if !uc.userRepo.Exists(input.UserID) {
		return LikeVideoOutput{UserExists: false}
	}

	if !uc.videoRepo.Exists(input.VideoID) {
		return LikeVideoOutput{UserExists: true, VideoExists: false}
	}

	// Re-use the existing LikeRepository but store video_id in post_id column.
	// Trade-off: avoids creating a separate video_likes table while the schema
	// is still evolving.  When the like table grows, split into dedicated tables.
	if uc.likeRepo.Exists(input.UserID, input.VideoID) {
		return LikeVideoOutput{UserExists: true, VideoExists: true, AlreadyLiked: true}
	}

	like := entities.Like{
		ID:        uuid.New().String(),
		UserID:    input.UserID,
		PostID:    input.VideoID,
		CreatedAt: time.Now().UTC().Format(time.RFC3339Nano),
	}

	if err := uc.likeRepo.Create(like); err != nil {
		return LikeVideoOutput{UserExists: true, VideoExists: true, Error: err}
	}

	if err := uc.videoRepo.IncrementLikes(input.VideoID); err != nil {
		return LikeVideoOutput{UserExists: true, VideoExists: true, Error: err}
	}

	return LikeVideoOutput{
		Like:        like,
		UserExists:  true,
		VideoExists: true,
	}
}
