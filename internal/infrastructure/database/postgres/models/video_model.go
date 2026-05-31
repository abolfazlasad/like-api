package models

import (
	"like-api/internal/domain/entities"
	"time"
)

// VideoModel is the GORM representation of a video/reel.
// created_at is indexed to support efficient cursor-based pagination.
type VideoModel struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)"`
	UserID      string    `gorm:"not null;index;type:varchar(36)"`
	Title       string    `gorm:"not null;type:varchar(255)"`
	Description string    `gorm:"type:text"`
	VideoURL    string    `gorm:"not null;type:varchar(2048)"`
	LikesCount  int       `gorm:"default:0;not null"`
	ViewsCount  int       `gorm:"default:0;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index"`
}

func (VideoModel) TableName() string { return "videos" }

func (m VideoModel) ToEntity() entities.Video {
	return entities.Video{
		ID:          m.ID,
		UserID:      m.UserID,
		Title:       m.Title,
		Description: m.Description,
		VideoURL:    m.VideoURL,
		LikesCount:  m.LikesCount,
		ViewsCount:  m.ViewsCount,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
	}
}

func VideoModelFromEntity(v entities.Video) VideoModel {
	var createdAt time.Time
	if v.CreatedAt != "" {
		createdAt, _ = time.Parse(time.RFC3339, v.CreatedAt)
	}
	return VideoModel{
		ID:          v.ID,
		UserID:      v.UserID,
		Title:       v.Title,
		Description: v.Description,
		VideoURL:    v.VideoURL,
		LikesCount:  v.LikesCount,
		ViewsCount:  v.ViewsCount,
		CreatedAt:   createdAt,
	}
}
