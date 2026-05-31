package videousecase

import (
	"like-api/internal/domain/repositories"
)

type VideoStats struct {
	VideoID        string  `json:"video_id"`
	Views          int     `json:"views"`
	Likes          int     `json:"likes"`
	EngagementRate float64 `json:"engagement_rate"`
}

type GetVideoStatsInput struct {
	VideoID string
}

type GetVideoStatsOutput struct {
	Stats       VideoStats
	Error       error
	VideoExists bool
}

type GetVideoStatsUseCase interface {
	Execute(input GetVideoStatsInput) GetVideoStatsOutput
}

type getVideoStatsUseCaseImpl struct {
	videoRepo repositories.VideoRepository
}

func NewGetVideoStatsUseCase(videoRepo repositories.VideoRepository) GetVideoStatsUseCase {
	return &getVideoStatsUseCaseImpl{videoRepo: videoRepo}
}

// Execute calculates basic engagement metrics.
//
// Engagement rate = likes / views * 100 (percentage).
// Returns 0 when views == 0 to avoid division by zero.
func (uc *getVideoStatsUseCaseImpl) Execute(input GetVideoStatsInput) GetVideoStatsOutput {
	video, err := uc.videoRepo.FindByID(input.VideoID)
	if err != nil {
		return GetVideoStatsOutput{VideoExists: false}
	}

	var engagementRate float64
	if video.ViewsCount > 0 {
		engagementRate = float64(video.LikesCount) / float64(video.ViewsCount) * 100
	}

	return GetVideoStatsOutput{
		VideoExists: true,
		Stats: VideoStats{
			VideoID:        video.ID,
			Views:          video.ViewsCount,
			Likes:          video.LikesCount,
			EngagementRate: engagementRate,
		},
	}
}
